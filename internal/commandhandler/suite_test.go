package commandhandler

import (
	"context"
	"net"
	"testing"

	"github.com/onsi/ginkgo/reporters"
	ggrpc "google.golang.org/grpc"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gitlab.figo.systems/platform/monoskope/monoskope/internal/eventstore"
	api "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/commandhandler"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/grpc"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/logger"
)

var (
	apiListener net.Listener
	log         logger.Logger
	grpcServer  *grpc.Server
	testEnv     *eventstore.EventStoreTestEnv
)

func TestCommandHandler(t *testing.T) {
	RegisterFailHandler(Fail)
	junitReporter := reporters.NewJUnitReporter("../../reports/commandhandler-junit.xml")
	RunSpecsWithDefaultAndCustomReporters(t, "commandhandler/integration", []Reporter{junitReporter})
}

var _ = BeforeSuite(func(done Done) {
	defer close(done)
	var err error
	log = logger.WithName("TestCommandHandler")

	By("bootstrapping test env")

	testEnv, err = eventstore.NewEventStoreTestEnv()
	Expect(err).To(Not(HaveOccurred()))

	esClient, err := testEnv.GetApiClient(context.Background())
	Expect(err).To(Not(HaveOccurred()))

	// Create server
	grpcServer = grpc.NewServer("command_handler_grpc", false)

	commandHandler := NewApiServer(esClient)
	grpcServer.RegisterService(func(s ggrpc.ServiceRegistrar) {
		api.RegisterCommandHandlerServer(s, commandHandler)
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
}, 60)

var _ = AfterSuite(func() {
	var err error
	By("tearing down the test environment")

	// Shutdown server
	grpcServer.Shutdown()
	err = apiListener.Close()
	Expect(err).To(BeNil())
})
