package storage

import (
	"gitlab.figo.systems/platform/monoskope/monoskope/internal/test"
	evs "gitlab.figo.systems/platform/monoskope/monoskope/pkg/event_sourcing"
)

const (
	testEventCreated      evs.EventType     = "TestEvent:Created"
	testEventChanged      evs.EventType     = "TestEvent:Changed"
	testEventDeleted      evs.EventType     = "TestEvent:Deleted"
	testAggregate         evs.AggregateType = "TestAggregate"
	testAggregateExtended evs.AggregateType = "TestAggregateExtended"
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
