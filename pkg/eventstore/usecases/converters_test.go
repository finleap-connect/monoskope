package usecases

import (
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/eventstore"
	api_es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/eventstore"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventstore/storage"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

var _ = Describe("Converters", func() {
	It("can convert to storage event from proto", func() {
		data, err := ptypes.MarshalAny(&api_es.TestEventData{Hello: "world"})
		Expect(err).ToNot(HaveOccurred())

		timestamp := time.Now().UTC()
		pe := &eventstore.Event{
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

		se := storage.NewEvent(
			storage.EventType("TestType"),
			storage.EventData("{\"@type\":\"type.googleapis.com/eventstore.TestEventData\",\"hello\":\"world\"}"),
			timestamp,
			storage.AggregateType("TestAggregateType"),
			aggregateId,
			0)
		pe, err := NewProtoFromEvent(se)
		Expect(err).ToNot(HaveOccurred())

		checkProtoStorageEventEquality(pe, se)
	})
})

func checkProtoStorageEventEquality(pe *eventstore.Event, se storage.Event) {
	Expect(pe.Type).To(Equal(string(se.EventType())))
	Expect(pe.Timestamp.AsTime()).To(Equal(se.Timestamp()))
	Expect(pe.AggregateId).To(Equal(se.AggregateID().String()))
	Expect(pe.AggregateType).To(Equal(string(se.AggregateType())))
	Expect(pe.AggregateVersion.GetValue()).To(Equal(se.AggregateVersion()))
}
