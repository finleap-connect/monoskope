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
)

type authServer struct {
	api         *http.Server
	engine      *gin.Engine
	log         logger.Logger
	shutdown    *util.ShutdownWaitGroup
	authHandler *auth.Handler
}

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

func (s *authServer) Serve(apiLis net.Listener) error {
	shutdown := s.shutdown
	// Start routine waiting for signals
	shutdown.RegisterSignalHandler(func() {
		// Stop the HTTP servers
		s.log.Info("metrics server shutting down")
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

func (s *authServer) registerViews(r *gin.Engine) {
	r.GET("/readyz", func(c *gin.Context) {
		c.String(http.StatusOK, "ready")
	})
	r.POST("/auth/*route", s.authN)
}

func (s *authServer) ensureValid(ctx context.Context, token string) (*auth.Claims, error) {
	// Perform the token validation here.
	claims, err := s.authHandler.Authorize(ctx, token)
	if err != nil {
		return nil, err
	}

	return claims, nil
}

func (s *authServer) authN(c *gin.Context) {
	route := c.Param("route")
	s.log.Info("Token validation requested...", "Route", route)

	token := defaultBearerTokenFromRequest(c.Request)
	var claims *auth.Claims

	for k := range c.Request.Header {
		s.log.Info("Request header.", "Key", k, "Value", c.Request.Header.Get(k))
		c.Writer.Header().Set(k, c.Request.Header.Get(k))
	}

	var err error
	if claims, err = s.ensureValid(c.Request.Context(), token); err != nil {
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
	c.Writer.WriteHeader(http.StatusOK)
}

func defaultBearerTokenFromRequest(r *http.Request) string {
	token := r.Header.Get(defaultAuthorizationHeader)
	split := strings.SplitN(token, " ", 2)
	if len(split) != 2 || !strings.EqualFold(strings.ToLower(split[0]), "bearer") {
		return ""
	}
	return split[1]
}
