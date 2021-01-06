package storage

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("storage/postgres", func() {
	ctx := context.Background()

	clearEs := func(es *postgresEventStore) {
		err := es.clear(ctx)
		Expect(err).ToNot(HaveOccurred())
	}

	createTestEventStore := func() *postgresEventStore {
		es, err := NewPostgresEventStore(&postgresStoreConfig{})
		Expect(err).ToNot(HaveOccurred())
		Expect(es).ToNot(BeNil())
		return es.(*postgresEventStore)
	}

	now := func() time.Time {
		return time.Now().UTC()
	}

	createTestEventData := func(something string) EventData {
		bytes, err := json.Marshal(&testEventData{Hello: something})
		Expect(err).ToNot(HaveOccurred())
		return EventData(bytes)
	}

	createTestEvents := func() []Event {
		aggregateId := uuid.New()

		return []Event{
			NewEvent(EventType(testEventCreated), createTestEventData("create"), now(), AggregateType(testAggregate), aggregateId, 0),
			NewEvent(EventType(testEventChanged), createTestEventData("change"), now(), AggregateType(testAggregate), aggregateId, 1),
			NewEvent(EventType(testEventDeleted), createTestEventData("delete"), now(), AggregateType(testAggregate), aggregateId, 2),
		}
	}

	It("can create new event store", func() {
		_ = createTestEventStore()
	})
	It("can append new events to the store", func() {
		es := createTestEventStore()
		defer clearEs(es)

		err := es.Save(ctx, createTestEvents())
		Expect(err).ToNot(HaveOccurred())
	})
	It("fails to append new events to the store when they are not of the same aggregate type", func() {
		es := createTestEventStore()
		defer clearEs(es)

		aggregateId := uuid.New()
		err := es.Save(ctx, []Event{
			NewEvent(testEventCreated, createTestEventData("create"), now(), testAggregate, aggregateId, 0),
			NewEvent(testEventChanged, createTestEventData("change"), now(), testAggregateExtended, aggregateId, 1),
		})
		Expect(err).To(HaveOccurred())
		Expect(err).To(Equal(eventStoreError{
			Err: ErrInvalidAggregateType,
		}))
	})
	It("fails to append new events to the store when they are not in the right aggregate version order", func() {
		es := createTestEventStore()
		defer clearEs(es)

		aggregateId := uuid.New()
		err := es.Save(ctx, []Event{
			NewEvent(testEventCreated, createTestEventData("create"), now(), testAggregate, aggregateId, 0),
			NewEvent(testEventChanged, createTestEventData("change"), now(), testAggregate, aggregateId, 2),
		})
		Expect(err).To(HaveOccurred())
		Expect(err).To(Equal(eventStoreError{
			Err: ErrIncorrectAggregateVersion,
		}))
	})
	It("fails to append new events to the store when the aggregate version does already exist", func() {
		es := createTestEventStore()
		defer clearEs(es)

		aggregateId := uuid.New()
		err := es.Save(ctx, []Event{
			NewEvent(testEventCreated, createTestEventData("create"), now(), testAggregate, aggregateId, 0),
			NewEvent(testEventChanged, createTestEventData("change"), now(), testAggregate, aggregateId, 1),
		})
		Expect(err).ToNot(HaveOccurred())

		err = es.Save(ctx, []Event{
			NewEvent(testEventChanged, createTestEventData("change"), now(), testAggregate, aggregateId, 1),
		})
		Expect(err).To(HaveOccurred())
		esErr := UnwrapEventStoreError(err)
		Expect(esErr).ToNot(BeNil())
		Expect(esErr.Cause()).To(Equal(ErrAggregateVersionAlreadyExists))
	})
	It("can load events from the store", func() {
		es := createTestEventStore()
		defer clearEs(es)

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
		storeEvents, err := es.Load(ctx, &StoreQuery{
			AggregateId: &aggregateId,
		})
		Expect(err).ToNot(HaveOccurred())
		Expect(storeEvents).ToNot(BeNil())
		Expect(storeEvents).ToNot(BeEmpty())
		Expect(len(storeEvents)).To(BeNumerically("==", expectedEventCount))
	})
	It("can filter events to load from the store by aggregate type", func() {
		es := createTestEventStore()
		defer clearEs(es)

		events := createTestEvents()
		err := es.Save(ctx, events)
		Expect(err).ToNot(HaveOccurred())

		aggregateType := AggregateType(testAggregate)
		storeEvents, err := es.Load(ctx, &StoreQuery{
			AggregateType: &aggregateType,
		})
		Expect(err).ToNot(HaveOccurred())
		Expect(storeEvents).ToNot(BeNil())
		Expect(storeEvents).ToNot(BeEmpty())
		Expect(len(storeEvents)).To(BeNumerically("==", len(events)))
	})
	It("can filter events to load from the store by aggregate version", func() {
		es := createTestEventStore()
		defer clearEs(es)

		events := createTestEvents()
		err := es.Save(ctx, events)
		Expect(err).ToNot(HaveOccurred())

		minVersion := uint64(1)
		maxVersion := uint64(1)
		storeEvents, err := es.Load(ctx, &StoreQuery{
			MinVersion: &minVersion,
		})
		Expect(err).ToNot(HaveOccurred())
		Expect(storeEvents).ToNot(BeNil())
		Expect(storeEvents).ToNot(BeEmpty())

		for _, ev := range storeEvents {
			Expect(ev.AggregateVersion()).To(BeNumerically(">=", 1))
		}

		storeEvents, err = es.Load(ctx, &StoreQuery{
			MaxVersion: &maxVersion,
		})
		Expect(err).ToNot(HaveOccurred())
		Expect(storeEvents).ToNot(BeNil())
		Expect(storeEvents).ToNot(BeEmpty())

		for _, ev := range storeEvents {
			Expect(ev.AggregateVersion()).To(BeNumerically("<=", 1))
		}

		storeEvents, err = es.Load(ctx, &StoreQuery{
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
