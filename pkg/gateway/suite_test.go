package gateway

import (
	"net"
	"net/http"
	"testing"

	"github.com/onsi/ginkgo/reporters"
	"golang.org/x/oauth2"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gitlab.figo.systems/platform/monoskope/monoskope/internal/test"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/logger"
)

const (
	anyLocalAddr = "127.0.0.1:0"
)

var (
	env *test.OAuthTestEnv

	gatewayApiListener net.Listener
	httpClient         *http.Client
	log                logger.Logger
	gatewayServer      *Server
)

func TestGateway(t *testing.T) {
	RegisterFailHandler(Fail)
	junitReporter := reporters.NewJUnitReporter("../../reports/gateway-junit.xml")
	RunSpecsWithDefaultAndCustomReporters(t, "gateway/integration", []Reporter{junitReporter})
}

var _ = BeforeSuite(func(done Done) {
	defer close(done)
	var err error
	log = logger.WithName("gateway")

	By("bootstrapping test env")
	env, err = test.SetupAuthTestEnv()
	Expect(err).ToNot(HaveOccurred())

	// Start gateway
	conf := &ServerConfig{
		KeepAlive:  false,
		AuthConfig: env.AuthConfig,
	}

	gatewayServer, err = NewServer(conf)
	Expect(err).ToNot(HaveOccurred())

	gatewayApiListener, err = net.Listen("tcp", anyLocalAddr)
	Expect(err).ToNot(HaveOccurred())
	go func() {
		err := gatewayServer.Serve(gatewayApiListener, nil)
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

	gatewayServer.shutdown.Expect()

	err = gatewayApiListener.Close()
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
	return toToken(test.AuthRootToken)
}
