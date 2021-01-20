package commands

import (
	"github.com/google/uuid"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/events"
)

const (
	TestCommandType   CommandType          = "TestCommandType"
	TestAggregateType events.AggregateType = "TestAggregateType"
)

// TestCommand is a command for tests.
type TestCommand struct {
	AggID uuid.UUID
}

func (c *TestCommand) AggregateID() uuid.UUID { return c.AggID }
func (c *TestCommand) AggregateType() events.AggregateType {
	return TestAggregateType
}
func (c *TestCommand) CommandType() CommandType { return TestCommandType }
