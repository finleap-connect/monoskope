package gateway

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"testing"

	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/reporters"
	. "github.com/onsi/gomega"
	"github.com/ory/dockertest/v3"
	dc "github.com/ory/dockertest/v3/docker"
	"gitlab.figo.systems/platform/monoskope/monoskope/internal/gateway/auth"
	monoctl_auth "gitlab.figo.systems/platform/monoskope/monoskope/internal/monoctl/auth"
	"gitlab.figo.systems/platform/monoskope/monoskope/internal/test"
	api_common "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/common"
	api_gw "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/gateway"
	api_gwauth "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/gateway/auth"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/grpc"
	"golang.org/x/oauth2"
	ggrpc "google.golang.org/grpc"
)

const (
	anyLocalAddr        = "127.0.0.1:0"
	RedirectURLHostname = "localhost"
	RedirectURLPort     = ":8000"
)

var (
	env *oAuthTestEnv

	apiListener net.Listener
	httpClient  *http.Client
	grpcServer  *grpc.Server
)

type oAuthTestEnv struct {
	*test.TestEnv
	DexWebEndpoint string
	AuthConfig     *auth.Config
}

func SetupAuthTestEnv(envName string) (*oAuthTestEnv, error) {
	env := &oAuthTestEnv{
		TestEnv: test.NewTestEnv(envName),
	}

	err := env.CreateDockerPool()
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
		Repository: "quay.io/dexidp/dex",
		Tag:        "v2.25.0",
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
	env.DexWebEndpoint = fmt.Sprintf("http://127.0.0.1:%s", dexContainer.GetPort("5556/tcp"))

	env.AuthConfig = &auth.Config{
		IssuerURL:      env.DexWebEndpoint,
		OfflineAsScope: true,
		ClientId:       "gateway",
		ClientSecret:   "app-secret",
		Nonce:          "secret-nonce",
	}

	return env, nil
}

func (env *oAuthTestEnv) Shutdown() error {
	return env.TestEnv.Shutdown()
}

func (env *oAuthTestEnv) NewOidcClientServer(ready chan<- string) (*monoctl_auth.Server, error) {
	serverConf := &monoctl_auth.Config{
		LocalServerBindAddress: []string{
			fmt.Sprintf("%s%s", RedirectURLHostname, RedirectURLPort),
		},
		RedirectURLHostname:  RedirectURLHostname,
		LocalServerReadyChan: ready,
	}
	server, err := monoctl_auth.NewServer(serverConf)
	if err != nil {
		return nil, err
	}
	return server, nil
}

func TestGateway(t *testing.T) {
	RegisterFailHandler(Fail)
	junitReporter := reporters.NewJUnitReporter("../../reports/gateway-junit.xml")
	RunSpecsWithDefaultAndCustomReporters(t, "gateway/integration", []Reporter{junitReporter})
}

var _ = BeforeSuite(func(done Done) {
	defer close(done)
	var err error

	By("bootstrapping test env")
	env, err = SetupAuthTestEnv("TestGateway")
	Expect(err).ToNot(HaveOccurred())

	// Start gateway
	authHandler, err := auth.NewHandler(env.AuthConfig)
	Expect(err).ToNot(HaveOccurred())
	authInterceptor, err := auth.NewInterceptor(authHandler)
	Expect(err).ToNot(HaveOccurred())

	gatewayApiServer := NewApiServer(env.AuthConfig, authHandler)

	// Create gRPC server and register implementation
	grpcServer = grpc.NewServerWithOpts("gateway-grpc", false,
		[]ggrpc.UnaryServerInterceptor{
			auth.UnaryServerInterceptor(authInterceptor.EnsureValid),
		},
		[]ggrpc.StreamServerInterceptor{
			auth.StreamServerInterceptor(authInterceptor.EnsureValid),
		})
	grpcServer.RegisterService(func(s ggrpc.ServiceRegistrar) {
		api_gw.RegisterGatewayServer(s, gatewayApiServer)
		api_gwauth.RegisterAuthServer(s, gatewayApiServer)
		api_common.RegisterServiceInformationServiceServer(s, gatewayApiServer)
	})

	apiListener, err = net.Listen("tcp", anyLocalAddr)
	Expect(err).ToNot(HaveOccurred())
	go func() {
		err := grpcServer.ServeFromListener(apiListener, nil)
		if err != nil {
			panic(err)
		}
	}()

	// Setup HTTP client
	httpClient = &http.Client{}
}, 60)

var _ = AfterSuite(func() {
	var err error
	By("tearing down the test environment")
	err = env.Shutdown()
	Expect(err).To(BeNil())

	grpcServer.Shutdown()

	err = apiListener.Close()
	Expect(err).To(BeNil())
})

func toToken(token string) *oauth2.Token {
	return &oauth2.Token{
		AccessToken: token,
	}
}

func invalidToken() *oauth2.Token {
	return toToken("some-invalid-token")
}
