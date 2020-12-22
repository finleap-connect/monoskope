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
	It("can unregister event data", func() {
		err := UnregisterEventData(TestEventExtended)
		Expect(err).ToNot(HaveOccurred())
	})
	It("can create event data", func() {
		data, err := CreateEventData(TestEvent)
		Expect(err).ToNot(HaveOccurred())
		Expect(data).ToNot(BeNil())
	})
})
