package config

import (
	"testing"

	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/reporters"
	. "github.com/onsi/gomega"
)

func TestConfig(t *testing.T) {
	RegisterFailHandler(Fail)
	junitReporter := reporters.NewJUnitReporter("../../../reports/monoctl-config-junit.xml")
	RunSpecsWithDefaultAndCustomReporters(t, "monoctl/conf", []Reporter{junitReporter})
}
