package repositories

import (
	"testing"

	"github.com/onsi/ginkgo/reporters"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	_ "gitlab.figo.systems/platform/monoskope/monoskope/internal/test"
)

func TestRepositories(t *testing.T) {
	RegisterFailHandler(Fail)
	junitReporter := reporters.NewJUnitReporter("../../../reports/event-sourcing-repositories-junit.xml")
	RunSpecsWithDefaultAndCustomReporters(t, "TestRepositories", []Reporter{junitReporter})
}

var _ = BeforeSuite(func(done Done) {
	defer close(done)
	By("bootstrapping test env")
}, 60)

var _ = AfterSuite(func() {
	By("tearing down the test environment")
})
