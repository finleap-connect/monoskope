package reactor

import (
	"context"
	"time"

	"github.com/cenkalti/backoff"
	"github.com/google/uuid"
	"gitlab.figo.systems/platform/monoskope/monoskope/internal/eventstore"
	"gitlab.figo.systems/platform/monoskope/monoskope/internal/gateway/auth"
	"gitlab.figo.systems/platform/monoskope/monoskope/internal/messagebus"
	esApi "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/eventsourcing"
	es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing"
	esMessaging "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing/messaging"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/logger"
)

const (
	reactorName  = "test-reactor"
	msgbusPrefix = "m8"
)

type testReactor struct {
	log      logger.Logger
	msgBus   es.EventBusConsumer
	esClient esApi.EventStoreClient

	observedEvents []es.Event
}

func NewTestReactor() (*testReactor, error) {
	r := &testReactor{
		log: logger.WithName("reactorEventHandler"),
	}

	err := r.init()
	if err != nil {
		return nil, err
	}

	return r, nil

}

func (*testReactor) init() error {
	return nil
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
	if err := r.msgBus.AddHandler(ctx, r, r.msgBus.Matcher().Any()); err != nil {
		return err
	}

	return nil

}

func (r *testReactor) Close() {
	r.msgBus.Close()
}

// HandleEvent handles a given event returns 0..* Events in reaction or an error
func (r *testReactor) HandleEvent(ctx context.Context, event es.Event) error {

	r.observedEvents = append(r.observedEvents, event)
	return nil
}

func (r *testReactor) GetObservedEvents() []es.Event {
	return r.observedEvents
}

func (r *testReactor) Emit(ctx context.Context, event es.Event) error {
	err := checkUserId(event)
	if err != nil {
		r.log.Error(err, "Event metadata do not contain user information.")
		return err
	}

	params := backoff.NewExponentialBackOff()
	params.MaxElapsedTime = 60 * time.Second
	err = backoff.Retry(func() error {
		if err := r.storeEvent(ctx, event); err != nil {
			r.log.Error(err, "Failed to send event to EventStore. Retrying...", "AggregateID", event.AggregateID(), "AggregateType", event.AggregateType(), "EventType", event.EventType())
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

func checkUserId(event es.Event) error {
	_, err := uuid.Parse(event.Metadata()[auth.HeaderAuthId])
	return err
}
