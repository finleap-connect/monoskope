package eventsourcing

import (
	"context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type metadataVal struct {
	Val string
}

var _ = Describe("MetadataManager", func() {
	ctx := context.Background()

	var notExistingKey = "notExistingKey"
	var existingKey = "existingKey"

	It("can't get not existing", func() {
		manager := NewMetadataManagerFromContext(ctx)
		val, ok := manager.Get(notExistingKey)
		Expect(val).To(BeNil())
		Expect(ok).To(BeFalse())
	})
	It("can set a value", func() {
		manager := NewMetadataManagerFromContext(ctx)

		val := &metadataVal{
			Val: "hello",
		}
		err := manager.SetObject(existingKey, val)
		Expect(err).To(Not(HaveOccurred()))

		valResult := &metadataVal{}
		err = manager.GetObject(existingKey, valResult)
		Expect(err).To(Not(HaveOccurred()))
		Expect(valResult.Val).To(Equal(valResult.Val))
	})
	It("can get from existing context", func() {
		manager := NewMetadataManagerFromContext(ctx)

		val := &metadataVal{
			Val: "hello",
		}
		err := manager.SetObject(existingKey, val)
		Expect(err).To(Not(HaveOccurred()))

		nuCtx := manager.GetContext()
		Expect(nuCtx).To(Not(BeNil()))

		nuManager := NewMetadataManagerFromContext(nuCtx)
		err = nuManager.SetObject(existingKey, val)
		Expect(err).To(Not(HaveOccurred()))

		valResult := &metadataVal{}
		err = nuManager.GetObject(existingKey, valResult)
		Expect(err).To(Not(HaveOccurred()))
		Expect(valResult.Val).To(Equal(valResult.Val))
	})
})
