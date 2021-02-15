package eventsourcing

import (
	"context"

	"github.com/google/uuid"
	cmdApi "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/eventsourcing/commands"
	"google.golang.org/protobuf/types/known/anypb"
)

const (
	testCommandType   CommandType   = "TestCommandType"
	testAggregateType AggregateType = "TestAggregateType"
)

// testCommand is a command for tests.
type testCommand struct {
	aggregateId uuid.UUID
	cmdApi.TestCommandData
}

func (c *testCommand) AggregateID() uuid.UUID { return c.aggregateId }
func (c *testCommand) AggregateType() AggregateType {
	return testAggregateType
}
func (c *testCommand) CommandType() CommandType { return testCommandType }
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
