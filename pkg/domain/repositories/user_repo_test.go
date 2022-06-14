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

var _ = Describe("domain/user_repo", func() {
	userId := uuid.New()
	adminUser := &projections.User{User: &projectionsApi.User{Id: userId.String(), Name: "admin", Email: "admin@monoskope.io"}, DomainProjection: projections.NewDomainProjection()}

	adminRoleBinding := projections.NewUserRoleBinding(uuid.New())
	adminRoleBinding.UserId = adminUser.Id
	adminRoleBinding.Role = string(roles.Admin)
	adminRoleBinding.Scope = string(scopes.System)

	It("can read/write projections", func() {
		inMemoryRoleRepo := es_repos.NewInMemoryRepository[*projections.UserRoleBinding]()
		err := inMemoryRoleRepo.Upsert(context.Background(), adminRoleBinding)
		Expect(err).NotTo(HaveOccurred())

		userRoleBindingRepo := NewUserRoleBindingRepository(inMemoryRoleRepo)

		inMemoryUserRepo := es_repos.NewInMemoryRepository[*projections.User]()
		userRepo := NewUserRepository(inMemoryUserRepo, userRoleBindingRepo)

		err = inMemoryUserRepo.Upsert(context.Background(), adminUser)
		Expect(err).NotTo(HaveOccurred())

		user, err := userRepo.ByEmail(context.Background(), adminUser.Email)
		Expect(err).NotTo(HaveOccurred())
		Expect(user).To(Equal(adminUser))
		Expect(user.Roles).ToNot(BeNil())
		Expect(len(user.Roles)).To(BeNumerically("==", 1))
		Expect(user.Roles[0]).To(Equal(adminRoleBinding.Proto()))
	})
})
