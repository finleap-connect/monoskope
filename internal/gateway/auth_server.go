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

	"gitlab.figo.systems/platform/monoskope/monoskope/internal/gateway/auth"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/logger"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/util"
)

const (
	defaultAuthorizationHeader = "Authorization"

	HeaderAuthName   = "x-auth-name"
	HeaderAuthEmail  = "x-auth-email"
	HeaderAuthIssuer = "x-auth-issuer"
)

// authServer is the AuthN/AuthZ decision API used as Ambassador Auth Service.
type authServer struct {
	api         *http.Server
	engine      *gin.Engine
	log         logger.Logger
	shutdown    *util.ShutdownWaitGroup
	authHandler *auth.Handler
}

// NewAuthServer creates a new instance of gateway.authServer.
func NewAuthServer(authHandler *auth.Handler) *authServer {
	engine := gin.Default()
	s := &authServer{
		api:         &http.Server{Handler: engine},
		engine:      engine,
		log:         logger.WithName("auth-server"),
		shutdown:    util.NewShutdownWaitGroup(),
		authHandler: authHandler,
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

// registerViews registers all routes necessary serving.
func (s *authServer) registerViews(r *gin.Engine) {
	r.GET("/readyz", func(c *gin.Context) {
		c.String(http.StatusOK, "ready")
	})
	r.POST("/auth/*route", s.auth)
}

// auth serves as handler for the auth route of the server.
func (s *authServer) auth(c *gin.Context) {
	route := c.Param("route")

	// Allow access to the Gateway in general
	if strings.HasPrefix(route, "/gateway.Gateway/") {
		s.log.Info("Allowing anonymous access.", "Route", route)
		c.Writer.WriteHeader(http.StatusOK)
		return
	}

	s.log.Info("Token validation requested...", "Route", route)
	var claims *auth.Claims
	var err error
	if claims, err = s.authHandler.Authorize(c.Request.Context(), defaultBearerTokenFromRequest(c.Request)); err != nil {
		s.log.Error(err, "Token validation failed.")
		c.String(http.StatusUnauthorized, fmt.Sprintf("error validating token: %v", err.Error()))
		return
	}
	s.log.Info("Token validation successful.", "Route", route, "User", claims.Email)

	for k := range c.Request.Header {
		// Avoid copying the original Content-Length header from the client
		if strings.ToLower(k) == "content-length" {
			continue
		}
		c.Writer.Header().Set(k, c.Request.Header.Get(k))
	}

	// Set headers with auth info
	c.Writer.Header().Set(HeaderAuthName, claims.Name)
	c.Writer.Header().Set(HeaderAuthEmail, claims.Email)
	c.Writer.Header().Set(HeaderAuthIssuer, claims.Issuer)

	c.Writer.WriteHeader(http.StatusOK)
}

// defaultBearerTokenFromRequest extracts the token from header
func defaultBearerTokenFromRequest(r *http.Request) string {
	token := r.Header.Get(defaultAuthorizationHeader)
	split := strings.SplitN(token, " ", 2)
	if len(split) != 2 || !strings.EqualFold(strings.ToLower(split[0]), "bearer") {
		return ""
	}
	return split[1]
}
