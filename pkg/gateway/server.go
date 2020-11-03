package gateway

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/golang/protobuf/ptypes/empty"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	api_gw "gitlab.figo.systems/platform/monoskope/monoskope/api/gateway"
	api_gwauth "gitlab.figo.systems/platform/monoskope/monoskope/api/gateway/auth"
	"gitlab.figo.systems/platform/monoskope/monoskope/internal/metadata"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/auth"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/grpcutil"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/logger"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/metrics"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/util"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Server struct {
	api_gw.UnimplementedGatewayServer
	// HTTP-server exposing the metrics
	http *http.Server
	// gRPC-server exposing both the API and health
	grpc *grpc.Server
	// Logger interface
	log logger.Logger
	//
	shutdown    *util.ShutdownWaitGroup
	authConfig  *auth.Config
	authHandler *auth.Handler
}

type ServerConfig struct {
	KeepAlive  bool
	AuthConfig *auth.Config
	TlsCert    *tls.Certificate
}

func NewServer(conf *ServerConfig) (*Server, error) {
	s := &Server{
		http:       metrics.NewServer(),
		log:        logger.WithName("gateway"),
		shutdown:   util.NewShutdownWaitGroup(),
		authConfig: conf.AuthConfig,
	}

	// Create interceptor for auth
	authHandler, err := auth.NewHandler(conf.AuthConfig)
	if err != nil {
		return nil, err
	}
	s.authHandler = authHandler

	authInterceptor, err := auth.NewInterceptor(authHandler)
	if err != nil {
		return nil, err
	}

	// Configure gRPC server
	opts := []grpc.ServerOption{ // add prometheus metrics interceptors
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			grpc_prometheus.StreamServerInterceptor,
			auth.StreamServerInterceptor(authInterceptor.EnsureValid),
		)),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_prometheus.UnaryServerInterceptor,
			auth.UnaryServerInterceptor(authInterceptor.EnsureValid),
		)),
	}
	if conf.KeepAlive {
		opts = append(opts, grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle: 5 * time.Minute,
			Time:              2 * time.Second,
		}))
	}
	if conf.TlsCert != nil {
		opts = append(opts, grpc.Creds(credentials.NewServerTLSFromCert(conf.TlsCert)))
	}
	s.grpc = grpc.NewServer(opts...)

	// Add user-authenticator service
	api_gw.RegisterGatewayServer(s.grpc, s)
	// Add grpc health check service
	healthpb.RegisterHealthServer(s.grpc, health.NewServer())
	// Register the metric interceptors with prometheus
	grpc_prometheus.Register(s.grpc)
	// Enable reflection API
	reflection.Register(s.grpc)
	return s, nil
}

func (s *Server) Serve(apiLis net.Listener, metricsLis net.Listener) error {
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

func (s *Server) GetServerInfo(context.Context, *empty.Empty) (*api_gw.ServerInformation, error) {
	return &api_gw.ServerInformation{
		Version: metadata.Version,
		Commit:  metadata.Commit,
	}, nil
}

func (s *Server) GetAuthInformation(ctx context.Context, state *api_gwauth.AuthState) (*api_gwauth.AuthInformation, error) {
	url, err := s.authHandler.GetAuthCodeURL(state, &auth.AuthCodeURLConfig{
		Scopes:        []string{"offline_access"},
		Clients:       []string{},
		OfflineAccess: true,
	})
	if err != nil {
		return nil, grpcutil.ErrInvalidArgument(err)
	}

	return &api_gwauth.AuthInformation{AuthCodeURL: url}, nil
}

func (s *Server) ExchangeAuthCode(ctx context.Context, code *api_gwauth.AuthCode) (*api_gwauth.UserInfo, error) {
	token, err := s.authHandler.Exchange(ctx, code.GetCode())
	if err != nil {
		return nil, grpcutil.ErrInvalidArgument(err)
	}

	err = s.authHandler.VerifyState(ctx, token, code.GetState())
	if err != nil {
		return nil, grpcutil.ErrInvalidArgument(err)
	}

	return &api_gwauth.UserInfo{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		Expiry:       timestamppb.New(token.Expiry),
	}, nil
}
