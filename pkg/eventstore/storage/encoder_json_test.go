package storage

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("jsonEncoder", func() {
	It("can marshal event data", func() {
		encoder := &jsonEncoder{}
		bytes, err := encoder.Marshal(testEventData)
		Expect(err).ToNot(HaveOccurred())
		Expect(bytes).ToNot(BeNil())
		Expect(bytes).ToNot(BeEmpty())
		Expect(jsonBytes).To(Equal(bytes))
	})
	It("can unmarshal event data", func() {
		encoder := &jsonEncoder{}
		decodedEventData, err := encoder.Unmarshal(TestEvent, jsonBytes)
		Expect(err).ToNot(HaveOccurred())
		Expect(decodedEventData).ToNot(BeNil())
		ed, ok := decodedEventData.(*TestEventData)
		Expect(ok).To(BeTrue())
		Expect(ed.Hello).To(Equal(testEventData.Hello))
	})
})
