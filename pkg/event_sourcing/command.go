package event_sourcing

import (
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/anypb"
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
	AggregateType() AggregateType

	// CommandType returns the type of the command.
	CommandType() CommandType

	// SetData sets type specific additional data.
	SetData(*anypb.Any) error

	// IsAuthorized checks if the given role/scope/resource allow execution.
	IsAuthorized(role Role, scope Scope, resource string) bool
}

// CommandType is the type of a command, used as its unique identifier.
type CommandType string

// String returns the string representation of a command type.
func (c CommandType) String() string {
	return string(c)
}
