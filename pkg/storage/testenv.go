package storage

import (
	"github.com/go-pg/pg"
	"gitlab.figo.systems/platform/monoskope/monoskope/internal/test"
)

const (
	testEventCreated      EventType     = "TestEvent:Created"
	testEventChanged      EventType     = "TestEvent:Changed"
	testEventDeleted      EventType     = "TestEvent:Deleted"
	testAggregate         AggregateType = "TestAggregate"
	testAggregateExtended AggregateType = "TestAggregateExtended"
)

type testEventData struct {
	Hello string `json:",omitempty"`
}

type eventStoreTestEnv struct {
	*test.TestEnv
	DB *pg.DB
}

func (env *eventStoreTestEnv) Shutdown() error {
	if env.DB != nil {
		defer env.DB.Close()
	}
	return env.TestEnv.Shutdown()
}
