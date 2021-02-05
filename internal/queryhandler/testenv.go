package queryhandler

import (
	"context"
	"net"

	"gitlab.figo.systems/platform/monoskope/monoskope/internal/eventstore"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/repositories"
	es_repos "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing/repositories"

	api "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/queryhandler"
	ggrpc "google.golang.org/grpc"

	"gitlab.figo.systems/platform/monoskope/monoskope/internal/test"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/grpc"
)

type TestEnv struct {
	*test.TestEnv
	apiListener       net.Listener
	grpcServer        *grpc.Server
	eventStoreTestEnv *eventstore.EventStoreTestEnv
}

func NewTestEnv() (*TestEnv, error) {
	var err error
	env := &TestEnv{
		TestEnv: test.NewTestEnv("QueryHandlerTestEnv"),
	}

	env.eventStoreTestEnv, err = eventstore.NewEventStoreTestEnv()
	if err != nil {
		return nil, err
	}

	esClient, err := env.eventStoreTestEnv.GetApiClient(context.Background())
	if err != nil {
		return nil, err
	}

	inMemoryUserRoleBindingRepo := es_repos.NewInMemoryRepository()
	userRoleBindingRepo := repositories.NewUserRoleBindingRepository(inMemoryUserRoleBindingRepo)

	inMemoryUserRepo := es_repos.NewInMemoryRepository()
	userRepo := repositories.NewUserRepository(inMemoryUserRepo, userRoleBindingRepo)

	// Create server
	env.grpcServer = grpc.NewServer("query_handler_grpc", false)

	env.grpcServer.RegisterService(func(s ggrpc.ServiceRegistrar) {
		api.RegisterUserServiceServer(s, NewUserServiceServer(esClient, userRepo))
		api.RegisterTenantServiceServer(s, NewTenantServiceServer(esClient))
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

func (env *TestEnv) GetApiAddr() string {
	return env.apiListener.Addr().String()
}

func (env *TestEnv) Shutdown() error {
	if err := env.eventStoreTestEnv.Shutdown(); err != nil {
		return err
	}
	// Shutdown server
	env.grpcServer.Shutdown()
	if err := env.apiListener.Close(); err != nil {
		return err
	}
	return env.TestEnv.Shutdown()
}
