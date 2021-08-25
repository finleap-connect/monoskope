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

func init() {
	es.DefaultCommandRegistry.RegisterCommand(NewDeleteUserCommand)
}

// DeleteUserCommand is a command for deleting a user.
type DeleteUserCommand struct {
	*es.BaseCommand
}

// NewDeleteUserCommand creates a DeleteUserCommand.
func NewDeleteUserCommand(id uuid.UUID) es.Command {
	return &DeleteUserCommand{
		BaseCommand: es.NewBaseCommand(id, aggregates.User, commands.DeleteUser),
	}
}

func (c *DeleteUserCommand) SetData(a *anypb.Any) error {
	return nil
}

// Policies returns the Role/Scope/Resource combination allowed to execute.
func (c *DeleteUserCommand) Policies(ctx context.Context) []es.Policy {
	return []es.Policy{
		es.NewPolicy().WithRole(roles.Admin).WithScope(scopes.System), // Allows system admins
	}
}
