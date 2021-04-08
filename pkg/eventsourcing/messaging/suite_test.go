package messaging

import (
	"testing"

	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/reporters"
	. "github.com/onsi/gomega"
	"gitlab.figo.systems/platform/monoskope/monoskope/internal/test"
)

var (
	baseTestEnv *test.TestEnv
	env         *TestEnv
)

func TestMessageBus(t *testing.T) {
	RegisterFailHandler(Fail)
	junitReporter := reporters.NewJUnitReporter("../../../reports/event-sourcing-messaging-junit.xml")
	RunSpecsWithDefaultAndCustomReporters(t, "eventsourcing/messaging", []Reporter{junitReporter})
}

var _ = BeforeSuite(func(done Done) {
	var err error
	defer close(done)

	By("bootstrapping test env")
	baseTestEnv = test.NewTestEnv("messaging-testenv")
	env, err = NewTestEnvWithParent(baseTestEnv)
	Expect(err).ToNot(HaveOccurred())
}, 60)

var _ = AfterSuite(func() {
	By("tearing down the test environment")

	Expect(env.Shutdown()).To(Not(HaveOccurred()))
	Expect(baseTestEnv.Shutdown()).To(Not(HaveOccurred()))
})
