package event_sourcing

import (
	"time"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/eventdata/test"
	api_es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/eventstore"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

var _ = Describe("EventData", func() {
	var (
		testEventType     EventType     = "TestEventType"
		testAggregateType AggregateType = "TestAggregateType"
	)

	checkProtoStorageEventEquality := func(pe *api_es.Event, se Event) {
		Expect(pe).ToNot(BeNil())
		Expect(se).ToNot(BeNil())
		Expect(pe.Type).To(Equal(se.EventType().String()))
		Expect(pe.Timestamp.AsTime()).To(Equal(se.Timestamp()))
		Expect(pe.AggregateId).To(Equal(se.AggregateID().String()))
		Expect(pe.AggregateType).To(Equal(se.AggregateType().String()))
		Expect(pe.AggregateVersion.GetValue()).To(Equal(se.AggregateVersion()))

		ed, err := ToEventDataFromAny(pe.Data)
		Expect(err).ToNot(HaveOccurred())
		Expect(se.Data()).To(Equal(ed))

		proto := &test.TestEventData{}
		err = se.Data().ToProto(proto)
		Expect(err).ToNot(HaveOccurred())

		a := &anypb.Any{}
		err = a.MarshalFrom(proto)
		Expect(err).ToNot(HaveOccurred())
		Expect(pe.Data.TypeUrl).To(Equal(a.TypeUrl))
		Expect(pe.Data.Value).To(Equal(a.Value))
	}

	It("can convert to storage event from proto", func() {
		timestamp := time.Now().UTC()
		pe := &api_es.Event{
			Type:             testEventType.String(),
			Timestamp:        timestamppb.New(timestamp),
			AggregateId:      uuid.New().String(),
			AggregateType:    testAggregateType.String(),
			AggregateVersion: wrapperspb.UInt64(0),
			Data:             &anypb.Any{},
		}

		proto := &test.TestEventData{Hello: "world"}
		err := pe.Data.MarshalFrom(proto)
		Expect(err).ToNot(HaveOccurred())

		se, err := NewEventFromProto(pe)
		Expect(err).ToNot(HaveOccurred())

		checkProtoStorageEventEquality(pe, se)
	})
	It("can convert to proto event from storage", func() {
		timestamp := time.Now().UTC()
		aggregateId := uuid.New()

		ed, err := ToEventDataFromProto(&test.TestEventData{Hello: "world"})
		Expect(err).ToNot(HaveOccurred())

		se := NewEvent(
			EventType("TestType"),
			ed,
			timestamp,
			AggregateType("TestAggregateType"),
			aggregateId,
			0)
		pe, err := NewProtoFromEvent(se)
		Expect(err).ToNot(HaveOccurred())

		checkProtoStorageEventEquality(pe, se)
	})
	It("fails to convert to storage query from proto filter for invalid aggregate id", func() {
		pe := &api_es.Event{
			Type:             testEventType.String(),
			Timestamp:        timestamppb.New(time.Now().UTC()),
			AggregateId:      "", // invalid id
			AggregateType:    testAggregateType.String(),
			AggregateVersion: wrapperspb.UInt64(0),
			Data:             &anypb.Any{},
		}

		proto := &test.TestEventData{Hello: "world"}
		err := pe.Data.MarshalFrom(proto)
		Expect(err).ToNot(HaveOccurred())

		se, err := NewEventFromProto(pe)
		Expect(err).To(HaveOccurred())
		Expect(se).To(BeNil())
		Expect(err).To(Equal(ErrCouldNotParseAggregateId))
	})
})
