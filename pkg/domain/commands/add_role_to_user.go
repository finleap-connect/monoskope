package commands

import (
	"github.com/google/uuid"
	apicmd "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/commands/user"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain"
	. "gitlab.figo.systems/platform/monoskope/monoskope/pkg/event_sourcing"
	"google.golang.org/protobuf/types/known/anypb"
)

// AddRoleToUser is a command for adding a role to a user.
type AddRoleToUserCommand struct {
	aggregateId uuid.UUID
	apicmd.AddRoleToUserCommand
}

func (c *AddRoleToUserCommand) AggregateID() uuid.UUID       { return c.aggregateId }
func (c *AddRoleToUserCommand) AggregateType() AggregateType { return domain.UserRoleBinding }
func (c *AddRoleToUserCommand) CommandType() CommandType     { return domain.AddRoleToUser }

func (c *AddRoleToUserCommand) SetData(a *anypb.Any) error {
	return a.UnmarshalTo(&c.AddRoleToUserCommand)
}
