package commands

import (
	"context"

	"github.com/google/uuid"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/constants/aggregates"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/constants/commands"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/constants/roles"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/constants/scopes"
	es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing"
	"google.golang.org/protobuf/types/known/anypb"
)

// DeleteUserRoleBindingCommand is a command for removing a role from a user.
type DeleteUserRoleBindingCommand struct {
	aggregateId uuid.UUID
}

func NewDeleteUserRoleBindingCommand() *DeleteUserRoleBindingCommand {
	return &DeleteUserRoleBindingCommand{
		aggregateId: uuid.New(),
	}
}

func (c *DeleteUserRoleBindingCommand) AggregateID() uuid.UUID { return c.aggregateId }
func (c *DeleteUserRoleBindingCommand) AggregateType() es.AggregateType {
	return aggregates.UserRoleBinding
}
func (c *DeleteUserRoleBindingCommand) CommandType() es.CommandType {
	return commands.DeleteUserRoleBinding
}
func (c *DeleteUserRoleBindingCommand) SetData(a *anypb.Any) error {
	return nil
}
func (c *DeleteUserRoleBindingCommand) Policies(ctx context.Context) []es.Policy {
	return []es.Policy{
		es.NewPolicy().WithRole(roles.Admin).WithScope(scopes.System),                         // System admin
		es.NewPolicy().WithRole(roles.Admin).WithScope(scopes.Tenant).WithResourceMatch(true), // Tenant admin
	}
}
