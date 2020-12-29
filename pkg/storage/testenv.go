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

type TestEventData struct {
	Hello string `json:",omitempty"`
}

type EventStoreTestEnv struct {
	*test.TestEnv
	DB *pg.DB
}

func (env *EventStoreTestEnv) Shutdown() error {
	if env.DB != nil {
		defer env.DB.Close()
	}
	return env.TestEnv.Shutdown()
}
