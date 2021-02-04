package event_sourcing

import (
	"context"

	"github.com/google/uuid"
	api_cmd "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/commands"
	"google.golang.org/protobuf/types/known/anypb"
)

const (
	TestCommandType   CommandType   = "TestCommandType"
	TestAggregateType AggregateType = "TestAggregateType"
	TestEventType     EventType     = "TestEventType"
)

// testCommand is a command for tests.
type testCommand struct {
	aggregateId uuid.UUID
	api_cmd.TestCommandData
}

func (c *testCommand) AggregateID() uuid.UUID { return c.aggregateId }
func (c *testCommand) AggregateType() AggregateType {
	return TestAggregateType
}
func (c *testCommand) CommandType() CommandType { return TestCommandType }
func (c *testCommand) SetData(a *anypb.Any) error {
	return a.UnmarshalTo(&c.TestCommandData)
}
func (c *testCommand) Policies(ctx context.Context) []Policy {
	return []Policy{}
}

// type testAggregate struct {
// 	*BaseAggregate
// }

// func newTestAggregate() *testAggregate {
// 	return &testAggregate{
// 		BaseAggregate: NewBaseAggregate(TestAggregateType, uuid.New()),
// 	}
// }

// func (a *testAggregate) HandleCommand(ctx context.Context, cmd Command) error {
// 	switch cmd := cmd.(type) {
// 	case *testCommand:
// 		ed, err := ToEventDataFromProto(&api_ev.TestEventData{
// 			Hello: cmd.TestCommandData.GetTest(),
// 		})
// 		if err != nil {
// 			return err
// 		}
// 		_ = a.AppendEvent(TestEventType, ed)
// 		return nil
// 	}
// 	return fmt.Errorf("couldn't handle command")
// }
