package commandhandler

import (
	"context"
	"net"

	"gitlab.figo.systems/platform/monoskope/monoskope/internal/eventstore"
	api "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/commandhandler"
	ggrpc "google.golang.org/grpc"

	"gitlab.figo.systems/platform/monoskope/monoskope/internal/test"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/grpc"
)

type CommandHandlerTestEnv struct {
	*test.TestEnv
	apiListener       net.Listener
	grpcServer        *grpc.Server
	eventStoreTestEnv *eventstore.EventStoreTestEnv
}

func NewCommandHandlerTestEnv() (*CommandHandlerTestEnv, error) {
	var err error
	env := &CommandHandlerTestEnv{
		TestEnv: test.NewTestEnv("CommandHandlerTestEnv"),
	}

	env.eventStoreTestEnv, err = eventstore.NewEventStoreTestEnv()
	if err != nil {
		return nil, err
	}

	esClient, err := env.eventStoreTestEnv.GetApiClient(context.Background())
	if err != nil {
		return nil, err
	}

	// Create server
	env.grpcServer = grpc.NewServer("command_handler_grpc", false)

	commandHandler := NewApiServer(esClient)
	env.grpcServer.RegisterService(func(s ggrpc.ServiceRegistrar) {
		api.RegisterCommandHandlerServer(s, commandHandler)
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

func (env *CommandHandlerTestEnv) GetApiAddr() string {
	return env.apiListener.Addr().String()
}

func (env *CommandHandlerTestEnv) Shutdown() error {
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
