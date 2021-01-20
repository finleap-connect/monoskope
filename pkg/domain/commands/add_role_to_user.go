package domain_commands

import (
	"github.com/google/uuid"
	apicmd "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/commands/user"
	domain_aggregates "gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/aggregates"
	evs "gitlab.figo.systems/platform/monoskope/monoskope/pkg/event_sourcing"
)

func init() {
	if err := evs.Registry.RegisterCommand(addRoleToUserCommandFactory); err != nil {
		panic(err)
	}
}

const (
	AddRoleToUser evs.CommandType = "AddRoleToUser"
)

// AddRoleToUser is a command for adding a role to a user.
type AddRoleToUserCommand struct {
	AggID uuid.UUID
	apicmd.AddRoleToUser
}

var addRoleToUserCommandFactory = func() evs.Command { return &AddRoleToUserCommand{} }

func (c *AddRoleToUserCommand) AggregateID() uuid.UUID { return c.AggID }
func (c *AddRoleToUserCommand) AggregateType() evs.AggregateType {
	return domain_aggregates.UserRoleBinding
}
func (c *AddRoleToUserCommand) CommandType() evs.CommandType { return AddRoleToUser }
