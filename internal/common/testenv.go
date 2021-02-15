package common

import (
	"net"

	api_common "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/domain/common"
	ggrpc "google.golang.org/grpc"

	"gitlab.figo.systems/platform/monoskope/monoskope/internal/test"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/grpc"
)

type TestEnv struct {
	*test.TestEnv
	apiListener net.Listener
	grpcServer  *grpc.Server
}

func NewCommandHandlerTestEnv() (*TestEnv, error) {
	var err error
	env := &TestEnv{
		TestEnv: test.NewTestEnv("CommandHandlerTestEnv"),
	}

	// Create server
	env.grpcServer = grpc.NewServer("command_handler_grpc", false)

	env.grpcServer.RegisterService(func(s ggrpc.ServiceRegistrar) {
		api_common.RegisterServiceInformationServiceServer(s, NewServiceInformationService())
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
	// Shutdown server
	env.grpcServer.Shutdown()
	if err := env.apiListener.Close(); err != nil {
		return err
	}
	return env.TestEnv.Shutdown()
}
