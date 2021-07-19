package rabbitmq

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

func TestRabbitmq(t *testing.T) {
	RegisterFailHandler(Fail)
	junitReporter := reporters.NewJUnitReporter("../../reports/jwt-junit.xml")
	RunSpecsWithDefaultAndCustomReporters(t, "Rabbitmq Suite", []Reporter{junitReporter})
}

var _ = BeforeSuite(func(done Done) {
	defer close(done)
	By("bootstrapping test env")

	var err error
	baseTestEnv = test.NewTestEnv("messaging-testenv")
	env, err = NewTestEnvWithParent(baseTestEnv)
	Expect(err).ToNot(HaveOccurred())
}, 60)

var _ = AfterSuite(func() {
	By("tearing down the test environment")
	Expect(env.Shutdown()).To(Not(HaveOccurred()))
	Expect(baseTestEnv.Shutdown()).To(Not(HaveOccurred()))
})
