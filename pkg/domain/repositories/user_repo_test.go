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

var _ = Describe("domain/user_repo", func() {
	userId := uuid.New()
	adminUser := &projections.User{User: &projectionsApi.User{Id: userId.String(), Name: "admin", Email: "admin@monoskope.io"}}

	adminRoleBinding := projections.NewUserRoleBinding(uuid.New())
	adminRoleBinding.UserId = adminUser.Id
	adminRoleBinding.Role = roles.Admin.String()
	adminRoleBinding.Scope = scopes.System.String()

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
		Expect(user.Roles[0]).To(Equal(adminRoleBinding.Proto()))
	})
})
