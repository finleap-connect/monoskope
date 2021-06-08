package eventhandler

import (
	"context"

	apiEs "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/eventsourcing"
	es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/logger"
)

type reactorEventHandler struct {
	log      logger.Logger
	esClient apiEs.EventStoreClient
	reactor  es.Reactor
	events   chan es.Event
}

// NewReactorEventHandler creates an EventHandler which automates storing Events in the EventStore when a Reactor has emitted any.
func NewReactorEventHandler(esClient apiEs.EventStoreClient, reactor es.Reactor) *reactorEventHandler {
	return &reactorEventHandler{
		log:      logger.WithName("reactorEventHandler"),
		esClient: esClient,
		reactor:  reactor,
		events:   make(chan es.Event),
	}
}

// HandleEvent implements the HandleEvent method of the es.EventHandler interface.
func (m *reactorEventHandler) HandleEvent(ctx context.Context, event es.Event) error {
	return m.reactor.HandleEvent(ctx, event, m.events)
}

func (m *reactorEventHandler) handle(ctx context.Context) error {
	var err error
	var stream apiEs.EventStore_StoreClient

	for ev := range m.events { // Read events from channel
		if stream == nil {
			// Create stream to send events to store.
			stream, err = m.esClient.Store(ctx)
			if err != nil {
				m.log.Error(err, "Failed to connect to EventStore. Retrying...")
				// Retry on failure
				continue
			}
		}

		// Convert to proto event
		protoEvent := es.NewProtoFromEvent(ev)

		// Send event to store
		err = stream.Send(protoEvent)
		if err != nil {
			return err
		}
	}

	if stream == nil {
		_, err := stream.CloseAndRecv()
		return err
	}

	return nil
}
