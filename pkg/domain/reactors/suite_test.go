package reactors

import (
	"testing"

	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/reporters"
	. "github.com/onsi/gomega"

	"gitlab.figo.systems/platform/monoskope/monoskope/internal/test"
	_ "gitlab.figo.systems/platform/monoskope/monoskope/internal/test"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/jwt"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/util"
)

func TestReactors(t *testing.T) {
	RegisterFailHandler(Fail)
	junitReporter := reporters.NewJUnitReporter("../../../reports/reactors-junit.xml")
	RunSpecsWithDefaultAndCustomReporters(t, "TestReactors", []Reporter{junitReporter})
}

var JwtTestEnv *jwt.TestEnv

var _ = BeforeSuite(func(done Done) {
	defer close(done)
	By("bootstrapping test env")

	testEnv, err := jwt.NewTestEnv(test.NewTestEnv("TestReactors"))
	util.PanicOnError(err)
	JwtTestEnv = testEnv

}, 60)

var _ = AfterSuite(func() {
	By("tearing down the test environment")
})
