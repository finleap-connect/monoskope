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

package gateway

import (
	"fmt"
	"net/http"
	"time"

	"github.com/PuerkitoBio/goquery"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/jwt"
	josejwt "gopkg.in/square/go-jose.v2/jwt"
)

var _ = Describe("Gateway Auth Server", func() {
	It("can retrieve openid conf", func() {
		res, err := env.HttpClient.Get(fmt.Sprintf("http://%s/.well-known/openid-configuration", localAddrAuthServer))
		Expect(err).NotTo(HaveOccurred())
		Expect(res.StatusCode).To(Equal(http.StatusOK))

		doc, err := goquery.NewDocumentFromReader(res.Body)
		Expect(err).NotTo(HaveOccurred())
		docText := doc.Text()
		Expect(docText).NotTo(BeEmpty())
	})
	It("can retrieve jwks", func() {
		res, err := env.HttpClient.Get(fmt.Sprintf("http://%s/keys", localAddrAuthServer))
		Expect(err).NotTo(HaveOccurred())
		Expect(res.StatusCode).To(Equal(http.StatusOK))

		doc, err := goquery.NewDocumentFromReader(res.Body)
		Expect(err).NotTo(HaveOccurred())
		Expect(doc.Text()).NotTo(BeEmpty())
	})
	It("can authenticate with JWT", func() {
		expectedValidity := time.Hour * 1
		token := jwt.NewAuthToken(&jwt.StandardClaims{Name: env.ExistingUser.Name, Email: env.ExistingUser.Email}, localAddrAPIServer, env.ExistingUser.Id, expectedValidity)
		signer := env.JwtTestEnv.CreateSigner()
		signedToken, err := signer.GenerateSignedToken(token)
		Expect(err).NotTo(HaveOccurred())

		req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("http://%s/auth/test", localAddrAuthServer), nil)
		Expect(err).NotTo(HaveOccurred())

		req.Header.Set(HeaderAuthorization, fmt.Sprintf("bearer %s", signedToken))
		res, err := env.HttpClient.Do(req)
		Expect(err).NotTo(HaveOccurred())
		Expect(res.StatusCode).To(Equal(http.StatusOK))
	})
	It("fails authentication with invalid JWT", func() {
		req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("http://%s/auth/test", localAddrAuthServer), nil)
		Expect(err).NotTo(HaveOccurred())

		req.Header.Set(HeaderAuthorization, fmt.Sprintf("bearer %s", "notavalidjwt"))
		res, err := env.HttpClient.Do(req)
		Expect(err).NotTo(HaveOccurred())
		Expect(res.StatusCode).To(Equal(http.StatusUnauthorized))
	})
	It("fails authentication with expired JWT", func() {
		expectedValidity := -30 * time.Minute
		token := jwt.NewAuthToken(&jwt.StandardClaims{Name: env.ExistingUser.Name, Email: env.ExistingUser.Email}, localAddrAPIServer, env.ExistingUser.Id, expectedValidity)
		token.NotBefore = josejwt.NewNumericDate(time.Now().UTC().Add(-1 * time.Hour))

		signer := env.JwtTestEnv.CreateSigner()
		signedToken, err := signer.GenerateSignedToken(token)
		Expect(err).NotTo(HaveOccurred())

		req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("http://%s/auth/test", localAddrAuthServer), nil)
		Expect(err).NotTo(HaveOccurred())

		req.Header.Set(HeaderAuthorization, fmt.Sprintf("bearer %s", signedToken))
		res, err := env.HttpClient.Do(req)
		Expect(err).NotTo(HaveOccurred())
		Expect(res.StatusCode).To(Equal(http.StatusUnauthorized))
	})
	It("fails authentication with not existing user", func() {
		expectedValidity := time.Hour * 1
		token := jwt.NewAuthToken(&jwt.StandardClaims{Name: env.NotExistingUser.Name, Email: env.NotExistingUser.Email}, localAddrAPIServer, env.NotExistingUser.Id, expectedValidity)
		signer := env.JwtTestEnv.CreateSigner()
		signedToken, err := signer.GenerateSignedToken(token)
		Expect(err).NotTo(HaveOccurred())

		req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("http://%s/auth/test", localAddrAuthServer), nil)
		Expect(err).NotTo(HaveOccurred())

		req.Header.Set(HeaderAuthorization, fmt.Sprintf("bearer %s", signedToken))
		res, err := env.HttpClient.Do(req)
		Expect(err).ToNot(HaveOccurred())
		Expect(res.StatusCode).To(Equal(http.StatusUnauthorized))
	})
})

var _ = Describe("Checks", func() {
	It("can do readiness checks", func() {
		res, err := env.HttpClient.Get(fmt.Sprintf("http://%s/readyz", localAddrAuthServer))
		Expect(err).NotTo(HaveOccurred())
		Expect(res.StatusCode).To(Equal(http.StatusOK))
	})
})
