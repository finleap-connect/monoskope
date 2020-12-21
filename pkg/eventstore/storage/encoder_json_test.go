package storage

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type TestEventData struct {
	Hello string `json:",omitempty"`
}

var (
	jsonString = "{\"Hello\":\"World\"}"
	jsonBytes  = []byte(jsonString)
	eventData  = TestEventData{Hello: "World"}
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
	// It("can unmarshal event data", func() {
	// 	encoder := &jsonEncoder{}
	// 	bytes, err := encoder.Unmarshal(EventType("TestEventData"), jsonBytes)
	// 	Expect(err).ToNot(HaveOccurred())
	// 	Expect(bytes).ToNot(BeNil())
	// 	Expect(bytes).ToNot(BeEmpty())
	// })
})
