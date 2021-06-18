package eventhandler

import (
	"context"
	"sync"
	"time"

	"github.com/cenkalti/backoff"
	apiEs "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/eventsourcing"
	es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/logger"
)

type reactorEventHandler struct {
	log       logger.Logger
	esClient  apiEs.EventStoreClient
	reactor   es.Reactor
	waitGroup sync.WaitGroup
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
	m.waitGroup.Add(1)
	eventsChannel := make(chan es.Event)
	go m.handle(ctx, eventsChannel)
	return m.reactor.HandleEvent(ctx, event, eventsChannel)
}

// Stop waits for all goroutines to finish
func (m *reactorEventHandler) Stop() {
	m.waitGroup.Wait()
}

func (m *reactorEventHandler) handle(ctx context.Context, events <-chan es.Event) {
	defer m.waitGroup.Done()
	for ev := range events { // Read events from channel
		params := backoff.NewExponentialBackOff()
		params.MaxElapsedTime = 60 * time.Second
		err := backoff.Retry(func() error {
			if err := m.storeEvent(ctx, ev); err != nil {
				m.log.Error(err, "Failed to send event to EventStore. Retrying...", "AggregateID", ev.AggregateID(), "AggregateType", ev.AggregateType(), "EventType", ev.EventType())
				return err
			}
			return nil
		}, params)

		if err != nil {
			m.log.Error(err, "Failed to send event to EventStore")
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
