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
	"github.com/finleap-connect/monoskope/pkg/eventsourcing"
	mock_eventsourcing "github.com/finleap-connect/monoskope/test/eventsourcing"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Pkg/Eventsourcing/Messaging/BusRabbmitMQ", func() {
	ctx := context.Background()
	testCount := 0
	eventCounter := 0
	expectedEventType := eventsourcing.EventType("TestEventXYZ")
	expectedAggregateType := eventsourcing.AggregateType("TestAggregateXYZ")
	var consumer eventsourcing.EventBusConsumer
	var publisher eventsourcing.EventBusPublisher
	var eventHandler *mock_eventsourcing.MockEventHandler
	var mockCtrl *gomock.Controller

	createTestEventData := func(something string) eventsourcing.EventData {
		return eventsourcing.ToEventDataFromProto(&testEd.TestEventData{Hello: something})
	}
	createEvent := func() eventsourcing.Event {
		data := createTestEventData(fmt.Sprintf("hello world %v!", eventCounter))
		event := eventsourcing.NewEvent(ctx, expectedEventType, data, time.Now().UTC(), expectedAggregateType, uuid.New(), uint64(eventCounter))
		eventCounter++
		return event
	}

	handleEventNamed := func(done chan interface{}, e eventsourcing.Event, name string) {
		Expect(e.AggregateType()).To(Equal(expectedAggregateType))
		Expect(e.EventType()).To(Equal(expectedEventType))

		testData := new(testEd.TestEventData)
		Expect(e.Data().ToProto(testData)).ToNot(HaveOccurred())

		env.Log.Info("Received event.", "Handler", name, "Event", e.String(), "Data", testData.Hello)
		close(done)
	}

	handleEvent := func(done chan interface{}, e eventsourcing.Event) {
		handleEventNamed(done, e, "default")
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

		// Setup mock event handler
		mockCtrl = gomock.NewController(GinkgoT())
		defer mockCtrl.Finish()
		eventHandler = mock_eventsourcing.NewMockEventHandler(mockCtrl)

		testCount++
	}, 10)

	AfterEach(func() {
		var err error

		err = consumer.Close()
		Expect(err).ToNot(HaveOccurred())

		err = publisher.Close()
		Expect(err).ToNot(HaveOccurred())
	}, 10)

	Context("RabbitMQ can be used to publish and receive events", func() {
		When("normal handler style", func() {
			It("can publish and receive an event matching aggregate type", func() {
				err := consumer.AddHandler(ctx, eventHandler, consumer.Matcher().MatchAggregateType(expectedAggregateType))
				Expect(err).ToNot(HaveOccurred())

				event := createEvent()
				done := make(chan interface{})
				eventHandler.EXPECT().HandleEvent(gomock.Any(), gomock.Any()).Do(func(_ context.Context, e eventsourcing.Event) { handleEvent(done, e) }).Return(nil)
				err = publisher.PublishEvent(ctx, event)
				Expect(err).ToNot(HaveOccurred())
				Eventually(done, 30).Should(BeClosed())
			})
			It("can publish and receive an event matching event type", func() {
				err := consumer.AddHandler(ctx, eventHandler, consumer.Matcher().MatchAggregateType(eventsourcing.AggregateType(expectedEventType)))
				Expect(err).ToNot(HaveOccurred())

				event := createEvent()
				done := make(chan interface{})
				eventHandler.EXPECT().HandleEvent(gomock.Any(), gomock.Any()).Do(func(_ context.Context, e eventsourcing.Event) { handleEvent(done, e) }).Return(nil)
				err = publisher.PublishEvent(ctx, event)
				Expect(err).ToNot(HaveOccurred())
				Eventually(done, 30).Should(BeClosed())
			})
			It("can publish and receive an event matching aggregate and event type", func() {
				err := consumer.AddHandler(ctx, eventHandler, consumer.Matcher().MatchAggregateType(expectedAggregateType), consumer.Matcher().MatchAggregateType(eventsourcing.AggregateType(expectedEventType)))
				Expect(err).ToNot(HaveOccurred())

				event := createEvent()
				done := make(chan interface{})
				eventHandler.EXPECT().HandleEvent(gomock.Any(), gomock.Any()).Do(func(_ context.Context, e eventsourcing.Event) { handleEvent(done, e) }).Return(nil)
				err = publisher.PublishEvent(ctx, event)
				Expect(err).ToNot(HaveOccurred())
				Eventually(done, 30).Should(BeClosed())
			})
		})
		When("normal worker style", func() {
			It("can publish and receive an event round robin among consumers", func() {
				eventHandlerA := mock_eventsourcing.NewMockEventHandler(mockCtrl)
				eventHandlerB := mock_eventsourcing.NewMockEventHandler(mockCtrl)

				// Add two workers
				err := consumer.AddWorker(ctx, eventHandlerA, "my-worker-group", consumer.Matcher().Any())
				Expect(err).ToNot(HaveOccurred())
				err = consumer.AddWorker(ctx, eventHandlerB, "my-worker-group", consumer.Matcher().Any())
				Expect(err).ToNot(HaveOccurred())

				eventA := createEvent()
				eventB := createEvent()
				doneA := make(chan interface{})
				doneB := make(chan interface{})

				eventHandlerA.EXPECT().HandleEvent(gomock.Any(), gomock.Any()).Do(func(_ context.Context, e eventsourcing.Event) { handleEventNamed(doneA, e, "a") }).Return(nil).AnyTimes()
				eventHandlerB.EXPECT().HandleEvent(gomock.Any(), gomock.Any()).Do(func(_ context.Context, e eventsourcing.Event) { handleEventNamed(doneB, e, "b") }).Return(nil).AnyTimes()

				err = publisher.PublishEvent(ctx, eventA)
				Expect(err).ToNot(HaveOccurred())

				err = publisher.PublishEvent(ctx, eventB)
				Expect(err).ToNot(HaveOccurred())

				Eventually(doneA, 30).Should(BeClosed())
				Eventually(doneB, 30).Should(BeClosed())
			})
		})
	})
})
