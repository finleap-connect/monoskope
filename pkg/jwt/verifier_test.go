// Copyright 2021 Monoskope Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
		signer := testEnv.CreateSigner()
		Expect(signer).ToNot(BeNil())

		claims := jwt.Claims{
			ID:      uuid.New().String(),
			Issuer:  "me",
			Subject: "you",
			Audience: jwt.Audience{
				"monoskope",
			},
			Expiry:   jwt.NewNumericDate(time.Now().Add(5 * time.Minute)),
			IssuedAt: jwt.NewNumericDate(time.Now()),
		}

		rawJWT, err := signer.GenerateSignedToken(claims)
		Expect(err).ToNot(HaveOccurred())
		Expect(rawJWT).ToNot(BeEmpty())
		testEnv.Log.Info("JWT created.", "JWT", rawJWT)

		verifier, err := testEnv.CreateVerifier()
		Expect(err).ToNot(HaveOccurred())
		Expect(verifier).ToNot(BeNil())
		defer verifier.Close()

		claimsFromJWT := jwt.Claims{}
		err = verifier.Verify(rawJWT, &claimsFromJWT)
		Expect(err).ToNot(HaveOccurred())
		Expect(claims).To(Equal(claimsFromJWT))

		err = testEnv.RotateCertificate()
		Expect(err).ToNot(HaveOccurred())
		time.Sleep(500 * time.Millisecond)

		claimsFromJWT = jwt.Claims{}
		err = verifier.Verify(rawJWT, &claimsFromJWT)
		Expect(err).To(HaveOccurred())

		rawJWT, err = signer.GenerateSignedToken(claims)
		Expect(err).ToNot(HaveOccurred())
		Expect(rawJWT).ToNot(BeEmpty())
		testEnv.Log.Info("JWT created.", "JWT", rawJWT)

		claimsFromJWT = jwt.Claims{}
		err = verifier.Verify(rawJWT, &claimsFromJWT)
		Expect(err).ToNot(HaveOccurred())
		Expect(claims).To(Equal(claimsFromJWT))

		jwks := verifier.JWKS()
		Expect(jwks).ToNot(BeNil())
	})
})
