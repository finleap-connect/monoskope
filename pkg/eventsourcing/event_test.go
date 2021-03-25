package eventsourcing

import (
	"context"
	"time"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	esApi "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/eventsourcing"
	testEd "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/eventsourcing/eventdata"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing/errors"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

var _ = Describe("EventData", func() {
	var (
		testEventType     EventType     = "TestEventType"
		testAggregateType AggregateType = "TestAggregateType"
	)

	checkProtoStorageEventEquality := func(pe *esApi.Event, se Event) {
		Expect(pe).ToNot(BeNil())
		Expect(se).ToNot(BeNil())
		Expect(pe.Type).To(Equal(se.EventType().String()))
		Expect(pe.Timestamp.AsTime()).To(Equal(se.Timestamp()))
		Expect(pe.AggregateId).To(Equal(se.AggregateID().String()))
		Expect(pe.AggregateType).To(Equal(se.AggregateType().String()))
		Expect(pe.AggregateVersion.GetValue()).To(Equal(se.AggregateVersion()))
		Expect(se.Data()).To(Equal(EventData(pe.Data)))

		proto := &testEd.TestEventData{}
		_, err := ToEventDataFromProto(proto)
		Expect(err).ToNot(HaveOccurred())
	}

	It("can convert to storage event from proto", func() {
		proto := &testEd.TestEventData{Hello: "world"}
		ed, err := ToEventDataFromProto(proto)
		Expect(err).ToNot(HaveOccurred())

		timestamp := time.Now().UTC()
		pe := &esApi.Event{
			Type:             testEventType.String(),
			Timestamp:        timestamppb.New(timestamp),
			AggregateId:      uuid.New().String(),
			AggregateType:    testAggregateType.String(),
			AggregateVersion: wrapperspb.UInt64(0),
			Data:             ed,
		}

		se, err := NewEventFromProto(pe)
		Expect(err).ToNot(HaveOccurred())

		checkProtoStorageEventEquality(pe, se)
	})
	It("can convert to proto event from storage", func() {
		timestamp := time.Now().UTC()
		aggregateId := uuid.New()

		ed, err := ToEventDataFromProto(&testEd.TestEventData{Hello: "world"})
		Expect(err).ToNot(HaveOccurred())

		se := NewEvent(
			context.Background(),
			EventType("TestType"),
			ed,
			timestamp,
			AggregateType("TestAggregateType"),
			aggregateId,
			0)
		pe := NewProtoFromEvent(se)

		checkProtoStorageEventEquality(pe, se)
	})
	It("fails to convert to storage query from proto filter for invalid aggregate id", func() {
		proto := &testEd.TestEventData{Hello: "world"}
		ed, err := ToEventDataFromProto(proto)
		Expect(err).ToNot(HaveOccurred())

		pe := &esApi.Event{
			Type:             testEventType.String(),
			Timestamp:        timestamppb.New(time.Now().UTC()),
			AggregateId:      "", // invalid id
			AggregateType:    testAggregateType.String(),
			AggregateVersion: wrapperspb.UInt64(0),
			Data:             ed,
		}

		se, err := NewEventFromProto(pe)
		Expect(err).To(HaveOccurred())
		Expect(se).To(BeNil())
		Expect(err).To(Equal(errors.ErrCouldNotParseAggregateId))
	})
})
