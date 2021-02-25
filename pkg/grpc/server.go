package grpc

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/logger"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/metrics"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/util"
	"google.golang.org/grpc"
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
func NewServer(name string, keepAlive bool) *Server {
	return NewServerWithOpts(name, keepAlive, []grpc.UnaryServerInterceptor{}, []grpc.StreamServerInterceptor{})
}

// NewServerWithOpts returns a new configured instance of Server with additional interceptros specified
func NewServerWithOpts(name string, keepAlive bool, unaryServerInterceptors []grpc.UnaryServerInterceptor, streamServerInterceptors []grpc.StreamServerInterceptor) *Server {
	s := &Server{
		http:     metrics.NewServer(),
		log:      logger.WithName(name),
		shutdown: util.NewShutdownWaitGroup(),
	}

	unaryServerInterceptors = append(unaryServerInterceptors, grpc_prometheus.UnaryServerInterceptor)    // add prometheus metrics interceptors
	streamServerInterceptors = append(streamServerInterceptors, grpc_prometheus.StreamServerInterceptor) // add prometheus metrics interceptors

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
