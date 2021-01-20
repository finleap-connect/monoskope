package event_sourcing

import (
	"github.com/google/uuid"
)

const (
	TestCommandType   CommandType   = "TestCommandType"
	TestAggregateType AggregateType = "TestAggregateType"
)

// TestCommand is a command for tests.
type TestCommand struct {
	AggID uuid.UUID
}

func (c *TestCommand) AggregateID() uuid.UUID { return c.AggID }
func (c *TestCommand) AggregateType() AggregateType {
	return TestAggregateType
}
func (c *TestCommand) CommandType() CommandType { return TestCommandType }
