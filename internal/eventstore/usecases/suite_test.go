package usecases

import (
	"testing"

	"github.com/onsi/ginkgo/reporters"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestUsecases(t *testing.T) {
	RegisterFailHandler(Fail)
	junitReporter := reporters.NewJUnitReporter("../../../reports/eventstore-usecases-junit.xml")
	RunSpecsWithDefaultAndCustomReporters(t, "eventstore/usecases", []Reporter{junitReporter})
}
