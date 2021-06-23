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
		signer := testEnv.CreateSigner()
		Expect(signer).ToNot(BeNil())

		rawJWT, err := signer.GenerateSignedToken(jwt.Claims{
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
		testEnv.Log.Info("JWT created.", "JWT", rawJWT)
	})
})
