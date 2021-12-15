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
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"crypto/x509"
	"encoding/pem"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/finleap-connect/monoskope/internal/gateway/auth"
	"github.com/finleap-connect/monoskope/pkg/api/gateway"
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

// authServer is the AuthN/AuthZ decision API used as Ambassador Auth Service.
type authServer struct {
	api        *http.Server
	engine     *gin.Engine
	log        logger.Logger
	shutdown   *util.ShutdownWaitGroup
	authClient *auth.Client
	authServer *auth.Server
	userRepo   repositories.ReadOnlyUserRepository
	url        string
}

// NewAuthServer creates a new instance of gateway.authServer.
func NewAuthServer(url string, client *auth.Client, server *auth.Server, userRepo repositories.ReadOnlyUserRepository) *authServer {
	engine := gin.Default()
	s := &authServer{
		api:        &http.Server{Handler: engine},
		engine:     engine,
		log:        logger.WithName("auth-server"),
		shutdown:   util.NewShutdownWaitGroup(),
		authClient: client,
		authServer: server,
		userRepo:   userRepo,
		url:        url,
	}
	engine.Use(gin.Recovery())
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
	r.GET("/readyz", func(c *gin.Context) {
		c.String(http.StatusOK, "ready")
	})
	r.POST("/auth/*route", s.auth)
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
	c.JSON(http.StatusOK, s.authServer.Keys())
}

// auth serves as handler for the auth route of the server.
func (s *authServer) auth(c *gin.Context) {
	route := c.Param("route")
	s.log.Info("Authenticating request...", "route", route)

	// Print headers
	for k := range c.Request.Header {
		s.log.V(logger.DebugLevel).Info("Metadata provided.", "Key", k, "Value", c.Request.Header.Get(k))
	}

	var authToken *jwt.AuthToken
	if authToken = s.tokenValidationFromContext(c); authToken != nil {
		s.writeSuccess(c, authToken)
		return
	}
	if authToken = s.certValidation(c); authToken != nil {
		s.writeSuccess(c, authToken)
		return
	}

	c.String(http.StatusUnauthorized, "authorization failed")
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
	if _, ok := s.retrieveUserId(c, authToken.Email); !ok {
		s.log.Info("Token validation failed. User does not exist.", "Email", authToken.Email)
		return nil
	}

	// Validate scopes
	route := c.Param("route")
	scopes := strings.Split(authToken.Scope, " ")

	if containsString(scopes, gateway.AuthorizationScope_API.String()) {
		return authToken
	}

	s.log.Info("Token validation failed. Token has not correct scopes for route.", "Route", route, "Scopes", authToken.Scope)
	return nil
}

func containsString(slice []string, value string) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}

// tokenValidation validates the token provided within the authorization flow
func (s *authServer) tokenValidation(ctx context.Context, token string) *jwt.AuthToken {
	s.log.Info("Validating token...")

	if token == "" {
		s.log.Info("Token validation failed.", "error", "token is empty")
		return nil
	}

	authToken := &jwt.AuthToken{}
	if err := s.authServer.Authorize(ctx, token, authToken); err != nil {
		s.log.Info("Token validation failed.", "error", err.Error())
		return nil
	}
	if err := authToken.Validate(s.url); err != nil {
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
		}, s.url, userId, time.Minute*5)
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
	c.Writer.Header().Set(auth.HeaderAuthIssuer, claims.Issuer)

	c.Writer.WriteHeader(http.StatusOK)
}

// defaultBearerTokenFromRequest extracts the token from header
func defaultBearerTokenFromRequest(r *http.Request) string {
	token := r.Header.Get(HeaderAuthorization)
	split := strings.SplitN(token, " ", 2)
	if len(split) != 2 || !strings.EqualFold(strings.ToLower(split[0]), "bearer") {
		return ""
	}
	return split[1]
}

// clientCertificateFromRequest extracts the client certificate from header
func clientCertificateFromRequest(r *http.Request) (*x509.Certificate, error) {
	pemData := r.Header.Get(auth.HeaderForwardedClientCert)
	if pemData == "" {
		return nil, errors.New("cert header is empty")
	}

	decodedValue, err := url.QueryUnescape(pemData)
	if err != nil {
		return nil, errors.New("could not unescape pem data from header")
	}

	block, _ := pem.Decode([]byte(decodedValue))
	if block == nil {
		return nil, errors.New("decoding pem failed")
	}

	return x509.ParseCertificate(block.Bytes)
}
