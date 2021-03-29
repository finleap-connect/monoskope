package commands

import (
	"context"

	"github.com/google/uuid"
	cmdData "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/domain/commanddata"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/constants/aggregates"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/constants/commands"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/constants/roles"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/constants/scopes"
	es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing"
	"google.golang.org/protobuf/types/known/anypb"
)

// AddRoleToUser is a command for adding a role to a user.
type CreateUserRoleBindingCommand struct {
	*es.BaseCommand
	cmdData.CreateUserRoleBindingCommandData
}

func NewCreateUserRoleBindingCommand(id uuid.UUID) es.Command {
	return &CreateUserRoleBindingCommand{
		BaseCommand: es.NewBaseCommand(id, aggregates.UserRoleBinding, commands.CreateUserRoleBinding),
	}
}
func (c *CreateUserRoleBindingCommand) SetData(a *anypb.Any) error {
	return a.UnmarshalTo(&c.CreateUserRoleBindingCommandData)
}
func (c *CreateUserRoleBindingCommand) Policies(ctx context.Context) []es.Policy {
	return []es.Policy{
		es.NewPolicy().WithRole(roles.Admin).WithScope(scopes.System),                         // System admin
		es.NewPolicy().WithRole(roles.Admin).WithScope(scopes.Tenant).WithResourceMatch(true), // Tenant admin
	}
}
