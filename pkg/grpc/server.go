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

package grpc

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	grpc_validator_wrapper "github.com/finleap-connect/monoskope/pkg/grpc/middleware/validator"
	"go.uber.org/zap/zapgrpc"

	"github.com/finleap-connect/monoskope/pkg/logger"
	"github.com/finleap-connect/monoskope/pkg/metrics"
	"github.com/finleap-connect/monoskope/pkg/util"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
)

// Server is the implementation of an API Server
type Server struct {
	// HTTP-server exposing the metrics
	http *http.Server
	// gRPC-server exposing both the API and health
	grpc *grpc.Server
	// Logger interface
	log logger.Logger
	//
	shutdown *util.ShutdownWaitGroup
}

// NewServer returns a new configured instance of Server
func NewServer(name string, keepAlive bool, opt ...grpc.ServerOption) *Server {
	return NewServerWithOpts(name, keepAlive, []grpc.UnaryServerInterceptor{}, []grpc.StreamServerInterceptor{}, opt...)
}

// NewServerWithOpts returns a new configured instance of Server with additional interceptors specified
func NewServerWithOpts(name string, keepAlive bool, unaryServerInterceptors []grpc.UnaryServerInterceptor, streamServerInterceptors []grpc.StreamServerInterceptor, opt ...grpc.ServerOption) *Server {
	s := &Server{
		http:     metrics.NewServer(),
		log:      logger.WithName(name),
		shutdown: util.NewShutdownWaitGroup(),
	}
	grpclog.SetLoggerV2(zapgrpc.NewLogger(logger.GetZapLogger()))

	// Add default interceptors
	unaryServerInterceptors = append(unaryServerInterceptors,
		grpc_ctxtags.UnaryServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
		grpc_prometheus.UnaryServerInterceptor, // add prometheus metrics interceptors
		grpc_recovery.UnaryServerInterceptor(), // add recovery from panics
		// own wrapper is used to unpack nested messages
		//grpc_validator.UnaryServerInterceptor(), // add message validator
		grpc_validator_wrapper.UnaryServerInterceptor(), // add message validator wrapper
		otelgrpc.UnaryServerInterceptor(),
	)
	streamServerInterceptors = append(streamServerInterceptors,
		grpc_ctxtags.StreamServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
		grpc_prometheus.StreamServerInterceptor, // add prometheus metrics interceptors
		grpc_recovery.StreamServerInterceptor(), // add recovery from panics
		// own wrapper is used to unpack nested messages
		//grpc_validator.StreamServerInterceptor(), // add message validator
		grpc_validator_wrapper.StreamServerInterceptor(), // add message validator wrapper
		otelgrpc.StreamServerInterceptor(),
	)

	// Configure gRPC server
	opts := []grpc.ServerOption{
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(streamServerInterceptors...)),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(unaryServerInterceptors...)),
	}
	if keepAlive {
		opts = append(opts, grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle: 5 * time.Minute,
			Time:              2 * time.Second,
		}))
	}
	s.grpc = grpc.NewServer(opts...)

	// Add grpc health check service
	healthpb.RegisterHealthServer(s.grpc, health.NewServer())

	// Register the metric interceptors with prometheus
	grpc_prometheus.Register(s.grpc)
	grpc_prometheus.EnableHandlingTimeHistogram()

	// Enable reflection API
	reflection.Register(s.grpc)

	return s
}

// RegisterService registers your gRPC service implementation with the server
func (s *Server) RegisterService(f func(grpc.ServiceRegistrar)) {
	f(s.grpc)
}

// Serve starts the api listeners of the Server
func (s *Server) Serve(apiAddr, metricsAddr string) error {
	// Setup grpc listener
	apiLis, err := net.Listen("tcp", apiAddr)
	if err != nil {
		return err
	}
	defer apiLis.Close()

	// Setup metrics listener
	metricsLis, err := net.Listen("tcp", metricsAddr)
	if err != nil {
		return err
	}
	defer metricsLis.Close()

	return s.ServeFromListener(apiLis, metricsLis)
}

// ServeFromListener starts the api listeners of the Server
func (s *Server) ServeFromListener(apiLis net.Listener, metricsLis net.Listener) error {
	shutdown := s.shutdown

	if metricsLis != nil {
		// Start the http server in a different goroutine
		shutdown.Add(1)

		go func() {
			s.log.Info("starting to serve prometheus metrics", "addr", metricsLis.Addr())
			err := s.http.Serve(metricsLis)
			// If shutdown is expected, we don't care about the error,
			// but if we do not expect shutdown, we panic!
			if !shutdown.IsExpected() && err != nil {
				panic(fmt.Sprintf("shutdown unexpected: %v", err))
			}
			s.log.Info("http server stopped")
			shutdown.Done() // Notify workgroup
		}()
	}

	// Start routine waiting for signals
	shutdown.RegisterSignalHandler(func() {
		// Stop the HTTP server
		s.log.Info("http server shutting down")
		if err := s.http.Shutdown(context.Background()); err != nil {
			s.log.Error(err, "http server shutdown problem")
		}

		// And the gRPC server
		s.log.Info("grpc server stopping gracefully")
		s.grpc.GracefulStop()
	})

	s.log.Info("starting to serve grpc", "addr", apiLis.Addr())
	err := s.grpc.Serve(apiLis)
	s.log.Info("grpc server stopped")

	// Check if we are expecting shutdown
	if !shutdown.IsExpected() {
		panic(fmt.Sprintf("shutdown unexpected, grpc serve returned: %v", err))
	}
	// Wait for both shutdown signals and close the channel
	if ok := shutdown.WaitOrTimeout(30 * time.Second); !ok {
		panic("shutting down gracefully exceeded 30 seconds")
	}
	return err // Return the error, if grpc stopped gracefully there is no error
}

// Tell the server to shutdown
func (s *Server) Shutdown() {
	s.shutdown.Expect()
}
