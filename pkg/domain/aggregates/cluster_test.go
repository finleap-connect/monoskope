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

var _ = Describe("Unit Test for Cluster Aggregate", func() {
	var (
		expectedJWT = "thisisnotajwt"
	)

	It("should set the data from a command to the resultant event", func() {
		ctx := createSysAdminCtx()
		agg := NewClusterAggregate(NewTestAggregateManager())

		reply, err := createCluster(ctx, agg)
		Expect(err).NotTo(HaveOccurred())
		Expect(reply.Id).ToNot(Equal(uuid.Nil))
		Expect(reply.Version).To(Equal(uint64(0)))

		event := agg.UncommittedEvents()[0]

		Expect(event.EventType()).To(Equal(events.ClusterCreatedV2))

		data := new(eventdata.ClusterCreatedV2)
		err = event.Data().ToProto(data)
		Expect(err).NotTo(HaveOccurred())

		Expect(data.DisplayName).To(Equal(expectedClusterDisplayName))
		Expect(data.Name).To(Equal(expectedClusterName))
		Expect(data.ApiServerAddress).To(Equal(expectedClusterApiServerAddress))
		Expect(data.CaCertificateBundle).To(Equal(expectedClusterCACertBundle))

	})

	It("should apply the data from an event to the aggregate", func() {

		ctx := createSysAdminCtx()
		agg := NewClusterAggregate(NewTestAggregateManager())

		ed := es.ToEventDataFromProto(&eventdata.ClusterCreatedV2{
			DisplayName:         expectedClusterDisplayName,
			Name:                expectedClusterName,
			ApiServerAddress:    expectedClusterApiServerAddress,
			CaCertificateBundle: expectedClusterCACertBundle,
		})
		esEvent := es.NewEvent(ctx, events.ClusterCreatedV2, ed, time.Now().UTC(),
			agg.Type(), agg.ID(), agg.Version())

		err := agg.ApplyEvent(esEvent)
		Expect(err).NotTo(HaveOccurred())

		Expect(agg.(*ClusterAggregate).displayName).To(Equal(expectedClusterDisplayName))
		Expect(agg.(*ClusterAggregate).name).To(Equal(expectedClusterName))
		Expect(agg.(*ClusterAggregate).apiServerAddr).To(Equal(expectedClusterApiServerAddress))
		Expect(agg.(*ClusterAggregate).caCertBundle).To(Equal(expectedClusterCACertBundle))
	})

	It("should write the jwt from a BootstrapTokenCreated event to the aggregate", func() {
		ctx := createSysAdminCtx()
		agg := NewClusterAggregate(NewTestAggregateManager())

		ed := es.ToEventDataFromProto(&eventdata.ClusterBootstrapTokenCreated{
			Jwt: expectedJWT,
		})
		esEvent := es.NewEvent(ctx, events.ClusterBootstrapTokenCreated, ed, time.Now().UTC(),
			agg.Type(), agg.ID(), agg.Version())

		err := agg.ApplyEvent(esEvent)
		Expect(err).NotTo(HaveOccurred())

		Expect(agg.(*ClusterAggregate).bootstrapToken).To(Equal(expectedJWT))
	})

	It("should update the cluster", func() {
		ctx := createSysAdminCtx()
		agg := NewClusterAggregate(NewTestAggregateManager())

		ed := es.ToEventDataFromProto(&eventdata.ClusterBootstrapTokenCreated{
			Jwt: expectedJWT,
		})
		esEvent := es.NewEvent(ctx, events.ClusterBootstrapTokenCreated, ed, time.Now().UTC(),
			agg.Type(), agg.ID(), agg.Version())

		err := agg.ApplyEvent(esEvent)
		Expect(err).NotTo(HaveOccurred())

		Expect(agg.(*ClusterAggregate).bootstrapToken).To(Equal(expectedJWT))
	})
})
