package storage

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

type TestEventDataExtened struct {
	Hello string `json:",omitempty"`
	World string `json:",omitempty"`
}
