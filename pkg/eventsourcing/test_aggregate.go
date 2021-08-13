package eventsourcing

import (
	"context"
	"fmt"

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
		BaseAggregate: NewBaseAggregate(testAggregateType),
	}
}

func (a *testAggregate) HandleCommand(ctx context.Context, cmd Command) (*CommandReply, error) {
	switch cmd := cmd.(type) {
	case *testCommand:
		ed := ToEventDataFromProto(&testEd.TestEventData{
			Hello: cmd.TestCommandData.GetTest(),
		})
		agg := a.AppendEvent(ctx, testEventType, ed)
		ret := &CommandReply{
			Id:      agg.AggregateID(),
			Version: agg.AggregateVersion(),
		}
		return ret, nil
	}
	return nil, fmt.Errorf("couldn't handle command")
}

func (a *testAggregate) ApplyEvent(ev Event) error {
	return nil
}
