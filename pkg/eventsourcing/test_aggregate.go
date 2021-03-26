package eventsourcing

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	testEd "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/eventsourcing/eventdata"
)

const (
	testAggregateType AggregateType = "TestAggregateType"
	testEventType     EventType     = "TestEventType"
)

type testAggregate struct {
	*BaseAggregate
	Test string
}

func newTestAggregate() *testAggregate {
	return &testAggregate{
		BaseAggregate: NewBaseAggregate(testAggregateType, uuid.New()),
	}
}

func (a *testAggregate) HandleCommand(ctx context.Context, cmd Command) error {
	switch cmd := cmd.(type) {
	case *testCommand:
		ed := ToEventDataFromProto(&testEd.TestEventData{
			Hello: cmd.TestCommandData.GetTest(),
		})
		_ = a.AppendEvent(ctx, testEventType, ed)
		return nil
	}
	return fmt.Errorf("couldn't handle command")
}

func (a *testAggregate) ApplyEvent(ev Event) error {
	return nil
}
