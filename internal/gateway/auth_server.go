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
	"strings"
	"time"

	"google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/protobuf/types/known/wrapperspb"

	envoy_core "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	envoy_auth "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	envoy_type "github.com/envoyproxy/go-control-plane/envoy/type/v3"

	"github.com/finleap-connect/monoskope/internal/gateway/auth"
	"github.com/finleap-connect/monoskope/pkg/api/gateway"
	m8roles "github.com/finleap-connect/monoskope/pkg/domain/constants/roles"
	m8scopes "github.com/finleap-connect/monoskope/pkg/domain/constants/scopes"
	"github.com/finleap-connect/monoskope/pkg/domain/repositories"
	"github.com/finleap-connect/monoskope/pkg/jwt"
	"github.com/finleap-connect/monoskope/pkg/logger"
	"google.golang.org/grpc/codes"
)

const (
	body_unauthorized    = "unauthorized"
	body_unauthenticated = "unauthenticated"
)

// authServer implements the AuthN/AuthZ decision API used as Ambassador Auth Service.
type authServer struct {
	envoy_auth.UnimplementedAuthorizationServer
	log        logger.Logger
	oidcServer *auth.Server
	userRepo   repositories.ReadOnlyUserRepository
	issuerURL  string
}

// NewAuthServer creates a new instance of gateway.authServer.
func NewAuthServer(issuerURL string, oidcServer *auth.Server, userRepo repositories.ReadOnlyUserRepository) *authServer {
	s := &authServer{
		log:        logger.WithName("auth-server"),
		oidcServer: oidcServer,
		userRepo:   userRepo,
		issuerURL:  issuerURL,
	}
	return s
}

// Check request object.
func (s *authServer) Check(ctx context.Context, req *envoy_auth.CheckRequest) (*envoy_auth.CheckResponse, error) {
	var err error
	var authToken *jwt.AuthToken
	path := req.GetAttributes().GetRequest().GetHttp().GetPath()
	authenticated := false
	authorized := false

	// Logging
	s.log.Info("Authenticating request...", "path", path)
	// Print headers
	s.log.V(logger.DebugLevel).Info("Request received.", "Value", req)

	// Authenticate user
	if !authenticated {
		authToken = s.tokenValidationFromContext(ctx, req) // via JWT
		authenticated = authToken != nil
	}
	if !authenticated {
		authToken = s.certValidation(ctx, req) // via client certificate validation
		authenticated = authToken != nil
	}
	if !authenticated {
		return s.createUnauthorizedResponse(body_unauthenticated), nil
	}

	// Get message body for policy evaluation in the future
	// body := req.Attributes.Request.Http.RawBody
	// s.log.V(logger.DebugLevel).Info("Message body received.", "body", body)

	// Authorize user
	authorized, err = s.validatePolicies(ctx, req)
	if err != nil {
		s.log.Error(err, "Error checking authorization of user.")
		return s.createUnauthorizedResponse(body_unauthorized), err
	}
	if authorized {
		return s.createAuthorizedResponse(authToken), nil
	}
	return s.createUnauthorizedResponse(body_unauthorized), nil
}

// validatePolicies validates the configured policies using OPA
func (s *authServer) validatePolicies(ctx context.Context, req *envoy_auth.CheckRequest) (bool, error) {
	// TODO: Implement
	return true, nil
}

func (s *authServer) retrieveUserId(ctx context.Context, email string) (string, bool) {
	user, err := s.userRepo.ByEmail(ctx, email)
	if err != nil {
		return "", false
	}
	return user.Id, true
}

// tokenValidationFromContext validates the token provided within the authorization flow from gin context
func (s *authServer) tokenValidationFromContext(ctx context.Context, req *envoy_auth.CheckRequest) *jwt.AuthToken {
	authToken := s.tokenValidation(ctx, defaultBearerTokenFromHeaders(req.Attributes.Request.Http.Headers))
	if authToken == nil {
		return nil
	}

	// Check user actually exists in m8
	user, err := s.userRepo.ByEmail(ctx, authToken.Email)
	if err != nil && !authToken.IsAPIToken {
		s.log.Info("Token validation failed. User does not exist.", "Email", authToken.Email)
		return nil
	}

	// Validate scopes
	route := req.GetAttributes().GetRequest().GetHttp().GetPath()
	scopes := strings.Split(authToken.Scope, " ")

	// Validation for API Token Endpoint
	// TODO: This is a temporary solution until authorization has been replaced with Open Policy Agent
	if strings.HasPrefix(route, "/"+gateway.APIToken_ServiceDesc.ServiceName) {
		if !authToken.IsAPIToken {
			for _, role := range user.Roles {
				if role.Role == m8roles.Admin.String() && role.Scope == m8scopes.System.String() { // Only system admins can issue API tokens
					return authToken
				}
			}
			s.log.Info("Token validation failed. Only system admins can call that route.", "Route", route, "Scopes", authToken.Scope)
			return nil
		} else { // API Tokens can't be used to issue new ones
			s.log.Info("Token validation failed. Token can not be used for route.", "Route", route, "Scopes", authToken.Scope)
			return nil
		}
	}

	// SCIM API Access
	if strings.HasPrefix(route, "/scim") {
		if containsString(scopes, gateway.AuthorizationScope_WRITE_SCIM.String()) {
			return authToken
		}
	}

	// General API access
	if containsString(scopes, gateway.AuthorizationScope_API.String()) {
		return authToken
	}

	s.log.Info("Token validation failed. Token has not correct scopes for route.", "Route", route, "Scopes", authToken.Scope)
	return nil
}

