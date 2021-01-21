package usecases

import (
	"time"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	api_es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/eventstore"
	evs "gitlab.figo.systems/platform/monoskope/monoskope/pkg/event_sourcing"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

var _ = Describe("Converters", func() {
	It("can convert to storage event from proto", func() {
		data, err := evs.ToEventData(&api_es.TestEventData{Hello: "world"})
		Expect(err).ToNot(HaveOccurred())

		timestamp := time.Now().UTC()
		pe := &api_es.Event{
			Type:             "TestEventType",
			Timestamp:        timestamppb.New(timestamp),
			AggregateId:      uuid.New().String(),
			AggregateType:    "TestAggregateType",
			AggregateVersion: wrapperspb.UInt64(0),
			Data:             data,
		}
		se, err := NewEventFromProto(pe)
		Expect(err).ToNot(HaveOccurred())

		checkProtoStorageEventEquality(pe, se)
	})
	It("can convert to proto event from storage", func() {
		timestamp := time.Now().UTC()
		aggregateId := uuid.New()

		se := evs.NewEvent(
			evs.EventType("TestType"),
			evs.EventData("{\"hello\":\"world\"}"),
			timestamp,
			evs.AggregateType("TestAggregateType"),
			aggregateId,
			0)
		pe, err := NewProtoFromEvent(se)
		Expect(err).ToNot(HaveOccurred())

		checkProtoStorageEventEquality(pe, se)
	})
	It("can convert to storage query from proto filter", func() {
		aggregateId := uuid.New()
		aggregateType := evs.AggregateType("TestAggregateType")
		maxTimestamp := time.Now().UTC()
		minTimestamp := maxTimestamp.Add(-1 * time.Hour)

		pf := &api_es.EventFilter{
			ByAggregate:  &api_es.EventFilter_AggregateId{AggregateId: wrapperspb.String(aggregateId.String())},
			MinVersion:   wrapperspb.UInt64(1),
			MaxVersion:   wrapperspb.UInt64(4),
			MinTimestamp: timestamppb.New(minTimestamp),
			MaxTimestamp: timestamppb.New(maxTimestamp),
		}
		q, err := NewStoreQueryFromProto(pf)
		Expect(err).ToNot(HaveOccurred())
		Expect(q).ToNot(BeNil())
		Expect(q.AggregateId).To(Equal(&aggregateId))

		pf.ByAggregate = &api_es.EventFilter_AggregateType{AggregateType: wrapperspb.String(aggregateType.String())}
		q, err = NewStoreQueryFromProto(pf)
		Expect(err).ToNot(HaveOccurred())
		Expect(q).ToNot(BeNil())
		Expect(q.AggregateId).To(BeNil())
		Expect(q.AggregateType).To(Equal(&aggregateType))
		Expect(*q.MinVersion).To(Equal(pf.MinVersion.GetValue()))
		Expect(*q.MaxVersion).To(Equal(pf.MaxVersion.GetValue()))
		Expect(q.MinTimestamp).To(Equal(&minTimestamp))
		Expect(q.MaxTimestamp).To(Equal(&maxTimestamp))
	})
	It("fails to convert to storage query from proto filter for invalid aggregate id", func() {
		data, err := evs.ToEventData(&api_es.TestEventData{Hello: "world"})
		Expect(err).ToNot(HaveOccurred())

		pe := &api_es.Event{
			Type:             "TestEventType",
			Timestamp:        timestamppb.New(time.Now().UTC()),
			AggregateId:      "", // invalid id
			AggregateType:    "TestAggregateType",
			AggregateVersion: wrapperspb.UInt64(0),
			Data:             data,
		}
		se, err := NewEventFromProto(pe)
		Expect(err).To(HaveOccurred())
		Expect(se).To(BeNil())
		Expect(err).To(Equal(ErrCouldNotParseAggregateId))
	})
})

func checkProtoStorageEventEquality(pe *api_es.Event, se evs.Event) {
	Expect(pe).ToNot(BeNil())
	Expect(se).ToNot(BeNil())
	Expect(pe.Type).To(Equal(se.EventType().String()))
	Expect(pe.Timestamp.AsTime()).To(Equal(se.Timestamp()))
	Expect(pe.AggregateId).To(Equal(se.AggregateID().String()))
	Expect(pe.AggregateType).To(Equal(se.AggregateType().String()))
	Expect(pe.AggregateVersion.GetValue()).To(Equal(se.AggregateVersion()))
	Expect(evs.EventData(pe.GetData())).To(Equal(se.Data()))
}
