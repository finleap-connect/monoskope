package repositories

import (
	"context"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/event_sourcing"
)

var _ = Describe("repositories/in_memory", func() {
	testProjection := newTestProjection(uuid.New())
	testReadWrite := func(repo es.Repository) {
		err := repo.Upsert(context.Background(), testProjection)
		Expect(err).NotTo(HaveOccurred())

		err = repo.Upsert(context.Background(), testProjection)
		Expect(err).NotTo(HaveOccurred())

		projection, err := repo.ById(context.Background(), testProjection.GetId())
		Expect(err).NotTo(HaveOccurred())
		Expect(projection).To(Equal(testProjection))

		projections, err := repo.All(context.Background())
		Expect(err).NotTo(HaveOccurred())
		Expect(projections).ToNot(BeNil())
		Expect(len(projections)).To(BeNumerically("==", 1))
		Expect(projections[0]).To(Equal(testProjection))
	}
	testRemove := func(repo es.Repository) {
		err := repo.Upsert(context.Background(), testProjection)
		Expect(err).NotTo(HaveOccurred())

		err = repo.Remove(context.Background(), testProjection.GetId())
		Expect(err).NotTo(HaveOccurred())

		projections, err := repo.All(context.Background())
		Expect(err).NotTo(HaveOccurred())
		Expect(projections).ToNot(BeNil())
		Expect(len(projections)).To(BeNumerically("==", 0))
	}

	It("can read/write projections", func() {
		testReadWrite(NewInMemoryRepository())
	})
	It("can remove projections", func() {
		testRemove(NewInMemoryRepository())
	})
})
