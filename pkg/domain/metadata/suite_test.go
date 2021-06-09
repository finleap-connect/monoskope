package metadata

import (
	"testing"

	"github.com/onsi/ginkgo/reporters"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestQueryHandler(t *testing.T) {
	RegisterFailHandler(Fail)
	junitReporter := reporters.NewJUnitReporter("../../../reports/domain-metadata-junit.xml")
	RunSpecsWithDefaultAndCustomReporters(t, "pkg/domain/metadata", []Reporter{junitReporter})
}

var _ = BeforeSuite(func(done Done) {
	defer close(done)
	By("bootstrapping test env")
}, 60)

var _ = AfterSuite(func() {
	By("tearing down the test environment")
})
