package messaging

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/storage"
)

var _ = Describe("messaging/rabbitmq", func() {
	var consumer EventBusConsumer
	var publisher EventBusPublisher
	ctx := context.Background()
	eventCounter := 0
	testCount := 0

	createEvent := func() storage.Event {
		eventType := storage.EventType("TestEvent")
		aggregateType := storage.AggregateType("TestAggregate")
		data := storage.EventData(fmt.Sprintf("test-%v", eventCounter))
		event := storage.NewEvent(eventType, data, time.Now().UTC(), aggregateType, uuid.New(), uint64(eventCounter))
		eventCounter++
		return event
	}

	publishEvent := func(event storage.Event) {
		defer GinkgoRecover()
		ctxWithTimeout, cancelFunc := context.WithTimeout(ctx, 10*time.Second)
		defer cancelFunc()
		err := publisher.PublishEvent(ctxWithTimeout, event)
		Expect(err).ToNot(HaveOccurred())
	}

	receiveEvent := func(receiveChan <-chan storage.Event, event storage.Event) {
		select {
		case eventFromBus := <-receiveChan:
			env.Log.Info("Received event.")
			Expect(eventFromBus).ToNot(BeNil())
			Expect(eventFromBus).To(Equal(event))
		case <-time.After(5 * time.Second):
			env.Log.Info("Timeout when receiving event.")
			Expect(fmt.Errorf("timeout waiting for receiving event")).ToNot(HaveOccurred())
		}
	}

	createReceiver := func(matchers ...EventMatcher) chan storage.Event {
		receiveChan := make(chan storage.Event)
		receiver := func(e storage.Event) error {
			env.Log.Info("Received event.")
			receiveChan <- e
			return nil
		}

		ctxWithTimeout, cancelFunc := context.WithTimeout(ctx, 10*time.Second)
		defer cancelFunc()
		err := consumer.AddReceiver(ctxWithTimeout, receiver, matchers...)
		Expect(err).ToNot(HaveOccurred())

		return receiveChan
	}

	testPubSub := func(eventCount int, matchers ...EventMatcher) {
		recChanA := createReceiver(matchers...)
		defer close(recChanA)
		recChanB := createReceiver(matchers...)
		defer close(recChanB)

		for i := 0; i < eventCount; i++ {
			event := createEvent()
			go publishEvent(event)
			receiveEvent(recChanA, event)
			receiveEvent(recChanB, event)
		}
	}

	BeforeEach(func() {
		var err error

		conf := NewRabbitEventBusConfig(fmt.Sprintf("test-%v", testCount), env.amqpURL)

		// init publisher
		publisher, err = NewRabbitEventBusPublisher(conf)
		Expect(err).ToNot(HaveOccurred())
		ctxWithTimeout, cancelFunc := context.WithTimeout(ctx, 30*time.Second)
		defer cancelFunc()
		err = publisher.Connect(ctxWithTimeout)
		Expect(err).ToNot(HaveOccurred())

		// init consumer
		consumer, err = NewRabbitEventBusConsumer(conf)
		Expect(err).ToNot(HaveOccurred())
		ctxWithTimeout, cancelFunc = context.WithTimeout(ctx, 30*time.Second)
		defer cancelFunc()
		err = consumer.Connect(ctxWithTimeout)
		Expect(err).ToNot(HaveOccurred())

		testCount++
	})
	AfterEach(func() {
		var err error

		err = consumer.Close()
		Expect(err).ToNot(HaveOccurred())

		err = publisher.Close()
		Expect(err).ToNot(HaveOccurred())
	})
	It("can publish and receive events", func() {
		testPubSub(3, consumer.Matcher().Any())
	})
	It("can publish and receive an event matching aggregate type", func() {
		event := createEvent()
		testPubSub(1, consumer.Matcher().MatchAggregateType(event.AggregateType()))
	})
	It("can publish and receive an event matching event type", func() {
		event := createEvent()
		testPubSub(1, consumer.Matcher().MatchEventType(event.EventType()))
	})
	It("can publish and receive an event matching aggregate type and event type", func() {
		event := createEvent()
		testPubSub(1, consumer.Matcher().MatchAggregateType(event.AggregateType()).MatchEventType(event.EventType()))
	})
})
