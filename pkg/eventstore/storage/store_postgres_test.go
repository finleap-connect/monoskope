package storage

import (
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
			NewEvent(typeTestEventCreated, createTestEventData("create"), time.Now().UTC(), typeTestAggregate, aggregateId, 0),
			NewEvent(typeTestEventChanged, createTestEventData("change"), time.Now().UTC(), typeTestAggregate, aggregateId, 1),
			NewEvent(typeTestEventDeleted, createTestEventData("delete"), time.Now().UTC(), typeTestAggregate, aggregateId, 2),
		})
		Expect(err).ToNot(HaveOccurred())
	})
	It("fails to append new events to the store when they are not of the same aggregate type", func() {
		es := createTestEventStore()
		defer func() { _ = es.clear(ctx) }()
		aggregateId := uuid.New()

		err := es.Save(ctx, []Event{
			NewEvent(typeTestEventCreated, createTestEventData("create"), time.Now().UTC(), typeTestAggregate, aggregateId, 0),
			NewEvent(typeTestEventChanged, createTestEventData("change"), time.Now().UTC(), typeTestAggregateExtended, aggregateId, 1),
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
			NewEvent(typeTestEventCreated, createTestEventData("create"), time.Now().UTC(), typeTestAggregate, aggregateId, 0),
			NewEvent(typeTestEventChanged, createTestEventData("change"), time.Now().UTC(), typeTestAggregate, aggregateId, 2),
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
			NewEvent(typeTestEventCreated, createTestEventData("create"), time.Now().UTC(), typeTestAggregate, aggregateId, 0),
			NewEvent(typeTestEventChanged, createTestEventData("change"), time.Now().UTC(), typeTestAggregate, aggregateId, 1),
		})
		Expect(err).ToNot(HaveOccurred())

		err = es.Save(ctx, []Event{
			NewEvent(typeTestEventChanged, createTestEventData("change"), time.Now().UTC(), typeTestAggregate, aggregateId, 1),
		})
		Expect(err).To(HaveOccurred())
		esErr := UnwrapEventStoreError(err)
		Expect(esErr).ToNot(BeNil())
		Expect(esErr.Err).To(Equal(ErrAggregateVersionAlreadyExists))
	})
})

func createTestEventStore() *EventStore {
	es, err := NewPostgresEventStore(env.DB, jsonEncoder{})
	Expect(err).ToNot(HaveOccurred())
	Expect(es).ToNot(BeNil())
	return es
}
