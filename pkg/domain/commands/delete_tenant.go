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

// DeleteTenantCommand is a command for deleting a tenant.
type DeleteTenantCommand struct {
	*es.BaseCommand
}

// NewDeleteTenantCommand creates a DeleteTenantCommand.
func NewDeleteTenantCommand(id uuid.UUID) es.Command {
	return &DeleteTenantCommand{
		BaseCommand: es.NewBaseCommand(id, aggregates.Tenant, commands.DeleteTenant),
	}
}

func (c *DeleteTenantCommand) SetData(a *anypb.Any) error {
	return nil
}

// Policies returns the Role/Scope/Resource combination allowed to execute.
func (c *DeleteTenantCommand) Policies(ctx context.Context) []es.Policy {
	return []es.Policy{
		es.NewPolicy().WithRole(roles.Admin).WithScope(scopes.System), // Allows system admins to delete a tenant
	}
}
