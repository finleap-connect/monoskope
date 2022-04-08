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
	"strings"
	"time"

	"github.com/finleap-connect/monoskope/internal/gateway/auth"
	"github.com/finleap-connect/monoskope/pkg/api/gateway"
	m8roles "github.com/finleap-connect/monoskope/pkg/domain/constants/roles"
	m8scopes "github.com/finleap-connect/monoskope/pkg/domain/constants/scopes"
	"github.com/finleap-connect/monoskope/pkg/domain/repositories"
	grpcUtil "github.com/finleap-connect/monoskope/pkg/grpc"
	"github.com/finleap-connect/monoskope/pkg/jwt"
	"github.com/finleap-connect/monoskope/pkg/logger"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"github.com/open-policy-agent/opa/rego"
	"golang.org/x/exp/slices"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

type policyRoles struct {
	Name     string
	Scope    string
	Resource string
}

type policyUser struct {
	Id    string
	Name  string
	Roles []policyRoles
}

type policyInput struct {
	User        policyUser
	Path        string
	CommandType string
}

// authServer implements the AuthN/AuthZ decision API used as Ambassador Auth Service.
type authServer struct {
	gateway.UnimplementedGatewayAuthZServer
	log           logger.Logger
	oidcServer    *auth.Server
	userRepo      repositories.ReadOnlyUserRepository
	issuerURL     string
	preparedQuery *rego.PreparedEvalQuery
}

func NewAuthServerClient(ctx context.Context, gatewayAddr string) (*grpc.ClientConn, gateway.GatewayAuthZClient, error) {
	conn, err := grpcUtil.
		NewGrpcConnectionFactoryWithDefaults(gatewayAddr).
		ConnectWithTimeout(ctx, 10*time.Second)
	if err != nil {
		return nil, nil, err
	}

	return conn, gateway.NewGatewayAuthZClient(conn), nil
}

// NewAuthServer creates a new instance of gateway.authServer.
func NewAuthServer(ctx context.Context, issuerURL string, oidcServer *auth.Server, userRepo repositories.ReadOnlyUserRepository, policiesPath string) (*authServer, error) {
	s := &authServer{
		log:        logger.WithName("auth-server"),
		oidcServer: oidcServer,
		userRepo:   userRepo,
		issuerURL:  issuerURL,
	}

	query, err := rego.New(
		rego.Query("data.m8.authz.authorized"),
		rego.Load([]string{policiesPath}, nil),
	).PrepareForEval(ctx)
	if err != nil {
		return nil, err
	}
	s.preparedQuery = &query

	return s, nil
}

// Check request object.
func (s *authServer) Check(ctx context.Context, req *gateway.CheckRequest) (*gateway.CheckResponse, error) {
	var err error
	var authToken *jwt.AuthToken
	path := req.FullMethodName
	authenticated := false
	authorized := false

	s.log.Info("Authenticating request...", "path", path, "req", req)

	// check authentication
	// via JWT
	authToken, err = s.tokenValidationFromContext(ctx, req)
	if err != nil {
		return nil, err
	}
	authenticated = err == nil

	// via client certificate validation if not authenticated already
	if !authenticated {
		authToken, err = s.certValidation(ctx, req)
		if err != nil {
			return nil, err
		}
		authenticated = err == nil
	}

	if !authenticated {
		// authentication failed
		return nil, status.Error(codes.Unauthenticated, "authentication failed")
	}

	s.log.Info("Request authenticated. Checking authorization...", "path", path, "req", req)

	// check authorization
	authorized, err = s.validatePolicies(ctx, req, authToken)
	if err != nil {
		// authorization failed with error
		s.log.Error(err, "Error checking authorization of user.")
		return nil, status.Error(codes.PermissionDenied, "authorization failed")
	}

	if authorized {
		// authorization successful
		return s.createAuthorizedResponse(authToken), nil
	}
	return nil, status.Error(codes.PermissionDenied, "authorization failed")
}

// validatePolicies validates the configured policies using OPA
func (s *authServer) validatePolicies(ctx context.Context, req *gateway.CheckRequest, authToken *jwt.AuthToken) (bool, error) {
	user, err := s.userRepo.ByEmail(ctx, authToken.Email)
	if err != nil {
		s.log.Error(err, "Policy evaluation failed. User does not exist.", "email", authToken.Email)
		return false, err
	}

	input := policyInput{
		User: policyUser{
			Id:   user.Id,
			Name: user.Name,
		},
		Path: req.FullMethodName,
	}

	input.User.Roles = make([]policyRoles, 0)
	for _, role := range user.Roles {
		input.User.Roles = append(input.User.Roles, policyRoles{
			Name:     role.Role,
			Scope:    role.Scope,
			Resource: role.Resource,
		})
	}

	results, err := s.preparedQuery.Eval(ctx, rego.EvalInput(input))
	if err != nil {
		s.log.Error(err, "Policy evaluation failed.", "email", authToken.Email)
		return false, err
	}
	if !results.Allowed() {
		s.log.Info("Policy evaluation failed.", "email", authToken.Email, "results", results)
		return false, nil
	}
	s.log.Info("Policy evaluation succeeded.", "email", authToken.Email, "results", results)
	return results.Allowed(), nil
}

