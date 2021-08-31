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

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	projectionsApi "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/domain/projections"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/constants/roles"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/constants/scopes"
	projections "gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/projections"
	es_repos "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing/repositories"
)

var _ = Describe("domain/tenant_user_repo_test", func() {
	tenantId := uuid.New()
	adminUserId := uuid.New()
	otherUserId := uuid.New()
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

	It("can read/write projections", func() {
		inMemoryRoleRepo := es_repos.NewInMemoryRepository()
		err := inMemoryRoleRepo.Upsert(context.Background(), adminRoleBinding)
		Expect(err).NotTo(HaveOccurred())
		err = inMemoryRoleRepo.Upsert(context.Background(), otherUserRoleBinding)
		Expect(err).NotTo(HaveOccurred())

		userRoleBindingRepo := NewUserRoleBindingRepository(inMemoryRoleRepo)
		inMemoryUserRepo := es_repos.NewInMemoryRepository()
		userRepo := NewUserRepository(inMemoryUserRepo, userRoleBindingRepo)

		err = inMemoryUserRepo.Upsert(context.Background(), adminUser)
		Expect(err).NotTo(HaveOccurred())
		err = inMemoryUserRepo.Upsert(context.Background(), otherUser)
		Expect(err).NotTo(HaveOccurred())

		tenantUserRepo := NewTenantUserRepository(userRepo, userRoleBindingRepo)
		users, err := tenantUserRepo.GetTenantUsersById(context.Background(), tenantId)
		Expect(err).NotTo(HaveOccurred())
		Expect(users).NotTo(BeEmpty())
		Expect(len(users)).To(BeNumerically("==", 2))
	})
})
