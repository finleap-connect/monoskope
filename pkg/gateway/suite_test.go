package gateway

import (
	"net"
	"net/http"
	"testing"

	"github.com/onsi/ginkgo/reporters"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/logger"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/util"
)

const (
	anyLocalAddr = "127.0.0.1:0"
	redirectURL  = "http://localhost:6555/oauth/callback"
)

var (
	env      *util.OAuthTestEnv
	authCode string

	gatewayApiListener net.Listener
	httpClient         *http.Client
	log                logger.Logger
	gatewayServer      *Server
	httpServer         *http.Server
)

func TestGateway(t *testing.T) {
	RegisterFailHandler(Fail)
	junitReporter := reporters.NewJUnitReporter("../../reports/gateway-junit.xml")
	RunSpecsWithDefaultAndCustomReporters(t, "Gateway", []Reporter{junitReporter})
}

var _ = BeforeSuite(func(done Done) {
	defer close(done)
	var err error
	log = logger.WithName("gateway")

	By("bootstrapping test env")
	env, err = util.SetupAuthTestEnv()
	Expect(err).ToNot(HaveOccurred())

	// Start gateway
	gatewayApiListener, err = net.Listen("tcp", anyLocalAddr)
	Expect(err).ToNot(HaveOccurred())

	conf := &ServerConfig{
		KeepAlive:             false,
		AuthServerInterceptor: env.AuthInterceptor,
		TlsCert:               env.GatewayTlsCert,
	}

	gatewayServer, err = NewServer(conf)
	Expect(err).ToNot(HaveOccurred())

	go func() {
		err := gatewayServer.Serve(gatewayApiListener, nil)
		if err != nil {
			panic(err)
		}
	}()

	// Setup HTTP client
	Expect(err).ToNot(HaveOccurred())
	httpClient = &http.Client{}

	mux := http.NewServeMux()
	mux.HandleFunc("/oauth/callback", callback)
	httpServer = &http.Server{
		Addr:    ":6555",
		Handler: mux,
	}
	go func() {
		_ = httpServer.ListenAndServe()
	}()
}, 60)

func callback(rw http.ResponseWriter, r *http.Request) {
	log.Info("received auth callback")
	err := r.ParseForm()
	if err != nil {
		return
	}
	// Authorization redirect callback from OAuth2 auth flow.
	if errMsg := r.Form.Get("error"); errMsg != "" {
		log.Error(err, errMsg)
		return
	}
	authCode = r.Form.Get("code")
}

var _ = AfterSuite(func() {
	var err error
	By("tearing down the test environment")
	err = env.Shutdown()
	Expect(err).To(BeNil())

	gatewayServer.shutdown.Expect()

	err = gatewayApiListener.Close()
	Expect(err).To(BeNil())

	_ = httpServer.Close()
})
