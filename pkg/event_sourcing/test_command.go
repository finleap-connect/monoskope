package event_sourcing

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	api_cmd "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/commands"
	api_ev "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/eventdata/test"
	"google.golang.org/protobuf/types/known/anypb"
)

const (
	TestCommandType   CommandType   = "TestCommandType"
	TestAggregateType AggregateType = "TestAggregateType"
	TestEventType     EventType     = "TestEventType"
)

// TestCommand is a command for tests.
type TestCommand struct {
	AggID uuid.UUID
	api_cmd.TestCommandData
}

func (c *TestCommand) AggregateID() uuid.UUID { return c.AggID }
func (c *TestCommand) AggregateType() AggregateType {
	return TestAggregateType
}
func (c *TestCommand) CommandType() CommandType { return TestCommandType }
func (c *TestCommand) SetData(a *anypb.Any) error {
	return a.UnmarshalTo(&c.TestCommandData)
}

type TestAggregate struct {
	AggregateBase
}

func (a *TestAggregate) HandleCommand(ctx context.Context, cmd Command) error {
	switch cmd := cmd.(type) {
	case *TestCommand:
		ed, err := ToEventDataFromProto(&api_ev.TestEventData{
			Hello: cmd.TestCommandData.GetTest(),
		})
		if err != nil {
			return err
		}
		_ = a.AppendEvent(TestEventType, ed)
		return nil
	}
	return fmt.Errorf("couldn't handle command")
}
