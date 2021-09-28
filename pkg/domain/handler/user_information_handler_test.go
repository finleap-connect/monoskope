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

package handler

import (
	"context"

	cmddata "github.com/finleap-connect/monoskope/pkg/api/domain/commanddata"
	projectionsApi "github.com/finleap-connect/monoskope/pkg/api/domain/projections"
	cmd "github.com/finleap-connect/monoskope/pkg/domain/commands"
	"github.com/finleap-connect/monoskope/pkg/domain/constants/aggregates"
	"github.com/finleap-connect/monoskope/pkg/domain/constants/commands"
	"github.com/finleap-connect/monoskope/pkg/domain/constants/roles"
	"github.com/finleap-connect/monoskope/pkg/domain/constants/scopes"
	metadata "github.com/finleap-connect/monoskope/pkg/domain/metadata"
	projections "github.com/finleap-connect/monoskope/pkg/domain/projections"
	"github.com/finleap-connect/monoskope/pkg/domain/repositories"
	es "github.com/finleap-connect/monoskope/pkg/eventsourcing"
	es_repos "github.com/finleap-connect/monoskope/pkg/eventsourcing/repositories"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("domain/handler", func() {
	adminUser := &projections.User{User: &projectionsApi.User{
		Id:    uuid.New().String(),
		Name:  "admin",
		Email: "admin@monoskope.io",
	}}
	someUser := &projections.User{User: &projectionsApi.User{
		Id:    uuid.New().String(),
		Name:  "some.user",
		Email: "some.user@monoskope.io",
	}}

	adminRoleBinding := projections.NewUserRoleBinding(uuid.New())
	adminRoleBinding.UserId = adminUser.Id
	adminRoleBinding.Role = roles.Admin.String()
	adminRoleBinding.Scope = scopes.System.String()

	inMemoryRoleRepo := es_repos.NewInMemoryRepository()
	err := inMemoryRoleRepo.Upsert(context.Background(), adminRoleBinding)
	Expect(err).NotTo(HaveOccurred())

	userRoleBindingRepo := repositories.NewUserRoleBindingRepository(inMemoryRoleRepo)

	inMemoryUserRepo := es_repos.NewInMemoryRepository()
	userRepo := repositories.NewUserRepository(inMemoryUserRepo, userRoleBindingRepo)

	userBase := es.NewBaseCommand(uuid.New(), aggregates.User, commands.CreateUser)
	tenantBase := es.NewBaseCommand(uuid.New(), aggregates.Tenant, commands.CreateTenant)
	roleBindingBase := es.NewBaseCommand(uuid.New(), aggregates.UserRoleBinding, commands.CreateUserRoleBinding)

	handler := NewUserInformationHandler(userRepo)

	It("system admin can create users", func() {
		err = inMemoryUserRepo.Upsert(context.Background(), adminUser)
		Expect(err).NotTo(HaveOccurred())

		manager, err := metadata.NewDomainMetadataManager(context.Background())
		Expect(err).ToNot(HaveOccurred())

		manager.SetUserInformation(&metadata.UserInformation{Email: adminUser.Email})

		_, err = handler.HandleCommand(manager.GetContext(), &cmd.CreateUserCommand{
			BaseCommand:           userBase,
			CreateUserCommandData: cmddata.CreateUserCommandData{Email: someUser.Email},
		})
		Expect(err).ToNot(HaveOccurred())
	})
	It("system admin can create tenants", func() {
		manager, err := metadata.NewDomainMetadataManager(context.Background())
		Expect(err).ToNot(HaveOccurred())

		manager.SetUserInformation(&metadata.UserInformation{Email: adminUser.Email})

		_, err = handler.HandleCommand(manager.GetContext(), &cmd.CreateTenantCommand{
			BaseCommand:             tenantBase,
			CreateTenantCommandData: cmddata.CreateTenantCommandData{Name: "dieter", Prefix: "dt"},
		})
		Expect(err).ToNot(HaveOccurred())
	})
	It("admin can create rolebinding for any user", func() {
		manager, err := metadata.NewDomainMetadataManager(context.Background())
		Expect(err).ToNot(HaveOccurred())

		manager.SetUserInformation(&metadata.UserInformation{Email: adminUser.Email})
		_, err = handler.HandleCommand(manager.GetContext(), &cmd.CreateUserRoleBindingCommand{
			BaseCommand: roleBindingBase,
			CreateUserRoleBindingCommandData: cmddata.CreateUserRoleBindingCommandData{
				UserId: someUser.Id,
				Role:   roles.Admin.String(),
				Scope:  scopes.System.String(),
			},
		})
		Expect(err).ToNot(HaveOccurred())
	})
	It("superuser can create admin rolebinding for any user", func() {
		manager, err := metadata.NewDomainMetadataManager(context.Background())
		Expect(err).ToNot(HaveOccurred())

		manager.SetUserInformation(&metadata.UserInformation{Email: adminUser.Email})
		command := &cmd.CreateUserRoleBindingCommand{
			BaseCommand: roleBindingBase,
			CreateUserRoleBindingCommandData: cmddata.CreateUserRoleBindingCommandData{
				UserId: someUser.Id,
				Role:   roles.Admin.String(),
				Scope:  scopes.System.String(),
			},
		}
		_, err = handler.HandleCommand(manager.GetContext(), command)
		Expect(err).ToNot(HaveOccurred())
	})
})
