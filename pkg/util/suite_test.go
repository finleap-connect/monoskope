package util

import (
	"testing"

	"github.com/onsi/ginkgo/reporters"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestUtil(t *testing.T) {
	RegisterFailHandler(Fail)
	junitReporter := reporters.NewJUnitReporter("../../reports/util-junit.xml")
	RunSpecsWithDefaultAndCustomReporters(t, "util", []Reporter{junitReporter})
}
