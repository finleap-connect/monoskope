package storage

import (
	"context"
	"time"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/eventdata/test"
	evs "gitlab.figo.systems/platform/monoskope/monoskope/pkg/event_sourcing"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/event_sourcing/errors"
)

var _ = Describe("storage/inmemory", func() {
	ctx := context.Background()

	now := func() time.Time {
		return time.Now().UTC()
	}
	clearInMemoryEs := func(es *InMemoryEventStore) {
		es.clear(ctx)
	}
	createInMemoryTestEventStore := func() *InMemoryEventStore {
		es := NewInMemoryEventStore()
		Expect(es).ToNot(BeNil())
		return es.(*InMemoryEventStore)
	}
	createTestEventData := func(something string) evs.EventData {
		ed, err := evs.ToEventDataFromProto(&test.TestEventData{Hello: something})
		Expect(err).ToNot(HaveOccurred())
		return ed
	}
	createTestEvents := func() []evs.Event {
		aggregateId := uuid.New()

		return []evs.Event{
			evs.NewEvent(evs.EventType(testEventCreated), createTestEventData("create"), now(), evs.AggregateType(testAggregate), aggregateId, 0),
			evs.NewEvent(evs.EventType(testEventChanged), createTestEventData("change"), now(), evs.AggregateType(testAggregate), aggregateId, 1),
			evs.NewEvent(evs.EventType(testEventDeleted), createTestEventData("delete"), now(), evs.AggregateType(testAggregate), aggregateId, 2),
		}
	}

	It("can create new event store", func() {
		_ = createInMemoryTestEventStore()
	})
	It("can append new events to the store", func() {
		es := createInMemoryTestEventStore()
		defer clearInMemoryEs(es)

		err := es.Save(ctx, createTestEvents())
		Expect(err).ToNot(HaveOccurred())
	})
	It("fails to append new events to the store when they are not of the same aggregate type", func() {
		es := createInMemoryTestEventStore()
		defer clearInMemoryEs(es)

		aggregateId := uuid.New()
		err := es.Save(ctx, []evs.Event{
			evs.NewEvent(testEventCreated, createTestEventData("create"), now(), testAggregate, aggregateId, 0),
			evs.NewEvent(testEventChanged, createTestEventData("change"), now(), testAggregateExtended, aggregateId, 1),
		})
		Expect(err).To(HaveOccurred())
		Expect(err).To(Equal(errors.ErrInvalidAggregateType))
	})
	It("fails to append new events to the store when they are not in the right aggregate version order", func() {
		es := createInMemoryTestEventStore()
		defer clearInMemoryEs(es)

		aggregateId := uuid.New()
		err := es.Save(ctx, []evs.Event{
			evs.NewEvent(testEventCreated, createTestEventData("create"), now(), testAggregate, aggregateId, 0),
			evs.NewEvent(testEventChanged, createTestEventData("change"), now(), testAggregate, aggregateId, 2),
		})
		Expect(err).To(HaveOccurred())
		Expect(err).To(Equal(errors.ErrIncorrectAggregateVersion))

	})
	It("fails to append new events to the store when the aggregate version does already exist", func() {
		es := createInMemoryTestEventStore()
		defer clearInMemoryEs(es)

		aggregateId := uuid.New()
		err := es.Save(ctx, []evs.Event{
			evs.NewEvent(testEventCreated, createTestEventData("create"), now(), testAggregate, aggregateId, 0),
			evs.NewEvent(testEventChanged, createTestEventData("change"), now(), testAggregate, aggregateId, 1),
		})
		Expect(err).ToNot(HaveOccurred())

		err = es.Save(ctx, []evs.Event{
			evs.NewEvent(testEventChanged, createTestEventData("change"), now(), testAggregate, aggregateId, 1),
		})
		Expect(err).To(HaveOccurred())
		Expect(err).To(Equal(errors.ErrAggregateVersionAlreadyExists))
	})
	It("can load events from the store", func() {
		es := createInMemoryTestEventStore()
		defer clearInMemoryEs(es)

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
		storeEvents, err := es.Load(ctx, &evs.StoreQuery{
			AggregateId: &aggregateId,
		})
		Expect(err).ToNot(HaveOccurred())
		Expect(storeEvents).ToNot(BeNil())
		Expect(storeEvents).ToNot(BeEmpty())
		Expect(len(storeEvents)).To(BeNumerically("==", expectedEventCount))
	})
	It("can filter events to load from the store by aggregate type", func() {
		es := createInMemoryTestEventStore()
		defer clearInMemoryEs(es)

		ev := createTestEvents()
		err := es.Save(ctx, ev)
		Expect(err).ToNot(HaveOccurred())

		aggregateType := evs.AggregateType(testAggregate)
		storeEvents, err := es.Load(ctx, &evs.StoreQuery{
			AggregateType: &aggregateType,
		})
		Expect(err).ToNot(HaveOccurred())
		Expect(storeEvents).ToNot(BeNil())
		Expect(storeEvents).ToNot(BeEmpty())
		Expect(len(storeEvents)).To(BeNumerically("==", len(ev)))
	})
	It("can filter events to load from the store by aggregate version", func() {
		es := createInMemoryTestEventStore()
		defer clearInMemoryEs(es)

		events := createTestEvents()
		err := es.Save(ctx, events)
		Expect(err).ToNot(HaveOccurred())

		minVersion := uint64(1)
		maxVersion := uint64(1)
		storeEvents, err := es.Load(ctx, &evs.StoreQuery{
			MinVersion: &minVersion,
		})
		Expect(err).ToNot(HaveOccurred())
		Expect(storeEvents).ToNot(BeNil())
		Expect(storeEvents).ToNot(BeEmpty())

		for _, ev := range storeEvents {
			Expect(ev.AggregateVersion()).To(BeNumerically(">=", 1))
		}

		storeEvents, err = es.Load(ctx, &evs.StoreQuery{
			MaxVersion: &maxVersion,
		})
		Expect(err).ToNot(HaveOccurred())
		Expect(storeEvents).ToNot(BeNil())
		Expect(storeEvents).ToNot(BeEmpty())

		for _, ev := range storeEvents {
			Expect(ev.AggregateVersion()).To(BeNumerically("<=", 1))
		}

		storeEvents, err = es.Load(ctx, &evs.StoreQuery{
			MinVersion: &minVersion,
			MaxVersion: &maxVersion,
		})
		Expect(err).ToNot(HaveOccurred())
		Expect(storeEvents).ToNot(BeNil())
		Expect(storeEvents).ToNot(BeEmpty())

		for _, ev := range storeEvents {
			Expect(ev.AggregateVersion()).To(BeNumerically("==", 1))
		}
	})
})
