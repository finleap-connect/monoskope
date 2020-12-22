package storage

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	testEventData = createTestEventData("World")
	jsonString    = "{\"Hello\":\"World\"}"
	jsonBytes     = []byte(jsonString)
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
		decodedEventData, err := encoder.Unmarshal(testEventCreated, jsonBytes)
		Expect(err).ToNot(HaveOccurred())
		Expect(decodedEventData).ToNot(BeNil())
		ed, ok := decodedEventData.(*TestEventData)
		Expect(ok).To(BeTrue())
		Expect(ed.Hello).To(Equal(testEventData.Hello))
	})
	It("ignores empty event data when marshalling", func() {
		encoder := &jsonEncoder{}
		bytes, err := encoder.Marshal(nil)
		Expect(err).ToNot(HaveOccurred())
		Expect(bytes).To(BeNil())
	})
	It("ignores empty event data when unmarshalling", func() {
		encoder := &jsonEncoder{}
		decodedEventData, err := encoder.Unmarshal(testEventCreated, []byte{})
		Expect(err).ToNot(HaveOccurred())
		Expect(decodedEventData).To(BeNil())
	})
})
