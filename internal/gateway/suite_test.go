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

package gateway

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/finleap-connect/monoskope/internal/common"
	"github.com/finleap-connect/monoskope/internal/gateway/auth"
	"github.com/finleap-connect/monoskope/internal/test"
	api_common "github.com/finleap-connect/monoskope/pkg/api/domain/common"
	api "github.com/finleap-connect/monoskope/pkg/api/gateway"
	clientAuth "github.com/finleap-connect/monoskope/pkg/auth"
	"github.com/finleap-connect/monoskope/pkg/domain/constants/roles"
	"github.com/finleap-connect/monoskope/pkg/domain/constants/scopes"
	"github.com/finleap-connect/monoskope/pkg/domain/projections"
	"github.com/finleap-connect/monoskope/pkg/domain/repositories"
	es_repos "github.com/finleap-connect/monoskope/pkg/eventsourcing/repositories"
	"github.com/finleap-connect/monoskope/pkg/grpc"
	"github.com/finleap-connect/monoskope/pkg/jwt"
	"github.com/google/uuid"
	"github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
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

var (
	env *oAuthTestEnv
)

type oAuthTestEnv struct {
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
	ClusterRepo                   repositories.ClusterRepository
	AdminUser                     *projections.User
	ExistingUser                  *projections.User
	NotExistingUser               *projections.User
	PoliciesPath                  string
}

func SetupAuthTestEnv(envName string) (*oAuthTestEnv, error) {
	env := &oAuthTestEnv{
		TestEnv: test.NewTestEnv(envName),
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
		ClientSecret:     "app-secret",
		Nonce:            "secret-nonce",
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

	return env, nil
}

func (env *oAuthTestEnv) NewOidcClientServer(ready chan<- string) (*clientAuth.Server, error) {
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

func (env *oAuthTestEnv) Shutdown() error {
	return env.TestEnv.Shutdown()
}

func TestGateway(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "gateway/integration")
}

var _ = BeforeSuite(func() {
	done := make(chan interface{})

	go func() {
		defer ginkgo.GinkgoRecover()

		var err error
		ctx := context.Background()

		By("bootstrapping test env")
		env, err = SetupAuthTestEnv("TestGateway")
		Expect(err).ToNot(HaveOccurred())

		signer := env.JwtTestEnv.CreateSigner()
		verifier, err := env.JwtTestEnv.CreateVerifier()
		Expect(err).ToNot(HaveOccurred())

		// Start gateway
		authClient := auth.NewClient(env.ClientAuthConfig)
		Expect(err).ToNot(HaveOccurred())
		authServer := auth.NewServer(env.ServerAuthConfig, signer, verifier)
		Expect(err).ToNot(HaveOccurred())

		// Setup OIDC
		err = authClient.SetupOIDC(ctx)
		Expect(err).ToNot(HaveOccurred())

		// Setup user repo
		adminUser := projections.NewUserProjection(uuid.New()).(*projections.User)
		adminUser.Name = "admin"
		adminUser.Email = "admin@monoskope.io"

		env.AdminUser = adminUser
		adminRoleBinding := projections.NewUserRoleBinding(uuid.New())
		adminRoleBinding.UserId = env.AdminUser.Id
		adminRoleBinding.Role = roles.Admin.String()
		adminRoleBinding.Scope = scopes.System.String()

		existingUser := projections.NewUserProjection(uuid.New()).(*projections.User)
		existingUser.Name = "someone"
		existingUser.Email = "someone@monoskope.io"

		notExistingUser := projections.NewUserProjection(uuid.New()).(*projections.User)
		notExistingUser.Name = "nobody"
		notExistingUser.Email = "nobody@monoskope.io"

		env.ExistingUser = existingUser
		env.NotExistingUser = notExistingUser

		inMemoryUserRepo := es_repos.NewInMemoryRepository()
		inMemoryUserRoleBindingRepo := es_repos.NewInMemoryRepository()
		Expect(inMemoryUserRepo.Upsert(ctx, env.AdminUser)).ToNot(HaveOccurred())
		Expect(inMemoryUserRepo.Upsert(ctx, env.ExistingUser)).ToNot(HaveOccurred())
		Expect(inMemoryUserRoleBindingRepo.Upsert(ctx, adminRoleBinding)).ToNot(HaveOccurred())

		err = inMemoryUserRepo.Upsert(ctx, env.AdminUser)
		Expect(err).ToNot(HaveOccurred())

		// Setup cluster repo
		clusterId := uuid.New()
		testCluster := projections.NewClusterProjection(clusterId).(*projections.Cluster)
		testCluster.Name = "test-cluster"
		testCluster.DisplayName = "Test Cluster"
		testCluster.ApiServerAddress = "https://somecluster.io"
		testCluster.CaCertBundle = []byte("some-bundle")

		inMemoryClusterRepo := es_repos.NewInMemoryRepository()
		err = inMemoryClusterRepo.Upsert(ctx, testCluster)
		Expect(err).ToNot(HaveOccurred())

		userRoleBindingRepo := repositories.NewUserRoleBindingRepository(inMemoryUserRoleBindingRepo)
		userRepo := repositories.NewUserRepository(inMemoryUserRepo, repositories.NewUserRoleBindingRepository(inMemoryUserRoleBindingRepo))
		env.ClusterRepo = repositories.NewClusterRepository(inMemoryClusterRepo)
		gatewayApiServer := NewGatewayAPIServer(env.ClientAuthConfig, authClient, authServer, userRepo)
		authApiServer := NewClusterAuthAPIServer("https://localhost", signer, env.ClusterRepo, map[string]time.Duration{
			"default": time.Hour * 1,
		})

		gatewayAuthServer, errAuthServer := NewAuthServer(ctx, localAddrAPIServer, authServer, env.PoliciesPath, userRoleBindingRepo)
		Expect(errAuthServer).ToNot(HaveOccurred())

		// Create gRPC server and register implementation
		env.GrpcServer = grpc.NewServer("gateway-grpc", false)
		env.GrpcServer.RegisterService(func(s ggrpc.ServiceRegistrar) {
			api.RegisterGatewayServer(s, gatewayApiServer)
			api.RegisterClusterAuthServer(s, authApiServer)
			api_common.RegisterServiceInformationServiceServer(s, common.NewServiceInformationService())
			api.RegisterGatewayAuthServer(s, gatewayAuthServer)
		})

		env.ApiListenerAPIServer, err = net.Listen("tcp", localAddrAPIServer)
		Expect(err).ToNot(HaveOccurred())
		go func() {
			err := env.GrpcServer.ServeFromListener(env.ApiListenerAPIServer, nil)
			if err != nil {
				panic(err)
			}
		}()

		env.LocalOIDCProviderServer = NewOIDCProviderServer(authServer)
		env.ApiListenerOIDCProviderServer, err = net.Listen("tcp", localAddrOIDCProviderServer)
		Expect(err).ToNot(HaveOccurred())
		go func() {
			err := env.LocalOIDCProviderServer.ServeFromListener(env.ApiListenerOIDCProviderServer)
			if err != nil {
				panic(err)
			}
		}()

		// Setup HTTP client
		env.HttpClient = &http.Client{}
		close(done)
	}()

	Eventually(done, 60).Should(BeClosed())
})

var _ = AfterSuite(func() {
	var err error
	By("tearing down the test environment")
	err = env.Shutdown()
	Expect(err).To(BeNil())

	env.GrpcServer.Shutdown()
	env.LocalOIDCProviderServer.Shutdown()

	defer env.ApiListenerAPIServer.Close()
	defer env.ApiListenerOIDCProviderServer.Close()
})
