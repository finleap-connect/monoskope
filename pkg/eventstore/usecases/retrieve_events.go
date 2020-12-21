package usecases

import (
	api_es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/eventstore"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventstore/storage"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/logger"
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
	// Retrieve events from Event Store
	u.log.Info("Retrieving events from the store...")
	events, err := u.store.Load(u.ctx, NewStoreQuery(u.filter))
	if err != nil {
		return err
	}

	// Send events to client
	u.log.Info("Sending events to client...")
	for _, e := range events {
		err := u.stream.Send(NewProtoEvent(e))
		if err != nil {
			return err
		}
	}

	return nil
}
