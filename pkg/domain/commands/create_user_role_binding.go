package commands

import (
	"context"

	"github.com/google/uuid"
	cmd "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/commands/user"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/constants/aggregates"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/constants/commands"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/constants/roles"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/constants/scopes"
	es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing"
	"google.golang.org/protobuf/types/known/anypb"
)

// AddRoleToUser is a command for adding a role to a user.
type CreateUserRoleBindingCommand struct {
	aggregateId uuid.UUID
	cmd.CreateUserRoleBindingCommand
}

func (c *CreateUserRoleBindingCommand) AggregateID() uuid.UUID { return c.aggregateId }
func (c *CreateUserRoleBindingCommand) AggregateType() es.AggregateType {
	return aggregates.UserRoleBinding
}
func (c *CreateUserRoleBindingCommand) CommandType() es.CommandType {
	return commands.CreateUserRoleBinding
}
func (c *CreateUserRoleBindingCommand) SetData(a *anypb.Any) error {
	return a.UnmarshalTo(&c.CreateUserRoleBindingCommand)
}
func (c *CreateUserRoleBindingCommand) Policies(ctx context.Context) []es.Policy {
	return []es.Policy{
		es.NewPolicy().WithRole(roles.Admin).WithScope(scopes.System),                          // System admin
		es.NewPolicy().WithRole(roles.Admin).WithScope(scopes.Tenant).WithResource(c.Resource), // Tenant admin
	}
}
