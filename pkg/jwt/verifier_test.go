package jwt

import (
	"time"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gopkg.in/square/go-jose.v2/jwt"
)

var _ = Describe("jwt/verifier", func() {
	It("can verify a JWT", func() {
		signer := NewSigner(testEnv.privateKeyFile)
		Expect(signer).ToNot(BeNil())

		claims := jwt.Claims{
			ID:      uuid.New().String(),
			Issuer:  "me",
			Subject: "you",
			Audience: jwt.Audience{
				"monoskope",
			},
			Expiry:   jwt.NewNumericDate(time.Now().Add(1 * time.Minute)),
			IssuedAt: jwt.NewNumericDate(time.Now()),
		}

		rawJWT, err := signer.GenerateSignedToken(claims)
		Expect(err).ToNot(HaveOccurred())
		Expect(rawJWT).ToNot(BeEmpty())
		testEnv.Log.Info("JWT created.", "JWT", rawJWT)

		verifier := NewVerifier(testEnv.publicKeyFile)
		Expect(verifier).ToNot(BeNil())

		claimsFromJWT := jwt.Claims{}
		err = verifier.Verify(rawJWT, &claimsFromJWT)
		Expect(err).ToNot(HaveOccurred())
		Expect(claims).To(Equal(claimsFromJWT))
	})
})
