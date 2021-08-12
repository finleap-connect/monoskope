package eventhandler

import (
	"context"
	"io"
	"sync"
	"time"

	apiEs "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/eventsourcing"
	es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/logger"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type eventStoreRefreshEventHandler struct {
	log             logger.Logger
	esClient        apiEs.EventStoreClient
	handler         es.EventHandler
	lastVersion     uint64
	mutex           sync.Mutex
	ticker          *time.Ticker
	refreshInterval time.Duration
	aggregateType   es.AggregateType
}

// NewEventStoreRefreshMiddleware creates an EventHandler which automates periodic querying of the EventStore to keep up-to-date.
func NewEventStoreRefreshMiddleware(esClient apiEs.EventStoreClient, refreshInterval time.Duration) es.EventHandlerMiddleware {
	return func(h es.EventHandler) es.EventHandler {
		return &eventStoreRefreshEventHandler{
			log:             logger.WithName("refresh-middleware"),
			esClient:        esClient,
			refreshInterval: refreshInterval,
			handler:         h,
		}
	}
}

// HandleEvent implements the HandleEvent method of the es.EventHandler interface.
func (m *eventStoreRefreshEventHandler) HandleEvent(ctx context.Context, event es.Event) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	err := m.handler.HandleEvent(ctx, event)
	if err == nil {
		m.aggregateType = event.AggregateType()
		m.lastVersion = event.AggregateVersion()
		m.resetTicker(ctx)
	}
	return err
}

// resetTicker starts a new ticker if not existing or resets the timer for the existing ticker.
func (m *eventStoreRefreshEventHandler) resetTicker(ctx context.Context) {
	if m.ticker == nil {
		m.ticker = time.NewTicker(m.refreshInterval)
		go func() {
			for range m.ticker.C {
				err := m.applyEventsFromStore(ctx)
				if err != nil {
					m.log.Error(err, "Failed to apply event from store.")
				}
			}
		}()
	}
	m.ticker.Reset(m.refreshInterval)
}

func (m *eventStoreRefreshEventHandler) applyEventsFromStore(ctx context.Context) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Retrieve events from store
	eventStream, err := m.esClient.Retrieve(ctx, &apiEs.EventFilter{
		MinVersion:    wrapperspb.UInt64(m.lastVersion + 1),
		AggregateType: wrapperspb.String(m.aggregateType.String()),
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
		event, err := es.NewEventFromProto(protoEvent)
		if err != nil {
			return err
		}

		m.log.Info("Applying event which wasn't received via bus from store.", "event", event.String())

		// Let the next handler in the chain handle the event
		err = m.handler.HandleEvent(ctx, event)
		if err != nil {
			return err
		}
		m.lastVersion = event.AggregateVersion()
	}

	return nil
}
