package domain_commands

import (
	"github.com/google/uuid"
	apicmd "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/commands/user"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/commands"
	domain_aggregates "gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/aggregates"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/events"
)

func init() {
	if err := commands.Registry.RegisterCommand(addRoleToUserCommandFactory); err != nil {
		panic(err)
	}
}

const (
	AddRoleToUser commands.CommandType = "AddRoleToUser"
)

// AddRoleToUser is a command for adding a role to a user.
type AddRoleToUserCommand struct {
	AggID uuid.UUID
	apicmd.AddRoleToUser
}

var addRoleToUserCommandFactory = func() commands.Command { return &AddRoleToUserCommand{} }

func (c *AddRoleToUserCommand) AggregateID() uuid.UUID { return c.AggID }
func (c *AddRoleToUserCommand) AggregateType() events.AggregateType {
	return domain_aggregates.UserRoleBinding
}
func (c *AddRoleToUserCommand) CommandType() commands.CommandType { return AddRoleToUser }
