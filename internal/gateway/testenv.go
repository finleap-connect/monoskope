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

package gateway

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/finleap-connect/monoskope/internal/common"
	"github.com/finleap-connect/monoskope/internal/eventstore"
	"github.com/finleap-connect/monoskope/internal/gateway/auth"
	"github.com/finleap-connect/monoskope/internal/messagebus"
	"github.com/finleap-connect/monoskope/internal/test"
	apiCommon "github.com/finleap-connect/monoskope/pkg/api/domain/common"
	esApi "github.com/finleap-connect/monoskope/pkg/api/eventsourcing"
	api "github.com/finleap-connect/monoskope/pkg/api/gateway"
	clientAuth "github.com/finleap-connect/monoskope/pkg/auth"
	"github.com/finleap-connect/monoskope/pkg/domain"
	"github.com/finleap-connect/monoskope/pkg/domain/mock"
	"github.com/finleap-connect/monoskope/pkg/domain/repositories"
	es "github.com/finleap-connect/monoskope/pkg/eventsourcing"
	esMessaging "github.com/finleap-connect/monoskope/pkg/eventsourcing/messaging"
	"github.com/finleap-connect/monoskope/pkg/grpc"
	"github.com/finleap-connect/monoskope/pkg/jwt"
	"github.com/ory/dockertest/v3"
	dc "github.com/ory/dockertest/v3/docker"
	ggrpc "google.golang.org/grpc"
)

const (
	localAddrAPIServer          = "127.0.0.1:9090"
	localAddrOIDCProviderServer = "127.0.0.1:9091"
	RedirectURLHostname         = "localhost"
	RedirectURLPort             = ":8000"
)

type TestEnv struct {
	*test.TestEnv
	JwtTestEnv                    *jwt.TestEnv
	ClientAuthConfig              *auth.ClientConfig
	ServerAuthConfig              *auth.ServerConfig
	IdentityProviderURL           string
	ApiListenerAPIServer          net.Listener
	ApiListenerOIDCProviderServer net.Listener
	HttpClient                    *http.Client
	GrpcServer                    *grpc.Server
	LocalOIDCProviderServer       *oidcProviderServer
	PoliciesPath                  string
	eventStoreTestEnv             *eventstore.TestEnv
	ebConsumer                    es.EventBusConsumer
	esConn                        *ggrpc.ClientConn
	esClient                      esApi.EventStoreClient
}

