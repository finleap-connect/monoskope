package storage

import "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventstore/storage/test"

const (
	testEventCreated      EventType     = "TestEvent:Created"
	testEventChanged      EventType     = "TestEvent:Changed"
	testEventDeleted      EventType     = "TestEvent:Deleted"
	testEventExtended     EventType     = "TestEventExtended:Created"
	testAggregate         AggregateType = "TestAggregate"
	testAggregateExtended AggregateType = "TestAggregateExtended"
)

func initTestDomain() error {
	return RegisterEventData(testEventCreated, func() EventData { return &test.TestEventData{} })
}

func createTestEventData(something string) *test.TestEventData {
	return &test.TestEventData{Hello: something}
}
