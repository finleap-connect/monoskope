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

	It("can publish an event", func() {
		event := storage.NewEvent(storage.EventType("TestEvent"), storage.EventData("test"), time.Now().UTC(), storage.AggregateType("TestAggregate"), uuid.New(), 0)
		err := env.Publisher.PublishEvent(ctx, event)
		Expect(err).ToNot(HaveOccurred())
	})
	It("can publish and receive an event", func() {
		event := storage.NewEvent(storage.EventType("TestEvent"), storage.EventData("test"), time.Now().UTC(), storage.AggregateType("TestAggregate"), uuid.New(), 0)
		eventsFromBus := make(chan storage.Event)

		err := env.Consumer.AddReceiver(func(e storage.Event) error {
			Expect(e).NotTo(BeNil())
			Expect(event).To(Equal(e))
			eventsFromBus <- e
			return nil
		}, env.Consumer.Matcher().Any())
		Expect(err).ToNot(HaveOccurred())

		err = env.Publisher.PublishEvent(ctx, event)
		Expect(err).ToNot(HaveOccurred())

		fromBus := <-eventsFromBus
		Expect(fromBus).ToNot(BeNil())
		Expect(fromBus).To(Equal(event))
	})
})
