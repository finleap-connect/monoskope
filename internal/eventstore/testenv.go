package eventstore

import (
	"context"
	"net"

	api "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/eventstore"
	ggrpc "google.golang.org/grpc"

	"gitlab.figo.systems/platform/monoskope/monoskope/internal/test"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing/messaging"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing/storage"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/grpc"
)

type EventStoreTestEnv struct {
	*test.TestEnv
	apiListener net.Listener
	grpcServer  *grpc.Server
}

func NewEventStoreTestEnv() (*EventStoreTestEnv, error) {
	var err error
	env := &EventStoreTestEnv{
		TestEnv: test.NewTestEnv("TestEventStoreEnv"),
	}

	// Create server
	env.grpcServer = grpc.NewServer("event_store_grpc", false)

	eventStore, err := NewApiServer(storage.NewInMemoryEventStore(), messaging.NewMockEventBusPublisher())
	if err != nil {
		return nil, err
	}

	env.grpcServer.RegisterService(func(s ggrpc.ServiceRegistrar) {
		api.RegisterEventStoreServer(s, eventStore)
	})
	env.grpcServer.RegisterOnShutdown(func() {
		eventStore.Shutdown()
	})

	env.apiListener, err = net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return nil, err
	}

	// Start server
	go func() {
		err := env.grpcServer.Serve(env.apiListener, nil)
		if err != nil {
			panic(err)
		}
	}()

	return env, nil
}

func (env *EventStoreTestEnv) GetApiAddr() string {
	return env.apiListener.Addr().String()
}

func (env *EventStoreTestEnv) GetApiClient(ctx context.Context) (api.EventStoreClient, error) {
	conn, err := grpc.
		NewGrpcConnectionFactory(env.GetApiAddr()).
		WithInsecure().
		WithBlock().
		Build(ctx)
	if err != nil {
		return nil, err
	}
	return api.NewEventStoreClient(conn), nil
}

func (env *EventStoreTestEnv) Shutdown() error {
	// Shutdown server
	env.grpcServer.Shutdown()
	if err := env.apiListener.Close(); err != nil {
		return err
	}
	return env.TestEnv.Shutdown()
}
