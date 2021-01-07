package storage

import (
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
	*postgresStoreConfig
}

func (env *eventStoreTestEnv) Shutdown() error {
	return env.TestEnv.Shutdown()
}
