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

package scimserver

import (
	"context"
	"net"
	"net/http"
	"os"

	"github.com/elimity-com/scim"
	"github.com/finleap-connect/monoskope/internal/commandhandler"
	"github.com/finleap-connect/monoskope/internal/eventstore"
	"github.com/finleap-connect/monoskope/internal/gateway"
	"github.com/finleap-connect/monoskope/internal/queryhandler"
	esApi "github.com/finleap-connect/monoskope/pkg/api/eventsourcing"

	domainApi "github.com/finleap-connect/monoskope/pkg/api/domain"
	ggrpc "google.golang.org/grpc"

	"github.com/finleap-connect/monoskope/internal/test"
)

type TestEnv struct {
	*test.TestEnv
	apiListener           net.Listener
	eventStoreTestEnv     *eventstore.TestEnv
	commandHandlerTestEnv *commandhandler.TestEnv
	queryHandlerTestEnv   *queryhandler.TestEnv
	gatewayTestEnv        *gateway.TestEnv
	userServiceConn       *ggrpc.ClientConn
	userSvcClient         domainApi.UserClient
	commandHandlerConn    *ggrpc.ClientConn
	commandHandlerClient  esApi.CommandHandlerClient
	scimServer            scim.Server
}

func NewTestEnv(testEnv *test.TestEnv) (*TestEnv, error) {
	var err error
	ctx := context.Background()

	env := &TestEnv{
		TestEnv: testEnv,
	}

	os.Setenv("SUPER_USERS", "")

	env.gatewayTestEnv, err = gateway.NewTestEnvWithParent(testEnv)
	if err != nil {
		return nil, err
	}

	env.eventStoreTestEnv, err = eventstore.NewTestEnvWithParent(testEnv)
	if err != nil {
		return nil, err
	}

	env.queryHandlerTestEnv, err = queryhandler.NewTestEnvWithParent(testEnv, env.eventStoreTestEnv, env.gatewayTestEnv)
	if err != nil {
		return nil, err
	}

	env.commandHandlerTestEnv, err = commandhandler.NewTestEnv(env.eventStoreTestEnv, env.gatewayTestEnv)
	if err != nil {
		return nil, err
	}

	env.userServiceConn, env.userSvcClient, err = queryhandler.NewUserClient(ctx, env.queryHandlerTestEnv.GetApiAddr())
	if err != nil {
		return nil, err
	}

	env.commandHandlerConn, env.commandHandlerClient, err = commandhandler.NewCommandHandlerClientWithAuthForward(ctx, env.commandHandlerTestEnv.GetApiAddr())
	if err != nil {
		return nil, err
	}

	env.apiListener, err = net.Listen("tcp", "127.0.0.1:5080")
	if err != nil {
		return nil, err
	}

	providerConfig := NewProvierConfig()
	userHandler := NewUserHandler(env.commandHandlerClient, env.userSvcClient)
	groupHandler := NewGroupHandler(env.commandHandlerClient, env.userSvcClient)
	env.scimServer = NewServer(providerConfig, userHandler, groupHandler)

	// Start server
	go func() {
		_ = http.Serve(env.apiListener, env.scimServer)
	}()

	return env, nil
}

func (env *TestEnv) GetApiAddr() string {
	return env.apiListener.Addr().String()
}

func (env *TestEnv) Shutdown() error {
	// Shutdown server
	if err := env.apiListener.Close(); err != nil {
		return err
	}

	if err := env.queryHandlerTestEnv.Shutdown(); err != nil {
		return err
	}

	if err := env.commandHandlerTestEnv.Shutdown(); err != nil {
		return err
	}

	if err := env.eventStoreTestEnv.Shutdown(); err != nil {
		return err
	}

	if err := env.gatewayTestEnv.Shutdown(); err != nil {
		return err
	}
	return nil
}
