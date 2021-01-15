package eventstore

import (
	"net"
	"net/http"
	"testing"

	"github.com/onsi/ginkgo/reporters"
	"google.golang.org/grpc"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	api "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/eventstore"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/grpcutil"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/logger"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/messaging"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/storage"
)

var (
	apiListener net.Listener
	httpClient  *http.Client
	log         logger.Logger
	grpcServer  *grpcutil.Server
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
	grpcServer = grpcutil.NewServer("event_store_grpc", false)

	eventStore, err := NewApiServer(storage.NewInMemoryEventStore(), messaging.NewMockEventBusPublisher())
	Expect(err).ToNot(HaveOccurred())

	grpcServer.RegisterService(func(s grpc.ServiceRegistrar) {
		api.RegisterEventStoreServer(s, eventStore)
	})
	grpcServer.RegisterOnShutdown(func() {
		eventStore.Shutdown()
	})

	apiListener, err = net.Listen("tcp", "127.0.0.1:0")
	Expect(err).ToNot(HaveOccurred())

	// Start server
	go func() {
		err := grpcServer.Serve(apiListener, nil)
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
	grpcServer.Shutdown()
	err = apiListener.Close()
	Expect(err).To(BeNil())
})
