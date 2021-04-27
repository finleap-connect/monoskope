package queryhandler

import (
	"context"
	"net"

	"gitlab.figo.systems/platform/monoskope/monoskope/internal/eventstore"
	"gitlab.figo.systems/platform/monoskope/monoskope/internal/messagebus"
	esApi "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/eventsourcing"
	es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing"

	api "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/domain"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain"
	esMessaging "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing/messaging"
	ggrpc "google.golang.org/grpc"

	"gitlab.figo.systems/platform/monoskope/monoskope/internal/test"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/grpc"
)

type TestEnv struct {
	*test.TestEnv
	apiListener       net.Listener
	grpcServer        *grpc.Server
	eventStoreTestEnv *eventstore.TestEnv
	ebConsumer        es.EventBusConsumer
	esConn            *ggrpc.ClientConn
	esClient          esApi.EventStoreClient
}

func NewTestEnvWithParent(testeEnv *test.TestEnv, eventStoreTestEnv *eventstore.TestEnv) (*TestEnv, error) {
	var err error
	ctx := context.Background()

	env := &TestEnv{
		TestEnv:           testeEnv,
		eventStoreTestEnv: eventStoreTestEnv,
	}

	rabbitConf := esMessaging.NewRabbitEventBusConfig("queryhandler", env.eventStoreTestEnv.GetMessagingTestEnv().AmqpURL, "")
	env.ebConsumer, err = messagebus.NewEventBusConsumerFromConfig(rabbitConf)
	if err != nil {
		return nil, err
	}

	env.esConn, env.esClient, err = eventstore.NewEventStoreClient(ctx, env.eventStoreTestEnv.GetApiAddr())
	if err != nil {
		return nil, err
	}

	// Setup domain
	qhDomain, err := domain.NewQueryHandlerDomain(context.Background(), env.ebConsumer, env.esClient)
	if err != nil {
		return nil, err
	}

	// Create server
	grpcServer := grpc.NewServer("queryhandler_grpc", false)
	env.grpcServer = grpcServer
	grpcServer.RegisterService(func(s ggrpc.ServiceRegistrar) {
		api.RegisterUserServer(s, NewUserServer(qhDomain.UserRepository))
		api.RegisterTenantServer(s, NewTenantServer(qhDomain.TenantRepository))
	})

	env.apiListener, err = net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return nil, err
	}

	// Start server
	go func() {
		err := grpcServer.ServeFromListener(env.apiListener, nil)
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
	if err := env.esConn.Close(); err != nil {
		return err
	}

	if err := env.ebConsumer.Close(); err != nil {
		return err
	}

	// Shutdown server
	env.grpcServer.Shutdown()
	if err := env.apiListener.Close(); err != nil {
		return err
	}

	return nil
}
