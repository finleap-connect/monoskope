package projectors

import (
	"context"
	"time"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/domain/eventdata"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/constants/aggregates"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/constants/events"
	metadata "gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/metadata"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/projections"
	es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing"
)

var (
	expectedName                = "the one cluster"
	expectedLabel               = "one-cluster"
	expectedApiServerAddress    = "one.example.com"
	expectedClusterCACertBundle = []byte("This should be a certificate")
	expectedM8CA                = []byte("m8 CA")
	expectedClusterCertificate  = []byte("This should also be a certificate")
	expectedJWT                 = "thisisnotajwt"
)

var _ = Describe("domain/cluster_repo", func() {
	ctx := context.Background()
	userId := uuid.New()

	mdManager, err := metadata.NewDomainMetadataManager(ctx)
	Expect(err).ToNot(HaveOccurred())
	mdManager.SetUserInformation(&metadata.UserInformation{
		Id:     userId,
		Name:   "admin",
		Email:  "admin@monoskope.io",
		Issuer: "monoskope",
	})
	ctx = mdManager.GetContext()

	It("can handle ClusterCreated events", func() {
		clusterProjector := NewClusterProjector()
		clusterProjection := clusterProjector.NewProjection(uuid.New())
		protoClusterCreatedEventData := &eventdata.ClusterCreated{
			Name:                expectedName,
			Label:               expectedLabel,
			ApiServerAddress:    expectedApiServerAddress,
			CaCertificateBundle: expectedClusterCACertBundle,
		}
		clusterCreatedEventData := es.ToEventDataFromProto(protoClusterCreatedEventData)
		event := es.NewEvent(ctx, events.ClusterCreated, clusterCreatedEventData, time.Now().UTC(), aggregates.Cluster, uuid.New(), 1)

		clusterProjection, err := clusterProjector.Project(context.Background(), event, clusterProjection)
		Expect(err).NotTo(HaveOccurred())

		Expect(clusterProjection.Version()).To(Equal(uint64(1)))

		cluster, ok := clusterProjection.(*projections.Cluster)
		Expect(ok).To(BeTrue())

		Expect(cluster.GetName()).To(Equal(expectedName))
		Expect(cluster.GetLabel()).To(Equal(expectedLabel))
		Expect(cluster.GetApiServerAddress()).To(Equal(expectedApiServerAddress))
		Expect(cluster.GetClusterCACertBundle()).To(Equal(expectedClusterCACertBundle))

		dp := cluster.DomainProjection
		Expect(dp.Created).ToNot(BeNil())
	})

	It("can handle ClusterBootstrapTokenCreated events", func() {
		clusterProjector := NewClusterProjector()
		clusterProjection := clusterProjector.NewProjection(uuid.New())
		protoClusterCreatedEventData := &eventdata.ClusterCreated{
			Name:                expectedName,
			Label:               expectedLabel,
			ApiServerAddress:    expectedApiServerAddress,
			CaCertificateBundle: expectedClusterCACertBundle,
		}
		clusterCreatedEventData := es.ToEventDataFromProto(protoClusterCreatedEventData)
		clusterProjection, err := clusterProjector.Project(context.Background(), es.NewEvent(ctx, events.ClusterCreated, clusterCreatedEventData, time.Now().UTC(), aggregates.Cluster, uuid.New(), 1), clusterProjection)
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
		clusterProjection, err = clusterProjector.Project(context.Background(), es.NewEvent(ctx, events.ClusterBootstrapTokenCreated, tokenCreatedEventData, time.Now().UTC(), aggregates.Cluster, uuid.New(), 1), clusterProjection)
		Expect(err).NotTo(HaveOccurred())
		Expect(clusterProjection.Version()).To(Equal(uint64(2)))
		cluster, ok = clusterProjection.(*projections.Cluster)
		Expect(ok).To(BeTrue())
		Expect(cluster.GetBootstrapToken()).To(Equal(expectedJWT))

		dp := cluster.DomainProjection
		Expect(dp.LastModified).ToNot(BeNil())
	})

	It("can handle ClusterOperatorCertificateIssued events", func() {
		clusterProjector := NewClusterProjector()
		clusterProjection := clusterProjector.NewProjection(uuid.New())
		protoClusterCreatedEventData := &eventdata.ClusterCreated{
			Name:                expectedName,
			Label:               expectedLabel,
			ApiServerAddress:    expectedApiServerAddress,
			CaCertificateBundle: expectedClusterCACertBundle,
		}
		clusterCreatedEventData := es.ToEventDataFromProto(protoClusterCreatedEventData)
		clusterProjection, err := clusterProjector.Project(context.Background(), es.NewEvent(ctx, events.ClusterCreated, clusterCreatedEventData, time.Now().UTC(), aggregates.Cluster, uuid.New(), 1), clusterProjection)
		Expect(err).NotTo(HaveOccurred())

		Expect(clusterProjection.Version()).To(Equal(uint64(1)))

		cluster, ok := clusterProjection.(*projections.Cluster)
		Expect(ok).To(BeTrue())
		Expect(cluster.GetName()).To(Equal(expectedName))
		Expect(cluster.GetLabel()).To(Equal(expectedLabel))
		Expect(cluster.GetApiServerAddress()).To(Equal(expectedApiServerAddress))
		Expect(cluster.GetClusterCACertBundle()).To(Equal(expectedClusterCACertBundle))

		protoClusterOperatorCertificateIssuedEventData := &eventdata.ClusterCertificateIssued{
			Ca:          expectedM8CA,
			Certificate: expectedClusterCertificate,
		}
		clusterOperatorCertificateIssuedEventData := es.ToEventDataFromProto(protoClusterOperatorCertificateIssuedEventData)
		clusterProjection, err = clusterProjector.Project(context.Background(), es.NewEvent(ctx, events.ClusterOperatorCertificateIssued, clusterOperatorCertificateIssuedEventData, time.Now().UTC(), aggregates.Cluster, uuid.New(), 1), clusterProjection)
		Expect(err).NotTo(HaveOccurred())
		Expect(clusterProjection.Version()).To(Equal(uint64(2)))
		cluster, ok = clusterProjection.(*projections.Cluster)
		Expect(ok).To(BeTrue())
		certs := cluster.GetClusterCertificates()
		Expect(certs.GetCa()).To(Equal(expectedM8CA))
		Expect(certs.GetCertificate()).To(Equal(expectedClusterCertificate))

		dp := cluster.DomainProjection
		Expect(dp.LastModified).ToNot(BeNil())
	})

})
