package reactor

import (
	"context"

	"gitlab.figo.systems/platform/monoskope/monoskope/internal/eventstore"
	"gitlab.figo.systems/platform/monoskope/monoskope/internal/messagebus"
	es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing"
	esMessaging "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing/messaging"
)

const (
	reactorName  = "test-reactor"
	msgbusPrefix = "m8"
)

type testReactor struct {
	msgBus es.EventBusConsumer

	observedEvents []es.Event
}

func NewTestReactor() (*testReactor, error) {
	r := &testReactor{}

	err := r.init()
	if err != nil {
		return nil, err
	}

	return r, nil

}

func (*testReactor) init() error {
	return nil
}

func (r *testReactor) Setup(ctx context.Context, env *eventstore.TestEnv) error {
	var err error

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

func (r *testReactor) Emit(event es.Event) error {
	return nil
}
