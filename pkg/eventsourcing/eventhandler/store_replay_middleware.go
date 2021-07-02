package eventhandler

import (
	"context"
	"errors"
	"io"

	apiEs "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/eventsourcing"
	es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type eventStoreReplayEventHandler struct {
	esClient apiEs.EventStoreClient
	handler  es.EventHandler
}

// NewEventStoreReplayEventHandler creates an EventHandler which automates querying the EventStore in case of gaps in AggregateVersion found in other EventHandlers later in the chain of EventHandlers.
func NewEventStoreReplayEventHandler(esClient apiEs.EventStoreClient) *eventStoreReplayEventHandler {
	return &eventStoreReplayEventHandler{
		esClient: esClient,
	}
}

func (m *eventStoreReplayEventHandler) AsMiddleware(h es.EventHandler) es.EventHandler {
	return &eventStoreReplayEventHandler{
		esClient: m.esClient,
		handler:  h,
	}
}

// HandleEvent implements the HandleEvent method of the es.EventHandler interface.
func (m *eventStoreReplayEventHandler) HandleEvent(ctx context.Context, event es.Event) error {
	var outdatedError *ProjectionOutdatedError
	if err := m.handler.HandleEvent(ctx, event); errors.As(err, &outdatedError) {
		// If the next handler in the chain tells that the projection is outdated
		if err := m.applyEventsFromStore(ctx, event, outdatedError.ProjectionVersion); err != nil {
			return err
		}
		return err
	} else {
		return err
	}
}

func (m *eventStoreReplayEventHandler) applyEventsFromStore(ctx context.Context, event es.Event, projectionVersion uint64) error {
	// Retrieve events from store
	eventStream, err := m.esClient.Retrieve(ctx, &apiEs.EventFilter{
		AggregateType: wrapperspb.String(event.AggregateID().String()),
		MaxVersion:    wrapperspb.UInt64(event.AggregateVersion()),
		MinVersion:    wrapperspb.UInt64(projectionVersion),
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
