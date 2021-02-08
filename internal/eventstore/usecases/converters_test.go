package usecases

import (
	"time"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	api_es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/eventstore"
	evs "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

var _ = Describe("Converters", func() {
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
})
