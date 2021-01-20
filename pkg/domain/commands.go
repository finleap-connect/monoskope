package domain

import (
	"github.com/google/uuid"
	apicmd_user "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/commands/user"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/commands"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/events"
)

func InitCommands(regsitry commands.CommandRegistry) error {
	if err := regsitry.RegisterCommand(addRoleToUserCommandFactory); err != nil {
		return err
	}
	return nil
}

const (
	AddRoleToUser commands.CommandType = "AddRoleToUser"
)

// AddRoleToUser is a command for adding a role to a user.
type AddRoleToUserCommand struct {
	AggID uuid.UUID
	apicmd_user.AddRoleToUser
}

var addRoleToUserCommandFactory = func() commands.Command { return &AddRoleToUserCommand{} }

func (c *AddRoleToUserCommand) AggregateID() uuid.UUID              { return c.AggID }
func (c *AddRoleToUserCommand) AggregateType() events.AggregateType { return UserRoleBinding }
func (c *AddRoleToUserCommand) CommandType() commands.CommandType   { return AddRoleToUser }
