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
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/finleap-connect/monoskope/internal/gateway/auth"
	"github.com/finleap-connect/monoskope/internal/telemetry"
	"github.com/finleap-connect/monoskope/pkg/logger"
	"github.com/finleap-connect/monoskope/pkg/util"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

type OpenIdConfiguration struct {
	Issuer  string `json:"issuer"`
	JwksURL string `json:"jwks_uri"`
}

// authServer implements the AuthN/AuthZ decision API used as Ambassador Auth Service.
type oidcProviderServer struct {
	api        *http.Server
	engine     *gin.Engine
	log        logger.Logger
	shutdown   *util.ShutdownWaitGroup
	oidcServer *auth.Server
}

// NewOIDCProviderServer creates a basic OIDC provider server
func NewOIDCProviderServer(oidcServer *auth.Server) *oidcProviderServer {
	r := gin.Default()
	s := &oidcProviderServer{
		api:        &http.Server{Handler: r},
		engine:     r,
		log:        logger.WithName("auth-server"),
		shutdown:   util.NewShutdownWaitGroup(),
		oidcServer: oidcServer,
	}
	r.Use(cors.Default())
	r.Use(otelgin.Middleware("auth-server"))
	return s
}

// Serve tells the server to start listening on the specified address.
func (s *oidcProviderServer) Serve(apiAddr string) error {
	// Setup grpc listener
	apiLis, err := net.Listen("tcp", apiAddr)
	if err != nil {
		return err
	}
	defer apiLis.Close()

	return s.ServeFromListener(apiLis)
}

// ServeFromListener tells the server to start listening on given listener.
func (s *oidcProviderServer) ServeFromListener(apiLis net.Listener) error {
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
func (s *oidcProviderServer) Shutdown() {
	s.shutdown.Expect()
}

func (s *oidcProviderServer) discovery(c *gin.Context) {
	_, span := telemetry.GetSpan(c.Request.Context(), "getOpenIDConfiguration")
	defer span.End()
	c.JSON(http.StatusOK, &OpenIdConfiguration{
		Issuer:  fmt.Sprintf("https://%s", c.Request.Host),
		JwksURL: fmt.Sprintf("https://%s%s", c.Request.Host, "/keys"),
	})
}

func (s *oidcProviderServer) keys(c *gin.Context) {
	_, span := telemetry.GetSpan(c.Request.Context(), "getKeys")
	defer span.End()
	c.Writer.Header().Set("Cache-Control", fmt.Sprintf("max-age=%d, must-revalidate", int(60*60*24)))
	c.JSON(http.StatusOK, s.oidcServer.Keys())
}

// registerViews registers all routes necessary serving.
func (s *oidcProviderServer) registerViews(r *gin.Engine) {
	// readiness check
	r.GET("/readyz", func(c *gin.Context) {
		c.String(http.StatusOK, "ready")
	})

	// OIDC
	r.GET("/.well-known/openid-configuration", s.discovery)
	r.GET("/keys", s.keys)
}
