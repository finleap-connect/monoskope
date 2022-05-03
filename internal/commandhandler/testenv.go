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

	"github.com/finleap-connect/monoskope/internal/eventstore"
	"github.com/finleap-connect/monoskope/internal/gateway"
	esApi "github.com/finleap-connect/monoskope/pkg/api/eventsourcing"
	gwApi "github.com/finleap-connect/monoskope/pkg/api/gateway"
	"github.com/finleap-connect/monoskope/pkg/domain"
	es "github.com/finleap-connect/monoskope/pkg/eventsourcing"
	"github.com/finleap-connect/monoskope/pkg/grpc/middleware/auth"
	ggrpc "google.golang.org/grpc"

	"github.com/finleap-connect/monoskope/internal/test"
	"github.com/finleap-connect/monoskope/pkg/grpc"
)

type TestEnv struct {
	*test.TestEnv
	apiListener        net.Listener
	grpcServer         *grpc.Server
	eventStoreTestEnv  *eventstore.TestEnv
	gatewayTestEnv     *gateway.TestEnv
	esConn             *ggrpc.ClientConn
	esClient           esApi.EventStoreClient
	gatewayServiceConn *ggrpc.ClientConn
	gatewaySvcClient   gwApi.GatewayAuthClient
}

func NewTestEnv(eventStoreTestEnv *eventstore.TestEnv, gatewayTestEnv *gateway.TestEnv) (*TestEnv, error) {
	var err error
	ctx := context.Background()

	env := &TestEnv{
		TestEnv:           test.NewTestEnv("CommandHandlerTestEnv"),
		eventStoreTestEnv: eventStoreTestEnv,
		gatewayTestEnv:    gatewayTestEnv,
	}

	env.esConn, env.esClient, err = eventstore.NewEventStoreClient(ctx, env.eventStoreTestEnv.GetApiAddr())
	if err != nil {
		return nil, err
	}

	env.gatewayServiceConn, env.gatewaySvcClient, err = gateway.NewInsecureAuthServerClient(ctx, env.gatewayTestEnv.GetApiAddr())
	if err != nil {
		return nil, err
	}

	authMiddleware := auth.NewAuthMiddleware(env.gatewaySvcClient, []string{"/grpc.health.v1.Health/Check"})

	err = domain.SetupCommandHandlerDomain(ctx, env.esClient)
	if err != nil {
		return nil, err
	}

	// Create server
	env.grpcServer = grpc.NewServerWithOpts("commandhandler-grpc", false,
		[]ggrpc.UnaryServerInterceptor{
			authMiddleware.UnaryServerInterceptor(),
		}, []ggrpc.StreamServerInterceptor{
			authMiddleware.StreamServerInterceptor(),
		},
	)

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

	if err := env.gatewayServiceConn.Close(); err != nil {
		return err
	}

	// Shutdown server
	env.grpcServer.Shutdown()
	if err := env.apiListener.Close(); err != nil {
		return err
	}
	return env.TestEnv.Shutdown()
}
