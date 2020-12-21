package storage

import (
	"testing"

	"github.com/onsi/ginkgo/reporters"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/logger"
)

type TestEventData struct {
	Hello string `json:",omitempty"`
}

var (
	TestEvent = EventType("TestEvent")
	log       logger.Logger

	jsonString = "{\"Hello\":\"World\"}"
	jsonBytes  = []byte(jsonString)
	eventData  = TestEventData{Hello: "World"}
)

func TestEventStoreStorage(t *testing.T) {
	RegisterFailHandler(Fail)
	junitReporter := reporters.NewJUnitReporter("../../../reports/eventstore-storage-junit.xml")
	RunSpecsWithDefaultAndCustomReporters(t, "eventstore/storage", []Reporter{junitReporter})
}

var _ = BeforeSuite(func(done Done) {
	defer close(done)
	log = logger.WithName("TestEventStoreStorage")

	By("bootstrapping test env")
}, 60)

var _ = AfterSuite(func() {
	By("tearing down the test environment")
})
