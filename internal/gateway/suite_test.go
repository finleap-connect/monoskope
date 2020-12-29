package gateway

import (
	"fmt"
	"net"
	"net/http"
	"testing"

	"github.com/onsi/ginkgo/reporters"
	"golang.org/x/oauth2"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/ory/dockertest/v3"
	dc "github.com/ory/dockertest/v3/docker"
	gw_auth "gitlab.figo.systems/platform/monoskope/monoskope/internal/gateway/auth"
	monoctl_auth "gitlab.figo.systems/platform/monoskope/monoskope/internal/monoctl/auth"
	"gitlab.figo.systems/platform/monoskope/monoskope/internal/test"
)

const (
	anyLocalAddr        = "127.0.0.1:0"
	AuthRootToken       = "super-secret-root-token"
	RedirectURLHostname = "localhost"
	RedirectURLPort     = ":8000"
)

var (
	env *OAuthTestEnv

	apiListener net.Listener
	httpClient  *http.Client
	server      *Server
)

type OAuthTestEnv struct {
	*test.TestEnv
	DexWebEndpoint string
	AuthConfig     *gw_auth.Config
}

func SetupAuthTestEnv(envName string) (*OAuthTestEnv, error) {
	env := &OAuthTestEnv{
		TestEnv: test.SetupGeneralTestEnv(envName),
	}

	err := env.CreateDockerPool()
	if err != nil {
		_ = env.Shutdown()
		return nil, err
	}

	dexContainer, err := env.RunWithOptions(&dockertest.RunOptions{
		Name:       "dex",
		Repository: "quay.io/dexidp/dex",
		Tag:        "v2.25.0",
		PortBindings: map[dc.Port][]dc.PortBinding{
			"5556": {{HostPort: "5556"}},
		},
		ExposedPorts: []string{"5556", "5000"},
		Cmd:          []string{"serve", "/etc/dex/cfg/config.yaml"},
		Mounts:       []string{fmt.Sprintf("%s:/etc/dex/cfg", test.DexConfigPath)},
	})
	if err != nil {
		_ = env.Shutdown()
		return nil, err
	}
	env.DexWebEndpoint = fmt.Sprintf("http://127.0.0.1:%s", dexContainer.GetPort("5556/tcp"))

	rootToken := AuthRootToken
	env.AuthConfig = &gw_auth.Config{
		IssuerURL:      env.DexWebEndpoint,
		OfflineAsScope: true,
		RootToken:      &rootToken,
		ClientId:       "gateway",
		ClientSecret:   "app-secret",
		Nonce:          "secret-nonce",
	}

	return env, nil
}

func (env *OAuthTestEnv) Shutdown() error {
	return env.TestEnv.Shutdown()
}

func (env *OAuthTestEnv) NewOidcClientServer(ready chan<- string) (*monoctl_auth.Server, error) {
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
	conf := &ServerConfig{
		KeepAlive:  false,
		AuthConfig: env.AuthConfig,
	}

	server, err = NewServer(conf)
	Expect(err).ToNot(HaveOccurred())

	apiListener, err = net.Listen("tcp", anyLocalAddr)
	Expect(err).ToNot(HaveOccurred())
	go func() {
		err := server.Serve(apiListener, nil)
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

	server.shutdown.Expect()

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

func rootToken() *oauth2.Token {
	return toToken(AuthRootToken)
}
