package commandhandler

import (
	"net"
	"testing"

	"github.com/onsi/ginkgo/reporters"
	"google.golang.org/grpc"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	api "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/commandhandler"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/grpcutil"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/logger"
)

var (
	apiListener net.Listener
	log         logger.Logger
	grpcServer  *grpcutil.Server
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

	// Create server
	conf := grpcutil.NewServerConfig("command_handler")
	grpcServer = grpcutil.NewServer(conf)

	commandHandler, err := NewApiServer()
	Expect(err).ToNot(HaveOccurred())

	grpcServer.RegisterService(func(s grpc.ServiceRegistrar) {
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
