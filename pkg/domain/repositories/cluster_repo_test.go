// Copyright 2021 Monoskope Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package repositories

import (
	"context"
	"time"

	"github.com/finleap-connect/monoskope/pkg/api/domain/common"
	projectionsApi "github.com/finleap-connect/monoskope/pkg/api/domain/projections"
	projections "github.com/finleap-connect/monoskope/pkg/domain/projections"
	es_repos "github.com/finleap-connect/monoskope/pkg/eventsourcing/repositories"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	timestamp "google.golang.org/protobuf/types/known/timestamppb"
)

var (
	expectedClusterName = "the one cluster"
	// expectedClusterLabel        = "one-cluster"
	// expectedApiServerAddress    = "one.example.com"
	// expectedClusterCACertBundle = []byte("This should be a certificate")
)

var _ = Describe("domain/cluster_repo", func() {

	clusterId := uuid.New()

	userId := uuid.New()
	adminUser := &projections.User{User: &projectionsApi.User{Id: userId.String(), Name: "admin", Email: "admin@monoskope.io"}}

	adminRoleBinding := projections.NewUserRoleBinding(uuid.New())
	adminRoleBinding.UserId = adminUser.Id
	adminRoleBinding.Role = common.Role_admin.String()
	adminRoleBinding.Scope = common.Scope_system.String()

	newCluster := projections.NewClusterProjection(clusterId).(*projections.Cluster)
	newCluster.Name = expectedClusterName
	newCluster.Created = timestamp.New(time.Now())

	It("can retrieve cluster by name", func() {
		inMemClusterRepo := es_repos.NewInMemoryRepository()
		clusterRepo := NewClusterRepository(inMemClusterRepo)

		err := inMemClusterRepo.Upsert(context.Background(), newCluster)
		Expect(err).NotTo(HaveOccurred())
		cluster, err := clusterRepo.ByClusterName(context.Background(), expectedClusterName)
		Expect(err).NotTo(HaveOccurred())

		Expect(cluster.Name).To(Equal(expectedClusterName))
		Expect(cluster.Created).NotTo(BeNil())
	})

	It("can retrieve cluster by ID", func() {
		inMemClusterRepo := es_repos.NewInMemoryRepository()
		clusterRepo := NewClusterRepository(inMemClusterRepo)

		err := inMemClusterRepo.Upsert(context.Background(), newCluster)
		Expect(err).NotTo(HaveOccurred())
		cluster, err := clusterRepo.ByClusterId(context.Background(), clusterId.String())
		Expect(err).NotTo(HaveOccurred())

		Expect(cluster.Name).To(Equal(expectedClusterName))
		Expect(cluster.Created).NotTo(BeNil())
	})
})
