package eventsourcing

import (
	"context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("MetadataManager", func() {
	ctx := context.Background()

	var notExistingKey EventMetadataKey
	var existingKey EventMetadataKey

	It("can't get not existing", func() {
		manager := NewMetadataManagerFromContext(ctx)
		val, ok := manager.Get(notExistingKey)
		Expect(val).To(BeNil())
		Expect(ok).To(BeFalse())
	})
	It("can set a value", func() {
		manager := NewMetadataManagerFromContext(ctx)
		val, err := manager.
			Set(existingKey, true).
			GetBool(existingKey)

		Expect(err).To(Not(HaveOccurred()))
		Expect(val).To(BeTrue())
	})
	It("can get from existing context", func() {
		manager := NewMetadataManagerFromContext(ctx)
		nuCtx := manager.
			Set(existingKey, "test").
			GetContext()
		Expect(nuCtx).To(Not(BeNil()))

		nuManager := NewMetadataManagerFromContext(nuCtx)

		val, err := nuManager.
			Set(existingKey, "test").
			GetString(existingKey)

		Expect(err).To(Not(HaveOccurred()))
		Expect(val).To(Not(BeNil()))
		Expect(val).To(Equal("test"))
	})
})
