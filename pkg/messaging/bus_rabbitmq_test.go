package messaging

import (
	"context"
	"time"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/storage"
)

var _ = Describe("messaging/rabbitmq", func() {
	ctx := context.Background()

	createEvent := func() storage.Event {
		return storage.NewEvent(storage.EventType("TestEvent"), storage.EventData("test"), time.Now().UTC(), storage.AggregateType("TestAggregate"), uuid.New(), 0)
	}

	publishEvent := func(event storage.Event) {
		err := env.Publisher.PublishEvent(ctx, event)
		Expect(err).ToNot(HaveOccurred())
	}

	createReceiver := func() (chan storage.Event, EventReceiver) {
		eventsFromBus := make(chan storage.Event)
		receiver := func(e storage.Event) error {
			Expect(e).NotTo(BeNil())
			eventsFromBus <- e
			return nil
		}
		return eventsFromBus, receiver
	}

	It("can publish an event", func() {
		event := createEvent()
		publishEvent(event)
	})
	It("can publish and receive an event", func() {
		eventsFromBus, receiver := createReceiver()
		err := env.Consumer.AddReceiver(receiver, env.Consumer.Matcher().Any())
		Expect(err).ToNot(HaveOccurred())

		event := createEvent()
		publishEvent(event)
		fromBus := <-eventsFromBus
		Expect(fromBus).ToNot(BeNil())
		Expect(fromBus).To(Equal(event))
	})
	It("can publish and receive an event matching aggregate type", func() {
		event := createEvent()

		eventsFromBus, receiver := createReceiver()
		err := env.Consumer.AddReceiver(receiver, env.Consumer.Matcher().MatchAggregateType(event.AggregateType()))
		Expect(err).ToNot(HaveOccurred())

		publishEvent(event)
		fromBus := <-eventsFromBus
		Expect(fromBus).ToNot(BeNil())
		Expect(fromBus).To(Equal(event))
	})
	It("can publish and receive an event matching event type", func() {
		event := createEvent()

		eventsFromBus, receiver := createReceiver()
		err := env.Consumer.AddReceiver(receiver, env.Consumer.Matcher().MatchEventType(event.EventType()))
		Expect(err).ToNot(HaveOccurred())

		publishEvent(event)
		fromBus := <-eventsFromBus
		Expect(fromBus).ToNot(BeNil())
		Expect(fromBus).To(Equal(event))
	})
	It("can publish and receive an event matching aggregate type and event type", func() {
		event := createEvent()

		eventsFromBus, receiver := createReceiver()
		err := env.Consumer.AddReceiver(receiver, env.Consumer.Matcher().MatchAggregateType(event.AggregateType()).MatchEventType(event.EventType()))
		Expect(err).ToNot(HaveOccurred())

		publishEvent(event)
		fromBus := <-eventsFromBus
		Expect(fromBus).ToNot(BeNil())
		Expect(fromBus).To(Equal(event))
	})
})
