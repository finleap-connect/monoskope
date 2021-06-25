package eventhandler

import (
	"context"

	"github.com/google/uuid"
	"gitlab.figo.systems/platform/monoskope/monoskope/internal/gateway"
	apiEs "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/eventsourcing"
	es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/logger"
)

type reactorEventHandler struct {
	log      logger.Logger
	esClient apiEs.EventStoreClient
	reactor  es.Reactor
}

// NewReactorEventHandler creates an EventHandler which automates storing Events in the EventStore when a Reactor has emitted any.
func NewReactorEventHandler(esClient apiEs.EventStoreClient, reactor es.Reactor) *reactorEventHandler {
	return &reactorEventHandler{
		log:      logger.WithName("reactorEventHandler"),
		esClient: esClient,
		reactor:  reactor,
	}
}

// HandleEvent implements the HandleEvent method of the es.EventHandler interface.
func (m *reactorEventHandler) HandleEvent(ctx context.Context, event es.Event) error {
	eventsChannel := make(chan es.Event)
	go m.handle(ctx, eventsChannel)
	return m.reactor.HandleEvent(ctx, event, eventsChannel)
}

func (m *reactorEventHandler) handle(ctx context.Context, events <-chan es.Event) {
	for ev := range events { // Read events from channel
		for {
			err := checkUserId(ev)
			if err != nil {
				m.log.Error(err, "Event metadata do not contain user information.")
				break
			}
			err = m.storeEvent(ctx, ev)
			if err != nil {
				m.log.Error(err, "Failed to send event to EventStore. Retrying...")
			} else {
				break
			}
		}
	}
}

func (m *reactorEventHandler) storeEvent(ctx context.Context, event es.Event) error {
	// Create stream to send events to store.
	stream, err := m.esClient.Store(ctx)
	if err != nil {
		m.log.Error(err, "Failed to connect to EventStore.")
		return err
	}

	// Convert to proto event
	protoEvent := es.NewProtoFromEvent(event)

	// Send event to store
	err = stream.Send(protoEvent)
	if err != nil {
		m.log.Error(err, "Failed to send event.")
		return err
	}

	// Close connection
	_, err = stream.CloseAndRecv()
	if err != nil {
		m.log.Error(err, "Failed to close connection with EventStore.")
	}

	return nil
}

func checkUserId(event es.Event) error {
	_, err := uuid.Parse(event.Metadata()[gateway.HeaderAuthId])
	return err
}
