package backup

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

func TestBackup(t *testing.T) {
	RegisterFailHandler(Fail)
	junitReporter := reporters.NewJUnitReporter("../../../reports/eventstore-backup-junit.xml")
	RunSpecsWithDefaultAndCustomReporters(t, "eventstore/backup", []Reporter{junitReporter})
}

var _ = BeforeSuite(func(done Done) {
	defer close(done)
	var err error

	By("bootstrapping test env")

	baseTestEnv = test.NewTestEnv("integration-testenv")
	testEnv, err = NewTestEnvWithParent(baseTestEnv)
	Expect(err).To(Not(HaveOccurred()))
}, 60)

var _ = AfterSuite(func() {
	By("tearing down the test environment")

	Expect(testEnv.Shutdown()).To(Not(HaveOccurred()))
	Expect(baseTestEnv.Shutdown()).To(Not(HaveOccurred()))
})
