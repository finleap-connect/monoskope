package eventhandler

import (
	"context"
	"io"

	api_es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/eventstore"
	es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/event_sourcing"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/event_sourcing/errors"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type EventStoreMiddleware struct {
	eventStoreClient api_es.EventStoreClient
	handler          es.EventHandler
}

func NewEventStoreMiddleware(eventStoreClient api_es.EventStoreClient, handler es.EventHandler) es.EventHandler {
	return &EventStoreMiddleware{
		eventStoreClient: eventStoreClient,
		handler:          handler,
	}
}

// HandleEvent implements the HandleEvent method of the es.EventHandler interface.
func (m *EventStoreMiddleware) HandleEvent(ctx context.Context, event es.Event) error {
	err := m.handler.HandleEvent(ctx, event)

	// If the next handler in the chain tells that the projection is outdated
	if err == errors.ErrProjectionOutdated {
		err = m.applyEventsFromStore(ctx, event)
		if err != nil {
			return err
		}
	}

	return m.handler.HandleEvent(ctx, event)
}

func (m *EventStoreMiddleware) applyEventsFromStore(ctx context.Context, event es.Event) error {
	// Retrieve events from store
	eventStream, err := m.eventStoreClient.Retrieve(ctx, &api_es.EventFilter{
		ByAggregate: &api_es.EventFilter_AggregateType{
			AggregateType: wrapperspb.String(event.AggregateID().String()),
		},
	})
	if err != nil {
		return err
	}

	for {
		// Read next
		protoEvent, err := eventStream.Recv()

		if err != nil {
			if err == io.EOF {
				// End of stream
				break
			} else {
				return err
			}
		}

		// Convert event from api to es
		esEvent, err := es.NewEventFromProto(protoEvent)
		if err != nil {
			return err
		}

		// Let the next handler in the chain handle the event
		err = m.handler.HandleEvent(ctx, esEvent)
		if err != nil {
			return err
		}
	}

	return nil
}
