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
		select {
		case eventFromBus := <-receiveChan:
			Expect(eventFromBus).ToNot(BeNil())
			Expect(eventFromBus).To(Equal(event))
		case <-time.After(10 * time.Second):
			Expect(fmt.Errorf("timeout waiting for receiving event")).ToNot(HaveOccurred())
		}
	}

	createReceiver := func() (<-chan storage.Event, EventReceiver) {
		receiveChan := make(chan storage.Event)
		receiver := func(e storage.Event) error {
			receiveChan <- e
			return nil
		}
		return receiveChan, receiver
	}
	It("can publish and receive an event", func() {
		receiveChan, receiver := createReceiver()
		err := env.Consumer.AddReceiver(receiver, env.Consumer.Matcher().Any())
		Expect(err).ToNot(HaveOccurred())

		event := createEvent()
		publishEvent(event)
		receiveEvent(receiveChan, event)
	})
	It("can publish and receive an event matching aggregate type", func() {
		event := createEvent()

		receiveChan, receiver := createReceiver()
		err := env.Consumer.AddReceiver(receiver, env.Consumer.Matcher().MatchAggregateType(event.AggregateType()))
		Expect(err).ToNot(HaveOccurred())

		publishEvent(event)
		receiveEvent(receiveChan, event)
	})
	It("can publish and receive an event matching event type", func() {
		event := createEvent()

		receiveChan, receiver := createReceiver()
		err := env.Consumer.AddReceiver(receiver, env.Consumer.Matcher().MatchEventType(event.EventType()))
		Expect(err).ToNot(HaveOccurred())

		publishEvent(event)
		receiveEvent(receiveChan, event)
	})
	It("can publish and receive an event matching aggregate type and event type", func() {
		event := createEvent()

		receiveChan, receiver := createReceiver()
		err := env.Consumer.AddReceiver(receiver, env.Consumer.Matcher().MatchAggregateType(event.AggregateType()).MatchEventType(event.EventType()))
		Expect(err).ToNot(HaveOccurred())

		publishEvent(event)
		receiveEvent(receiveChan, event)
	})
})
