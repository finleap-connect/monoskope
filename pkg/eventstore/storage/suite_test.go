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
)

var (
	env *EventStoreTestEnv
	ctx = context.Background()
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
	env = &EventStoreTestEnv{
		TestEnv: test.SetupGeneralTestEnv("TestEventStoreStorage"),
	}

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
