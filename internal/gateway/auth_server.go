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
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/finleap-connect/monoskope/internal/gateway/auth"
	"github.com/finleap-connect/monoskope/pkg/api/gateway"
	m8roles "github.com/finleap-connect/monoskope/pkg/domain/constants/roles"
	m8scopes "github.com/finleap-connect/monoskope/pkg/domain/constants/scopes"
	"github.com/finleap-connect/monoskope/pkg/domain/repositories"
	"github.com/finleap-connect/monoskope/pkg/jwt"
	"github.com/finleap-connect/monoskope/pkg/logger"
	"github.com/finleap-connect/monoskope/pkg/util"
)

const (
	HeaderAuthorization = "Authorization"
)

type OpenIdConfiguration struct {
	Issuer  string `json:"issuer"`
	JwksURL string `json:"jwks_uri"`
}

// authServer implements the AuthN/AuthZ decision API used as Ambassador Auth Service.
type authServer struct {
	api        *http.Server
	engine     *gin.Engine
	log        logger.Logger
	shutdown   *util.ShutdownWaitGroup
	oidcClient *auth.Client
	oidcServer *auth.Server
	userRepo   repositories.ReadOnlyUserRepository
	issuerURL  string
}

// NewAuthServer creates a new instance of gateway.authServer.
func NewAuthServer(issuerURL string, oidcClient *auth.Client, oidcServer *auth.Server, userRepo repositories.ReadOnlyUserRepository) *authServer {
	engine := gin.Default()
	s := &authServer{
		api:        &http.Server{Handler: engine},
		engine:     engine,
		log:        logger.WithName("auth-server"),
		shutdown:   util.NewShutdownWaitGroup(),
		oidcClient: oidcClient,
		oidcServer: oidcServer,
		userRepo:   userRepo,
		issuerURL:  issuerURL,
	}
	engine.Use(cors.Default())
	return s
}

// Serve tells the server to start listening on the specified address.
func (s *authServer) Serve(apiAddr string) error {
	// Setup grpc listener
	apiLis, err := net.Listen("tcp", apiAddr)
	if err != nil {
		return err
	}
	defer apiLis.Close()

	return s.ServeFromListener(apiLis)
}

// Serve tells the server to start listening on given listener.
func (s *authServer) ServeFromListener(apiLis net.Listener) error {
	shutdown := s.shutdown
	// Start routine waiting for signals
	shutdown.RegisterSignalHandler(func() {
		// Stop the HTTP servers
		s.log.Info("api server shutting down")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := s.api.Shutdown(ctx); err != nil {
			s.log.Error(err, "api server shutdown problem")
		}
	})
	//
	s.registerViews(s.engine)
	//
	s.log.Info("starting to serve RESTful API", "addr", apiLis.Addr())
	err := s.api.Serve(apiLis)
	if !shutdown.IsExpected() && err != nil {
		panic(fmt.Sprintf("shutdown unexpected: %v", err))
	}
	s.log.Info("api server stopped")
	// Check if we are expecting shutdown
	// Wait for both shutdown signals and close the channel
	if ok := shutdown.WaitOrTimeout(30 * time.Second); !ok {
		panic("shutting down gracefully exceeded 30 seconds")
	}
	return nil
}

// Tell the server to shutdown
func (s *authServer) Shutdown() {
	s.shutdown.Expect()
}

// registerViews registers all routes necessary serving.
func (s *authServer) registerViews(r *gin.Engine) {
	// readiness check
	r.GET("/readyz", func(c *gin.Context) {
		c.String(http.StatusOK, "ready")
	})

	// auth route
	r.GET("/auth/*route", s.auth)
	r.PUT("/auth/*route", s.auth)
	r.POST("/auth/*route", s.auth)
	r.PATCH("/auth/*route", s.auth)
	r.DELETE("/auth/*route", s.auth)
	r.OPTIONS("/auth/*route", s.auth)

	// OIDC
	r.GET("/.well-known/openid-configuration", s.discovery)
	r.GET("/keys", s.keys)
}

