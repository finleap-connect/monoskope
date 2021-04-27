package eventsourcing

import (
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("aggregate_registry", func() {
	It("can register and unregistered aggregates", func() {
		registry := NewAggregateRegistry()
		registry.RegisterAggregate(func(id uuid.UUID) Aggregate { return newTestAggregate() })
	})
	It("can't register the same aggregate twice", func() {
		registry := NewAggregateRegistry()
		registry.RegisterAggregate(func(id uuid.UUID) Aggregate { return newTestAggregate() })

		defer func() {
			Expect(recover()).To(HaveOccurred())
		}()
		registry.RegisterAggregate(func(id uuid.UUID) Aggregate { return newTestAggregate() })
	})
	It("can't create aggregates which are not registered", func() {
		registry := NewAggregateRegistry()
		aggregate, err := registry.CreateAggregate(testAggregateType, uuid.Nil)
		Expect(err).To(HaveOccurred())
		Expect(aggregate).To(BeNil())
	})
	It("can create aggregates which are registered", func() {
		registry := NewAggregateRegistry()
		registry.RegisterAggregate(func(id uuid.UUID) Aggregate { return newTestAggregate() })

		aggregate, err := registry.CreateAggregate(testAggregateType, uuid.Nil)
		Expect(err).ToNot(HaveOccurred())
		Expect(aggregate).ToNot(BeNil())
	})
})
