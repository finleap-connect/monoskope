package aggregates

import (
	"time"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/domain/eventdata"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/constants/events"
	es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing"
)

var (
	expectedTenantName = "the one tenant"
	expectedPrefix     = "tenant-one"
)

var _ = Describe("Unit Test for the Tenant Aggregate", func() {
	It("should set the data from a command to the resultant event", func() {
		ctx := createSysAdminCtx()
		inID := uuid.Nil
		agg := NewTenantAggregate(NewTestAggregateManager())

		reply, err := createTenant(ctx, agg)
		Expect(err).NotTo(HaveOccurred())
		Expect(reply.Id).ToNot(Equal(inID))
		Expect(reply.Version).To(Equal(uint64(0)))

		event := agg.UncommittedEvents()[0]

		Expect(event.EventType()).To(Equal(events.TenantCreated))
		Expect(event.AggregateID()).ToNot(Equal(inID))

		data := &eventdata.TenantCreated{}
		err = event.Data().ToProto(data)
		Expect(err).NotTo(HaveOccurred())

		Expect(data.Name).To(Equal(expectedTenantName))

	})

	It("should apply the data from an event to the aggregate", func() {
		ctx := createSysAdminCtx()
		agg := NewTenantAggregate(NewTestAggregateManager())

		ed := es.ToEventDataFromProto(&eventdata.TenantCreated{
			Name:   expectedTenantName,
			Prefix: expectedPrefix,
		})
		esEvent := es.NewEvent(ctx, events.TenantCreated, ed, time.Now().UTC(),
			agg.Type(), agg.ID(), agg.Version())

		err := agg.ApplyEvent(esEvent)
		Expect(err).NotTo(HaveOccurred())

		Expect(agg.(*TenantAggregate).Name).To(Equal(expectedTenantName))
		Expect(agg.(*TenantAggregate).Prefix).To(Equal(expectedPrefix))

	})
})
