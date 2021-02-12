package storage

import (
	"fmt"
	"os"

	"github.com/go-pg/pg/v10"
	"github.com/ory/dockertest/v3"
	"gitlab.figo.systems/platform/monoskope/monoskope/internal/test"
	evs "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing"
)

const (
	testEventCreated      evs.EventType     = "TestEvent:Created"
	testEventChanged      evs.EventType     = "TestEvent:Changed"
	testEventDeleted      evs.EventType     = "TestEvent:Deleted"
	testAggregate         evs.AggregateType = "TestAggregate"
	testAggregateExtended evs.AggregateType = "TestAggregateExtended"
)

type TestEnv struct {
	*test.TestEnv
	*postgresStoreConfig
}

func (t *TestEnv) GetStoreConfig() *postgresStoreConfig {
	return t.postgresStoreConfig
}

func NewTestEnv() (*TestEnv, error) {
	env := &TestEnv{
		TestEnv: test.NewTestEnv("StorageTestEnv"),
	}

	if err := env.CreateDockerPool(); err != nil {
		return nil, err
	}

	if v := os.Getenv("DB_URL"); v != "" {
		// create test db
		err := env.Retry(func() error {
			options, err := pg.ParseURL(v)
			if err != nil {
				return err
			}
			testDb := pg.Connect(options)
			_, err = testDb.Exec("CREATE DATABASE IF NOT EXISTS test")
			return err
		})
		if err != nil {
			return nil, err
		}

		conf, err := NewPostgresStoreConfig(v)
		if err != nil {
			return nil, err
		}
		env.postgresStoreConfig = conf
	} else {
		// Start single node crdb
		container, err := env.Run(&dockertest.RunOptions{
			Name:       "cockroach",
			Repository: "gitlab.figo.systems/platform/dependency_proxy/containers/cockroachdb/cockroach",
			Tag:        "v20.2.2",
			Cmd: []string{
				"start-single-node", "--insecure",
			},
		})
		if err != nil {
			return nil, err
		}

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
		if err != nil {
			return nil, err
		}

		conf, err := NewPostgresStoreConfig(fmt.Sprintf("postgres://root@127.0.0.1:%s/test?sslmode=disable", container.GetPort("26257/tcp")))
		if err != nil {
			return nil, err
		}
		env.postgresStoreConfig = conf
	}

	return env, nil
}

func (env *TestEnv) Shutdown() error {
	return env.TestEnv.Shutdown()
}
