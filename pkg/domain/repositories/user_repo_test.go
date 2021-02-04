package repositories

import (
	"context"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/constants"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/projections"
	es_repos "gitlab.figo.systems/platform/monoskope/monoskope/pkg/event_sourcing/repositories"
)

var _ = Describe("domain/user_repo", func() {
	adminUser := projections.NewUser(uuid.New(), "admin", "admin@monoskope.io", []*projections.UserRoleBinding{
		projections.NewUserRoleBinding(uuid.New(), constants.Admin, constants.System, ""),
	})

	It("can read/write projections", func() {
		inMemoryRepo := es_repos.NewInMemoryRepository()
		userRepo := NewUserRepository(inMemoryRepo)
		err := userRepo.Upsert(context.Background(), adminUser)
		Expect(err).NotTo(HaveOccurred())

		user, err := userRepo.ByEmail(context.Background(), adminUser.Email())
		Expect(err).NotTo(HaveOccurred())
		Expect(user).To(Equal(adminUser))
	})
})
