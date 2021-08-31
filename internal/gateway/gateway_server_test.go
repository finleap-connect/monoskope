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
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/PuerkitoBio/goquery"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	api "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/gateway"
	"golang.org/x/sync/errgroup"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

var _ = Describe("Gateway", func() {
	var (
		ctx = context.Background()
	)
	It("can retrieve auth url", func() {
		conn, err := CreateInsecureConnection(ctx, env.ApiListenerAPIServer.Addr().String())
		Expect(err).ToNot(HaveOccurred())
		defer conn.Close()
		gwc := api.NewGatewayClient(conn)

		authInfo, err := gwc.GetAuthInformation(context.Background(), &api.AuthState{CallbackUrl: "http://localhost:8000"})
		Expect(err).ToNot(HaveOccurred())
		Expect(authInfo).ToNot(BeNil())
		env.Log.Info("AuthCodeURL: " + authInfo.AuthCodeUrl)
	})
	It("can go through oidc-flow with existing user", func() {
		conn, err := CreateInsecureConnection(ctx, env.ApiListenerAPIServer.Addr().String())
		Expect(err).ToNot(HaveOccurred())
		defer conn.Close()
		gwcAuth := api.NewGatewayClient(conn)

		ready := make(chan string, 1)
		oidcClientServer, err := env.NewOidcClientServer(ready)
		Expect(err).ToNot(HaveOccurred())
		defer oidcClientServer.Close()

		env.Log.Info("oidc redirect uri: " + oidcClientServer.RedirectURI)
		authInfo, err := gwcAuth.GetAuthInformation(context.Background(), &api.AuthState{CallbackUrl: oidcClientServer.RedirectURI})
		Expect(err).ToNot(HaveOccurred())
		Expect(authInfo).ToNot(BeNil())

		var innerErr error
		res, err := env.HttpClient.Get(authInfo.AuthCodeUrl)
		Expect(err).NotTo(HaveOccurred())
		doc, err := goquery.NewDocumentFromReader(res.Body)
		Expect(err).NotTo(HaveOccurred())
		path, ok := doc.Find("form").Attr("action")
		Expect(ok).To(BeTrue())
		formAction := fmt.Sprintf("%s%s", env.IdentityProviderURL, path)

		var authCode string
		var statusCode int
		var eg errgroup.Group
		eg.Go(func() error {
			defer GinkgoRecover()
			var innerErr error
			authCode, innerErr = oidcClientServer.ReceiveCodeViaLocalServer(ctx, authInfo.AuthCodeUrl, authInfo.State)
			return innerErr
		})
		eg.Go(func() error {
			defer GinkgoRecover()
			env.Log.Info("wait for oidc client server to get ready...")
			<-ready
			res, err = env.HttpClient.PostForm(formAction, url.Values{
				"login": {"admin@monoskope.io"}, "password": {"password"},
			})
			if err == nil {
				statusCode = res.StatusCode
			}
			return innerErr
		})
		Expect(eg.Wait()).NotTo(HaveOccurred())
		Expect(statusCode).To(Equal(http.StatusOK))

		authResponse, err := gwcAuth.ExchangeAuthCode(context.Background(), &api.AuthCode{Code: authCode, State: authInfo.GetState(), CallbackUrl: oidcClientServer.RedirectURI})
		Expect(err).ToNot(HaveOccurred())
		Expect(authResponse).ToNot(BeNil())
		Expect(authResponse.GetAccessToken()).ToNot(Equal(""))
		Expect(authResponse.GetUsername()).ToNot(Equal(""))
		env.Log.Info("Received user info", "AccessToken", authResponse.GetAccessToken(), "Expiry", authResponse.GetExpiry().AsTime())
	})
})

var _ = Describe("HealthCheck", func() {
	var (
		ctx = context.Background()
	)
	It("can do health checks", func() {
		conn, err := CreateInsecureConnection(ctx, env.ApiListenerAPIServer.Addr().String())
		Expect(err).ToNot(HaveOccurred())
		defer conn.Close()

		hc := healthpb.NewHealthClient(conn)
		res, err := hc.Check(ctx, &healthpb.HealthCheckRequest{})
		Expect(err).ToNot(HaveOccurred())
		Expect(res.GetStatus()).To(Equal(healthpb.HealthCheckResponse_SERVING))
	})
})
