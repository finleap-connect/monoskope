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
)

const shortDuration = 1 * time.Second

var _ = Describe("Unit Test for Cluster Aggregate", func() {
	It("should set the data from a command to the resultant event", func() {

		baseCtx, cancel := context.WithTimeout(context.Background(), shortDuration)
		defer cancel()

		metaMgr, err := meta.NewDomainMetadataManager(baseCtx)
		Expect(err).NotTo(HaveOccurred())
		metaMgr, err = meta.NewDomainMetadataManager(metaMgr.GetContext())
		Expect(err).NotTo(HaveOccurred())

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

		ctx := metaMgr.GetContext()

		var (
			expectedName                = "the one cluster"
			expectedLabel               = "one-cluster"
			expectedApiServerAddress    = "one.example.com"
			expectedClusterCACertBundle = []byte("This should be a certificate")
		)

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

	/*
		It("should apply the data from an event to the aggregate", func() {

			ctx, cancel := context.WithTimeout(context.Background(), shortDuration)
			defer cancel()

			agg := NewClusterAggregate(uuid.New())

			esCommand := cmd.NewCreateClusterCommand(uuid.New())

			err := agg.HandleCommand(ctx, esCommand)
			Expect(err).NotTo(HaveOccurred())

			Expect(ClusterAggregate(agg).GetName()).To(Equal("the one cluster"))

		})

	*/
})
