package eventstore

import (
	"context"
	"net"

	api "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/eventsourcing"
	ggrpc "google.golang.org/grpc"

	"gitlab.figo.systems/platform/monoskope/monoskope/internal/test"
	es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing/messaging"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing/storage"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/grpc"
)

type TestEnv struct {
	*test.TestEnv
	apiListener      net.Listener
	MetricsListener  net.Listener
	grpcServer       *grpc.Server
	messagingTestEnv *messaging.TestEnv
	storageTestEnv   *storage.TestEnv
	publisher        es.EventBusPublisher
}

func (t *TestEnv) GetMessagingTestEnv() *messaging.TestEnv {
	return t.messagingTestEnv
}

func NewTestEnvWithParent(testEnv *test.TestEnv) (*TestEnv, error) {
	var err error

	env := &TestEnv{
		TestEnv: testEnv,
	}

	env.messagingTestEnv, err = messaging.NewTestEnvWithParent(testEnv)
	if err != nil {
		return nil, err
	}

	conf := messaging.NewRabbitEventBusConfig("eventstore", env.messagingTestEnv.AmqpURL, "")
	env.publisher, err = messaging.NewRabbitEventBusPublisher(conf)
	if err != nil {
		return nil, err
	}

	err = env.publisher.Open(context.Background())
	if err != nil {
		return nil, err
	}

	env.storageTestEnv, err = storage.NewTestEnvWithParent(testEnv)
	if err != nil {
		return nil, err
	}

	err = env.storageTestEnv.Store.Connect(context.Background())
	if err != nil {
		return nil, err
	}

	// Create server
	env.grpcServer = grpc.NewServer("eventstore_grpc", false)
	env.grpcServer.RegisterService(func(s ggrpc.ServiceRegistrar) {
		api.RegisterEventStoreServer(s, NewApiServer(env.storageTestEnv.Store, env.publisher))
	})

	env.apiListener, err = net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return nil, err
	}

	env.MetricsListener, err = net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return nil, err
	}

	// Start server
	go func() {
		err := env.grpcServer.ServeFromListener(env.apiListener, env.MetricsListener)
		if err != nil {
			panic(err)
		}
	}()

	return env, nil
}

func (env *TestEnv) GetApiAddr() string {
	return env.apiListener.Addr().String()
}

func (env *TestEnv) Shutdown() error {
	if err := env.publisher.Close(); err != nil {
		return err
	}

	if err := env.messagingTestEnv.Shutdown(); err != nil {
		return err
	}

	if err := env.storageTestEnv.Shutdown(); err != nil {
		return err
	}

	// Shutdown server
	env.grpcServer.Shutdown()
	if err := env.apiListener.Close(); err != nil {
		return err
	}
	return nil
}
