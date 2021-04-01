package storage

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

func TestEventStoreStorage(t *testing.T) {
	RegisterFailHandler(Fail)
	junitReporter := reporters.NewJUnitReporter("../../../reports/event-sourcing-storage-junit.xml")
	RunSpecsWithDefaultAndCustomReporters(t, "eventsourcing/storage", []Reporter{junitReporter})
}

var _ = BeforeSuite(func(done Done) {
	var err error
	defer close(done)

	By("bootstrapping test env")
	baseTestEnv = test.NewTestEnv("storage-testenv")
	env, err = NewTestEnvWithParent(baseTestEnv)
	Expect(err).ToNot(HaveOccurred())
}, 180)

var _ = AfterSuite(func() {
	By("tearing down the test environment")

	Expect(env.Shutdown()).To(Not(HaveOccurred()))
	Expect(baseTestEnv.Shutdown()).To(Not(HaveOccurred()))
})
