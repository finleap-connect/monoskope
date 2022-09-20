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
	"context"
	"time"

	"github.com/finleap-connect/monoskope/internal/gateway/auth"
	cmdData "github.com/finleap-connect/monoskope/pkg/api/domain/commanddata"
	"github.com/finleap-connect/monoskope/pkg/api/gateway"
	cmd "github.com/finleap-connect/monoskope/pkg/domain/commands"
	commandTypes "github.com/finleap-connect/monoskope/pkg/domain/constants/commands"
	"github.com/finleap-connect/monoskope/pkg/domain/constants/roles"
	"github.com/finleap-connect/monoskope/pkg/domain/constants/scopes"
	"github.com/finleap-connect/monoskope/pkg/domain/mock"
	"github.com/finleap-connect/monoskope/pkg/domain/projections"
	"github.com/finleap-connect/monoskope/pkg/jwt"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
	jose_jwt "gopkg.in/square/go-jose.v2/jwt"
)

var _ = Describe("Gateway Auth Server", func() {
	var (
		ctx            = context.Background()
		expectedUserId = uuid.New()
		expectedRole   = roles.User
		expectedScope  = scopes.Tenant
	)

	getTokenForUser := func(user *projections.User) string {
		expectedValidity := time.Hour * 1
		token := auth.NewAuthToken(&jwt.StandardClaims{Name: user.Name, Email: user.Email}, localAddrAPIServer, user.Id, expectedValidity)
		signer := testEnv.JwtTestEnv.CreateSigner()
		signedToken, err := signer.GenerateSignedToken(token)
		Expect(err).NotTo(HaveOccurred())
		return signedToken
	}

	getCreateUserRoleBindingCmd := func() *anypb.Any {
		command := cmd.NewCommandWithData(uuid.Nil, commandTypes.CreateUserRoleBinding,
			&cmdData.CreateUserRoleBindingCommandData{
				UserId:   expectedUserId.String(),
				Role:     string(expectedRole),
				Scope:    string(expectedScope),
				Resource: wrapperspb.String(mock.TestTenant.Id),
			},
		)

		a := &anypb.Any{}
		err := a.MarshalFrom(command)
		Expect(err).ToNot(HaveOccurred())
		return a
	}
	It("admin can auth with JWT", func() {
		conn, err := CreateInsecureConnection(ctx, testEnv.ApiListenerAPIServer.Addr().String())
		Expect(err).ToNot(HaveOccurred())
		defer conn.Close()
		authClient := gateway.NewGatewayAuthClient(conn)

		bytes, err := protojson.Marshal(getCreateUserRoleBindingCmd())
		Expect(err).ToNot(HaveOccurred())

		resp, err := authClient.Check(ctx, &gateway.CheckRequest{
			FullMethodName: "/eventsourcing.CommandHandler/",
			AccessToken:    getTokenForUser(mock.TestAdminUser),
			Request:        bytes,
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(resp).ToNot(BeNil())
		Expect(resp.Tags).ToNot(BeNil())
	})

	It("can authorize with JWT as tenant admin", func() {
		conn, err := CreateInsecureConnection(ctx, testEnv.ApiListenerAPIServer.Addr().String())
		Expect(err).ToNot(HaveOccurred())
		defer conn.Close()
		authClient := gateway.NewGatewayAuthClient(conn)

		bytes, err := protojson.Marshal(getCreateUserRoleBindingCmd())
		Expect(err).ToNot(HaveOccurred())

		resp, err := authClient.Check(ctx, &gateway.CheckRequest{
			FullMethodName: "/eventsourcing.CommandHandler/Execute",
			AccessToken:    getTokenForUser(mock.TestTenantAdminUser),
			Request:        bytes,
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(resp).ToNot(BeNil())
		Expect(resp.Tags).ToNot(BeNil())
	})

	It("regular user can't authenticate with JWT for commandhandler", func() {
		conn, err := CreateInsecureConnection(ctx, testEnv.ApiListenerAPIServer.Addr().String())
		Expect(err).ToNot(HaveOccurred())
		defer conn.Close()
		authClient := gateway.NewGatewayAuthClient(conn)

		resp, err := authClient.Check(ctx, &gateway.CheckRequest{
			FullMethodName: "/eventsourcing.CommandHandler/Execute",
			AccessToken:    getTokenForUser(mock.TestExistingUser),
		})
		Expect(err).To(HaveOccurred())
		Expect(resp).To(BeNil())
		status, ok := status.FromError(err)
		Expect(ok).To(BeTrue())
		Expect(status).NotTo(BeNil())
		Expect(status.Code()).To(Equal(codes.PermissionDenied))
	})

	It("fails authentication with invalid JWT", func() {
		conn, err := CreateInsecureConnection(ctx, testEnv.ApiListenerAPIServer.Addr().String())
		Expect(err).ToNot(HaveOccurred())
		defer conn.Close()
		authClient := gateway.NewGatewayAuthClient(conn)

		resp, err := authClient.Check(ctx, &gateway.CheckRequest{
			FullMethodName: "/gateway.Gateway/",
			AccessToken:    "invalidjwt",
		})
		Expect(err).To(HaveOccurred())
		Expect(resp).To(BeNil())
		status, ok := status.FromError(err)
		Expect(ok).To(BeTrue())
		Expect(status).NotTo(BeNil())
		Expect(status.Code()).To(Equal(codes.Unauthenticated))
	})

	It("fails authentication with expired JWT", func() {
		expectedValidity := -30 * time.Minute
		token := auth.NewAuthToken(&jwt.StandardClaims{Name: mock.TestExistingUser.Name, Email: mock.TestExistingUser.Email}, localAddrAPIServer, mock.TestExistingUser.Id, expectedValidity)
		token.NotBefore = jose_jwt.NewNumericDate(time.Now().UTC().Add(-1 * time.Hour))

		signer := testEnv.JwtTestEnv.CreateSigner()
		signedToken, err := signer.GenerateSignedToken(token)
		Expect(err).NotTo(HaveOccurred())

		conn, err := CreateInsecureConnection(ctx, testEnv.ApiListenerAPIServer.Addr().String())
		Expect(err).ToNot(HaveOccurred())
		defer conn.Close()
		authClient := gateway.NewGatewayAuthClient(conn)

		resp, err := authClient.Check(ctx, &gateway.CheckRequest{
			FullMethodName: "/gateway.Gateway/",
			AccessToken:    signedToken,
		})
		Expect(err).To(HaveOccurred())
		Expect(resp).To(BeNil())
		status, ok := status.FromError(err)
		Expect(ok).To(BeTrue())
		Expect(status).NotTo(BeNil())
		Expect(status.Code()).To(Equal(codes.Unauthenticated))
	})
	It("can not authenticate with JWT for wrong scope", func() {
		expectedValidity := time.Hour * 1
		token := auth.NewAuthToken(&jwt.StandardClaims{Name: mock.TestExistingUser.Name, Email: mock.TestExistingUser.Email}, localAddrAPIServer, mock.TestExistingUser.Id, expectedValidity)
		signer := testEnv.JwtTestEnv.CreateSigner()
		signedToken, err := signer.GenerateSignedToken(token)
		Expect(err).NotTo(HaveOccurred())

		conn, err := CreateInsecureConnection(ctx, testEnv.ApiListenerAPIServer.Addr().String())
		Expect(err).ToNot(HaveOccurred())
		defer conn.Close()
		authClient := gateway.NewGatewayAuthClient(conn)

		resp, err := authClient.Check(ctx, &gateway.CheckRequest{
			FullMethodName: "/scim/Users",
			AccessToken:    signedToken,
		})
		Expect(err).To(HaveOccurred())
		Expect(resp).To(BeNil())
		status, ok := status.FromError(err)
		Expect(ok).To(BeTrue())
		Expect(status).NotTo(BeNil())
		Expect(status.Code()).To(Equal(codes.PermissionDenied))
	})
	It("can authenticate with JWT for correct scope", func() {
		expectedValidity := time.Hour * 1
		token := auth.NewApiToken(&jwt.StandardClaims{Name: mock.TestNoneExistingUser.Name}, localAddrAPIServer, mock.TestNoneExistingUser.Id, expectedValidity, []gateway.AuthorizationScope{
			gateway.AuthorizationScope_WRITE_SCIM,
		})
		signer := testEnv.JwtTestEnv.CreateSigner()
		signedToken, err := signer.GenerateSignedToken(token)
		Expect(err).NotTo(HaveOccurred())

		conn, err := CreateInsecureConnection(ctx, testEnv.ApiListenerAPIServer.Addr().String())
		Expect(err).ToNot(HaveOccurred())
		defer conn.Close()
		authClient := gateway.NewGatewayAuthClient(conn)

		resp, err := authClient.Check(ctx, &gateway.CheckRequest{
			FullMethodName: "/scim/Users",
			AccessToken:    signedToken,
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(resp).ToNot(BeNil())
		Expect(resp.Tags).ToNot(BeNil())
	})

})
