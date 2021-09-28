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

package reactor

import (
	"context"
	"time"

	"github.com/cenkalti/backoff"
	"github.com/finleap-connect/monoskope/internal/eventstore"
	"github.com/finleap-connect/monoskope/internal/messagebus"
	esApi "github.com/finleap-connect/monoskope/pkg/api/eventsourcing"
	es "github.com/finleap-connect/monoskope/pkg/eventsourcing"
	esMessaging "github.com/finleap-connect/monoskope/pkg/eventsourcing/messaging"
	"github.com/finleap-connect/monoskope/pkg/logger"
)

type testReactor struct {
	log      logger.Logger
	msgBus   es.EventBusConsumer
	esClient esApi.EventStoreClient

	observedEvents []es.Event
}

func NewTestReactor() *testReactor {
	r := &testReactor{
		log: logger.WithName("reactorEventHandler"),
	}
	return r
}

func (r *testReactor) Setup(ctx context.Context, env *eventstore.TestEnv, client esApi.EventStoreClient) error {
	var err error

	r.esClient = client

	rabbitConf, err := esMessaging.NewRabbitEventBusConfig("queryhandler", env.GetMessagingTestEnv().AmqpURL, "")
	if err != nil {
		return err
	}

	r.msgBus, err = messagebus.NewEventBusConsumerFromConfig(rabbitConf)
	if err != nil {
		return err
	}

	// Register event handler with event bus to consume all events
	return r.msgBus.AddHandler(ctx, r, r.msgBus.Matcher().Any())
}

func (r *testReactor) Close() {
	r.msgBus.Close()
}

// HandleEvent handles a given event returns 0..* Events in reaction or an error
func (r *testReactor) HandleEvent(ctx context.Context, event es.Event) error {
	r.log.Info("Event observed.", "Event", event.String())
	r.observedEvents = append(r.observedEvents, event)
	return nil
}

func (r *testReactor) GetObservedEvents() []es.Event {
	return r.observedEvents
}

func (r *testReactor) Emit(ctx context.Context, event es.Event) error {
	params := backoff.NewExponentialBackOff()
	params.MaxElapsedTime = 60 * time.Second
	err := backoff.Retry(func() error {
		if err := r.storeEvent(ctx, event); err != nil {
			r.log.Error(err, "Failed to send event to EventStore. Retrying...", "Event", event.String())
			return err
		}
		return nil
	}, params)
	if err != nil {
		r.log.Error(err, "Failed to send event to EventStore")
	}

	return nil
}

func (r *testReactor) storeEvent(ctx context.Context, event es.Event) error {
	// Create stream to send events to store.
	stream, err := r.esClient.Store(ctx)
	if err != nil {
		r.log.Error(err, "Failed to connect to EventStore.")
		return err
	}

	// Convert to proto event
	protoEvent := es.NewProtoFromEvent(event)

	// Send event to store
	err = stream.Send(protoEvent)
	if err != nil {
		r.log.Error(err, "Failed to send event.")
		return err
	}

	// Close connection
	_, err = stream.CloseAndRecv()
	if err != nil {
		r.log.Error(err, "Failed to close connection with EventStore.")
	}

	return nil
}
