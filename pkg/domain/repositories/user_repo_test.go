package repositories

import (
	"context"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	projections "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/queryhandler"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/constants/roles"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/constants/scopes"
	es_repos "gitlab.figo.systems/platform/monoskope/monoskope/pkg/event_sourcing/repositories"
)

var _ = Describe("domain/user_repo", func() {
	userId := uuid.New()
	adminUser := &projections.User{Id: userId.String(), Name: "admin", Email: "admin@monoskope.io"}
	adminRoleBinding := &projections.UserRoleBinding{Id: uuid.New().String(), UserId: adminUser.Id, Role: roles.Admin.String(), Scope: scopes.System.String(), Resource: ""}

	It("can read/write projections", func() {
		inMemoryRoleRepo := es_repos.NewInMemoryRepository()
		err := inMemoryRoleRepo.Upsert(context.Background(), adminRoleBinding)
		Expect(err).NotTo(HaveOccurred())

		userRoleBindingRepo := NewUserRoleBindingRepository(inMemoryRoleRepo)

		inMemoryUserRepo := es_repos.NewInMemoryRepository()
		userRepo := NewUserRepository(inMemoryUserRepo, userRoleBindingRepo)

		err = inMemoryUserRepo.Upsert(context.Background(), adminUser)
		Expect(err).NotTo(HaveOccurred())

		user, err := userRepo.ByEmail(context.Background(), adminUser.Email)
		Expect(err).NotTo(HaveOccurred())
		Expect(user).To(Equal(adminUser))
		Expect(user.Roles).ToNot(BeNil())
		Expect(len(user.Roles)).To(BeNumerically("==", 1))
		Expect(user.Roles[0]).To(Equal(adminRoleBinding))
	})
})
