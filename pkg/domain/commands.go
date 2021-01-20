package domain

import (
	"github.com/google/uuid"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/commands"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/events"
)

const (
	AddRoleToUser commands.CommandType = "AddRoleToUser"
)

// AddRoleToUser is a command for adding a role to a user.
type AddRoleToUserCommand struct {
	ID      uuid.UUID
	Role    string
	Context string
}

func (c AddRoleToUserCommand) AggregateID() uuid.UUID              { return c.ID }
func (c AddRoleToUserCommand) AggregateType() events.AggregateType { return UserRoleBinding }
func (c AddRoleToUserCommand) CommandType() commands.CommandType   { return AddRoleToUser }
