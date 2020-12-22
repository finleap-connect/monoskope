package storage

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("eventfactory", func() {
	It("can register event data", func() {
		err := RegisterEventData(TestEventExtended, func() EventData { return &TestEventDataExtened{} })
		Expect(err).ToNot(HaveOccurred())
	})
	It("fails to register empty event data", func() {
		err := RegisterEventData(EventType(""), func() EventData { return &TestEventDataExtened{} })
		Expect(err).To(HaveOccurred())
	})
	It("fails to register event data more than once", func() {
		err := RegisterEventData(TestEvent, func() EventData { return &TestEventData{} })
		Expect(err).To(HaveOccurred())
	})
	It("can unregister event data", func() {
		err := UnregisterEventData(TestEventExtended)
		Expect(err).ToNot(HaveOccurred())
	})
	It("fails to unregister empty event data", func() {
		err := UnregisterEventData(EventType(""))
		Expect(err).To(HaveOccurred())
	})
	It("fails to unregister unknown event data", func() {
		err := UnregisterEventData(EventType("foobar"))
		Expect(err).To(HaveOccurred())
	})
	It("can create event data", func() {
		data, err := CreateEventData(TestEvent)
		Expect(err).ToNot(HaveOccurred())
		Expect(data).ToNot(BeNil())
	})
	It("fails to create event data for empty type", func() {
		data, err := CreateEventData(EmptyEventType)
		Expect(err).To(HaveOccurred())
		Expect(data).To(BeNil())
	})
	It("fails to create event data for unknown type", func() {
		data, err := CreateEventData(EventType("foobar"))
		Expect(err).To(HaveOccurred())
		Expect(data).To(BeNil())
	})
})
