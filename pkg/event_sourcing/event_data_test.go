package event_sourcing

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	api "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/commands"
	"google.golang.org/protobuf/types/known/anypb"
)

var _ = Describe("EventData", func() {
	getProto := func() *api.TestCommandData {
		return &api.TestCommandData{Test: "Hello world!"}
	}
	eventDataFromProto := func() EventData {
		eventData, err := ToEventDataFromProto(getProto())
		Expect(err).To(Not(HaveOccurred()))
		Expect(eventData).To(Not(BeNil()))
		return eventData
	}

	It("can create from proto", func() {
		_ = eventDataFromProto()
	})
	It("can create from any", func() {
		proto := getProto()
		any := anypb.Any{}
		err := any.MarshalFrom(proto)
		Expect(err).To(Not(HaveOccurred()))

		eventData, err := ToEventDataFromAny(&any)
		Expect(err).To(Not(HaveOccurred()))
		Expect(eventData).To(Not(BeNil()))
	})
	It("can unmarshall to any", func() {
		eventData := eventDataFromProto()
		any, err := eventData.ToAny()
		Expect(err).To(Not(HaveOccurred()))
		Expect(any).To(Not(BeNil()))
	})
	It("can unmarshall to proto", func() {
		eventData := eventDataFromProto()
		proto := &api.TestCommandData{}
		err := eventData.ToProto(proto)
		Expect(err).To(Not(HaveOccurred()))
		Expect(proto.GetTest()).To(Equal(getProto().GetTest()))
	})
})
