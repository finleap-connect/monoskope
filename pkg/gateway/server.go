package gateway

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/golang/protobuf/ptypes/empty"
	promgrpc "github.com/grpc-ecosystem/go-grpc-prometheus"

	"gitlab.figo.systems/platform/monoskope/monoskope/api"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/grpcutil"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/logger"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/metrics"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/util"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	api.UnimplementedGatewayServer

	// HTTP-server exposing the metrics
	http *http.Server
	// gRPC-server exposing both the API and health
	grpc *grpc.Server
	// Logger interface
	log logger.Logger
	//
	shutdown *util.ShutdownWaitGroup
}

func NewServer(keepAlive bool) *Server {
	// Configure gRPC server
	opts := []grpc.ServerOption{ // add prometheus metrics interceptors
		grpc.StreamInterceptor(promgrpc.StreamServerInterceptor),
		grpc.UnaryInterceptor(promgrpc.UnaryServerInterceptor),
		grpc.UnaryInterceptor(ensureValidToken),
	}
	if keepAlive {
		opts = append(opts, grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle: 5 * time.Minute,
			Time:              2 * time.Second,
		}))
	}
	s := &Server{
		http:     metrics.NewServer(),
		grpc:     grpc.NewServer(opts...),
		log:      logger.WithName("gateway"),
		shutdown: util.NewShutdownWaitGroup(),
	}
	// Add user-authenticator service
	api.RegisterGatewayServer(s.grpc, s)
	// Add grpc health check service
	healthpb.RegisterHealthServer(s.grpc, health.NewServer())
	// Register the metric interceptors with prometheus
	promgrpc.Register(s.grpc)
	// Enable reflection API
	reflection.Register(s.grpc)
	return s
}

func (s *Server) Serve(apiLis net.Listener, metricsLis net.Listener) error {
	shutdown := s.shutdown
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
	//
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

// valid validates the authorization.
func valid(authorization []string) bool {
	if len(authorization) < 1 {
		return false
	}
	token := strings.TrimPrefix(authorization[0], "Bearer ")
	// Perform the token validation here. For the sake of this example, the code
	// here forgoes any of the usual OAuth2 token validation and instead checks
	// for a token matching an arbitrary string.
	return token == "some-secret-token"
}

// ensureValidToken ensures a valid token exists within a request's metadata. If
// the token is missing or invalid, the interceptor blocks execution of the
// handler and returns an error. Otherwise, the interceptor invokes the unary
// handler.
func ensureValidToken(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, grpcutil.ErrMissingMetadata
	}
	// The keys within metadata.MD are normalized to lowercase.
	// See: https://godoc.org/google.golang.org/grpc/metadata#New
	if !valid(md["authorization"]) {
		return nil, grpcutil.ErrInvalidToken
	}
	// Continue execution of handler after ensuring a valid token.
	return handler(ctx, req)
}

func (s *Server) GetServerInfo(context.Context, *empty.Empty) (*api.ServerInformation, error) {
	return nil, nil
}
