package repositories

import (
	"context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/projections"
	es_repos "gitlab.figo.systems/platform/monoskope/monoskope/pkg/event_sourcing/repositories"
)

var _ = Describe("domain/user_repo", func() {
	adminUser := &projections.User{Email: "admin@monoskope.io", Name: "admin"}
	It("can read/write projections", func() {
		inMemoryRepo := es_repos.NewInMemoryRepository()
		userRepo := NewUserRepository(inMemoryRepo)
		err := userRepo.Upsert(context.Background(), adminUser)
		Expect(err).NotTo(HaveOccurred())

		user, err := userRepo.ByEmail(context.Background(), adminUser.Email)
		Expect(err).NotTo(HaveOccurred())
		Expect(user).To(Equal(adminUser))
	})
})
