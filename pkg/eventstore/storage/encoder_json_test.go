package storage

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	storage_test "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventstore/storage/test"
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
		decodedEventData, err := encoder.Unmarshal(typeTestEventCreated, jsonBytes)
		Expect(err).ToNot(HaveOccurred())
		Expect(decodedEventData).ToNot(BeNil())
		ed, ok := decodedEventData.(*storage_test.TestEventData)
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
		decodedEventData, err := encoder.Unmarshal(typeTestEventCreated, []byte{})
		Expect(err).ToNot(HaveOccurred())
		Expect(decodedEventData).To(BeNil())
	})
})
