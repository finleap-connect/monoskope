package metadata

import (
	"context"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gitlab.figo.systems/platform/monoskope/monoskope/internal/gateway"
)

var _ = Describe("Managing Metadata", func() {
	It("should set user information", func() {

		expectedUserId := uuid.New()

		ctx := context.Background()
		mdManager, err := NewDomainMetadataManager(ctx)
		Expect(err).ToNot(HaveOccurred())

		mdManager.SetUserInformation(&UserInformation{
			Id:     expectedUserId,
			Name:   "admin",
			Email:  "admin@monoskope.io",
			Issuer: "monoskope",
		})

		Expect(mdManager.GetMetadata()[gateway.HeaderAuthId]).To(Equal(expectedUserId.String()))

	})
})
