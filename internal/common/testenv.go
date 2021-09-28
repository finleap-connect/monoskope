// Copyright 2021 Monoskope Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package common

import (
	"net"

	api_common "github.com/finleap-connect/monoskope/pkg/api/domain/common"
	ggrpc "google.golang.org/grpc"

	"github.com/finleap-connect/monoskope/internal/test"
	"github.com/finleap-connect/monoskope/pkg/grpc"
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
