package backup

import (
	"testing"

	"github.com/onsi/ginkgo/reporters"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	testEnv *TestEnv
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

	testEnv, err = NewTestEnv()
	Expect(err).To(Not(HaveOccurred()))
}, 60)

var _ = AfterSuite(func() {
	var err error
	By("tearing down the test environment")

	err = testEnv.Shutdown()
	Expect(err).To(Not(HaveOccurred()))
})
