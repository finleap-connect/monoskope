// Copyright 2021 Monoskope Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package eventhandler

import (
	"context"
	"io"
	"sync"
	"time"

	apiEs "github.com/finleap-connect/monoskope/pkg/api/eventsourcing"
	es "github.com/finleap-connect/monoskope/pkg/eventsourcing"
	"github.com/finleap-connect/monoskope/pkg/logger"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type eventStoreRefreshEventHandler struct {
	log             logger.Logger
	esClient        apiEs.EventStoreClient
	handler         es.EventHandler
	mutex           sync.Mutex
	ticker          *time.Ticker
	refreshInterval time.Duration
	aggregateType   es.AggregateType
	lastTimestamp   time.Time
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
		m.lastTimestamp = event.Timestamp()
		if m.aggregateType == "" {
			m.aggregateType = event.AggregateType()
		}

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
		MinTimestamp:  timestamppb.New(m.lastTimestamp),
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

		m.log.V(logger.DebugLevel).Info("Applying event which wasn't received via bus from store.", "event", event.String())

		// Let the next handler in the chain handle the event
		err = m.handler.HandleEvent(ctx, event)
		if err != nil {
			return err
		}
		m.lastTimestamp = event.Timestamp()
	}

	return nil
}
