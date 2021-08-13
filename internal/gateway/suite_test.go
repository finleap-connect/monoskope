package gateway

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/reporters"
	. "github.com/onsi/gomega"
	"github.com/ory/dockertest/v3"
	dc "github.com/ory/dockertest/v3/docker"
	"gitlab.figo.systems/platform/monoskope/monoskope/internal/common"
	"gitlab.figo.systems/platform/monoskope/monoskope/internal/gateway/auth"
	"gitlab.figo.systems/platform/monoskope/monoskope/internal/test"
	api_common "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/domain/common"
	projectionsApi "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/domain/projections"
	api "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/gateway"
	clientAuth "gitlab.figo.systems/platform/monoskope/monoskope/pkg/auth"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/constants/roles"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/constants/scopes"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/projections"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/repositories"
	es_repos "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing/repositories"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/grpc"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/jwt"
	ggrpc "google.golang.org/grpc"
)

const (
	localAddrAPIServer  = "127.0.0.1:9090"
	localAddrAuthServer = "127.0.0.1:9091"
	RedirectURLHostname = "localhost"
	RedirectURLPort     = ":8000"
)

var (
	env *oAuthTestEnv
)

type oAuthTestEnv struct {
	*test.TestEnv
	JwtTestEnv            *jwt.TestEnv
	AuthConfig            *auth.Config
	IdentityProviderURL   string
	ApiListenerAPIServer  net.Listener
	ApiListenerAuthServer net.Listener
	HttpClient            *http.Client
	GrpcServer            *grpc.Server
	LocalAuthServer       *authServer
	ClusterRepo           repositories.ClusterRepository
	AdminUser             *projections.User
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

	if !test.IsRunningInCI() {
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

		env.AuthConfig = &auth.Config{
			IdentityProviderName: "dex",
			IdentityProvider:     env.IdentityProviderURL,
			OfflineAsScope:       true,
			ClientId:             "gateway",
			ClientSecret:         "app-secret",
			Nonce:                "secret-nonce",
			Scopes: []string{
				"openid",
				"profile",
				"email",
				"federated:id",
			},
			RedirectURIs: []string{
				fmt.Sprintf("http://%s%s", RedirectURLHostname, RedirectURLPort),
			},
		}
	}
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
	if test.IsRunningInCI() {
		return
	}
	RegisterFailHandler(Fail)
	junitReporter := reporters.NewJUnitReporter("../../reports/gateway-junit.xml")
	RunSpecsWithDefaultAndCustomReporters(t, "gateway/integration", []Reporter{junitReporter})
}

var _ = BeforeSuite(func(done Done) {
	defer close(done)
	var err error
	ctx := context.Background()

	By("bootstrapping test env")
	env, err = SetupAuthTestEnv("TestGateway")
	Expect(err).ToNot(HaveOccurred())

	signer := env.JwtTestEnv.CreateSigner()
	verifier, err := env.JwtTestEnv.CreateVerifier(10 * time.Minute)
	Expect(err).ToNot(HaveOccurred())

	// Start gateway
	authHandler := auth.NewHandler(env.AuthConfig, signer, verifier)
	Expect(err).ToNot(HaveOccurred())

	// Setup OIDC
	err = authHandler.SetupOIDC(ctx)
	Expect(err).ToNot(HaveOccurred())

	// Setup user repo
	userId := uuid.New()
	env.AdminUser = &projections.User{User: &projectionsApi.User{Id: userId.String(), Name: "admin", Email: "admin@monoskope.io"}}
	adminRoleBinding := projections.NewUserRoleBinding(uuid.New())
	adminRoleBinding.UserId = env.AdminUser.Id
	adminRoleBinding.Role = roles.Admin.String()
	adminRoleBinding.Scope = scopes.System.String()

	inMemoryUserRepo := es_repos.NewInMemoryRepository()
	inMemoryUserRoleBindingRepo := es_repos.NewInMemoryRepository()
	err = inMemoryUserRepo.Upsert(ctx, env.AdminUser)
	Expect(err).ToNot(HaveOccurred())
	err = inMemoryUserRoleBindingRepo.Upsert(ctx, adminRoleBinding)
	Expect(err).ToNot(HaveOccurred())

	// Setup cluster repo
	clusterId := uuid.New()
	testCluster := projections.NewClusterProjection(clusterId).(*projections.Cluster)
	testCluster.Name = "test-cluster"
	testCluster.Label = "test-cluster"
	testCluster.ApiServerAddress = "https://somecluster.io"
	testCluster.CaCertBundle = []byte("some-bundle")

	inMemoryClusterRepo := es_repos.NewInMemoryRepository()
	err = inMemoryClusterRepo.Upsert(ctx, testCluster)
	Expect(err).ToNot(HaveOccurred())

	userRepo := repositories.NewUserRepository(inMemoryUserRepo, repositories.NewUserRoleBindingRepository(inMemoryUserRoleBindingRepo))
	env.ClusterRepo = repositories.NewClusterRepository(inMemoryClusterRepo)
	gatewayApiServer := NewGatewayAPIServer(env.AuthConfig, authHandler, userRepo)
	authApiServer := NewClusterAuthAPIServer("https://localhost", signer, userRepo, env.ClusterRepo, time.Hour*1)

	// Create gRPC server and register implementation
	env.GrpcServer = grpc.NewServer("gateway-grpc", false)
	env.GrpcServer.RegisterService(func(s ggrpc.ServiceRegistrar) {
		api.RegisterGatewayServer(s, gatewayApiServer)
		api.RegisterClusterAuthServer(s, authApiServer)
		api_common.RegisterServiceInformationServiceServer(s, common.NewServiceInformationService())
	})

	env.ApiListenerAPIServer, err = net.Listen("tcp", localAddrAPIServer)
	Expect(err).ToNot(HaveOccurred())
	go func() {
		err := env.GrpcServer.ServeFromListener(env.ApiListenerAPIServer, nil)
		if err != nil {
			panic(err)
		}
	}()

	env.LocalAuthServer = NewAuthServer(localAddrAPIServer, authHandler, userRepo)
	env.ApiListenerAuthServer, err = net.Listen("tcp", localAddrAuthServer)
	Expect(err).ToNot(HaveOccurred())
	go func() {
		err := env.LocalAuthServer.ServeFromListener(env.ApiListenerAuthServer)
		if err != nil {
			panic(err)
		}
	}()

	// Setup HTTP client
	env.HttpClient = &http.Client{}
}, 60)

var _ = AfterSuite(func() {
	var err error
	By("tearing down the test environment")
	err = env.Shutdown()
	Expect(err).To(BeNil())

	env.GrpcServer.Shutdown()
	env.LocalAuthServer.Shutdown()

	defer env.ApiListenerAPIServer.Close()
	defer env.ApiListenerAuthServer.Close()
})