func (s *authServer) discovery(c *gin.Context) {
	c.JSON(http.StatusOK, &OpenIdConfiguration{
		Issuer:  fmt.Sprintf("https://%s", c.Request.Host),
		JwksURL: fmt.Sprintf("https://%s%s", c.Request.Host, "/keys"),
	})
}

func (s *authServer) keys(c *gin.Context) {
	c.Writer.Header().Set("Cache-Control", fmt.Sprintf("max-age=%d, must-revalidate", int(60*60*24)))
	c.JSON(http.StatusOK, s.oidcServer.Keys())
}

// auth serves as handler for the auth route of the server.
func (s *authServer) auth(c *gin.Context) {
	var err error
	var authToken *jwt.AuthToken
	route := c.Param("route")
	authenticated := false
	authorized := false

	// Logging
	s.log.Info("Authenticating request...", "route", route)
	for k := range c.Request.Header {
		// Print headers
		s.log.V(logger.DebugLevel).Info("Metadata provided.", "Key", k, "Value", c.Request.Header.Get(k))
	}

	// Authenticate user
	if !authenticated {
		authToken = s.tokenValidationFromContext(c) // via JWT
		authenticated = authToken != nil
	}
	if !authenticated {
		authToken = s.certValidation(c) // via client certificate validation
		authenticated = authToken != nil
	}
	if !authenticated {
		c.String(http.StatusUnauthorized, "authentication failed")
		return
	}

	// Authorize user
	authorized, err = s.validatePolicies(c)
	if err != nil {
		s.log.Error(err, "Error checking authorization of user.")
	}
	if authorized {
		s.writeSuccess(c, authToken)
		return
	}
	c.String(http.StatusUnauthorized, "authorization failed")
}

// validatePolicies validates the configured policies using OPA
func (s *authServer) validatePolicies(c *gin.Context) (bool, error) {
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
func (s *authServer) tokenValidationFromContext(c *gin.Context) *jwt.AuthToken {
	authToken := s.tokenValidation(c.Request.Context(), defaultBearerTokenFromRequest(c.Request))
	if authToken == nil {
		return nil
	}

	// Check user actually exists in m8

	user, err := s.userRepo.ByEmail(c, authToken.Email)
	if err != nil && !authToken.IsAPIToken {
		s.log.Info("Token validation failed. User does not exist.", "Email", authToken.Email)
		return nil
	}

	// Validate scopes
	route := c.Param("route")
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
func (s *authServer) certValidation(c *gin.Context) *jwt.AuthToken {
	s.log.Info("Validating client certificate...")

	cert, err := clientCertificateFromRequest(c.Request)
	if err != nil {
		s.log.Info("Certificate validation failed.", "error", err.Error())
		return nil
	}

	if userId, ok := s.retrieveUserId(c, cert.EmailAddresses[0]); !ok {
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

// writeSuccess writes all request headers back to the response along with information got from the upstream IdP and sends an http status ok
func (s *authServer) writeSuccess(c *gin.Context, claims *jwt.AuthToken) {
	// Copy request headers to response
	for k := range c.Request.Header {
		// Avoid copying the original Content-Length header from the client
		if strings.ToLower(k) == "content-length" {
			continue
		}
		c.Writer.Header().Set(k, c.Request.Header.Get(k))
	}

	// Set headers with auth info
	c.Writer.Header().Set(auth.HeaderAuthId, claims.Subject)
	c.Writer.Header().Set(auth.HeaderAuthName, claims.Name)
	c.Writer.Header().Set(auth.HeaderAuthEmail, claims.Email)
	c.Writer.Header().Set(auth.HeaderAuthNotBefore, claims.NotBefore.Time().Format(auth.HeaderAuthNotBeforeFormat))

	c.Writer.WriteHeader(http.StatusOK)
}
