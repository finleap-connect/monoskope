package queryhandler

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

func TestQueryHandler(t *testing.T) {
	RegisterFailHandler(Fail)
	junitReporter := reporters.NewJUnitReporter("../../reports/queryhandler-junit.xml")
	RunSpecsWithDefaultAndCustomReporters(t, "queryhandler/integration", []Reporter{junitReporter})
}

var _ = BeforeSuite(func(done Done) {
	defer close(done)
	var err error

	By("bootstrapping test env")
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
