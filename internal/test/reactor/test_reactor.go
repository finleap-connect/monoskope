package reactor

import (
	"context"

	"gitlab.figo.systems/platform/monoskope/monoskope/internal/messagebus"
	esApi "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/eventsourcing"
	es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing"
	ggrpc "google.golang.org/grpc"
)

const (
	reactorName    = "test-reactor"
	msgbusPrefix   = "m8"
	eventStoreAddr = ":8001"
)

type testReactor struct {
	esConnection *ggrpc.ClientConn
	esClient     esApi.EventStoreClient
	msgBus       es.EventBusConsumer

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

func (r *testReactor) Setup(ctx context.Context, esConn *ggrpc.ClientConn, esClient esApi.EventStoreClient) error {
	var err error

	// Create EventStore client
	r.esConnection = esConn
	r.esClient = esClient

	r.msgBus, err = messagebus.NewEventBusConsumer(reactorName, msgbusPrefix)
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
	r.esConnection.Close()
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
