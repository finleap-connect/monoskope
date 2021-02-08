package commandhandler

import (
	"testing"

	"github.com/onsi/ginkgo/reporters"
	"gitlab.figo.systems/platform/monoskope/monoskope/internal/eventstore"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	testEnv *TestEnv
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
	eventStoreTestEnv, err := eventstore.NewTestEnv()
	Expect(err).To(Not(HaveOccurred()))

	testEnv, err = NewTestEnv(eventStoreTestEnv)
	Expect(err).To(Not(HaveOccurred()))
}, 60)

var _ = AfterSuite(func() {
	var err error
	By("tearing down the test environment")

	err = testEnv.Shutdown()
	Expect(err).To(Not(HaveOccurred()))

	err = testEnv.eventStoreTestEnv.Shutdown()
	Expect(err).To(Not(HaveOccurred()))
})
