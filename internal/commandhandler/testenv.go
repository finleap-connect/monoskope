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

package commandhandler

import (
	"context"
	"net"

	"gitlab.figo.systems/platform/monoskope/monoskope/internal/eventstore"
	"gitlab.figo.systems/platform/monoskope/monoskope/internal/queryhandler"
	domainApi "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/domain"
	esApi "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/eventsourcing"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain"
	es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing"
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
	userSvcClient       domainApi.UserClient
	tenantServiceConn   *ggrpc.ClientConn
	tenantSvcClient     domainApi.TenantClient
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

	env.userServiceConn, env.userSvcClient, err = queryhandler.NewUserClient(ctx, env.queryHandlerTestEnv.GetApiAddr())
	if err != nil {
		return nil, err
	}

	env.tenantServiceConn, env.tenantSvcClient, err = queryhandler.NewTenantClient(ctx, env.queryHandlerTestEnv.GetApiAddr())
	if err != nil {
		return nil, err
	}

	err = domain.SetupCommandHandlerDomain(ctx, env.userSvcClient, env.esClient)
	if err != nil {
		return nil, err
	}

	// Create server
	env.grpcServer = grpc.NewServer("commandhandler_grpc", false)

	commandHandler := NewApiServer(es.DefaultCommandRegistry)
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