// tokenValidation validates the token provided within the authorization flow
func (s *authServer) tokenValidation(ctx context.Context, token string) *jwt.AuthToken {
	s.log.Info("Validating token...")

	if token == "" {
		s.log.Info("Token validation failed.", "error", "token is empty")
		return nil
	}

	authToken := &jwt.AuthToken{}
	if err := s.oidcServer.Authorize(ctx, token, authToken); err != nil {
		s.log.Info("Token validation failed.", "error", err.Error())
		return nil
	}
	if err := authToken.Validate(s.issuerURL); err != nil {
		s.log.Info("Token validation failed.", "error", err.Error())
		return nil
	}

	s.log.Info("Token validation successful", "subject", authToken.Subject, "email", authToken.Email, "scope", authToken.Scope)

	return authToken
}

// tokenValidation validates the client certificate provided within the forwarded client secret header
func (s *authServer) certValidation(ctx context.Context, req *envoy_auth.CheckRequest) *jwt.AuthToken {
	s.log.Info("Validating client certificate...")

	cert, err := clientCertificateFromHeaders(req.GetAttributes().GetRequest().GetHttp().GetHeaders())
	if err != nil {
		s.log.Info("Certificate validation failed.", "error", err.Error())
		return nil
	}

	if userId, ok := s.retrieveUserId(ctx, cert.EmailAddresses[0]); !ok {
		s.log.Info("Certificate validation failed. User does not exist.", "Email", cert.EmailAddresses[0])
		return nil
	} else {
		claims := auth.NewAuthToken(&jwt.StandardClaims{
			Name:  cert.Subject.CommonName,
			Email: cert.EmailAddresses[0],
		}, s.issuerURL, userId, time.Minute*5)
		claims.Subject = userId
		claims.Issuer = cert.Issuer.CommonName
		s.log.Info("Client certificate validation successful.", "User", claims.Email)
		return claims
	}
}

func (s *authServer) createAuthorizedResponse(authToken *jwt.AuthToken) *envoy_auth.CheckResponse {
	// Set headers with auth info
	return &envoy_auth.CheckResponse{
		Status: &status.Status{Code: int32(codes.OK)},
		HttpResponse: &envoy_auth.CheckResponse_OkResponse{
			OkResponse: &envoy_auth.OkHttpResponse{
				Headers: []*envoy_core.HeaderValueOption{
					{Header: &envoy_core.HeaderValue{Key: auth.HeaderAuthId, Value: authToken.Subject}, Append: wrapperspb.Bool(true)},
					{Header: &envoy_core.HeaderValue{Key: auth.HeaderAuthName, Value: authToken.Name}, Append: wrapperspb.Bool(true)},
					{Header: &envoy_core.HeaderValue{Key: auth.HeaderAuthEmail, Value: authToken.Email}, Append: wrapperspb.Bool(true)},
					{Header: &envoy_core.HeaderValue{Key: auth.HeaderAuthNotBefore, Value: authToken.NotBefore.Time().Format(auth.HeaderAuthNotBeforeFormat)}, Append: wrapperspb.Bool(true)},
				},
			},
		},
	}
}

func (s *authServer) createUnauthorizedResponse(body string) *envoy_auth.CheckResponse {
	return &envoy_auth.CheckResponse{
		Status: &status.Status{Code: int32(codes.Unauthenticated)},
		HttpResponse: &envoy_auth.CheckResponse_DeniedResponse{
			DeniedResponse: &envoy_auth.DeniedHttpResponse{
				Status: &envoy_type.HttpStatus{
					Code: envoy_type.StatusCode_Unauthorized,
				},
				Body: body,
			},
		},
	}
}
