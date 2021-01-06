package eventstore

import (
	"net"
	"net/http"
	"testing"

	"github.com/onsi/ginkgo/reporters"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/logger"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/messaging"
)

const (
	anyLocalAddr = "127.0.0.1:0"
)

var (
	apiListener net.Listener
	httpClient  *http.Client
	log         logger.Logger
	testServer  *server
)

func TestEventStore(t *testing.T) {
	RegisterFailHandler(Fail)
	junitReporter := reporters.NewJUnitReporter("../../reports/eventstore-junit.xml")
	RunSpecsWithDefaultAndCustomReporters(t, "eventstore/integration", []Reporter{junitReporter})
}

var _ = BeforeSuite(func(done Done) {
	defer close(done)
	var err error
	log = logger.WithName("TestEventStore")

	By("bootstrapping test env")

	// Create server
	conf := NewServerConfig()
	conf.Bus = messaging.NewMockEventBusPublisher()

	testServer, err = NewServer(conf)
	Expect(err).ToNot(HaveOccurred())
	apiListener, err = net.Listen("tcp", anyLocalAddr)
	Expect(err).ToNot(HaveOccurred())

	// Start server
	go func() {
		err := testServer.Serve(apiListener, nil)
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

	// Shutdown server
	testServer.shutdown.Expect()
	err = apiListener.Close()
	Expect(err).To(BeNil())
})
