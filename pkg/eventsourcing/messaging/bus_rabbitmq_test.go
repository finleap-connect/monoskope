package messaging

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/eventdata/test"
	evs "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing"
)

type testEventHandler struct {
	event evs.Event
}

func (t testEventHandler) HandleEvent(ctx context.Context, e evs.Event) error {
	defer GinkgoRecover()
	env.Log.Info("Received event.")
	Expect(e).ToNot(BeNil())
	Expect(e).To(Equal(t.event))
	return nil
}

var _ = Describe("messaging/rabbitmq", func() {
	var consumer evs.EventBusConsumer
	var publisher evs.EventBusPublisher
	ctx := context.Background()
	eventCounter := 0
	testCount := 0

	createTestEventData := func(something string) evs.EventData {
		ed, err := evs.ToEventDataFromProto(&test.TestEventData{Hello: something})
		Expect(err).ToNot(HaveOccurred())
		return ed
	}
	createEvent := func() evs.Event {
		eventType := evs.EventType("TestEvent")
		aggregateType := evs.AggregateType("TestAggregate")
		data := createTestEventData("world!")
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

	createHandler := func(event evs.Event, matchers ...evs.EventMatcher) {
		handler := &testEventHandler{event: event}
		ctxWithTimeout, cancel := context.WithTimeout(ctx, 20*time.Second)
		defer cancel()
		err := consumer.AddHandler(ctxWithTimeout, handler, matchers...)
		Expect(err).ToNot(HaveOccurred())
	}

	testPubSub := func(matchers ...evs.EventMatcher) {
		event := createEvent()
		createHandler(event, matchers...)
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
