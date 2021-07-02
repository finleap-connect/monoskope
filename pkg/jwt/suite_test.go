package jwt

import (
	"testing"

	"github.com/onsi/ginkgo/reporters"
	"gitlab.figo.systems/platform/monoskope/monoskope/internal/test"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var testEnv *TestEnv

func TestJWT(t *testing.T) {
	RegisterFailHandler(Fail)
	junitReporter := reporters.NewJUnitReporter("../../reports/jwt-junit.xml")
	RunSpecsWithDefaultAndCustomReporters(t, "jwt", []Reporter{junitReporter})
}

var _ = BeforeSuite(func(done Done) {
	defer close(done)
	By("bootstrapping test env")

	var err error
	baseTestEnv := test.NewTestEnv("integration-testenv")
	testEnv, err = NewTestEnv(baseTestEnv)
	Expect(err).To(Not(HaveOccurred()))
}, 60)

var _ = AfterSuite(func() {
	By("tearing down the test environment")
})
