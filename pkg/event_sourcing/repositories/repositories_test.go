package repositories

import (
	"context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/event_sourcing"
)

var _ = Describe("repositories/in_memory", func() {
	It("can read", func() {
		repo := NewInMemoryRepository()
		err := repo.Save(context.Background(), es.NewTestProjection())
		Expect(err).NotTo(HaveOccurred())
	})
})
