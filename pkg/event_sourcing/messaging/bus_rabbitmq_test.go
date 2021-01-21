package messaging

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	evs "gitlab.figo.systems/platform/monoskope/monoskope/pkg/event_sourcing"
)

var _ = Describe("messaging/rabbitmq", func() {
	var consumer EventBusConsumer
	var publisher EventBusPublisher
	ctx := context.Background()
	eventCounter := 0
	testCount := 0

	createEvent := func() evs.Event {
		eventType := evs.EventType("TestEvent")
		aggregateType := evs.AggregateType("TestAggregate")
		data := evs.NewEventData()
		event := evs.NewEvent(eventType, data, time.Now().UTC(), aggregateType, uuid.New(), uint64(eventCounter))
		eventCounter++
		return event
	}

	publishEvent := func(event evs.Event) {
		ctxWithTimeout, cancel := context.WithTimeout(ctx, 20*time.Second)
		defer cancel()
		err := publisher.PublishEvent(ctxWithTimeout, event)
		Expect(err).ToNot(HaveOccurred())
	}

	createReceiver := func(event evs.Event, matchers ...EventMatcher) {
		receiver := func(e evs.Event) (err error) {
			defer ginkgo.GinkgoRecover()
			env.Log.Info("Received event.")
			Expect(e).ToNot(BeNil())
			Expect(e).To(Equal(event))
			return nil
		}

		ctxWithTimeout, cancel := context.WithTimeout(ctx, 20*time.Second)
		defer cancel()
		err := consumer.AddReceiver(ctxWithTimeout, receiver, matchers...)
		Expect(err).ToNot(HaveOccurred())
	}

	testPubSub := func(matchers ...EventMatcher) {
		event := createEvent()
		createReceiver(event, matchers...)
		publishEvent(event)
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
	}, 10)
	AfterEach(func() {
		var err error

		err = consumer.Close()
		Expect(err).ToNot(HaveOccurred())

		err = publisher.Close()
		Expect(err).ToNot(HaveOccurred())
	}, 10)
	It("can publish and receive events", func() {
		testPubSub(consumer.Matcher().Any())
	})
	It("can publish and receive an event matching aggregate type", func() {
		event := createEvent()
		testPubSub(consumer.Matcher().MatchAggregateType(event.AggregateType()))
	})
	It("can publish and receive an event matching event type", func() {
		event := createEvent()
		testPubSub(consumer.Matcher().MatchEventType(event.EventType()))
	})
	It("can publish and receive an event matching aggregate type and event type", func() {
		event := createEvent()
		testPubSub(consumer.Matcher().MatchAggregateType(event.AggregateType()).MatchEventType(event.EventType()))
	})
})
