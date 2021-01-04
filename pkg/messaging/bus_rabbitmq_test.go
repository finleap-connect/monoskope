package messaging

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/storage"
)

var _ = Describe("messaging/rabbitmq", func() {
	var wg sync.WaitGroup
	ctx := context.Background()
	eventCounter := 0

	createEvent := func() storage.Event {
		eventType := storage.EventType("TestEvent")
		aggregateType := storage.AggregateType("TestAggregate")
		data := storage.EventData(fmt.Sprintf("test-%v", eventCounter))
		event := storage.NewEvent(eventType, data, time.Now().UTC(), aggregateType, uuid.New(), uint64(eventCounter))
		eventCounter++
		return event
	}

	publishEvent := func(event storage.Event) {
		err := env.Publisher.PublishEvent(ctx, event)
		Expect(err).ToNot(HaveOccurred())
	}

	receiveEvent := func(receiveChan <-chan storage.Event, event storage.Event) {
		defer GinkgoRecover()
		defer wg.Done()
		select {
		case eventFromBus := <-receiveChan:
			Expect(eventFromBus).ToNot(BeNil())
			Expect(eventFromBus).To(Equal(event))
		case <-time.After(10 * time.Second):
			Expect(fmt.Errorf("timeout waiting for receiving event")).ToNot(HaveOccurred())
		}
	}

	createReceiver := func(matchers ...EventMatcher) <-chan storage.Event {
		receiveChan := make(chan storage.Event)
		receiver := func(e storage.Event) error {
			receiveChan <- e
			return nil
		}
		err := env.Consumer.AddReceiver(receiver, matchers...)
		Expect(err).ToNot(HaveOccurred())
		return receiveChan
	}

	failOnErr := func() {
		defer GinkgoRecover()
		env.Consumer.AddErrorHandler(func(mbe MessageBusError) {
			Expect(mbe).ToNot(HaveOccurred())
		})
	}

	testPubSub := func(matchers ...EventMatcher) {
		go failOnErr()

		recChanA := createReceiver(matchers...)
		recChanB := createReceiver(matchers...)
		event := createEvent()

		wg.Add(2)
		go receiveEvent(recChanA, event)
		go receiveEvent(recChanB, event)
		publishEvent(event)
		wg.Wait()
	}

	It("can publish and receive an event", func() {
		testPubSub(env.Consumer.Matcher().Any())
	})
	It("can publish and receive an event matching aggregate type", func() {
		event := createEvent()
		testPubSub(env.Consumer.Matcher().MatchAggregateType(event.AggregateType()))
	})
	It("can publish and receive an event matching event type", func() {
		event := createEvent()
		testPubSub(env.Consumer.Matcher().MatchEventType(event.EventType()))
	})
	It("can publish and receive an event matching aggregate type and event type", func() {
		event := createEvent()
		testPubSub(env.Consumer.Matcher().MatchAggregateType(event.AggregateType()).MatchEventType(event.EventType()))
	})
})
