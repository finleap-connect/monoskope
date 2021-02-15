package queryhandler

import (
	"context"
	"net"

	"gitlab.figo.systems/platform/monoskope/monoskope/internal/eventstore"
	"gitlab.figo.systems/platform/monoskope/monoskope/internal/util"
	esApi "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/eventsourcing"
	es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing"

	api "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/domain"
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

func NewTestEnv(eventStoreTestEnv *eventstore.TestEnv) (*TestEnv, error) {
	var err error
	env := &TestEnv{
		TestEnv:           test.NewTestEnv("QueryHandlerTestEnv"),
		eventStoreTestEnv: eventStoreTestEnv,
	}

	rabbitConf := esMessaging.NewRabbitEventBusConfig("queryhandler", env.eventStoreTestEnv.GetMessagingTestEnv().AmqpURL, "")
	env.ebConsumer, err = util.NewEventBusConsumerFromConfig(rabbitConf)
	if err != nil {
		return nil, err
	}

	env.esConn, env.esClient, err = util.NewEventStoreClient(env.eventStoreTestEnv.GetApiAddr())
	if err != nil {
		return nil, err
	}

	// Setup domain
	userRepo, err := util.SetupQueryHandlerDomain(context.Background(), env.ebConsumer, env.esClient)
	if err != nil {
		return nil, err
	}

	// Create server
	grpcServer := grpc.NewServer("queryhandler_grpc", false)
	env.grpcServer = grpcServer
	grpcServer.RegisterService(func(s ggrpc.ServiceRegistrar) {
		api.RegisterUserServiceServer(s, NewUserServiceServer(userRepo))
		api.RegisterTenantServiceServer(s, NewTenantServiceServer())
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
	return env.TestEnv.Shutdown()
}
