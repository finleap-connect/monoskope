package commandhandler

import (
	"context"
	"net"

	"gitlab.figo.systems/platform/monoskope/monoskope/internal/eventstore"
	"gitlab.figo.systems/platform/monoskope/monoskope/internal/queryhandler"
	domainApi "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/domain"
	esApi "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/eventsourcing"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain"
	ggrpc "google.golang.org/grpc"

	"gitlab.figo.systems/platform/monoskope/monoskope/internal/test"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/grpc"
)

type TestEnv struct {
	*test.TestEnv
	apiListener         net.Listener
	grpcServer          *grpc.Server
	eventStoreTestEnv   *eventstore.TestEnv
	queryHandlerTestEnv *queryhandler.TestEnv
	esConn              *ggrpc.ClientConn
	esClient            esApi.EventStoreClient
	userServiceConn     *ggrpc.ClientConn
	userSvcClient       domainApi.UserServiceClient
	tenantServiceConn   *ggrpc.ClientConn
	tenantSvcClient     domainApi.TenantServiceClient
}

func NewTestEnv(eventStoreTestEnv *eventstore.TestEnv, queryHandlerTestEnv *queryhandler.TestEnv) (*TestEnv, error) {
	var err error
	ctx := context.Background()

	env := &TestEnv{
		TestEnv:             test.NewTestEnv("CommandHandlerTestEnv"),
		eventStoreTestEnv:   eventStoreTestEnv,
		queryHandlerTestEnv: queryHandlerTestEnv,
	}

	env.esConn, env.esClient, err = eventstore.NewEventStoreClient(ctx, env.eventStoreTestEnv.GetApiAddr())
	if err != nil {
		return nil, err
	}

	env.userServiceConn, env.userSvcClient, err = queryhandler.NewUserServiceClient(ctx, env.queryHandlerTestEnv.GetApiAddr())
	if err != nil {
		return nil, err
	}

	env.tenantServiceConn, env.tenantSvcClient, err = queryhandler.NewTenantServiceClient(ctx, env.queryHandlerTestEnv.GetApiAddr())
	if err != nil {
		return nil, err
	}

	cmdRegistry, err := domain.SetupCommandHandlerDomain(ctx, env.userSvcClient, env.tenantSvcClient, env.esClient)
	if err != nil {
		return nil, err
	}

	// Create server
	env.grpcServer = grpc.NewServer("commandhandler_grpc", false)

	commandHandler := NewApiServer(cmdRegistry)
	env.grpcServer.RegisterService(func(s ggrpc.ServiceRegistrar) {
		esApi.RegisterCommandHandlerServer(s, commandHandler)
	})

	env.apiListener, err = net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return nil, err
	}

	// Start server
	go func() {
		err := env.grpcServer.ServeFromListener(env.apiListener, nil)
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

	if err := env.userServiceConn.Close(); err != nil {
		return err
	}

	if err := env.tenantServiceConn.Close(); err != nil {
		return err
	}

	// Shutdown server
	env.grpcServer.Shutdown()
	if err := env.apiListener.Close(); err != nil {
		return err
	}
	return env.TestEnv.Shutdown()
}
