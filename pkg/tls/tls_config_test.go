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

package tls

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("pkg/tls/tls_config", func() {
	It("can consume", func() {
		err := testEnv.CreateCertificate()
		Expect(err).ToNot(HaveOccurred())

		loader, err := NewTLSConfigLoader()
		Expect(err).ToNot(HaveOccurred())

		err = loader.SetServerCACertificate(testEnv.caCertFile)
		Expect(err).ToNot(HaveOccurred())

		err = loader.SetClientCertificate(testEnv.certFile, testEnv.certKeyFile)
		Expect(err).ToNot(HaveOccurred())

		Expect(err).ToNot(HaveOccurred())
		defer loader.Stop()

		err = loader.Watch()
		Expect(err).ToNot(HaveOccurred())

		// set up the httptest.Server using our certificate signed by our CA
		server := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, "success!")
		}))
		server.TLS = &tls.Config{
			Certificates: []tls.Certificate{testEnv.cert},
		}
		server.StartTLS()
		defer server.Close()

		// communicate with the server using an http.Client configured to trust our CA
		transport := &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs: loader.GetRootCAs(),
			},
		}
		http := http.Client{
			Transport: transport,
		}
		resp, err := http.Get(server.URL)
		Expect(err).ToNot(HaveOccurred())

		// verify the response
		respBodyBytes, err := ioutil.ReadAll(resp.Body)
		Expect(err).ToNot(HaveOccurred())

		body := strings.TrimSpace(string(respBodyBytes[:]))
		Expect(body).To(Equal("success!"))
	})
})
