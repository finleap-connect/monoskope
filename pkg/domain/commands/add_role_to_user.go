package commands

import (
	"github.com/google/uuid"
	apicmd "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/commands/user"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain"
	. "gitlab.figo.systems/platform/monoskope/monoskope/pkg/event_sourcing"
	"google.golang.org/protobuf/types/known/anypb"
)

// AddRoleToUser is a command for adding a role to a user.
type CreateUserRoleBindingCommand struct {
	aggregateId uuid.UUID
	apicmd.AddRoleToUserCommand
}

func (c *CreateUserRoleBindingCommand) AggregateID() uuid.UUID       { return c.aggregateId }
func (c *CreateUserRoleBindingCommand) AggregateType() AggregateType { return domain.UserRoleBinding }
func (c *CreateUserRoleBindingCommand) CommandType() CommandType     { return domain.CreateUserRoleBinding }

func (c *CreateUserRoleBindingCommand) SetData(a *anypb.Any) error {
	return a.UnmarshalTo(&c.AddRoleToUserCommand)
}
