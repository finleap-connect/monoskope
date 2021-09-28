// Copyright 2021 Monoskope Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package messaging

import (
	"context"
	"fmt"
	"time"

	testEd "github.com/finleap-connect/monoskope/pkg/api/eventsourcing/eventdata"
	evs "github.com/finleap-connect/monoskope/pkg/eventsourcing"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type testEventHandler struct {
	event evs.Event
}

func (t testEventHandler) HandleEvent(ctx context.Context, e evs.Event) error {
	defer GinkgoRecover()
	env.Log.Info("Received event.")
	Expect(e).ToNot(BeNil())
	Expect(e.AggregateID()).To(BeEquivalentTo(t.event.AggregateID()))
	Expect(e.AggregateType()).To(BeEquivalentTo(t.event.AggregateType()))
	Expect(e.AggregateVersion()).To(BeEquivalentTo(t.event.AggregateVersion()))
	return nil
}

var _ = Describe("messaging/rabbitmq", func() {
	var consumer evs.EventBusConsumer
	var publisher evs.EventBusPublisher
	ctx := context.Background()
	eventCounter := 0
	testCount := 0

	createTestEventData := func(something string) evs.EventData {
		return evs.ToEventDataFromProto(&testEd.TestEventData{Hello: something})
	}
	createEvent := func() evs.Event {
		eventType := evs.EventType("TestEvent")
		aggregateType := evs.AggregateType("TestAggregate")
		data := createTestEventData("world!")
		event := evs.NewEvent(ctx, eventType, data, time.Now().UTC(), aggregateType, uuid.New(), uint64(eventCounter))
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

		conf, err := NewRabbitEventBusConfig(fmt.Sprintf("test-%v", testCount), env.AmqpURL, "")
		Expect(err).ToNot(HaveOccurred())

		// init publisher
		publisher, err = NewRabbitEventBusPublisher(conf)
		Expect(err).ToNot(HaveOccurred())

		// init consumer
		consumer, err = NewRabbitEventBusConsumer(conf)
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
