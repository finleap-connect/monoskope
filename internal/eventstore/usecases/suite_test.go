package usecases

import (
	"testing"

	"github.com/onsi/ginkgo/reporters"
	evs "gitlab.figo.systems/platform/monoskope/monoskope/pkg/event_sourcing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestUsecases(t *testing.T) {
	RegisterFailHandler(Fail)
	junitReporter := reporters.NewJUnitReporter("../../../reports/eventstore-usecases-junit.xml")
	RunSpecsWithDefaultAndCustomReporters(t, "eventstore/usecases", []Reporter{junitReporter})
}

const (
	testEventType     evs.EventType     = "TestEventType"
	testAggregateType evs.AggregateType = "TestAggregateType"
)
