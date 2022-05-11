// Copyright 2022 Monoskope Authors
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

	"github.com/PuerkitoBio/goquery"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Gateway Auth Server", func() {
	It("can retrieve openid conf", func() {
		res, err := testEnv.HttpClient.Get(fmt.Sprintf("http://%s/.well-known/openid-configuration", localAddrOIDCProviderServer))
		Expect(err).NotTo(HaveOccurred())
		Expect(res.StatusCode).To(Equal(http.StatusOK))

		doc, err := goquery.NewDocumentFromReader(res.Body)
		Expect(err).NotTo(HaveOccurred())
		docText := doc.Text()
		Expect(docText).NotTo(BeEmpty())
	})
	It("can retrieve jwks", func() {
		res, err := testEnv.HttpClient.Get(fmt.Sprintf("http://%s/keys", localAddrOIDCProviderServer))
		Expect(err).NotTo(HaveOccurred())
		Expect(res.StatusCode).To(Equal(http.StatusOK))

		doc, err := goquery.NewDocumentFromReader(res.Body)
		Expect(err).NotTo(HaveOccurred())
		Expect(doc.Text()).NotTo(BeEmpty())
	})
})

var _ = Describe("Checks", func() {
	It("can do readiness checks", func() {
		res, err := testEnv.HttpClient.Get(fmt.Sprintf("http://%s/readyz", localAddrOIDCProviderServer))
		Expect(err).NotTo(HaveOccurred())
		Expect(res.StatusCode).To(Equal(http.StatusOK))
	})
})
