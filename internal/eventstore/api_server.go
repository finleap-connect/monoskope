package eventstore

import (
	"gitlab.figo.systems/platform/monoskope/monoskope/internal/eventstore/metrics"
	"gitlab.figo.systems/platform/monoskope/monoskope/internal/eventstore/usecases"
	esApi "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/eventsourcing"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/errors"
	es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/logger"
)

// apiServer is the implementation of the EventStore API
type apiServer struct {
	esApi.UnimplementedEventStoreServer
	// Logger interface
	log     logger.Logger
	store   es.Store
	bus     es.EventBusPublisher
	metrics *metrics.EventStoreMetrics
}

// NewApiServer returns a new configured instance of apiServer
func NewApiServer(store es.Store, bus es.EventBusPublisher) *apiServer {
	s := &apiServer{
		log:     logger.WithName("server"),
		store:   store,
		bus:     bus,
		metrics: metrics.NewEventStoreMetrics(),
	}

	return s
}

// Store implements the API method for storing events
func (s *apiServer) Store(stream esApi.EventStore_StoreServer) error {
	// Perform the use case for storing events
	if err := usecases.NewStoreEventsUseCase(stream, s.store, s.bus, s.metrics).Run(); err != nil {
		return errors.TranslateToGrpcError(err)
	}
	return nil
}

// Retrieve implements the API method for retrieving events from the store
func (s *apiServer) Retrieve(filter *esApi.EventFilter, stream esApi.EventStore_RetrieveServer) error {
	// Perform the use case for storing events
	err := usecases.NewRetrieveEventsUseCase(stream, s.store, filter, s.metrics).Run()
	if err != nil {
		return errors.TranslateToGrpcError(err)
	}
	return nil
}
