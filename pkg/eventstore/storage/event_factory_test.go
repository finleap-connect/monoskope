package storage

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("eventfactory", func() {
	It("can register event data", func() {
		err := RegisterEventData(testEventExtended, func() EventData { return &TestEventDataExtened{} })
		Expect(err).ToNot(HaveOccurred())
	})
	It("fails to register empty event data", func() {
		err := RegisterEventData("", func() EventData { return &TestEventDataExtened{} })
		Expect(err).To(HaveOccurred())
	})
	It("fails to register event data more than once", func() {
		err := RegisterEventData(testEventCreated, func() EventData { return &TestEventData{} })
		Expect(err).To(HaveOccurred())
	})
	It("can unregister event data", func() {
		err := UnregisterEventData(testEventExtended)
		Expect(err).ToNot(HaveOccurred())
	})
	It("fails to unregister empty event data", func() {
		err := UnregisterEventData("")
		Expect(err).To(HaveOccurred())
	})
	It("fails to unregister unknown event data", func() {
		err := UnregisterEventData("foobar")
		Expect(err).To(HaveOccurred())
	})
	It("can create event data", func() {
		data, err := CreateEventData(testEventCreated)
		Expect(err).ToNot(HaveOccurred())
		Expect(data).ToNot(BeNil())
	})
	It("fails to create event data for empty type", func() {
		data, err := CreateEventData("")
		Expect(err).To(HaveOccurred())
		Expect(data).To(BeNil())
	})
	It("fails to create event data for unknown type", func() {
		data, err := CreateEventData("foobar")
		Expect(err).To(HaveOccurred())
		Expect(data).To(BeNil())
	})
})
