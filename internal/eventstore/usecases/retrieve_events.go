package usecases

import (
	api_es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/eventstore"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/logger"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/storage"
)

type RetrieveEventsUseCase struct {
	UseCaseBase

	store  storage.Store
	filter *api_es.EventFilter
	stream api_es.EventStore_RetrieveServer
}

// NewRetrieveEventsUseCase creates a new usecase which retrieves all events
// from the store which match the filter
func NewRetrieveEventsUseCase(stream api_es.EventStore_RetrieveServer, store storage.Store, filter *api_es.EventFilter) UseCase {
	useCase := &RetrieveEventsUseCase{
		UseCaseBase: UseCaseBase{
			log: logger.WithName("retrieve-events-use-case"),
			ctx: stream.Context(),
		},
		store:  store,
		filter: filter,
		stream: stream,
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
	u.log.Info("Retrieving events from the store...")
	events, err := u.store.Load(u.ctx, sq)
	if err != nil {
		return err
	}

	// Send events to client
	u.log.Info("Sending events to client...")
	for _, e := range events {
		protoEvent, err := NewProtoFromEvent(e)
		if err != nil {
			return err
		}

		err = u.stream.Send(protoEvent)
		if err != nil {
			return err
		}
	}

	return nil
}
