package repositories

import (
	"context"
	"time"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	projectionsApi "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/domain/projections"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/constants/roles"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/constants/scopes"
	projections "gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/projections"
	es_repos "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing/repositories"
	timestamp "google.golang.org/protobuf/types/known/timestamppb"
)

var (
	expectedClusterName = "the one cluster"
	// expectedClusterLabel        = "one-cluster"
	// expectedApiServerAddress    = "one.example.com"
	expectedClusterCACertBundle = []byte("This should be a certificate")
)

var _ = Describe("domain/cluster_repo", func() {

	clusterId := uuid.New()

	userId := uuid.New()
	adminUser := &projections.User{User: &projectionsApi.User{Id: userId.String(), Name: "admin", Email: "admin@monoskope.io"}}

	adminRoleBinding := projections.NewUserRoleBinding(uuid.New())
	adminRoleBinding.UserId = adminUser.Id
	adminRoleBinding.Role = roles.Admin.String()
	adminRoleBinding.Scope = scopes.System.String()

	newCluster := projections.NewClusterProjection(clusterId)
	newCluster.Name = expectedClusterName
	newCluster.ClusterCACertBundle = expectedClusterCACertBundle
	newCluster.Created = timestamp.New(time.Now())

	It("can retrieve cluster by name", func() {
		inMemoryRoleRepo := es_repos.NewInMemoryRepository()
		err := inMemoryRoleRepo.Upsert(context.Background(), adminRoleBinding)
		Expect(err).NotTo(HaveOccurred())

		userRoleBindingRepo := NewUserRoleBindingRepository(inMemoryRoleRepo)

		inMemoryUserRepo := es_repos.NewInMemoryRepository()
		userRepo := NewUserRepository(inMemoryUserRepo, userRoleBindingRepo)
		err = inMemoryUserRepo.Upsert(context.Background(), adminUser)
		Expect(err).NotTo(HaveOccurred())

		inMemClusterRepo := es_repos.NewInMemoryRepository()
		clusterRepo := NewClusterRepository(inMemClusterRepo, userRepo)

		err = inMemClusterRepo.Upsert(context.Background(), newCluster)
		Expect(err).NotTo(HaveOccurred())
		cluster, err := clusterRepo.ByClusterName(context.Background(), expectedClusterName)
		Expect(err).NotTo(HaveOccurred())

		Expect(cluster.Name).To(Equal(expectedClusterName))
		Expect(cluster.Created).NotTo(BeNil())
	})

	It("can retrieve cluster by ID", func() {
		inMemoryRoleRepo := es_repos.NewInMemoryRepository()
		err := inMemoryRoleRepo.Upsert(context.Background(), adminRoleBinding)
		Expect(err).NotTo(HaveOccurred())

		userRoleBindingRepo := NewUserRoleBindingRepository(inMemoryRoleRepo)

		inMemoryUserRepo := es_repos.NewInMemoryRepository()
		userRepo := NewUserRepository(inMemoryUserRepo, userRoleBindingRepo)
		err = inMemoryUserRepo.Upsert(context.Background(), adminUser)
		Expect(err).NotTo(HaveOccurred())

		inMemClusterRepo := es_repos.NewInMemoryRepository()
		clusterRepo := NewClusterRepository(inMemClusterRepo, userRepo)

		err = inMemClusterRepo.Upsert(context.Background(), newCluster)
		Expect(err).NotTo(HaveOccurred())
		cluster, err := clusterRepo.ByClusterId(context.Background(), clusterId.String())
		Expect(err).NotTo(HaveOccurred())

		Expect(cluster.Name).To(Equal(expectedClusterName))
		Expect(cluster.Created).NotTo(BeNil())
	})
	It("can query a cluster's ca bundle by cluster ID", func() {
		inMemoryRoleRepo := es_repos.NewInMemoryRepository()
		err := inMemoryRoleRepo.Upsert(context.Background(), adminRoleBinding)
		Expect(err).NotTo(HaveOccurred())

		userRoleBindingRepo := NewUserRoleBindingRepository(inMemoryRoleRepo)

		inMemoryUserRepo := es_repos.NewInMemoryRepository()
		userRepo := NewUserRepository(inMemoryUserRepo, userRoleBindingRepo)
		err = inMemoryUserRepo.Upsert(context.Background(), adminUser)
		Expect(err).NotTo(HaveOccurred())

		inMemClusterRepo := es_repos.NewInMemoryRepository()
		clusterRepo := NewClusterRepository(inMemClusterRepo, userRepo)

		err = inMemClusterRepo.Upsert(context.Background(), newCluster)
		Expect(err).NotTo(HaveOccurred())
		bundle, err := clusterRepo.GetCACertificateBundle(context.Background(), clusterId.String())
		Expect(err).NotTo(HaveOccurred())
		Expect(bundle).To(Equal(expectedClusterCACertBundle))
	})
})
