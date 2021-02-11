package handler

import (
	"context"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	projections "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/queryhandler"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/commands"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/constants/roles"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/constants/scopes"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/errors"
	metadata "gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/metadata"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/repositories"
	es_repos "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing/repositories"
)

var _ = Describe("domain/handler", func() {
	userId := uuid.New()
	adminUser := &projections.User{Id: userId.String(), Name: "admin", Email: "admin@monoskope.io"}
	adminRoleBinding := &projections.UserRoleBinding{Id: uuid.New().String(), UserId: adminUser.Id, Role: roles.Admin.String(), Scope: scopes.System.String(), Resource: ""}
	ctx := metadata.NewDomainMetadataManager(context.Background()).
		SetUserInformation(&metadata.UserInformation{Email: adminUser.Email}).GetContext()

	It("can handle events", func() {
		inMemoryRoleRepo := es_repos.NewInMemoryRepository()
		err := inMemoryRoleRepo.Upsert(ctx, adminRoleBinding)
		Expect(err).NotTo(HaveOccurred())

		userRoleBindingRepo := repositories.NewUserRoleBindingRepository(inMemoryRoleRepo)

		inMemoryUserRepo := es_repos.NewInMemoryRepository()
		userRepo := repositories.NewUserRepository(inMemoryUserRepo, userRoleBindingRepo)

		err = inMemoryUserRepo.Upsert(ctx, adminUser)
		Expect(err).NotTo(HaveOccurred())

		user, err := userRepo.ByEmail(ctx, adminUser.Email)
		Expect(err).NotTo(HaveOccurred())
		Expect(user).To(Equal(adminUser))
		Expect(user.Roles).ToNot(BeNil())
		Expect(len(user.Roles)).To(BeNumerically("==", 1))
		Expect(user.Roles[0]).To(Equal(adminRoleBinding))

		handler := NewAuthorizationHandler(userRepo)
		err = handler.HandleCommand(ctx, &commands.CreateUserCommand{})
		Expect(err).To(HaveOccurred())
		Expect(err).To(Equal(errors.ErrUnauthorized))
	})
})
