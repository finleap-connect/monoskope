package commands

import (
	"context"

	"github.com/google/uuid"
	cmd "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/commands/user"
	types "gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/constants"
	es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/event_sourcing"
	"google.golang.org/protobuf/types/known/anypb"
)

// AddRoleToUser is a command for adding a role to a user.
type CreateUserRoleBindingCommand struct {
	aggregateId uuid.UUID
	cmd.AddRoleToUserCommand
}

func (c *CreateUserRoleBindingCommand) AggregateID() uuid.UUID { return c.aggregateId }
func (c *CreateUserRoleBindingCommand) AggregateType() es.AggregateType {
	return types.UserRoleBindingType
}
func (c *CreateUserRoleBindingCommand) CommandType() es.CommandType {
	return types.CreateUserRoleBindingType
}
func (c *CreateUserRoleBindingCommand) SetData(a *anypb.Any) error {
	return a.UnmarshalTo(&c.AddRoleToUserCommand)
}
func (c *CreateUserRoleBindingCommand) Policies(ctx context.Context) []es.Policy {
	return []es.Policy{
		{Role: types.Admin, Scope: types.System},                                            // System admin
		{Role: types.Admin, Scope: types.Tenant, Resource: c.AddRoleToUserCommand.Resource}, // Tenant admin
	}
}
