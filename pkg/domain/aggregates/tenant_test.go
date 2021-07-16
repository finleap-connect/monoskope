package aggregates

import (
	"context"
	"time"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/domain/eventdata"
	cmd "gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/commands"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/constants/events"
	es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing"
)

var (
	expectedTenantName = "the one tenant"
	expectedPrefix     = "tenant-one"
)

var _ = Describe("Unit Test for the Tenant Aggregate", func() {

	It("should set the data from a command to the resultant event", func() {

		ctx, err := makeMetadataContextWithSystemAdminUser()
		Expect(err).NotTo(HaveOccurred())

		inID := uuid.Nil
		agg := NewTenantAggregate(inID, NewTestAggregateManager())

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

		ctx, err := makeMetadataContextWithSystemAdminUser()
		Expect(err).NotTo(HaveOccurred())

		agg := NewTenantAggregate(uuid.New(), NewTestAggregateManager())

		ed := es.ToEventDataFromProto(&eventdata.TenantCreated{
			Name:   expectedTenantName,
			Prefix: expectedPrefix,
		})
		esEvent := es.NewEvent(ctx, events.TenantCreated, ed, time.Now().UTC(),
			agg.Type(), agg.ID(), agg.Version())

		err = agg.ApplyEvent(esEvent)
		Expect(err).NotTo(HaveOccurred())

		Expect(agg.(*TenantAggregate).Name).To(Equal(expectedTenantName))
		Expect(agg.(*TenantAggregate).Prefix).To(Equal(expectedPrefix))

	})
})

func createTenant(ctx context.Context, agg es.Aggregate) (*es.CommandReply, error) {
	esCommand, ok := cmd.NewCreateTenantCommand(agg.ID()).(*cmd.CreateTenantCommand)
	Expect(ok).To(BeTrue())

	esCommand.CreateTenantCommandData.Name = expectedTenantName
	esCommand.CreateTenantCommandData.Prefix = expectedPrefix

	return agg.HandleCommand(ctx, esCommand)
}
