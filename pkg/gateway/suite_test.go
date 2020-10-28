package gateway

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"net/http/cookiejar"
	"testing"
	"time"

	"github.com/cenkalti/backoff"
	"github.com/onsi/ginkgo/reporters"
	"github.com/ory/dockertest/v3"
	dc "github.com/ory/dockertest/v3/docker"
	"golang.org/x/net/publicsuffix"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/examples/data"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gitlab.figo.systems/platform/monoskope/monoskope/internal/test"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/auth"
	auth_server "gitlab.figo.systems/platform/monoskope/monoskope/pkg/auth/server"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/logger"
)

const (
	anyLocalAddr = "127.0.0.1:0"
	redirectURL  = "http://localhost:6555/oauth/callback"
)

var (
	authRootToken = "super-secret-root-token"
	authCode      string
	pool          *dockertest.Pool
	dexContainer  *dockertest.Resource

	metricsLis                 net.Listener
	apiLis                     net.Listener
	dexConn                    *grpc.ClientConn
	server                     *Server
	authInterceptor            *auth_server.AuthServerInterceptor
	clientTransportCredentials credentials.TransportCredentials
	httpClient                 *http.Client
	dexWebEndpoint             string
	log                        logger.Logger
)

func TestGateway(t *testing.T) {
	RegisterFailHandler(Fail)
	junitReporter := reporters.NewJUnitReporter("../../reports/gateway-junit.xml")
	RunSpecsWithDefaultAndCustomReporters(t, "Gateway", []Reporter{junitReporter})
}

var _ = BeforeSuite(func(done Done) {
	var err error
	log = logger.WithName("gatewaysetup")

	By("bootstrapping gateway test env")
	pool, err = dockertest.NewPool("")
	Expect(err).ToNot(HaveOccurred())

	log.Info("spawn dex container")
	options := &dockertest.RunOptions{
		Repository: "quay.io/dexidp/dex",
		Tag:        "v2.25.0",
		PortBindings: map[dc.Port][]dc.PortBinding{
			"5556": {{HostPort: "5556"}},
		},
		ExposedPorts: []string{"5556", "5000"},
		Cmd:          []string{"serve", "/etc/dex/cfg/config.yaml"},
		Mounts:       []string{fmt.Sprintf("%s:/etc/dex/cfg", test.DexConfigPath)},
	}
	dexContainer, err = pool.RunWithOptions(options)
	Expect(err).ToNot(HaveOccurred())

	clientTransportCredentials, err = credentials.NewClientTLSFromFile(data.Path("x509/ca_cert.pem"), "x.test.example.com")
	Expect(err).ToNot(HaveOccurred())

	dexWebEndpoint = fmt.Sprintf("http://127.0.0.1:%s", dexContainer.GetPort("5556/tcp"))
	authConfig := &auth_server.Config{
		BaseConfig: auth.BaseConfig{
			IssuerURL:      dexWebEndpoint,
			OfflineAsScope: true,
		},
		RootToken:     &authRootToken,
		ValidClientId: "monoctl",
	}
	log.Info("dex issuer url: " + dexWebEndpoint)

	// Create interceptor for auth
	authInterceptor, err = auth_server.NewInterceptor(authConfig)
	Expect(err).ToNot(HaveOccurred())

	// Start gateway
	metricsLis, err = net.Listen("tcp", anyLocalAddr)
	Expect(err).ToNot(HaveOccurred())

	apiLis, err = net.Listen("tcp", anyLocalAddr)
	Expect(err).ToNot(HaveOccurred())

	cert, err := tls.LoadX509KeyPair(data.Path("x509/server_cert.pem"), data.Path("x509/server_key.pem"))
	Expect(err).ToNot(HaveOccurred())

	conf := &ServerConfig{
		KeepAlive:             false,
		AuthServerInterceptor: authInterceptor,
		TlsCert:               &cert,
	}

	ebo := backoff.NewExponentialBackOff()
	ebo.MaxElapsedTime = 5 * time.Second
	err = backoff.Retry(func() error {
		var err error
		server, err = NewServer(conf)
		return err
	}, ebo)
	Expect(err).ToNot(HaveOccurred())

	go func() {
		err := server.Serve(apiLis, metricsLis)
		if err != nil {
			panic(err)
		}
	}()

	// Setup HTTP client
	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	Expect(err).ToNot(HaveOccurred())
	httpClient = &http.Client{
		Jar: jar,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/oauth/callback", callback)
	server := &http.Server{
		Addr:    ":6555",
		Handler: mux,
	}
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			panic(err)
		}
	}()

	close(done)
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

	err = pool.Purge(dexContainer)
	Expect(err).ToNot(HaveOccurred())

	server.shutdown.Expect()

	if dexConn != nil {
		dexConn.Close()
	}
	if metricsLis != nil {
		metricsLis.Close()
	}
	if apiLis != nil {
		apiLis.Close()
	}
})
