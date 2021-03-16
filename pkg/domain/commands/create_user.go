package commands

import (
	"context"

	cmdData "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/domain/commanddata"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/constants/aggregates"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/constants/commands"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/constants/roles"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/constants/scopes"
	es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing"
	"google.golang.org/protobuf/types/known/anypb"
)

// CreateUserCommand is a command for creating a user.
type CreateUserCommand struct {
	*es.BaseCommand
	cmdData.CreateUserCommandData
}

func NewCreateUserCommand() *CreateUserCommand {
	return &CreateUserCommand{
		BaseCommand: es.NewBaseCommand(aggregates.User, commands.CreateUser),
	}
}

func (c *CreateUserCommand) SetData(a *anypb.Any) error {
	return a.UnmarshalTo(&c.CreateUserCommandData)
}

func (c *CreateUserCommand) Policies(ctx context.Context) []es.Policy {
	return []es.Policy{
		es.NewPolicy().WithSubject(c.GetEmail()),                      // Allows user to create themselves
		es.NewPolicy().WithRole(roles.Admin).WithScope(scopes.System), // Allows system admins to create any users
	}
}
