package eventstore

import (
	"io"

	"gitlab.figo.systems/platform/monoskope/monoskope/internal/eventstore/usecases"
	esApi "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/eventsourcing"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/errors"
	es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/logger"
	"google.golang.org/protobuf/types/known/emptypb"
)

// apiServer is the implementation of the EventStore API
type apiServer struct {
	esApi.UnimplementedEventStoreServer
	// Logger interface
	log   logger.Logger
	store es.Store
	bus   es.EventBusPublisher
}

// NewApiServer returns a new configured instance of apiServer
func NewApiServer(store es.Store, bus es.EventBusPublisher) *apiServer {
	s := &apiServer{
		log:   logger.WithName("server"),
		store: store,
		bus:   bus,
	}

	return s
}

// Store implements the API method for storing events
func (s *apiServer) Store(stream esApi.EventStore_StoreServer) error {
	var eventStream []*esApi.Event
	for {
		// Read next event
		event, err := stream.Recv()

		// End of stream
		if err == io.EOF {
			break
		}
		if err != nil { // Some other error
			return errors.TranslateToGrpcError(err)
		}

		// Append events to the stream
		eventStream = append(eventStream, event)
	}

	// Perform the use case for storing events
	err := usecases.NewStoreEventsUseCase(stream.Context(), s.store, s.bus, eventStream).Run()
	if err != nil {
		return errors.TranslateToGrpcError(err)
	}

	return stream.SendAndClose(&emptypb.Empty{})
}

// Retrieve implements the API method for retrieving events from the store
func (s *apiServer) Retrieve(filter *esApi.EventFilter, stream esApi.EventStore_RetrieveServer) error {
	// Perform the use case for storing events
	err := usecases.NewRetrieveEventsUseCase(stream, s.store, filter).Run()
	if err != nil {
		return errors.TranslateToGrpcError(err)
	}
	return nil
}
