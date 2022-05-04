// Copyright 2022 Monoskope Authors
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
	"errors"
	"io"
	"sync"

	apiEs "github.com/finleap-connect/monoskope/pkg/api/eventsourcing"
	es "github.com/finleap-connect/monoskope/pkg/eventsourcing"
	"github.com/finleap-connect/monoskope/pkg/logger"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type eventStoreReplayEventHandler struct {
	log      logger.Logger
	esClient apiEs.EventStoreClient
	handler  es.EventHandler
	mutex    sync.Mutex
}

// NewEventStoreReplayMiddleware creates an EventHandler which automates querying the EventStore in case of gaps in AggregateVersion found in other EventHandlers later in the chain of EventHandlers.
func NewEventStoreReplayMiddleware(esClient apiEs.EventStoreClient) es.EventHandlerMiddleware {
	return func(h es.EventHandler) es.EventHandler {
		return &eventStoreReplayEventHandler{
			log:      logger.WithName("replay-middleware"),
			esClient: esClient,
			handler:  h,
		}
	}
}

// HandleEvent implements the HandleEvent method of the es.EventHandler interface.
func (m *eventStoreReplayEventHandler) HandleEvent(ctx context.Context, event es.Event) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	var outdatedError *ProjectionOutdatedError
	if err := m.handler.HandleEvent(ctx, event); errors.As(err, &outdatedError) {
		// If the next handler in the chain tells that the projection is outdated
		m.log.Info("Gap in event stream found. Replaying missing events from store.", "event", event.String())
		if err := m.applyEventsFromStore(ctx, event, outdatedError.ProjectionVersion); err != nil {
			return err
		}
		return nil
	} else {
		return err
	}
}

func (m *eventStoreReplayEventHandler) applyEventsFromStore(ctx context.Context, event es.Event, projectionVersion uint64) error {
	// Retrieve events from store
	eventStream, err := m.esClient.Retrieve(ctx, &apiEs.EventFilter{
		AggregateId:   wrapperspb.String(event.AggregateID().String()),
		AggregateType: wrapperspb.String(event.AggregateType().String()),
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
