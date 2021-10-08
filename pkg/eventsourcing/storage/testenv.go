// Copyright 2021 Monoskope Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package storage

import (
	"context"
	"fmt"

	"github.com/finleap-connect/monoskope/internal/test"
	evs "github.com/finleap-connect/monoskope/pkg/eventsourcing"
	"github.com/go-pg/pg/v10"
	"github.com/ory/dockertest/v3"
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
	Store evs.EventStore
}

func NewTestEnvWithParent(parent *test.TestEnv) (*TestEnv, error) {
	env := &TestEnv{
		TestEnv: parent,
	}

	if err := env.CreateDockerPool(false); err != nil {
		return nil, err
	}

	// Start single node crdb
	container, err := env.Run(&dockertest.RunOptions{
		Name:       "cockroach",
		Repository: "cockroachdb/cockroach",
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

	store, err := NewPostgresEventStore(env.postgresStoreConfig)
	if err != nil {
		return nil, err
	}
	env.Store = store

	return env, nil
}

func (env *TestEnv) ClearStore(ctx context.Context) {
	if pgStore, ok := env.Store.(*postgresEventStore); ok {
		if err := pgStore.clear(ctx); err != nil {
			panic(err)
		}
	} else {
		panic("that thing is not a pgstore")
	}
}

func (env *TestEnv) Shutdown() error {
	if err := env.Store.Close(); err != nil {
		return err
	}

	return env.TestEnv.Shutdown()
}
