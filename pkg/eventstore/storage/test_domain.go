package storage

const (
	testEventCreated      EventType     = "TestEvent:Created"
	testEventChanged      EventType     = "TestEvent:Changed"
	testEventDeleted      EventType     = "TestEvent:Deleted"
	testEventExtended     EventType     = "TestEventExtended:Created"
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

func initTestDomain() error {
	return RegisterEventData(testEventCreated, func() EventData { return &TestEventData{} })
}

func createTestEventData(something string) *TestEventData {
	return &TestEventData{Hello: something}
}
