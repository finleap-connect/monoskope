// Copyright 2022 Monoskope Authors
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

	projectionsApi "github.com/finleap-connect/monoskope/pkg/api/domain/projections"
	"github.com/finleap-connect/monoskope/pkg/domain/constants/roles"
	"github.com/finleap-connect/monoskope/pkg/domain/constants/scopes"
	"github.com/finleap-connect/monoskope/pkg/domain/projections"
	es_repos "github.com/finleap-connect/monoskope/pkg/eventsourcing/repositories"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("pkg/domain/repositories/clusterAccessRepository", func() {
	tenantId := uuid.New()
	adminUserId := uuid.New()
	otherUserId := uuid.New()
	clusterId := uuid.New()

	adminUser := &projections.User{User: &projectionsApi.User{Id: adminUserId.String(), Name: "admin", Email: "admin@monoskope.io"}}
	otherUser := &projections.User{User: &projectionsApi.User{Id: otherUserId.String(), Name: "otheruser", Email: "otheruser@monoskope.io"}}

	adminRoleBinding := projections.NewUserRoleBinding(uuid.New())
	adminRoleBinding.UserId = adminUser.Id
	adminRoleBinding.Role = roles.Admin.String()
	adminRoleBinding.Scope = scopes.Tenant.String()
	adminRoleBinding.Resource = tenantId.String()

	otherUserRoleBinding := projections.NewUserRoleBinding(uuid.New())
	otherUserRoleBinding.UserId = otherUser.Id
	otherUserRoleBinding.Role = roles.User.String()
	otherUserRoleBinding.Scope = scopes.Tenant.String()
	otherUserRoleBinding.Resource = tenantId.String()

	cluster := projections.NewClusterProjection(clusterId).(*projections.Cluster)
	cluster.Name = "test-cluster"
	cluster.DisplayName = "Test Cluster"

	binding := projections.NewTenantClusterBindingProjection(uuid.New()).(*projections.TenantClusterBinding)
	binding.ClusterId = clusterId.String()
	binding.TenantId = tenantId.String()

	It("can read/write projections", func() {
		inMemoryRoleRepo := es_repos.NewInMemoryRepository()
		Expect(inMemoryRoleRepo.Upsert(context.Background(), adminRoleBinding)).NotTo(HaveOccurred())
		Expect(inMemoryRoleRepo.Upsert(context.Background(), otherUserRoleBinding)).NotTo(HaveOccurred())

		inMemoryClusterRepo := es_repos.NewInMemoryRepository()
		Expect(inMemoryClusterRepo.Upsert(context.Background(), cluster)).NotTo(HaveOccurred())

		clusterRepo := NewClusterRepository(inMemoryClusterRepo)
		userRoleBindingRepo := NewUserRoleBindingRepository(inMemoryRoleRepo)

		inMemoryTenantClusterBindingRepo := es_repos.NewInMemoryRepository()
		Expect(inMemoryTenantClusterBindingRepo.Upsert(context.Background(), binding)).NotTo(HaveOccurred())
		tenantClusterBindingRepo := NewTenantClusterBindingRepository(inMemoryTenantClusterBindingRepo)

		clusterAccessRepo := NewClusterAccessRepository(tenantClusterBindingRepo, clusterRepo, userRoleBindingRepo)

		clusters, err := clusterAccessRepo.GetClustersAccessibleByUserId(context.Background(), otherUserId)
		Expect(err).NotTo(HaveOccurred())
		Expect(clusters).NotTo(BeEmpty())
		Expect(len(clusters)).To(BeNumerically("==", 1))
		Expect(clusters[0].Id).To(Equal(clusterId.String()))

		clusters, err = clusterAccessRepo.GetClustersAccessibleByTenantId(context.Background(), tenantId)
		Expect(err).NotTo(HaveOccurred())
		Expect(clusters).NotTo(BeEmpty())
		Expect(len(clusters)).To(BeNumerically("==", 1))
		Expect(clusters[0].Id).To(Equal(clusterId.String()))
	})
})
