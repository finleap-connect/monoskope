package storage

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("storage/postgres", func() {
	It("can create new event store", func() {
		_ = createTestEventStore()
	})
	It("can append new events to the store", func() {
		es := createTestEventStore()
		defer func() { _ = es.clear(ctx) }()
		aggregateId := uuid.New()

		err := es.Save(ctx, []Event{
			NewEvent(EventType(testEventCreated), createTestEventData("create"), now(), AggregateType(testAggregate), aggregateId, 0),
			NewEvent(EventType(testEventChanged), createTestEventData("change"), now(), AggregateType(testAggregate), aggregateId, 1),
			NewEvent(EventType(testEventDeleted), createTestEventData("delete"), now(), AggregateType(testAggregate), aggregateId, 2),
		})
		Expect(err).ToNot(HaveOccurred())
	})
	It("fails to append new events to the store when they are not of the same aggregate type", func() {
		es := createTestEventStore()
		defer func() { _ = es.clear(ctx) }()
		aggregateId := uuid.New()

		err := es.Save(ctx, []Event{
			NewEvent(testEventCreated, createTestEventData("create"), now(), testAggregate, aggregateId, 0),
			NewEvent(testEventChanged, createTestEventData("change"), now(), testAggregateExtended, aggregateId, 1),
		})
		Expect(err).To(HaveOccurred())
		Expect(err).To(Equal(EventStoreError{
			Err: ErrInvalidAggregateType,
		}))
	})
	It("fails to append new events to the store when they are not in the right aggregate version order", func() {
		es := createTestEventStore()
		defer func() { _ = es.clear(ctx) }()
		aggregateId := uuid.New()

		err := es.Save(ctx, []Event{
			NewEvent(testEventCreated, createTestEventData("create"), now(), testAggregate, aggregateId, 0),
			NewEvent(testEventChanged, createTestEventData("change"), now(), testAggregate, aggregateId, 2),
		})
		Expect(err).To(HaveOccurred())
		Expect(err).To(Equal(EventStoreError{
			Err: ErrIncorrectAggregateVersion,
		}))
	})
	It("fails to append new events to the store when the aggregate version does already exist", func() {
		es := createTestEventStore()
		defer func() { _ = es.clear(ctx) }()
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
		Expect(esErr.Err).To(Equal(ErrAggregateVersionAlreadyExists))
	})
	It("can load events from the store", func() {
		es := createTestEventStore()
		defer func() { _ = es.clear(ctx) }()
		aggregateId := uuid.New()

		events := []Event{
			NewEvent(EventType(testEventCreated), createTestEventData("create"), now(), AggregateType(testAggregate), aggregateId, 0),
			NewEvent(EventType(testEventChanged), createTestEventData("change"), now(), AggregateType(testAggregate), aggregateId, 1),
			NewEvent(EventType(testEventDeleted), createTestEventData("delete"), now(), AggregateType(testAggregate), aggregateId, 2),
		}
		err := es.Save(ctx, events)
		Expect(err).ToNot(HaveOccurred())

		storeEvents, err := es.Load(ctx, &StoreQuery{
			AggregateId: &aggregateId,
		})
		Expect(err).ToNot(HaveOccurred())
		Expect(storeEvents).ToNot(BeNil())
		Expect(storeEvents).ToNot(BeEmpty())
		Expect(len(storeEvents)).To(BeNumerically("==", len(events)))
	})
})

func createTestEventStore() *EventStore {
	es, err := NewPostgresEventStore(env.DB)
	Expect(err).ToNot(HaveOccurred())
	Expect(es).ToNot(BeNil())
	return es.(*EventStore)
}

func now() time.Time {
	return time.Now().UTC()
}

func createTestEventData(something string) EventData {
	bytes, err := json.Marshal(&TestEventData{Hello: something})
	Expect(err).ToNot(HaveOccurred())
	return EventData(bytes)
}
