package internal

import (
	"testing"

	"github.com/onsi/ginkgo/reporters"
	"gitlab.figo.systems/platform/monoskope/monoskope/internal/test"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	baseTestEnv *test.TestEnv
	testEnv     *TestEnv
)

func TestQueryHandler(t *testing.T) {
	RegisterFailHandler(Fail)
	junitReporter := reporters.NewJUnitReporter("../reports/internal-junit.xml")
	RunSpecsWithDefaultAndCustomReporters(t, "integration", []Reporter{junitReporter})
}

var _ = BeforeSuite(func(done Done) {
	defer close(done)
	var err error

	By("bootstrapping test env")
	baseTestEnv = test.NewTestEnv("integration-testenv")
	testEnv, err = NewTestEnv(baseTestEnv)
	Expect(err).To(Not(HaveOccurred()))
}, 120)

var _ = AfterSuite(func() {
	By("tearing down the test environment")

	Expect(testEnv.Shutdown()).To(Not(HaveOccurred()))
	Expect(baseTestEnv.Shutdown()).To(Not(HaveOccurred()))
})
