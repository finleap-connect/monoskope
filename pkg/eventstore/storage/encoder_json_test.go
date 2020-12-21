package storage

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("jsonEncoder", func() {
	It("can marshal event data", func() {
		encoder := &jsonEncoder{}
		bytes, err := encoder.Marshal(eventData)
		Expect(err).ToNot(HaveOccurred())
		Expect(bytes).ToNot(BeNil())
		Expect(bytes).ToNot(BeEmpty())
		Expect(jsonBytes).To(Equal(bytes))
	})
	It("can unmarshal event data", func() {
		encoder := &jsonEncoder{}

		decodedEventData := &TestEventData{}
		err := encoder.Unmarshal(jsonBytes, decodedEventData)
		Expect(err).ToNot(HaveOccurred())
		Expect(decodedEventData).ToNot(BeNil())
		Expect(decodedEventData.Hello).To(Equal(eventData.Hello))
	})
})
