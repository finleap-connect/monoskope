package monoctl

import (
	"testing"

	"github.com/onsi/ginkgo/reporters"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestMonoctl(t *testing.T) {
	RegisterFailHandler(Fail)
	junitReporter := reporters.NewJUnitReporter("../../reports/monoctl-junit.xml")
	RunSpecsWithDefaultAndCustomReporters(t, "Monoctl", []Reporter{junitReporter})
}
