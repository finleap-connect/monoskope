package eventstore

import (
	"context"
	"io"
	"time"

	"gitlab.figo.systems/platform/monoskope/monoskope/internal/eventstore/usecases"
	api_es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/eventstore"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/logger"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/messaging"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/storage"
	"google.golang.org/protobuf/types/known/emptypb"
)

// apiServer is the implementation of the EventStore API
type apiServer struct {
	api_es.UnimplementedEventStoreServer
	// Logger interface
	log   logger.Logger
	store storage.Store
	bus   messaging.EventBusPublisher
}

// NewApiServer returns a new configured instance of apiServer
func NewApiServer(store storage.Store, bus messaging.EventBusPublisher) (*apiServer, error) {
	s := &apiServer{
		log:   logger.WithName("server"),
		store: store,
		bus:   bus,
	}

	s.log.Info("connecting to the message bus")
	ctx, cancelFunc := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancelFunc()
	msgbusErr := s.bus.Connect(ctx)
	if msgbusErr != nil {
		s.log.Error(msgbusErr, "failed connecting the message bus")
		return nil, msgbusErr.Cause()
	}

	s.log.Info("connecting to the storage backend")
	ctx, cancelFunc = context.WithTimeout(context.Background(), 30*time.Second)
	defer cancelFunc()
	err := s.store.Connect(ctx)
	if err != nil {
		s.log.Error(err, "failed connecting to the storage backend")
		return nil, err
	}

	return s, nil
}

// // Serve starts the api listeners of the Server
// func (s *server) Serve(apiLis net.Listener, metricsLis net.Listener) error {
// 	shutdown := s.shutdown

// 	if metricsLis != nil {
// 		// Start the http server in a different goroutine
// 		shutdown.Add(1)

// 		go func() {
// 			s.log.Info("starting to serve prometheus metrics", "addr", metricsLis.Addr())
// 			err := s.http.Serve(metricsLis)
// 			// If shutdown is expected, we don't care about the error,
// 			// but if we do not expect shutdown, we panic!
// 			if !shutdown.IsExpected() && err != nil {
// 				panic(fmt.Sprintf("shutdown unexpected: %v", err))
// 			}
// 			s.log.Info("http server stopped")
// 			shutdown.Done() // Notify workgroup
// 		}()
// 	}

// 	// Start routine waiting for signals
// 	shutdown.RegisterSignalHandler(func() {
// 		// Stop the HTTP server
// 		s.log.Info("http server shutting down")
// 		if err := s.http.Shutdown(context.Background()); err != nil {
// 			s.log.Error(err, "http server shutdown problem")
// 		}

// 		// And the gRPC server
// 		s.log.Info("grpc server stopping gracefully")
// 		s.grpc.GracefulStop()

// 		// And the bus
// 		s.log.Info("closing connection to message bus gracefully")
// 		if err := s.bus.Close(); err != nil {
// 			s.log.Error(err, "message bzs shutdown problem")
// 		}

// 		// And the store
// 		s.log.Info("closing connection to store gracefully")
// 		if err := s.store.Close(); err != nil {
// 			s.log.Error(err, "store shutdown problem")
// 		}
// 	})

// 	s.log.Info("connecting to the message bus")
// 	ctx, cancelFunc := context.WithTimeout(context.Background(), 30*time.Second)
// 	defer cancelFunc()
// 	msgbusErr := s.bus.Connect(ctx)
// 	if msgbusErr != nil {
// 		s.log.Error(msgbusErr, "failed connecting the message bus")
// 		return msgbusErr.Cause()
// 	}

// 	s.log.Info("connecting to the storage backend")
// 	ctx, cancelFunc = context.WithTimeout(context.Background(), 30*time.Second)
// 	defer cancelFunc()
// 	err := s.store.Connect(ctx)
// 	if err != nil {
// 		s.log.Error(err, "failed connecting to the storage backend")
// 		return err
// 	}

// 	s.log.Info("starting to serve grpc", "addr", apiLis.Addr())
// 	err = s.grpc.Serve(apiLis)
// 	s.log.Info("grpc server stopped")

// 	// Check if we are expecting shutdown
// 	if !shutdown.IsExpected() {
// 		panic(fmt.Sprintf("shutdown unexpected, grpc serve returned: %v", err))
// 	}
// 	// Wait for both shutdown signals and close the channel
// 	if ok := shutdown.WaitOrTimeout(30 * time.Second); !ok {
// 		panic("shutting down gracefully exceeded 30 seconds")
// 	}
// 	return err // Return the error, if grpc stopped gracefully there is no error
// }

// Store implements the API method for storing events
func (s *apiServer) Store(stream api_es.EventStore_StoreServer) error {
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
	err := usecases.NewStoreEventsUseCase(stream.Context(), s.store, s.bus, eventStream).Run()
	if err != nil {
		return err
	}

	return stream.SendAndClose(&emptypb.Empty{})
}

// Retrieve implements the API method for retrieving events from the store
func (s *apiServer) Retrieve(filter *api_es.EventFilter, stream api_es.EventStore_RetrieveServer) error {
	// Perform the use case for storing events
	err := usecases.NewRetrieveEventsUseCase(stream, s.store, filter).Run()
	if err != nil {
		return err
	}
	return nil
}
