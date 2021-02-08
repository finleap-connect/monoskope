package eventstore

import (
	"io"

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
func NewApiServer(store evs.Store, bus evs.EventBusPublisher) *apiServer {
	s := &apiServer{
		log:   logger.WithName("server"),
		store: store,
		bus:   bus,
	}

	return s
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
