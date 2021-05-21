package projectors

import (
	"github.com/onsi/ginkgo/reporters"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	_ "gitlab.figo.systems/platform/monoskope/monoskope/internal/test"
)

var _ = Describe("Testing with Ginkgo", func() {
	It("projectors", func() {

		RegisterFailHandler(Fail)
		junitReporter := reporters.NewJUnitReporter("../../../reports/domain-projectors-junit.xml")
		RunSpecsWithDefaultAndCustomReporters(GinkgoT(), "TestProjectors", []Reporter{junitReporter})
	})
})

var _ = BeforeSuite(func(done Done) {
	defer close(done)
	By("bootstrapping test env")
}, 60)

var _ = AfterSuite(func() {
	By("tearing down the test environment")
})
