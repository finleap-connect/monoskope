package aggregates

import (
	"context"
	"time"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/domain/eventdata"
	projections "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/domain/projections"
	cmd "gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/commands"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/constants/events"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/constants/roles"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/constants/scopes"
	meta "gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/metadata"
	es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing"
)

var (
	expectedName                = "the one cluster"
	expectedLabel               = "one-cluster"
	expectedApiServerAddress    = "one.example.com"
	expectedClusterCACertBundle = []byte("This should be a certificate")
)

var _ = Describe("Unit Test for Cluster Aggregate", func() {
	It("should set the data from a command to the resultant event", func() {

		ctx, err := makeMetadataContextWithSystemAdminUser()
		Expect(err).NotTo(HaveOccurred())

		agg := NewClusterAggregate(uuid.New())
		agg.IncrementVersion() // to make aggregate verification pass

		esCommand, ok := cmd.NewCreateClusterCommand(uuid.New()).(*cmd.CreateClusterCommand)
		Expect(ok).To(BeTrue())

		esCommand.CreateCluster.Name = expectedName
		esCommand.CreateCluster.Label = expectedLabel
		esCommand.CreateCluster.ApiServerAddress = expectedApiServerAddress
		esCommand.CreateCluster.ClusterCACertBundle = expectedClusterCACertBundle

		err = agg.HandleCommand(ctx, esCommand)
		Expect(err).NotTo(HaveOccurred())

		event := agg.UncommittedEvents()[0]

		Expect(event.EventType()).To(Equal(events.ClusterCreated))

		data := &eventdata.ClusterCreated{}
		err = event.Data().ToProto(data)
		Expect(err).NotTo(HaveOccurred())

		Expect(data.Name).To(Equal(expectedName))
		Expect(data.Label).To(Equal(expectedLabel))
		Expect(data.ApiServerAddress).To(Equal(expectedApiServerAddress))
		Expect(data.CaCertificateBundle).To(Equal(expectedClusterCACertBundle))

	})

	It("should apply the data from an event to the aggregate", func() {

		ctx, err := makeMetadataContextWithSystemAdminUser()
		Expect(err).NotTo(HaveOccurred())

		agg := NewClusterAggregate(uuid.New())
		agg.IncrementVersion()

		ed := es.ToEventDataFromProto(&eventdata.ClusterCreated{
			Name:                expectedName,
			Label:               expectedLabel,
			ApiServerAddress:    expectedApiServerAddress,
			CaCertificateBundle: expectedClusterCACertBundle,
		})
		esEvent := es.NewEvent(ctx, events.ClusterCreated, ed, time.Now().UTC(),
			agg.Type(), agg.ID(), agg.Version())

		err = agg.ApplyEvent(esEvent)
		Expect(err).NotTo(HaveOccurred())

		Expect(agg.(*ClusterAggregate).name).To(Equal(expectedName))
		Expect(agg.(*ClusterAggregate).label).To(Equal(expectedLabel))
		Expect(agg.(*ClusterAggregate).apiServerAddr).To(Equal(expectedApiServerAddress))
		Expect(agg.(*ClusterAggregate).caCertBundle).To(Equal(expectedClusterCACertBundle))

	})

})

func makeMetadataContextWithSystemAdminUser() (context.Context, error) {
	metaMgr, err := meta.NewDomainMetadataManager(context.Background())
	if err != nil {
		return nil, err
	}

	// forces the setting of the domain context
	metaMgr, err = meta.NewDomainMetadataManager(metaMgr.GetContext())
	if err != nil {
		return nil, err
	}

	metaMgr.SetUserInformation(&meta.UserInformation{
		Id:     uuid.New(),
		Name:   "admin",
		Email:  "admin@monoskope.io",
		Issuer: "monoskope",
	})

	metaMgr.SetRoleBindings([]*projections.UserRoleBinding{
		{
			Role:  roles.Admin.String(),
			Scope: scopes.System.String(),
		},
	})

	return metaMgr.GetContext(), nil

}
