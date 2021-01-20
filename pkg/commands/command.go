package commands

import (
	"github.com/google/uuid"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/events"
)

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
	AggregateType() events.AggregateType

	// CommandType returns the type of the command.
	CommandType() CommandType
}

// CommandType is the type of a command, used as its unique identifier.
type CommandType string

// String returns the string representation of a command type.
func (c CommandType) String() string {
	return string(c)
}
