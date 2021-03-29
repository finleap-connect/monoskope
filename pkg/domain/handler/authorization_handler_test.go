package handler

import (
	"context"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	cmddata "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/domain/commanddata"
	projectionsApi "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/domain/projections"
	cmd "gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/commands"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/constants/aggregates"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/constants/commands"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/constants/roles"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/constants/scopes"
	metadata "gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/metadata"
	projections "gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/projections"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/repositories"
	es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing"
	es_repos "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing/repositories"
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

	adminRoleBinding := &projections.UserRoleBinding{UserRoleBinding: projectionsApi.UserRoleBinding{
		Id:     uuid.New().String(),
		UserId: adminUser.Id,
		Role:   roles.Admin.String(),
		Scope:  scopes.System.String(),
	}}

	inMemoryRoleRepo := es_repos.NewInMemoryRepository()
	err := inMemoryRoleRepo.Upsert(context.Background(), adminRoleBinding)
	Expect(err).NotTo(HaveOccurred())

	userRoleBindingRepo := repositories.NewUserRoleBindingRepository(inMemoryRoleRepo)

	inMemoryUserRepo := es_repos.NewInMemoryRepository()
	userRepo := repositories.NewUserRepository(inMemoryUserRepo, userRoleBindingRepo)

	userBase := es.NewBaseCommand(uuid.New(), aggregates.User, commands.CreateUser)
	tenantBase := es.NewBaseCommand(uuid.New(), aggregates.Tenant, commands.CreateTenant)
	roleBindingBase := es.NewBaseCommand(uuid.New(), aggregates.UserRoleBinding, commands.CreateUserRoleBinding)

	handler := NewAuthorizationHandler(userRepo)

	It("system admin can create users", func() {
		err = inMemoryUserRepo.Upsert(context.Background(), adminUser)
		Expect(err).NotTo(HaveOccurred())

		manager, err := metadata.NewDomainMetadataManager(context.Background())
		Expect(err).ToNot(HaveOccurred())

		manager.SetUserInformation(&metadata.UserInformation{Email: adminUser.Email})

		err = handler.HandleCommand(manager.GetContext(), &cmd.CreateUserCommand{
			BaseCommand:           userBase,
			CreateUserCommandData: cmddata.CreateUserCommandData{Email: someUser.Email},
		})
		Expect(err).ToNot(HaveOccurred())
	})
	It("system admin can create tenants", func() {
		manager, err := metadata.NewDomainMetadataManager(context.Background())
		Expect(err).ToNot(HaveOccurred())

		manager.SetUserInformation(&metadata.UserInformation{Email: adminUser.Email})

		err = handler.HandleCommand(manager.GetContext(), &cmd.CreateTenantCommand{
			BaseCommand:             tenantBase,
			CreateTenantCommandData: cmddata.CreateTenantCommandData{Name: "dieter", Prefix: "dt"},
		})
		Expect(err).ToNot(HaveOccurred())
	})
	It("admin can create rolebinding for any user", func() {
		manager, err := metadata.NewDomainMetadataManager(context.Background())
		Expect(err).ToNot(HaveOccurred())

		manager.SetUserInformation(&metadata.UserInformation{Email: adminUser.Email})
		err = handler.HandleCommand(manager.GetContext(), &cmd.CreateUserRoleBindingCommand{
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
		err = handler.HandleCommand(manager.GetContext(), command)
		Expect(err).ToNot(HaveOccurred())
	})

})
