package eventstore

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"gitlab.figo.systems/platform/monoskope/monoskope/internal/eventstore/usecases"
	api_es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/eventstore"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/logger"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/metrics"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/storage"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/util"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
)

// Server is the implementation of the API server
type Server struct {
	api_es.UnimplementedEventStoreServer
	// HTTP-server exposing the metrics
	http *http.Server
	// gRPC-server exposing both the API and health
	grpc *grpc.Server
	// Logger interface
	log logger.Logger
	//
	shutdown *util.ShutdownWaitGroup
	store    storage.Store
}

// NewServer returns a new configured instance of Server
func NewServer(conf *ServerConfig) (*Server, error) {
	s := &Server{
		http:     metrics.NewServer(),
		log:      logger.WithName("server"),
		shutdown: util.NewShutdownWaitGroup(),
		store:    conf.Store,
	}

	// Configure gRPC server
	opts := []grpc.ServerOption{ // add prometheus metrics interceptors
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			grpc_prometheus.StreamServerInterceptor,
		)),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_prometheus.UnaryServerInterceptor,
		)),
	}
	if conf.KeepAlive {
		opts = append(opts, grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle: 5 * time.Minute,
			Time:              2 * time.Second,
		}))
	}
	s.grpc = grpc.NewServer(opts...)

	// Add actual api service
	api_es.RegisterEventStoreServer(s.grpc, s)

	// Add grpc health check service
	healthpb.RegisterHealthServer(s.grpc, health.NewServer())
	// Register the metric interceptors with prometheus
	grpc_prometheus.Register(s.grpc)
	// Enable reflection API
	reflection.Register(s.grpc)

	return s, nil
}

// Serve starts the api listeners of the Server
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

// Store implements the API method for storing events
func (s *Server) Store(stream api_es.EventStore_StoreServer) error {
	eventStream := make([]*api_es.Event, 0)
	for {
		// Read next event
		event, err := stream.Recv()

		// End of stream
		if err == io.EOF {
			break
		}
		if err != nil { // Some other error
			return err
		}

		// Append events to the stream
		eventStream = append(eventStream, event)
	}

	// Perform the use case for storing events
	err := usecases.NewStoreEventsUseCase(stream.Context(), s.store, eventStream).Run()
	if err != nil {
		return err
	}

	return stream.SendAndClose(&emptypb.Empty{})
}

// Retrieve implements the API method for retrieving events from the store
func (s *Server) Retrieve(filter *api_es.EventFilter, stream api_es.EventStore_RetrieveServer) error {
	// Perform the use case for storing events
	err := usecases.NewRetrieveEventsUseCase(stream, s.store, filter).Run()
	if err != nil {
		return err
	}
	return nil
}
