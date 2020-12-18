package eventstore

import (
	"net"
	"net/http"
	"testing"

	"github.com/onsi/ginkgo/reporters"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/logger"
)

const (
	anyLocalAddr = "127.0.0.1:0"
)

var (
	apiListener net.Listener
	httpClient  *http.Client
	log         logger.Logger
	server      *Server
)

func TestEventStore(t *testing.T) {
	RegisterFailHandler(Fail)
	junitReporter := reporters.NewJUnitReporter("../../reports/eventstore-junit.xml")
	RunSpecsWithDefaultAndCustomReporters(t, "eventstore integration tests", []Reporter{junitReporter})
}

var _ = BeforeSuite(func(done Done) {
	defer close(done)
	var err error
	log = logger.WithName("eventstore")

	By("bootstrapping test env")

	// Start server
	conf := &ServerConfig{
		KeepAlive: false,
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

	server.shutdown.Expect()

	err = apiListener.Close()
	Expect(err).To(BeNil())
})
