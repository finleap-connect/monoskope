package eventstore

import (
	"context"
	"io"
	"time"

	"gitlab.figo.systems/platform/monoskope/monoskope/internal/eventstore/usecases"
	api "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/eventstore"
	evs "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/logger"
	"google.golang.org/protobuf/types/known/emptypb"
)

// apiServer is the implementation of the EventStore API
type apiServer struct {
	api.UnimplementedEventStoreServer
	// Logger interface
	log   logger.Logger
	store evs.Store
	bus   evs.EventBusPublisher
}

// NewApiServer returns a new configured instance of apiServer
func NewApiServer(store evs.Store, bus evs.EventBusPublisher) (*apiServer, error) {
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
		return nil, msgbusErr
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

// Shutdown closes all underyling connections
func (s *apiServer) Shutdown() {
	// And the bus
	s.log.Info("closing connection to message bus gracefully")
	if err := s.bus.Close(); err != nil {
		s.log.Error(err, "message bzs shutdown problem")
	}

	// And the store
	s.log.Info("closing connection to store gracefully")
	if err := s.store.Close(); err != nil {
		s.log.Error(err, "store shutdown problem")
	}
}

// Store implements the API method for storing events
func (s *apiServer) Store(stream api.EventStore_StoreServer) error {
	var eventStream []*api.Event
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
func (s *apiServer) Retrieve(filter *api.EventFilter, stream api.EventStore_RetrieveServer) error {
	// Perform the use case for storing events
	err := usecases.NewRetrieveEventsUseCase(stream, s.store, filter).Run()
	if err != nil {
		return err
	}
	return nil
}
