package event_sourcing

import (
	"github.com/google/uuid"
	api "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/commands"
	"google.golang.org/protobuf/types/known/anypb"
)

const (
	TestCommandType   CommandType   = "TestCommandType"
	TestAggregateType AggregateType = "TestAggregateType"
)

// TestCommand is a command for tests.
type TestCommand struct {
	AggID uuid.UUID
	api.TestCommandData
}

func (c *TestCommand) AggregateID() uuid.UUID { return c.AggID }
func (c *TestCommand) AggregateType() AggregateType {
	return TestAggregateType
}
func (c *TestCommand) CommandType() CommandType { return TestCommandType }
func (c *TestCommand) SetData(a *anypb.Any) error {
	return a.UnmarshalTo(&c.TestCommandData)
}
