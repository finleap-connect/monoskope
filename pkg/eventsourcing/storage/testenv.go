package storage

import (
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

type eventStoreTestEnv struct {
	*test.TestEnv
	*postgresStoreConfig
}

func (env *eventStoreTestEnv) Shutdown() error {
	return env.TestEnv.Shutdown()
}
