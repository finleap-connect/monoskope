package storage

import (
	"gitlab.figo.systems/platform/monoskope/monoskope/internal/test"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/events"
)

const (
	testEventCreated      events.EventType     = "TestEvent:Created"
	testEventChanged      events.EventType     = "TestEvent:Changed"
	testEventDeleted      events.EventType     = "TestEvent:Deleted"
	testAggregate         events.AggregateType = "TestAggregate"
	testAggregateExtended events.AggregateType = "TestAggregateExtended"
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
