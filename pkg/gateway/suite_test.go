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
)

const (
	anyLocalAddr = "127.0.0.1:0"
)

var (
	env *test.OAuthTestEnv

	apiListener net.Listener
	httpClient  *http.Client
	server      *Server
)

func TestGateway(t *testing.T) {
	RegisterFailHandler(Fail)
	junitReporter := reporters.NewJUnitReporter("../../reports/gateway-junit.xml")
	RunSpecsWithDefaultAndCustomReporters(t, "gateway/integration", []Reporter{junitReporter})
}

var _ = BeforeSuite(func(done Done) {
	defer close(done)
	var err error

	By("bootstrapping test env")
	env, err = test.SetupAuthTestEnv("TestGateway")
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
	return toToken(test.AuthRootToken)
}
