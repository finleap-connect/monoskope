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
	projections "github.com/finleap-connect/monoskope/pkg/domain/projections"
	es_repos "github.com/finleap-connect/monoskope/pkg/eventsourcing/repositories"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("domain/tenant_user_repo_test", func() {
	tenantId := uuid.New()
	adminUserId := uuid.New()
	otherUserId := uuid.New()

	tenant := &projections.Tenant{
		Tenant: &projectionsApi.Tenant{
			Id:     tenantId.String(),
			Name:   "tenant-a",
			Prefix: "ta",
		},
		DomainProjection: projections.NewDomainProjection(),
	}

	adminUser := &projections.User{User: &projectionsApi.User{Id: adminUserId.String(), Name: "admin", Email: "admin@monoskope.io", Metadata: &projectionsApi.LifecycleMetadata{}}}
	otherUser := &projections.User{User: &projectionsApi.User{Id: otherUserId.String(), Name: "otheruser", Email: "otheruser@monoskope.io", Metadata: &projectionsApi.LifecycleMetadata{}}}

	adminRoleBinding := projections.NewUserRoleBinding(uuid.New())
	adminRoleBinding.UserId = adminUser.Id
	adminRoleBinding.Role = string(roles.Admin)
	adminRoleBinding.Scope = string(scopes.Tenant)
	adminRoleBinding.Resource = tenantId.String()

	otherUserRoleBinding := projections.NewUserRoleBinding(uuid.New())
	otherUserRoleBinding.UserId = otherUser.Id
	otherUserRoleBinding.Role = string(roles.User)
	otherUserRoleBinding.Scope = string(scopes.Tenant)
	otherUserRoleBinding.Resource = tenantId.String()

	It("can read/write projections", func() {
		inMemoryTenantRepo := es_repos.NewInMemoryRepository[*projections.Tenant]()
		tenantRepo := NewTenantRepository(inMemoryTenantRepo)
		err := tenantRepo.Upsert(context.Background(), tenant)
		Expect(err).NotTo(HaveOccurred())

		inMemoryRoleRepo := es_repos.NewInMemoryRepository[*projections.UserRoleBinding]()
		err = inMemoryRoleRepo.Upsert(context.Background(), adminRoleBinding)
		Expect(err).NotTo(HaveOccurred())
		err = inMemoryRoleRepo.Upsert(context.Background(), otherUserRoleBinding)
		Expect(err).NotTo(HaveOccurred())

		userRoleBindingRepo := NewUserRoleBindingRepository(inMemoryRoleRepo)
		inMemoryUserRepo := es_repos.NewInMemoryRepository[*projections.User]()
		userRepo := NewUserRepository(inMemoryUserRepo, userRoleBindingRepo)

		err = inMemoryUserRepo.Upsert(context.Background(), adminUser)
		Expect(err).NotTo(HaveOccurred())
		err = inMemoryUserRepo.Upsert(context.Background(), otherUser)
		Expect(err).NotTo(HaveOccurred())

		tenantUserRepo := NewTenantUserRepository(userRepo, userRoleBindingRepo, tenantRepo)
		users, err := tenantUserRepo.GetTenantUsersById(context.Background(), tenantId)
		Expect(err).NotTo(HaveOccurred())
		Expect(users).NotTo(BeEmpty())
		Expect(len(users)).To(BeNumerically("==", 2))
	})
})
