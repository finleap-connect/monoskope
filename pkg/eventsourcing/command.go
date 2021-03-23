package eventsourcing

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/anypb"
)

// CommandType is the type of a command, used as its unique identifier.
type CommandType string

// String returns the string representation of a command type.
func (c CommandType) String() string {
	return string(c)
}

// Command is a domain command that is executed by a CommandHandler.
//
// A command name should 1) be in present tense and 2) contain the intent
// (CreateTenant, AddRoleToUser).
//
// The command should contain all the data needed when handling it as fields.
type Command interface {
	// AggregateID returns the ID of the aggregate that the command should be
	// handled by.
	AggregateID() uuid.UUID

	// AggregateType returns the type of the aggregate that the command can be
	// handled by.
	AggregateType() AggregateType

	// CommandType returns the type of the command.
	CommandType() CommandType

	// SetData sets type specific additional data.
	SetData(*anypb.Any) error

	// Policies returns the Role/Scope/Resource combination allowed to execute.
	Policies(ctx context.Context) []Policy
}

// BaseCommand is the base implementation for all commands
type BaseCommand struct {
	aggregateID   uuid.UUID
	aggregateType AggregateType
	commandType   CommandType
}

// NewBaseCommand creates a command.
func NewBaseCommand(id uuid.UUID, aggregateType AggregateType, commandType CommandType) *BaseCommand {
	return &BaseCommand{
		aggregateID:   id,
		aggregateType: aggregateType,
		commandType:   commandType,
	}
}

// AggregateID returns the ID of the aggregate that the command should be
// handled by.
func (c *BaseCommand) AggregateID() uuid.UUID { return c.aggregateID }

// AggregateType returns the type of the aggregate that the command can be
// handled by.
func (c *BaseCommand) AggregateType() AggregateType { return c.aggregateType }

// CommandType returns the type of the command.
func (c *BaseCommand) CommandType() CommandType { return c.commandType }
