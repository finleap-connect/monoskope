package eventhandler

import (
	"context"
	"errors"
	"io"

	apiEs "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/eventsourcing"
	es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type eventStoreReplayMiddleware struct {
	esClient apiEs.EventStoreClient
	handler  es.EventHandler
}

// NewEventStoreReplayMiddleware creates an EventHandler which automates querying the EventStore in case of gaps in AggregateVersion found in other EventHandlers later in the chain of EventHandlers.
func NewEventStoreReplayMiddleware(esClient apiEs.EventStoreClient) *eventStoreReplayMiddleware {
	return &eventStoreReplayMiddleware{
		esClient: esClient,
	}
}

func (m *eventStoreReplayMiddleware) Middleware(h es.EventHandler) es.EventHandler {
	m.handler = h
	return m
}

// HandleEvent implements the HandleEvent method of the es.EventHandler interface.
func (m *eventStoreReplayMiddleware) HandleEvent(ctx context.Context, event es.Event) error {
	var outdatedError *ProjectionOutdatedError
	if err := m.handler.HandleEvent(ctx, event); errors.As(err, &outdatedError) {
		// If the next handler in the chain tells that the projection is outdated
		if err := m.applyEventsFromStore(ctx, event, outdatedError.ProjectionVersion); err != nil {
			return err
		}
	}

	return m.handler.HandleEvent(ctx, event)
}

func (m *eventStoreReplayMiddleware) applyEventsFromStore(ctx context.Context, event es.Event, projectionVersion uint64) error {
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
