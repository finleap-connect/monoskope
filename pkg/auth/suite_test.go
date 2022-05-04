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

package auth

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"testing"

	"github.com/finleap-connect/monoskope/pkg/logger"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const (
	RedirectURLHostname = "localhost"
	RedirectURLPort     = ":8000"
	AuthURLPort         = ":8050"
)

var (
	httpClient *http.Client
	httpServer *http.Server
	log        logger.Logger
	ctx        = context.Background()
)

func TestAuth(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "monoctl/auth")
}

var _ = BeforeSuite(func() {
	done := make(chan interface{})

	go func() {
		log = logger.WithName("TestAuth")

		By("bootstrapping test env")

		// Setup HTTP client
		httpClient = &http.Client{}

		mux := http.NewServeMux()
		mux.HandleFunc("/auth", auth)
		httpServer = &http.Server{
			Addr:    AuthURLPort,
			Handler: mux,
		}
		go func() {
			_ = httpServer.ListenAndServe()
		}()
		close(done)
	}()

	Eventually(done, 60).Should(BeClosed())
})

func getAuthCodeUrl(redirectURI, state string) string {
	return fmt.Sprintf("http://localhost%s/auth?callback=%s&state=%s", AuthURLPort, url.QueryEscape(redirectURI), state)
}

func auth(rw http.ResponseWriter, r *http.Request) {
	log.Info("received auth request")
	err := r.ParseForm()
	if err != nil {
		return
	}
	if errMsg := r.Form.Get("error"); errMsg != "" {
		log.Error(err, errMsg)
		return
	}
	callBackUrl := r.Form.Get("callback")
	state := r.Form.Get("state")
	http.Redirect(rw, r, fmt.Sprintf("%s?state=%s&code=my-fancy-auth-code", callBackUrl, state), http.StatusSeeOther)
}

var _ = AfterSuite(func() {
	By("tearing down the test environment")
	defer httpServer.Close()
})
