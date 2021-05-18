package eventhandler

import (
	"context"

	apiEs "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/eventsourcing"
	es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing"
)

type reactorEventHandler struct {
	esClient apiEs.EventStoreClient
	reactor  es.Reactor
}

// NewReactorEventHandler creates an EventHandler which automates storing Events in the EventStore when a Reactor has emitted any.
func NewReactorEventHandler(esClient apiEs.EventStoreClient, reactor es.Reactor) *reactorEventHandler {
	return &reactorEventHandler{
		esClient: esClient,
		reactor:  reactor,
	}
}

// HandleEvent implements the HandleEvent method of the es.EventHandler interface.
func (m *reactorEventHandler) HandleEvent(ctx context.Context, event es.Event) error {
	uncomittedEvents, err := m.reactor.HandleEvent(ctx, event)
	if err != nil {
		return err
	}

	// Create stream to send events to store.
	stream, err := m.esClient.Store(ctx)
	if err != nil {
		return err
	}

	for _, event := range uncomittedEvents {
		// Convert to proto event
		protoEvent := es.NewProtoFromEvent(event)

		// Send event to store
		err = stream.Send(protoEvent)
		if err != nil {
			return err
		}
	}
	_, err = stream.CloseAndRecv()

	return err
}
