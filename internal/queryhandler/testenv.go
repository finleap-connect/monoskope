// Copyright 2022 Monoskope Authors
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

package queryhandler

import (
	"context"
	"net"

	ef "github.com/finleap-connect/monoskope/pkg/audit/formatters/event"
	"github.com/finleap-connect/monoskope/pkg/grpc/middleware/auth"

	"github.com/finleap-connect/monoskope/internal/eventstore"
	"github.com/finleap-connect/monoskope/internal/gateway"
	"github.com/finleap-connect/monoskope/internal/messagebus"
	esApi "github.com/finleap-connect/monoskope/pkg/api/eventsourcing"
	es "github.com/finleap-connect/monoskope/pkg/eventsourcing"

	api "github.com/finleap-connect/monoskope/pkg/api/domain"
	gwApi "github.com/finleap-connect/monoskope/pkg/api/gateway"
	"github.com/finleap-connect/monoskope/pkg/domain"
	esMessaging "github.com/finleap-connect/monoskope/pkg/eventsourcing/messaging"
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
	ebConsumer         es.EventBusConsumer
	esConn             *ggrpc.ClientConn
	esClient           esApi.EventStoreClient
	gatewayServiceConn *ggrpc.ClientConn
	gatewaySvcClient   gwApi.GatewayAuthClient
}

func NewTestEnvWithParent(testeEnv *test.TestEnv, eventStoreTestEnv *eventstore.TestEnv, gatewayTestEnv *gateway.TestEnv) (*TestEnv, error) {
	var err error
	ctx := context.Background()

	env := &TestEnv{
		TestEnv:           testeEnv,
		eventStoreTestEnv: eventStoreTestEnv,
		gatewayTestEnv:    gatewayTestEnv,
	}

	rabbitConf, err := esMessaging.NewRabbitEventBusConfig("queryhandler", env.eventStoreTestEnv.GetMessagingTestEnv().AmqpURL, "")
	if err != nil {
		return nil, err
	}

	env.ebConsumer, err = messagebus.NewEventBusConsumerFromConfig(rabbitConf)
	if err != nil {
		return nil, err
	}

	env.esConn, env.esClient, err = eventstore.NewEventStoreClient(ctx, env.eventStoreTestEnv.GetApiAddr())
	if err != nil {
		return nil, err
	}

	env.gatewayServiceConn, env.gatewaySvcClient, err = gateway.NewInsecureAuthServerClient(ctx, env.gatewayTestEnv.GetApiAddr())
	if err != nil {
		return nil, err
	}

	// Setup domain
	qhDomain, err := domain.NewQueryHandlerDomain(ctx, env.ebConsumer, env.esClient)
	if err != nil {
		return nil, err
	}

	authMiddleware := auth.NewAuthMiddleware(env.gatewaySvcClient, []string{"/grpc.health.v1.Health/Check"})

	// Create server
	env.grpcServer = grpc.NewServerWithOpts("queryhandler_grpc-grpc", false,
		[]ggrpc.UnaryServerInterceptor{
			authMiddleware.UnaryServerInterceptor(),
		}, []ggrpc.StreamServerInterceptor{
			authMiddleware.StreamServerInterceptor(),
		},
	)

	env.grpcServer.RegisterService(func(s ggrpc.ServiceRegistrar) {
		api.RegisterUserServer(s, NewUserServer(qhDomain.UserRepository))
		api.RegisterTenantServer(s, NewTenantServer(qhDomain.TenantRepository, qhDomain.TenantUserRepository))
		api.RegisterClusterServer(s, NewClusterServer(qhDomain.ClusterRepository))
		api.RegisterClusterAccessServer(s, NewClusterAccessServer(qhDomain.ClusterAccessRepo, qhDomain.TenantClusterBindingRepository))
		api.RegisterAuditLogServer(s, NewAuditLogServer(env.esClient, ef.DefaultEventFormatterRegistry, qhDomain.UserRepository))
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
