package gateway

import (
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/cenkalti/backoff"
	"github.com/onsi/ginkgo/reporters"
	"github.com/ory/dockertest/v3"
	dc "github.com/ory/dockertest/v3/docker"
	"google.golang.org/grpc"

	dexpb "github.com/dexidp/dex/api"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gitlab.figo.systems/platform/monoskope/monoskope/internal/test"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/gateway/auth"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/logger"
)

const (
	anyLocalAddr = "127.0.0.1:0"
)

var (
	pool         *dockertest.Pool
	dexContainer *dockertest.Resource

	metricsLis      net.Listener
	apiLis          net.Listener
	dexConn         *grpc.ClientConn
	server          *Server
	authInterceptor *auth.AuthInterceptor
)

func TestGateway(t *testing.T) {
	RegisterFailHandler(Fail)
	junitReporter := reporters.NewJUnitReporter("../../reports/gateway-junit.xml")
	RunSpecsWithDefaultAndCustomReporters(t, "Gateway", []Reporter{junitReporter})
}

var _ = BeforeSuite(func(done Done) {
	var err error
	log := logger.WithName("gatewaysetup")

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

	authConfig := &auth.Config{
		IssuerURL:      fmt.Sprintf("http://127.0.0.1:%s", dexContainer.GetPort("5556/tcp")),
		OfflineAsScope: true,
	}
	log.Info("dex issuer url: " + authConfig.IssuerURL)

	opts := []grpc.DialOption{grpc.WithInsecure()}
	dexConn, err = grpc.Dial(fmt.Sprintf("127.0.0.1:%s", dexContainer.GetPort("5000/tcp")), opts...)
	Expect(err).ToNot(HaveOccurred())
	dexClient := dexpb.NewDexClient(dexConn)

	// Create interceptor for auth
	authInterceptor, err = auth.NewAuthInterceptor(dexClient, authConfig)
	Expect(err).ToNot(HaveOccurred())

	// Start gateway
	metricsLis, err = net.Listen("tcp", anyLocalAddr)
	Expect(err).ToNot(HaveOccurred())

	apiLis, err = net.Listen("tcp", anyLocalAddr)
	Expect(err).ToNot(HaveOccurred())

	ebo := backoff.NewExponentialBackOff()
	ebo.MaxElapsedTime = 5 * time.Second
	err = backoff.Retry(func() error {
		var err error
		server, err = NewServer(false, authInterceptor)
		return err
	}, ebo)
	Expect(err).ToNot(HaveOccurred())

	go func() {
		err := server.Serve(apiLis, metricsLis)
		if err != nil {
			panic(err)
		}
	}()

	close(done)
}, 60)

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
