package commands

import (
	"testing"

	"github.com/onsi/ginkgo/reporters"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	_ "gitlab.figo.systems/platform/monoskope/monoskope/internal/test"
)

func TestCommandHandler(t *testing.T) {
	RegisterFailHandler(Fail)
	junitReporter := reporters.NewJUnitReporter("../../reports/commands-junit.xml")
	RunSpecsWithDefaultAndCustomReporters(t, "commands", []Reporter{junitReporter})
}

var _ = BeforeSuite(func(done Done) {
	defer close(done)
	By("bootstrapping test env")
}, 60)

var _ = AfterSuite(func() {
	By("tearing down the test environment")
})
