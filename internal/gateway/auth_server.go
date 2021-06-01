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

	"gitlab.figo.systems/platform/monoskope/monoskope/internal/gateway/auth"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/repositories"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/jwt"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/logger"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/util"
)

const (
	HeaderAuthorization = "Authorization"

	HeaderAuthId              = "x-auth-id"
	HeaderAuthName            = "x-auth-name"
	HeaderAuthEmail           = "x-auth-email"
	HeaderAuthIssuer          = "x-auth-issuer"
	HeaderForwardedClientCert = "x-forwarded-client-cert"
)

// authServer is the AuthN/AuthZ decision API used as Ambassador Auth Service.
type authServer struct {
	api         *http.Server
	engine      *gin.Engine
	log         logger.Logger
	shutdown    *util.ShutdownWaitGroup
	authHandler *auth.Handler
	userRepo    repositories.ReadOnlyUserRepository
}

// NewAuthServer creates a new instance of gateway.authServer.
func NewAuthServer(authHandler *auth.Handler, userRepo repositories.ReadOnlyUserRepository) *authServer {
	engine := gin.Default()
	s := &authServer{
		api:         &http.Server{Handler: engine},
		engine:      engine,
		log:         logger.WithName("auth-server"),
		shutdown:    util.NewShutdownWaitGroup(),
		authHandler: authHandler,
		userRepo:    userRepo,
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
	c.JSON(http.StatusOK, s.authHandler.Discovery("/keys"))
}

func (s *authServer) keys(c *gin.Context) {
	c.Writer.Header().Set("Cache-Control", fmt.Sprintf("max-age=%d, must-revalidate", int(s.authHandler.KeyExpiration().Seconds())))
	c.JSON(http.StatusOK, s.authHandler.Keys())
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

func (s *authServer) retrieveUserId(c *gin.Context, email string) (string, bool) {
	user, err := s.userRepo.ByEmail(c, email)
	if err != nil {
		return "", false
	}
	return user.Id, true
}

// tokenValidationFromContext validates the token provided within the authorization flow from gin context
func (s *authServer) tokenValidationFromContext(c *gin.Context) *jwt.AuthToken {
	return s.tokenValidation(c.Request.Context(), defaultBearerTokenFromRequest(c.Request), jwt.AudienceMonoctl, jwt.AudienceM8Operator)
}

// tokenValidation validates the token provided within the authorization flow
func (s *authServer) tokenValidation(ctx context.Context, token string, expectedAudience ...string) *jwt.AuthToken {
	s.log.Info("Validating token...")

	if token == "" {
		s.log.Info("Token validation failed.", "error", "token is empty")
		return nil
	}

	authToken := &jwt.AuthToken{}
	if err := s.authHandler.Authorize(ctx, token, authToken); err != nil {
		s.log.Info("Token validation failed.", "error", err.Error())
		return nil
	} else if err := authToken.Validate(expectedAudience...); err != nil {
		s.log.Info("Token validation failed.", "error", err.Error())
		return nil
	}
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
		claims := jwt.NewAuthToken(&jwt.StandardClaims{
			Name:  cert.Subject.CommonName,
			Email: cert.EmailAddresses[0],
		}, userId, "mtls")
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
	c.Writer.Header().Set(HeaderAuthId, claims.Subject)
	c.Writer.Header().Set(HeaderAuthName, claims.Name)
	c.Writer.Header().Set(HeaderAuthEmail, claims.Email)
	c.Writer.Header().Set(HeaderAuthIssuer, claims.Issuer)

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
	pemData := r.Header.Get(HeaderForwardedClientCert)
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
