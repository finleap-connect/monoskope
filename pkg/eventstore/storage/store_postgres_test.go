package storage

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("storage/postgres", func() {
	It("can create new event store", func() {
		es, err := NewPostgresEventStore(env.DB)
		Expect(err).ToNot(HaveOccurred())
		Expect(es).ToNot(BeNil())
	})
})
