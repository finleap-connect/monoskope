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
	if err := RegisterEventData(testEventCreated, func() EventData { return &TestEventData{} }); err != nil {
		return err
	}
	if err := RegisterEventData(testEventChanged, func() EventData { return &TestEventData{} }); err != nil {
		return err
	}
	if err := RegisterEventData(testEventDeleted, func() EventData { return &TestEventData{} }); err != nil {
		return err
	}
	return nil
}

func createTestEventData(something string) *TestEventData {
	return &TestEventData{Hello: something}
}
