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

package auth

import (
	"fmt"
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"golang.org/x/sync/errgroup"
)

var _ = Describe("monoctl auth", func() {
	It("can go through fake oidc-flow", func() {
		ready := make(chan string, 1)
		oidcClientServer, err := newOidcClientServer(ready)
		Expect(err).ToNot(HaveOccurred())
		defer oidcClientServer.Close()

		redirectURI := oidcClientServer.RedirectURI
		log.Info("RedirectURI: " + redirectURI)
		authCodeUrl := getAuthCodeUrl(redirectURI, "")
		Expect(err).NotTo(HaveOccurred())
		log.Info("AuthCodeUrl: " + authCodeUrl)

		var authCode string
		var statusCode int
		var eg errgroup.Group
		var innerErr error
		eg.Go(func() error {
			defer GinkgoRecover()
			var innerErr error
			authCode, innerErr = oidcClientServer.ReceiveCodeViaLocalServer(ctx, authCodeUrl, "")
			return innerErr
		})
		eg.Go(func() error {
			defer GinkgoRecover()
			log.Info("wait for oidc client server to get ready...")
			<-ready
			res, err := httpClient.Get(authCodeUrl)
			if err == nil {
				statusCode = res.StatusCode
			}
			return innerErr
		})
		Expect(eg.Wait()).NotTo(HaveOccurred())
		Expect(statusCode).To(Equal(http.StatusOK))
		Expect(authCode).ToNot(BeNil())
	})
})

func newOidcClientServer(ready chan<- string) (*Server, error) {
	serverConf := &Config{
		LocalServerBindAddress: []string{
			fmt.Sprintf("%s%s", RedirectURLHostname, RedirectURLPort),
		},
		RedirectURLHostname:  RedirectURLHostname,
		LocalServerReadyChan: ready,
	}
	server, err := NewServer(serverConf)
	if err != nil {
		return nil, err
	}
	return server, nil
}
