package handler

import (
	"context"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	commandsApi "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/domain/commands"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/domain/common"
	projectionsApi "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/domain/projections"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/commands"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/constants/roles"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/constants/scopes"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/errors"
	metadata "gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/metadata"
	projections "gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/projections"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/repositories"
	es_repos "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing/repositories"
)

var _ = Describe("domain/handler", func() {
	adminUser := &projections.User{User: projectionsApi.User{
		Id:    uuid.New().String(),
		Name:  "admin",
		Email: "admin@monoskope.io",
	}}
	someUser := &projections.User{User: projectionsApi.User{
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

	err = inMemoryUserRepo.Upsert(context.Background(), adminUser)
	Expect(err).NotTo(HaveOccurred())

	err = inMemoryUserRepo.Upsert(context.Background(), someUser)
	Expect(err).NotTo(HaveOccurred())

	handler := NewAuthorizationHandler(userRepo)

	It("user can't create other users", func() {
		manager, err := metadata.NewDomainMetadataManager(context.Background())
		Expect(err).ToNot(HaveOccurred())

		err = manager.SetUserInformation(&metadata.UserInformation{Email: someUser.Email})
		Expect(err).ToNot(HaveOccurred())

		err = handler.HandleCommand(manager.GetContext(), &commands.CreateUserCommand{CreateUserCommand: commandsApi.CreateUserCommand{
			UserMetadata: &common.UserMetadata{
				Email: "janedoe@monoskope.io",
			},
		}})
		Expect(err).To(HaveOccurred())
		Expect(err).To(Equal(errors.ErrUnauthorized))
	})
	It("user can create themselves", func() {
		manager, err := metadata.NewDomainMetadataManager(context.Background())
		Expect(err).ToNot(HaveOccurred())

		err = manager.SetUserInformation(&metadata.UserInformation{Email: adminUser.Email})
		Expect(err).ToNot(HaveOccurred())

		err = handler.HandleCommand(manager.GetContext(), &commands.CreateUserCommand{CreateUserCommand: commandsApi.CreateUserCommand{
			UserMetadata: &common.UserMetadata{
				Email: adminUser.Email,
			},
		}})
		Expect(err).ToNot(HaveOccurred())
	})
	It("system admin can create users", func() {
		manager, err := metadata.NewDomainMetadataManager(context.Background())
		Expect(err).ToNot(HaveOccurred())

		err = manager.SetUserInformation(&metadata.UserInformation{Email: adminUser.Email})
		Expect(err).ToNot(HaveOccurred())

		err = handler.HandleCommand(manager.GetContext(), &commands.CreateUserCommand{CreateUserCommand: commandsApi.CreateUserCommand{
			UserMetadata: &common.UserMetadata{
				Email: someUser.Email,
			},
		}})
		Expect(err).ToNot(HaveOccurred())
	})
	It("user can't make themselfes admin", func() {
		manager, err := metadata.NewDomainMetadataManager(context.Background())
		Expect(err).ToNot(HaveOccurred())

		err = manager.SetUserInformation(&metadata.UserInformation{Email: someUser.Email})
		Expect(err).ToNot(HaveOccurred())

		err = handler.HandleCommand(manager.GetContext(), &commands.CreateUserRoleBindingCommand{
			CreateUserRoleBindingCommand: commandsApi.CreateUserRoleBindingCommand{
				UserId: someUser.Id,
				Role:   roles.Admin.String(),
				Scope:  scopes.System.String(),
			},
		})
		Expect(err).To(HaveOccurred())
		Expect(err).To(Equal(errors.ErrUnauthorized))
	})
	It("admin can create rolebinding for any user", func() {
		manager, err := metadata.NewDomainMetadataManager(context.Background())
		Expect(err).ToNot(HaveOccurred())

		err = manager.SetUserInformation(&metadata.UserInformation{Email: adminUser.Email})
		Expect(err).ToNot(HaveOccurred())

		err = handler.HandleCommand(manager.GetContext(), &commands.CreateUserRoleBindingCommand{
			CreateUserRoleBindingCommand: commandsApi.CreateUserRoleBindingCommand{
				UserId: someUser.Id,
				Role:   roles.Admin.String(),
				Scope:  scopes.System.String(),
			},
		})
		Expect(err).ToNot(HaveOccurred())
	})
})
