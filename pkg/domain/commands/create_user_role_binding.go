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
	superuserPolicies []es.Policy
}

func NewCreateUserRoleBindingCommand(id uuid.UUID) *CreateUserRoleBindingCommand {
	return &CreateUserRoleBindingCommand{
		BaseCommand:       es.NewBaseCommand(id, aggregates.UserRoleBinding, commands.CreateUserRoleBinding),
		superuserPolicies: []es.Policy{},
	}
}

func (c *CreateUserRoleBindingCommand) SetData(a *anypb.Any) error {
	return a.UnmarshalTo(&c.CreateUserRoleBindingCommandData)
}
func (c *CreateUserRoleBindingCommand) Policies(ctx context.Context) []es.Policy {
	return append(
		c.superuserPolicies,
		[]es.Policy{
			es.NewPolicy().WithRole(roles.Admin).WithScope(scopes.System),                          // System admin
			es.NewPolicy().WithRole(roles.Admin).WithScope(scopes.Tenant).WithResource(c.Resource), // Tenant admin
		}...,
	)
}

func (c *CreateUserRoleBindingCommand) DeclareSuperusers(users []string) {
	for _, user := range users {
		c.superuserPolicies = append(c.superuserPolicies, es.NewPolicy().WithSubject(user))
	}
}
