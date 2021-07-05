package usecases_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/reporters"
	. "github.com/onsi/gomega"
)

func TestUsecases(t *testing.T) {
	RegisterFailHandler(Fail)
	junitReporter := reporters.NewJUnitReporter("../../../reports/aggregates-junit.xml")
	RunSpecsWithDefaultAndCustomReporters(t, "Usecases Suite", []Reporter{junitReporter})
}
