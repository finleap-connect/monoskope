package eventsourcing

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("aggregate_registry", func() {
	It("can register and unregister aggregates", func() {
		registry := NewAggregateRegistry()
		err := registry.RegisterAggregate(func() Aggregate { return newTestAggregate() })
		Expect(err).ToNot(HaveOccurred())
	})
	It("can't register the same aggregate twice", func() {
		registry := NewAggregateRegistry()
		err := registry.RegisterAggregate(func() Aggregate { return newTestAggregate() })
		Expect(err).ToNot(HaveOccurred())

		err = registry.RegisterAggregate(func() Aggregate { return newTestAggregate() })
		Expect(err).To(HaveOccurred())
	})
	It("can't create aggregates which are not registered", func() {
		registry := NewAggregateRegistry()
		aggregate, err := registry.CreateAggregate(testAggregateType)
		Expect(err).To(HaveOccurred())
		Expect(aggregate).To(BeNil())
	})
	It("can create aggregates which are registered", func() {
		registry := NewAggregateRegistry()
		err := registry.RegisterAggregate(func() Aggregate { return newTestAggregate() })
		Expect(err).ToNot(HaveOccurred())

		aggregate, err := registry.CreateAggregate(testAggregateType)
		Expect(err).ToNot(HaveOccurred())
		Expect(aggregate).ToNot(BeNil())
	})
})
