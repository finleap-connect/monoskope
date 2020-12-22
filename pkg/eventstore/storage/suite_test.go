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

type TestEventDataExtened struct {
	Hello string `json:",omitempty"`
	World string `json:",omitempty"`
}

const (
	TestEvent         = EventType("TestEvent")
	TestEventExtended = EventType("TestEventExtended")
	jsonString        = "{\"Hello\":\"World\"}"
)

var (
	log logger.Logger

	jsonBytes     = []byte(jsonString)
	testEventData = TestEventData{Hello: "World"}
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

	// Register event data for test event
	err := RegisterEventData(TestEvent, func() EventData { return &TestEventData{} })
	Expect(err).ToNot(HaveOccurred())
}, 60)

var _ = AfterSuite(func() {
	By("tearing down the test environment")
})
