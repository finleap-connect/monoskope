package projectors

import (
	"context"
	"time"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/domain/eventdata"
	apiProjections "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/domain/projections"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/constants/aggregates"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/constants/events"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/projections"
	es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing"
)

var (
	expectedName                = "the one cluster"
	expectedLabel               = "one-cluster"
	expectedApiServerAddress    = "one.example.com"
	expectedClusterCACertBundle = []byte("This should be a certificate")
	expectedJWT                 = "thisisnotajwt"
)

var _ = Describe("domain/cluster_repo", func() {
	ctx := context.Background()
	userId := uuid.New()
	adminUser := &apiProjections.User{Id: userId.String(), Name: "admin", Email: "admin@monoskope.io"}

	It("can handle events", func() {
		clusterProjector := NewClusterProjector()
		clusterProjection := clusterProjector.NewProjection(uuid.New())
		protoClusterCreatedEventData := &eventdata.ClusterCreated{
			Name:                expectedName,
			Label:               expectedLabel,
			ApiServerAddress:    expectedApiServerAddress,
			CaCertificateBundle: expectedClusterCACertBundle,
		}
		clusterCreatedEventData := es.ToEventDataFromProto(protoClusterCreatedEventData)
		clusterProjection, err := clusterProjector.Project(context.Background(), es.NewEvent(ctx, events.ClusterCreated, clusterCreatedEventData, time.Now().UTC(), aggregates.Cluster, uuid.MustParse(adminUser.Id), 1), clusterProjection)
		Expect(err).NotTo(HaveOccurred())
		Expect(clusterProjection.Version()).To(Equal(uint64(1)))
		cluster, ok := clusterProjection.(*projections.Cluster)
		Expect(ok).To(BeTrue())
		Expect(cluster.GetName()).To(Equal(expectedName))
		Expect(cluster.GetLabel()).To(Equal(expectedLabel))
		Expect(cluster.GetApiServerAddress()).To(Equal(expectedApiServerAddress))
		Expect(cluster.GetClusterCACertBundle()).To(Equal(expectedClusterCACertBundle))

		protoTokenCreatedEventData := &eventdata.ClusterBootstrapTokenCreated{
			JWT: expectedJWT,
		}
		tokenCreatedEventData := es.ToEventDataFromProto(protoTokenCreatedEventData)
		clusterProjection, err = clusterProjector.Project(context.Background(), es.NewEvent(ctx, events.ClusterCreated, tokenCreatedEventData, time.Now().UTC(), aggregates.Cluster, uuid.MustParse(adminUser.Id), 1), clusterProjection)
		Expect(err).NotTo(HaveOccurred())
		Expect(clusterProjection.Version()).To(Equal(uint64(2)))
		cluster, ok = clusterProjection.(*projections.Cluster)
		Expect(ok).To(BeTrue())
		Expect(cluster.GetBootstrapToken()).To(Equal(expectedJWT))
	})
})
