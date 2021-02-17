package messaging

import (
	"testing"

	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/reporters"
	. "github.com/onsi/gomega"
)

var (
	env *TestEnv
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
	env, err = NewTestEnv()
	Expect(err).ToNot(HaveOccurred())
}, 60)

var _ = AfterSuite(func() {
	var err error
	By("tearing down the test environment")

	err = env.Shutdown()
	Expect(err).To(BeNil())
})
