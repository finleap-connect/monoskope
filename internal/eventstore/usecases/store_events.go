package usecases

import (
	"context"

	api_es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/eventstore"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/logger"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/storage"
)

type StoreEventsUseCase struct {
	UseCaseBase

	store  storage.Store
	events []*api_es.Event
}

// NewStoreEventsUseCase creates a new usecase which stores all events in the store
// and broadcasts these events via the message bus
func NewStoreEventsUseCase(ctx context.Context, store storage.Store, events []*api_es.Event) UseCase {
	useCase := &StoreEventsUseCase{
		UseCaseBase: UseCaseBase{
			log: logger.WithName("store-events-use-case"),
			ctx: ctx,
		},
		store:  store,
		events: events,
	}
	return useCase
}

func (u *StoreEventsUseCase) Run() error {
	// Convert from proto events to storage events
	storageEvents := make([]storage.Event, 0)
	for _, v := range u.events {
		ev, err := NewEventFromProto(v)
		if err != nil {
			return err
		}
		storageEvents = append(storageEvents, ev)
	}

	// Store events in Event Store
	u.log.Info("Saving events in the store...")
	err := u.store.Save(u.ctx, storageEvents)
	if err != nil {
		return err
	}

	// TODO: Send events via message bus
	u.log.Info("Sending events to the message bus...")

	return nil
}
