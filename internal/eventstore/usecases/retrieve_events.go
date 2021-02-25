package usecases

import (
	"gitlab.figo.systems/platform/monoskope/monoskope/internal/eventstore/metrics"
	esApi "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/eventsourcing"
	es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/logger"
)

type RetrieveEventsUseCase struct {
	UseCaseBase

	store   es.Store
	filter  *esApi.EventFilter
	stream  esApi.EventStore_RetrieveServer
	metrics *metrics.EventStoreMetrics
}

// NewRetrieveEventsUseCase creates a new usecase which retrieves all events
// from the store which match the filter
func NewRetrieveEventsUseCase(stream esApi.EventStore_RetrieveServer, store es.Store, filter *esApi.EventFilter, metrics *metrics.EventStoreMetrics) UseCase {
	useCase := &RetrieveEventsUseCase{
		UseCaseBase: UseCaseBase{
			log: logger.WithName("retrieve-events-use-case"),
			ctx: stream.Context(),
		},
		store:   store,
		filter:  filter,
		stream:  stream,
		metrics: metrics,
	}
	return useCase
}

func (u *RetrieveEventsUseCase) Run() error {
	// Convert filter
	sq, err := NewStoreQueryFromProto(u.filter)
	if err != nil {
		return err
	}

	// Retrieve events from Event Store
	u.log.Info("Retrieving events from the database...")
	events, err := u.store.Load(u.ctx, sq)
	if err != nil {
		return err
	}

	// Send events to client
	u.log.Info("Streaming events to client...")
	for _, e := range events {
		protoEvent, err := es.NewProtoFromEvent(e)
		if err != nil {
			return err
		}

		err = u.stream.Send(protoEvent)
		if err != nil {
			return err
		}

		// Count retrieved event
		u.metrics.RetrievedTotalCounter.WithLabelValues(protoEvent.Type, protoEvent.AggregateType).Inc()
	}

	return nil
}
