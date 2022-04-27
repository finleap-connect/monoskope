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

package storage

import (
	"context"
	"io"
	"time"

	testEd "github.com/finleap-connect/monoskope/pkg/api/eventsourcing/eventdata"
	evs "github.com/finleap-connect/monoskope/pkg/eventsourcing"
	"github.com/finleap-connect/monoskope/pkg/eventsourcing/errors"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type metadataVal struct {
	Val string
}

var _ = Describe("storage/postgres", func() {
	var userInformationKey = "userInformationKey"

	manager := evs.NewMetadataManagerFromContext(context.Background())
	err := manager.SetObject(userInformationKey, &metadataVal{Val: "admin"})
	Expect(err).ToNot(HaveOccurred())
	ctx := manager.GetContext()

	var es *postgresEventStore

	clearEs := func(es *postgresEventStore) {
		err := es.clear(ctx)
		Expect(err).ToNot(HaveOccurred())
	}

	createTestEventStore := func() *postgresEventStore {
		es, err := NewPostgresEventStore(env.postgresStoreConfig)
		Expect(err).ToNot(HaveOccurred())
		Expect(es).ToNot(BeNil())
		return es.(*postgresEventStore)
	}

	now := func() time.Time {
		return time.Now().UTC()
	}

	createTestEventData := func(something string) evs.EventData {
		return evs.ToEventDataFromProto(&testEd.TestEventData{Hello: something})
	}

	createTestEvents := func() []evs.Event {
		aggregateId := uuid.New()
		return []evs.Event{
			evs.NewEvent(ctx, testEventCreated, createTestEventData("create"), now(), testAggregate, aggregateId, 0),
			evs.NewEvent(ctx, testEventChanged, createTestEventData("change"), now(), testAggregate, aggregateId, 1),
			evs.NewEvent(ctx, testEventDeleted, createTestEventData("delete"), now(), testAggregate, aggregateId, 2),
		}
	}

	BeforeEach(func() {
		var err error
		es = createTestEventStore()

		ctxWithTimeout, cancelFunc := context.WithTimeout(ctx, 10*time.Second)
		defer cancelFunc()
		err = es.Open(ctxWithTimeout)
		Expect(err).ToNot(HaveOccurred())
	})
	AfterEach(func() {
		clearEs(es)
		err := es.Close()
		Expect(err).ToNot(HaveOccurred())
	})
	It("can append new events to the store", func() {
		err := es.Save(ctx, createTestEvents())
		Expect(err).ToNot(HaveOccurred())
	})
	It("fails to append new events to the store when they are not of the same aggregate type", func() {
		aggregateId := uuid.New()
		err := es.Save(ctx, []evs.Event{
			evs.NewEvent(ctx, testEventCreated, createTestEventData("create"), now(), testAggregate, aggregateId, 0),
			evs.NewEvent(ctx, testEventChanged, createTestEventData("change"), now(), testAggregateExtended, aggregateId, 1),
		})
		Expect(err).To(HaveOccurred())
		Expect(err).To(Equal(errors.ErrInvalidAggregateType))
	})
	It("fails to append new events to the store when they are not in the right aggregate version order", func() {
		aggregateId := uuid.New()
		err := es.Save(ctx, []evs.Event{
			evs.NewEvent(ctx, testEventCreated, createTestEventData("create"), now(), testAggregate, aggregateId, 0),
			evs.NewEvent(ctx, testEventChanged, createTestEventData("change"), now(), testAggregate, aggregateId, 2),
		})
		Expect(err).To(HaveOccurred())
		Expect(err).To(Equal(errors.ErrIncorrectAggregateVersion))
	})
	It("fails to append new events to the store when the aggregate version does already exist", func() {
		aggregateId := uuid.New()
		err := es.Save(ctx, []evs.Event{
			evs.NewEvent(ctx, testEventCreated, createTestEventData("create"), now(), testAggregate, aggregateId, 0),
			evs.NewEvent(ctx, testEventChanged, createTestEventData("change"), now(), testAggregate, aggregateId, 1),
		})
		Expect(err).ToNot(HaveOccurred())

		err = es.Save(ctx, []evs.Event{
			evs.NewEvent(ctx, testEventChanged, createTestEventData("change"), now(), testAggregate, aggregateId, 1),
		})
		Expect(err).To(HaveOccurred())
		Expect(err).To(Equal(errors.ErrAggregateVersionAlreadyExists))
	})
	It("can load events from the store", func() {
		// append some events to load later
		events := createTestEvents()
		err := es.Save(ctx, events)
		Expect(err).ToNot(HaveOccurred())
		expectedEventCount := len(events)

		// append some additional events
		events = createTestEvents()
		err = es.Save(ctx, events)
		Expect(err).ToNot(HaveOccurred())

		aggregateId := events[0].AggregateID()
		eventStream, err := es.Load(ctx, &evs.StoreQuery{
			AggregateId: &aggregateId,
		})
		Expect(err).ToNot(HaveOccurred())

		var storeEvents []evs.Event
		for {
			event, err := eventStream.Receive()
			if err == io.EOF {
				break
			}
			Expect(err).ToNot(HaveOccurred())
			storeEvents = append(storeEvents, event)
		}

		Expect(storeEvents).ToNot(BeNil())
		Expect(storeEvents).ToNot(BeEmpty())
		Expect(len(storeEvents)).To(BeNumerically("==", expectedEventCount))

		valResult := &metadataVal{}
		err = evs.NewMetadataManagerFromContext(context.Background()).SetMetadata(storeEvents[0].Metadata()).GetObject(userInformationKey, valResult)
		Expect(err).ToNot(HaveOccurred())
		Expect(valResult.Val).To(Equal("admin"))
	})
	It("can load events from the store by concatenating the filters with the logical OR", func() {
		events := createTestEvents()
		err := es.Save(ctx, events)
		Expect(err).ToNot(HaveOccurred())
		otherAggregateId := uuid.New()
		otherEvents := []evs.Event{
			evs.NewEvent(ctx, testEventCreated, createTestEventData("create"), now(), testAggregate, otherAggregateId, 0),
		}
		err = es.Save(ctx, otherEvents)
		Expect(err).ToNot(HaveOccurred())
		expectedEventCount := len(events) + len(otherEvents)

		aggregateId := events[0].AggregateID()
		eventStream, err := es.LoadOr(ctx, []*evs.StoreQuery{
			{AggregateId: &aggregateId},
			{AggregateId: &otherAggregateId},
		})
		Expect(err).ToNot(HaveOccurred())

		var storeEvents []evs.Event
		for {
			event, err := eventStream.Receive()
			if err == io.EOF {
				break
			}
			Expect(err).ToNot(HaveOccurred())
			storeEvents = append(storeEvents, event)
		}

		Expect(len(storeEvents)).To(BeNumerically("==", expectedEventCount))
	})
	It("can filter events to load from the store by aggregate type", func() {
		ev := createTestEvents()
		err := es.Save(ctx, ev)
		Expect(err).ToNot(HaveOccurred())

		aggregateType := testAggregate
		eventStream, err := es.Load(ctx, &evs.StoreQuery{
			AggregateType: &aggregateType,
		})
		Expect(err).ToNot(HaveOccurred())

		var storeEvents []evs.Event
		for {
			event, err := eventStream.Receive()
			if err == io.EOF {
				break
			}
			Expect(err).ToNot(HaveOccurred())
			storeEvents = append(storeEvents, event)
		}

		Expect(storeEvents).ToNot(BeNil())
		Expect(storeEvents).ToNot(BeEmpty())
		Expect(len(storeEvents)).To(BeNumerically("==", len(ev)))
	})
	It("can filter events to load from the store by aggregate version", func() {
		events := createTestEvents()
		err := es.Save(ctx, events)
		Expect(err).ToNot(HaveOccurred())

		minVersion := uint64(1)
		maxVersion := uint64(1)
		eventStream, err := es.Load(ctx, &evs.StoreQuery{
			MinVersion: &minVersion,
		})
		Expect(err).ToNot(HaveOccurred())

		var storeEvents []evs.Event
		for {
			event, err := eventStream.Receive()
			if err == io.EOF {
				break
			}
			Expect(err).ToNot(HaveOccurred())
			storeEvents = append(storeEvents, event)
		}

		Expect(storeEvents).ToNot(BeNil())
		Expect(storeEvents).ToNot(BeEmpty())

		for _, ev := range storeEvents {
			Expect(ev.AggregateVersion()).To(BeNumerically(">=", 1))
		}

		eventStream, err = es.Load(ctx, &evs.StoreQuery{
			MaxVersion: &maxVersion,
		})
		Expect(err).ToNot(HaveOccurred())

		storeEvents = make([]evs.Event, 0)
		for {
			event, err := eventStream.Receive()
			if err == io.EOF {
				break
			}
			Expect(err).ToNot(HaveOccurred())
			storeEvents = append(storeEvents, event)
		}

		Expect(storeEvents).ToNot(BeNil())
		Expect(storeEvents).ToNot(BeEmpty())

		for _, ev := range storeEvents {
			Expect(ev.AggregateVersion()).To(BeNumerically("<=", 1))
		}

		eventStream, err = es.Load(ctx, &evs.StoreQuery{
			MinVersion: &minVersion,
			MaxVersion: &maxVersion,
		})
		Expect(err).ToNot(HaveOccurred())

		storeEvents = make([]evs.Event, 0)
		for {
			event, err := eventStream.Receive()
			if err == io.EOF {
				break
			}
			Expect(err).ToNot(HaveOccurred())
			storeEvents = append(storeEvents, event)
		}

		Expect(storeEvents).ToNot(BeNil())
		Expect(storeEvents).ToNot(BeEmpty())

		for _, ev := range storeEvents {
			Expect(ev.AggregateVersion()).To(BeNumerically("==", 1))
		}
	})
})