func NewTestEnvWithParent(testeEnv *test.TestEnv, eventStoreTestEnv *eventstore.TestEnv, mockData bool) (*TestEnv, error) {
	ctx := context.Background()

	env := &TestEnv{
		TestEnv:           testeEnv,
		eventStoreTestEnv: eventStoreTestEnv,
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

	jwtTestEnv, err := jwt.NewTestEnv(env.TestEnv)
	if err != nil {
		_ = env.Shutdown()
		return nil, err
	}
	env.JwtTestEnv = jwtTestEnv

	err = env.CreateDockerPool(false)
	if err != nil {
		_ = env.Shutdown()
		return nil, err
	}

	dexConfigDir := os.Getenv("DEX_CONFIG")
	if dexConfigDir == "" {
		return nil, fmt.Errorf("DEX_CONFIG not specified")
	}
	env.Log.Info("Config for dex specified.", "DEX_CONFIG", dexConfigDir)

	dexContainer, err := env.Run(&dockertest.RunOptions{
		Name:       "dex",
		Repository: "dexidp/dex",
		Tag:        "v2.27.0",
		PortBindings: map[dc.Port][]dc.PortBinding{
			"5556": {{HostPort: "5556"}},
		},
		ExposedPorts: []string{"5556", "5000"},
		Cmd:          []string{"serve", "/etc/dex/cfg/config.yaml"},
		Mounts:       []string{fmt.Sprintf("%s:/etc/dex/cfg", dexConfigDir)},
	})

	if err != nil {
		_ = env.Shutdown()
		return nil, err
	}
	env.IdentityProviderURL = fmt.Sprintf("http://127.0.0.1:%s", dexContainer.GetPort("5556/tcp"))

	env.ClientAuthConfig = &auth.ClientConfig{
		IdentityProvider: env.IdentityProviderURL,
		OfflineAsScope:   true,
		ClientId:         "gateway",
		// deepcode ignore HardcodedPassword: just for test
		ClientSecret: "app-secret",
		Nonce:        "secret-nonce",
		Scopes: []string{
			"openid",
			"profile",
			"email",
		},
		RedirectURIs: []string{
			fmt.Sprintf("http://%s%s", RedirectURLHostname, RedirectURLPort),
		},
	}
	env.ServerAuthConfig = &auth.ServerConfig{
		URL:           localAddrOIDCProviderServer,
		TokenValidity: time.Minute * 1,
	}

	policiesPath := os.Getenv("POLICIES_PATH")
	if policiesPath == "" {
		return nil, fmt.Errorf("POLICIES_PATH not specified")
	}
	env.PoliciesPath = policiesPath

	signer := env.JwtTestEnv.CreateSigner()
	verifier, err := env.JwtTestEnv.CreateVerifier()
	if err != nil {
		return nil, err
	}

	// Start gateway
	authClient := auth.NewClient(env.ClientAuthConfig)
	if err != nil {
		return nil, err
	}
	authServer := auth.NewServer(env.ServerAuthConfig, signer, verifier)
	if err != nil {
		return nil, err
	}

	// Setup OIDC
	err = authClient.SetupOIDC(ctx)
	if err != nil {
		return nil, err
	}

	// Setup domain
	gwDomain, err := domain.NewGatewayDomain(ctx, env.ebConsumer, env.esClient)
	if err != nil {
		return nil, err
	}

	if mockData {
		err = gwDomain.UserRepository.Upsert(ctx, mock.TestAdminUser)
		if err != nil {
			return nil, err
		}
		err = gwDomain.UserRepository.Upsert(ctx, mock.TestExistingUser)
		if err != nil {
			return nil, err
		}
		err = gwDomain.UserRepository.Upsert(ctx, mock.TestTenantAdminUser)
		if err != nil {
			return nil, err
		}
		err = gwDomain.UserRoleBindingRepository.Upsert(ctx, mock.TestAdminUserRoleBinding)
		if err != nil {
			return nil, err
		}
		err = gwDomain.UserRoleBindingRepository.Upsert(ctx, mock.TestTenantAdminUserRoleBinding)
		if err != nil {
			return nil, err
		}
		err = gwDomain.ClusterRepository.Upsert(ctx, mock.TestCluster)
		if err != nil {
			return nil, err
		}
		err = gwDomain.TenantRepository.Upsert(ctx, mock.TestTenant)
		if err != nil {
			return nil, err
		}
		err = gwDomain.TenantClusterBindingRepository.Upsert(ctx, mock.TestTenantClusterBinding)
		if err != nil {
			return nil, err
		}
	}

	gatewayApiServer := NewGatewayAPIServer(env.ClientAuthConfig, authClient, authServer, gwDomain.UserRepository)
	authApiServer := NewClusterAuthAPIServer("https://localhost", signer, repositories.NewClusterAccessRepository(gwDomain.TenantClusterBindingRepository, gwDomain.ClusterRepository, gwDomain.UserRoleBindingRepository, gwDomain.TenantRepository), map[string]time.Duration{
		"default": time.Hour * 1,
	})

	gatewayAuthServer, errAuthServer := NewAuthServer(ctx, localAddrAPIServer, authServer, env.PoliciesPath, gwDomain.UserRoleBindingRepository)
	if errAuthServer != nil {
		return nil, errAuthServer
	}

	// Create gRPC server and register implementation
	env.GrpcServer = grpc.NewServer("gateway-grpc", false)
	env.GrpcServer.RegisterService(func(s ggrpc.ServiceRegistrar) {
		api.RegisterGatewayServer(s, gatewayApiServer)
		api.RegisterClusterAuthServer(s, authApiServer)
		apiCommon.RegisterServiceInformationServiceServer(s, common.NewServiceInformationService())
		api.RegisterGatewayAuthServer(s, gatewayAuthServer)
	})

	env.ApiListenerAPIServer, err = net.Listen("tcp", localAddrAPIServer)
	if err != nil {
		return nil, err
	}

	go func() {
		err := env.GrpcServer.ServeFromListener(env.ApiListenerAPIServer, nil)
		if err != nil {
			panic(err)
		}
	}()

	env.LocalOIDCProviderServer = NewOIDCProviderServer(authServer)
	env.ApiListenerOIDCProviderServer, err = net.Listen("tcp", localAddrOIDCProviderServer)
	if err != nil {
		return nil, err
	}

	go func() {
		err := env.LocalOIDCProviderServer.ServeFromListener(env.ApiListenerOIDCProviderServer)
		if err != nil {
			panic(err)
		}
	}()

	// Setup HTTP client
	env.HttpClient = &http.Client{}

	return env, nil
}

func (env *TestEnv) NewOidcClientServer(ready chan<- string) (*clientAuth.Server, error) {
	serverConf := &clientAuth.Config{
		LocalServerBindAddress: []string{
			fmt.Sprintf("%s%s", RedirectURLHostname, RedirectURLPort),
		},
		RedirectURLHostname:  RedirectURLHostname,
		LocalServerReadyChan: ready,
	}
	server, err := clientAuth.NewServer(serverConf)
	if err != nil {
		return nil, err
	}
	return server, nil
}

func (env *TestEnv) GetApiAddr() string {
	return env.ApiListenerAPIServer.Addr().String()
}

func (env *TestEnv) Shutdown() error {
	return env.TestEnv.Shutdown()
}
