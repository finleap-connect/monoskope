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
	"time"

	envoy_auth "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	envoy_type "github.com/envoyproxy/go-control-plane/envoy/type/v3"
	"github.com/finleap-connect/monoskope/internal/gateway/auth"
	"github.com/finleap-connect/monoskope/pkg/jwt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"google.golang.org/grpc/codes"
	josejwt "gopkg.in/square/go-jose.v2/jwt"
)

var _ = Describe("Gateway Auth Server", func() {
	var (
		ctx = context.Background()
	)

	It("can authenticate with JWT", func() {
		expectedValidity := time.Hour * 1
		token := auth.NewAuthToken(&jwt.StandardClaims{Name: env.AdminUser.Name, Email: env.AdminUser.Email}, localAddrAPIServer, env.AdminUser.Id, expectedValidity)
		signer := env.JwtTestEnv.CreateSigner()
		signedToken, err := signer.GenerateSignedToken(token)
		Expect(err).NotTo(HaveOccurred())

		conn, err := CreateInsecureConnection(ctx, env.ApiListenerAPIServer.Addr().String())
		Expect(err).ToNot(HaveOccurred())
		defer conn.Close()
		authClient := envoy_auth.NewAuthorizationClient(conn)

		resp, err := authClient.Check(ctx, &envoy_auth.CheckRequest{Attributes: &envoy_auth.AttributeContext{
			Request: &envoy_auth.AttributeContext_Request{
				Http: &envoy_auth.AttributeContext_HttpRequest{
					Path: "/domain.Cluster/",
					Headers: map[string]string{
						auth.HeaderAuthorization: fmt.Sprintf("Bearer %s", signedToken),
					},
				},
			},
		}})
		Expect(err).NotTo(HaveOccurred())
		Expect(resp.Status.Code).To(Equal(int32(codes.OK)))
	})
	It("fails authentication with invalid JWT", func() {
		conn, err := CreateInsecureConnection(ctx, env.ApiListenerAPIServer.Addr().String())
		Expect(err).ToNot(HaveOccurred())
		defer conn.Close()
		authClient := envoy_auth.NewAuthorizationClient(conn)

		resp, err := authClient.Check(ctx, &envoy_auth.CheckRequest{Attributes: &envoy_auth.AttributeContext{
			Request: &envoy_auth.AttributeContext_Request{
				Http: &envoy_auth.AttributeContext_HttpRequest{
					Headers: map[string]string{
						auth.HeaderAuthorization: fmt.Sprintf("Bearer %s", "notavalidjwt"),
					},
				},
			},
		}})
		Expect(err).NotTo(HaveOccurred())
		Expect(resp.Status.Code).To(Equal(int32(codes.Unauthenticated)))
		Expect(resp.GetDeniedResponse().Status.Code).To(Equal(envoy_type.StatusCode_Unauthorized))
	})
	It("fails authentication with expired JWT", func() {
		expectedValidity := -30 * time.Minute
		token := auth.NewAuthToken(&jwt.StandardClaims{Name: env.ExistingUser.Name, Email: env.ExistingUser.Email}, localAddrAPIServer, env.ExistingUser.Id, expectedValidity)
		token.NotBefore = josejwt.NewNumericDate(time.Now().UTC().Add(-1 * time.Hour))

		signer := env.JwtTestEnv.CreateSigner()
		signedToken, err := signer.GenerateSignedToken(token)
		Expect(err).NotTo(HaveOccurred())

		conn, err := CreateInsecureConnection(ctx, env.ApiListenerAPIServer.Addr().String())
		Expect(err).ToNot(HaveOccurred())
		defer conn.Close()
		authClient := envoy_auth.NewAuthorizationClient(conn)

		resp, err := authClient.Check(ctx, &envoy_auth.CheckRequest{Attributes: &envoy_auth.AttributeContext{
			Request: &envoy_auth.AttributeContext_Request{
				Http: &envoy_auth.AttributeContext_HttpRequest{
					Headers: map[string]string{
						auth.HeaderAuthorization: fmt.Sprintf("Bearer %s", signedToken),
					},
				},
			},
		}})
		Expect(err).NotTo(HaveOccurred())
		Expect(resp.Status.Code).To(Equal(int32(codes.Unauthenticated)))
		Expect(resp.GetDeniedResponse().Status.Code).To(Equal(envoy_type.StatusCode_Unauthorized))
	})
	It("fails authentication with not existing user", func() {
		expectedValidity := time.Hour * 1
		token := auth.NewAuthToken(&jwt.StandardClaims{Name: env.NotExistingUser.Name, Email: env.NotExistingUser.Email}, localAddrAPIServer, env.NotExistingUser.Id, expectedValidity)
		signer := env.JwtTestEnv.CreateSigner()
		signedToken, err := signer.GenerateSignedToken(token)
		Expect(err).NotTo(HaveOccurred())

		conn, err := CreateInsecureConnection(ctx, env.ApiListenerAPIServer.Addr().String())
		Expect(err).ToNot(HaveOccurred())
		defer conn.Close()
		authClient := envoy_auth.NewAuthorizationClient(conn)

		resp, err := authClient.Check(ctx, &envoy_auth.CheckRequest{Attributes: &envoy_auth.AttributeContext{
			Request: &envoy_auth.AttributeContext_Request{
				Http: &envoy_auth.AttributeContext_HttpRequest{
					Headers: map[string]string{
						auth.HeaderAuthorization: fmt.Sprintf("Bearer %s", signedToken),
					},
				},
			},
		}})
		Expect(err).NotTo(HaveOccurred())
		Expect(resp.Status.Code).To(Equal(int32(codes.Unauthenticated)))
		Expect(resp.GetDeniedResponse().Status.Code).To(Equal(envoy_type.StatusCode_Unauthorized))
	})
	// It("can not authenticate with JWT for wrong scope", func() {
	// 	token := auth.NewClusterBootstrapToken(&jwt.StandardClaims{Name: env.ExistingUser.Name, Email: env.ExistingUser.Email}, localAddrAPIServer, env.ExistingUser.Id)
	// 	signer := env.JwtTestEnv.CreateSigner()
	// 	signedToken, err := signer.GenerateSignedToken(token)
	// 	Expect(err).NotTo(HaveOccurred())

	// 	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("http://%s/auth/test", localAddrOIDCProviderServer), nil)
	// 	Expect(err).NotTo(HaveOccurred())

	// 	req.Header.Set(auth.HeaderAuthorization, fmt.Sprintf("bearer %s", signedToken))
	// 	res, err := env.HttpClient.Do(req)
	// 	Expect(err).ToNot(HaveOccurred())
	// 	Expect(res.StatusCode).To(Equal(http.StatusUnauthorized))
	// })
	// It("can authenticate with JWT for correct scope", func() {
	// 	expectedValidity := time.Hour * 1
	// 	token := auth.NewApiToken(&jwt.StandardClaims{Name: env.NotExistingUser.Name}, localAddrAPIServer, env.NotExistingUser.Id, expectedValidity, []gateway.AuthorizationScope{
	// 		gateway.AuthorizationScope_WRITE_SCIM,
	// 	})
	// 	signer := env.JwtTestEnv.CreateSigner()
	// 	signedToken, err := signer.GenerateSignedToken(token)
	// 	Expect(err).NotTo(HaveOccurred())

	// 	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("http://%s/auth/scim/Users", localAddrOIDCProviderServer), nil)
	// 	Expect(err).NotTo(HaveOccurred())

	// 	req.Header.Set(auth.HeaderAuthorization, fmt.Sprintf("bearer %s", signedToken))
	// 	res, err := env.HttpClient.Do(req)
	// 	Expect(err).ToNot(HaveOccurred())
	// 	Expect(res.StatusCode).To(Equal(http.StatusOK))
	// })
	// It("can authenticate with JWT for correct scope", func() {
	// 	expectedValidity := time.Hour * 1
	// 	token := auth.NewAuthToken(&jwt.StandardClaims{Name: env.AdminUser.Name, Email: env.AdminUser.Email}, localAddrAPIServer, env.ExistingUser.Id, expectedValidity)
	// 	signer := env.JwtTestEnv.CreateSigner()
	// 	signedToken, err := signer.GenerateSignedToken(token)
	// 	Expect(err).NotTo(HaveOccurred())

	// 	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("http://%s/auth/gateway.APIToken/", localAddrOIDCProviderServer), nil)
	// 	Expect(err).NotTo(HaveOccurred())

	// 	req.Header.Set(auth.HeaderAuthorization, fmt.Sprintf("bearer %s", signedToken))
	// 	res, err := env.HttpClient.Do(req)
	// 	Expect(err).ToNot(HaveOccurred())
	// 	Expect(res.StatusCode).To(Equal(http.StatusOK))
	// })
})
