package commandhandler

import (
	"testing"

	"github.com/onsi/ginkgo/reporters"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	testEnv *CommandHandlerTestEnv
)

func TestCommandHandler(t *testing.T) {
	RegisterFailHandler(Fail)
	junitReporter := reporters.NewJUnitReporter("../../reports/commandhandler-junit.xml")
	RunSpecsWithDefaultAndCustomReporters(t, "commandhandler/integration", []Reporter{junitReporter})
}

var _ = BeforeSuite(func(done Done) {
	defer close(done)

	By("bootstrapping test env")
	var err error
	testEnv, err = NewCommandHandlerTestEnv()
	Expect(err).To(Not(HaveOccurred()))
}, 60)

var _ = AfterSuite(func() {
	var err error
	By("tearing down the test environment")

	err = testEnv.Shutdown()
	Expect(err).To(Not(HaveOccurred()))
})
