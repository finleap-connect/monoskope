package usecases

import (
	"context"

	api_es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/eventstore"
	evs "gitlab.figo.systems/platform/monoskope/monoskope/pkg/event_sourcing"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/event_sourcing/messaging"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/event_sourcing/storage"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/logger"
)

type StoreEventsUseCase struct {
	UseCaseBase

	store  storage.Store
	bus    messaging.EventBusPublisher
	events []*api_es.Event
}

// NewStoreEventsUseCase creates a new usecase which stores all events in the store
// and broadcasts these events via the message bus
func NewStoreEventsUseCase(ctx context.Context, store storage.Store, bus messaging.EventBusPublisher, events []*api_es.Event) UseCase {
	useCase := &StoreEventsUseCase{
		UseCaseBase: UseCaseBase{
			log: logger.WithName("store-events-use-case"),
			ctx: ctx,
		},
		store:  store,
		bus:    bus,
		events: events,
	}
	return useCase
}

func (u *StoreEventsUseCase) Run() error {
	// Convert from proto events to storage events
	var storageEvents []evs.Event
	for _, v := range u.events {
		ev, err := evs.NewEventFromProto(v)
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

	u.log.Info("Sending events to the message bus...")
	for _, event := range storageEvents {
		err := u.bus.PublishEvent(u.ctx, event)
		if err != nil {
			return err
		}
	}

	return nil
}
