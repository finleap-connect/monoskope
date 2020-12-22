package storage

import (
	"context"
	"fmt"
	"testing"

	"github.com/go-pg/pg"
	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/reporters"
	. "github.com/onsi/gomega"
	"github.com/ory/dockertest/v3"
	"gitlab.figo.systems/platform/monoskope/monoskope/internal/test"
	storage_test "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventstore/storage/test"
)

const (
	typeTestEventCreated      = EventType("TestEvent:Created")
	typeTestEventChanged      = EventType("TestEvent:Changed")
	typeTestEventDeleted      = EventType("TestEvent:Deleted")
	typeTestEventExtended     = EventType("TestEventExtended:Created")
	typeTestAggregate         = AggregateType("TestAggregate")
	typeTestAggregateExtended = AggregateType("TestAggregateExtended")
	jsonString                = "{\"Hello\":\"World\"}"
)

var (
	env           *storage_test.EventStoreTestEnv
	jsonBytes     = []byte(jsonString)
	testEventData = createTestEventData("World")
	ctx           = context.Background()
)

func TestEventStoreStorage(t *testing.T) {
	RegisterFailHandler(Fail)
	junitReporter := reporters.NewJUnitReporter("../../../reports/eventstore-storage-junit.xml")
	RunSpecsWithDefaultAndCustomReporters(t, "eventstore/storage", []Reporter{junitReporter})
}

var _ = BeforeSuite(func(done Done) {
	var err error
	defer close(done)

	By("bootstrapping test env")
	env = &storage_test.EventStoreTestEnv{
		TestEnv: test.SetupGeneralTestEnv("TestEventStoreStorage"),
	}

	// Register event data for test event
	err = RegisterEventData(typeTestEventCreated, func() EventData { return &storage_test.TestEventData{} })
	Expect(err).ToNot(HaveOccurred())

	err = env.CreateDockerPool()
	Expect(err).ToNot(HaveOccurred())

	// Start single node crdb
	container, err := env.RunWithOptions(&dockertest.RunOptions{
		Name:       "cockroach",
		Repository: "gitlab.figo.systems/platform/dependency_proxy/containers/cockroachdb/cockroach",
		Tag:        "v20.2.2",
		Cmd: []string{
			"start-single-node", "--insecure",
		},
	})
	Expect(err).ToNot(HaveOccurred())

	// create test db
	err = env.Retry(func() error {
		testDb := pg.Connect(&pg.Options{
			Addr:     fmt.Sprintf("127.0.0.1:%s", container.GetPort("26257/tcp")),
			Database: "",
			User:     "root",
			Password: "",
		})
		_, err := testDb.Exec("CREATE DATABASE IF NOT EXISTS test")
		return err
	})
	Expect(err).ToNot(HaveOccurred())

	// create test db connection for tests
	env.DB = pg.Connect(&pg.Options{
		Addr:     fmt.Sprintf("127.0.0.1:%s", container.GetPort("26257/tcp")),
		Database: "test",
		User:     "root",
		Password: "",
	})
}, 60)

var _ = AfterSuite(func() {
	var err error
	By("tearing down the test environment")
	err = env.Shutdown()
	Expect(err).To(BeNil())
})

func createTestEventData(something string) *storage_test.TestEventData {
	return &storage_test.TestEventData{Hello: something}
}
