package eventhandler

import (
	"context"
	"io"
	"sync"
	"time"

	apiEs "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/eventsourcing"
	es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/logger"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type eventStoreRefreshEventHandler struct {
	log           logger.Logger
	esClient      apiEs.EventStoreClient
	handler       es.EventHandler
	lastTimestamp time.Time
	mutex         sync.Mutex
	ticker        *time.Ticker
}

// NewEventStoreRefreshEventHandler creates an EventHandler which automates periodic querying of the EventStore to keep up-to-date.
func NewEventStoreRefreshEventHandler(esClient apiEs.EventStoreClient) *eventStoreRefreshEventHandler {
	return &eventStoreRefreshEventHandler{
		esClient: esClient,
	}
}

func (m *eventStoreRefreshEventHandler) AsMiddleware(h es.EventHandler) es.EventHandler {
	return &eventStoreRefreshEventHandler{
		log:      logger.WithName("eventStoreRefreshEventHandler"),
		esClient: m.esClient,
		handler:  h,
	}
}

// HandleEvent implements the HandleEvent method of the es.EventHandler interface.
func (m *eventStoreRefreshEventHandler) HandleEvent(ctx context.Context, event es.Event) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	err := m.handler.HandleEvent(ctx, event)
	if err == nil {
		m.lastTimestamp = event.Timestamp()
		m.resetTicker(ctx)
	}
	return err
}

// resetTicker starts a new ticker if not existing or resets the timer for the existing ticker.
func (m *eventStoreRefreshEventHandler) resetTicker(ctx context.Context) {
	if m.ticker == nil {
		m.ticker = time.NewTicker(time.Minute * 1)
		go func() {
			for range m.ticker.C {
				err := m.applyEventsFromStore(ctx)
				if err != nil {
					m.log.Error(err, "Failed to apply event from store.")
				}
			}
		}()
	}
	m.ticker.Reset(time.Minute * 1)
}

func (m *eventStoreRefreshEventHandler) applyEventsFromStore(ctx context.Context) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Retrieve events from store
	eventStream, err := m.esClient.Retrieve(ctx, &apiEs.EventFilter{
		MinTimestamp: timestamppb.New(m.lastTimestamp),
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

		m.log.Info("Applying event which got lost from store.", "event", esEvent.String())

		// Let the next handler in the chain handle the event
		err = m.handler.HandleEvent(ctx, esEvent)
		if err != nil {
			return err
		}
		m.lastTimestamp = esEvent.Timestamp()
	}

	return nil
}
