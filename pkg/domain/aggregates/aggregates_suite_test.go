package aggregates_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/reporters"
	. "github.com/onsi/gomega"
)

func TestAggregates(t *testing.T) {
	RegisterFailHandler(Fail)
	junitReporter := reporters.NewJUnitReporter("../../../reports/aggregates-junit.xml")
	RunSpecsWithDefaultAndCustomReporters(t, "Aggregates Suite", []Reporter{junitReporter})
}
