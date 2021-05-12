package jwt

import (
	"time"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gopkg.in/square/go-jose.v2/jwt"
)

var _ = Describe("jwt/signer", func() {
	It("can sign a JWT", func() {
		signer, err := NewSigner(testEnv.privateKeyFile)
		Expect(err).ToNot(HaveOccurred())
		Expect(signer).ToNot(BeNil())

		rawJWT, err := Sign(signer, &jwt.Claims{
			ID:      uuid.New().String(),
			Issuer:  "me",
			Subject: "you",
			Audience: jwt.Audience{
				"monoskope",
			},
			Expiry:   jwt.NewNumericDate(time.Now().Add(1 * time.Minute)),
			IssuedAt: jwt.NewNumericDate(time.Now()),
		})
		Expect(err).ToNot(HaveOccurred())
		Expect(rawJWT).ToNot(BeEmpty())
	})
})