func (s *authServer) retrieveUserId(ctx context.Context, email string) (string, bool) {
	user, err := s.userRepo.ByEmail(ctx, email)
	if err != nil {
		return "", false
	}
	return user.Id, true
}

// tokenValidationFromContext validates the token provided within the authorization flow from gin context
func (s *authServer) tokenValidationFromContext(ctx context.Context, req *gateway.CheckRequest) (*jwt.AuthToken, error) {
	token, err := grpc_auth.AuthFromMD(ctx, "bearer")
	if err != nil {
		return nil, err
	}

	authToken, err := s.tokenValidation(ctx, token)
	if err != nil {
		return nil, err
	}

	// Check user actually exists in m8
	user, err := s.userRepo.ByEmail(ctx, authToken.Email)
	if err != nil && !authToken.IsAPIToken {
		s.log.Info("Token validation failed. User does not exist.", "Email", authToken.Email)
		return nil, err
	}

	// Validate scopes
	route := req.FullMethodName
	scopes := strings.Split(authToken.Scope, " ")

	// Validation for API Token Endpoint
	// TODO: This is a temporary solution until authorization has been replaced with Open Policy Agent
	if strings.HasPrefix(route, "/"+gateway.APIToken_ServiceDesc.ServiceName) {
		if !authToken.IsAPIToken {
			for _, role := range user.Roles {
				if role.Role == m8roles.Admin.String() && role.Scope == m8scopes.System.String() { // Only system admins can issue API tokens
					return authToken, nil
				}
			}
			s.log.Info("Token validation failed. Only system admins can call that route.", "Route", route, "Scopes", authToken.Scope)
			return nil, status.Error(codes.Unauthenticated, "token validation failed")
		} else { // API Tokens can't be used to issue new ones
			s.log.Info("Token validation failed. Token can not be used for route.", "Route", route, "Scopes", authToken.Scope)
			return nil, err
		}
	}

	// SCIM API Access
	if strings.HasPrefix(route, "/scim") {
		if slices.Contains(scopes, gateway.AuthorizationScope_WRITE_SCIM.String()) {
			return authToken, err
		}
	}

	// General API access
	if slices.Contains(scopes, gateway.AuthorizationScope_API.String()) {
		return authToken, err
	}

	s.log.Info("Token validation failed. Token has not correct scopes for route.", "Route", route, "Scopes", authToken.Scope)
	return nil, err
}

// tokenValidation validates the token provided within the authorization flow
func (s *authServer) tokenValidation(ctx context.Context, token string) (*jwt.AuthToken, error) {
	s.log.Info("Validating token...")

	if token == "" {
		s.log.Info("Token validation failed.", "error", "token is empty")
		return nil, nil
	}

	authToken := &jwt.AuthToken{}
	if err := s.oidcServer.Authorize(ctx, token, authToken); err != nil {
		s.log.Info("Token validation failed.", "error", err.Error())
		return nil, err
	}
	if err := authToken.Validate(s.issuerURL); err != nil {
		s.log.Info("Token validation failed.", "error", err.Error())
		return nil, err
	}

	s.log.Info("Token validation successful", "subject", authToken.Subject, "email", authToken.Email, "scope", authToken.Scope)

	return authToken, nil
}

// tokenValidation validates the client certificate provided within the forwarded client secret header
func (s *authServer) certValidation(ctx context.Context, req *gateway.CheckRequest) (*jwt.AuthToken, error) {
	s.log.Info("Validating client certificate...")

	p, ok := peer.FromContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "no peer found")
	}

	tlsAuth, ok := p.AuthInfo.(credentials.TLSInfo)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "unexpected peer transport credentials")
	}

	if len(tlsAuth.State.VerifiedChains) == 0 || len(tlsAuth.State.VerifiedChains[0]) == 0 {
		return nil, status.Error(codes.Unauthenticated, "could not verify peer certificate")
	}

	userName := tlsAuth.State.VerifiedChains[0][0].Subject.CommonName
	emailAddress := tlsAuth.State.VerifiedChains[0][0].EmailAddresses[0]
	if userId, ok := s.retrieveUserId(ctx, emailAddress); !ok {
		return nil, status.Error(codes.Unauthenticated, "invalid subject common name")
	} else {
		claims := auth.NewAuthToken(&jwt.StandardClaims{
			Name:  userName,
			Email: emailAddress,
		}, s.issuerURL, userId, time.Minute*5)
		claims.Subject = userId
		claims.Issuer = tlsAuth.State.VerifiedChains[0][0].Issuer.CommonName
		s.log.Info("Client certificate validation successful.", "User", claims.Email)
		return claims, nil
	}
}

func (s *authServer) createAuthorizedResponse(authToken *jwt.AuthToken) *gateway.CheckResponse {
	// Set headers with auth info
	return &gateway.CheckResponse{
		Tags: []*gateway.CheckResponse_CheckResponseTag{
			{Key: auth.HeaderAuthId, Value: authToken.Subject},
			{Key: auth.HeaderAuthName, Value: authToken.Name},
			{Key: auth.HeaderAuthEmail, Value: authToken.Email},
			{Key: auth.HeaderAuthNotBefore, Value: authToken.NotBefore.Time().Format(auth.HeaderAuthNotBeforeFormat)},
		},
	}
}
