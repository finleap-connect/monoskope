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

// DeleteTenantCommand is a command for deleting a tenant.
type DeleteTenantCommand struct {
	*es.BaseCommand
	cmdData.DeleteTenantCommandData
}

// NewDeleteTenantCommand creates a DeleteTenantCommand.
func NewDeleteTenantCommand(id uuid.UUID) *DeleteTenantCommand {
	return &DeleteTenantCommand{
		BaseCommand: es.NewBaseCommand(id, aggregates.Tenant, commands.DeleteTenant),
	}
}

func (c *DeleteTenantCommand) SetData(a *anypb.Any) error {
	return a.UnmarshalTo(&c.DeleteTenantCommandData)
}

// Policies returns the Role/Scope/Resource combination allowed to execute.
func (c *DeleteTenantCommand) Policies(ctx context.Context) []es.Policy {
	return []es.Policy{
		es.NewPolicy().WithRole(roles.Admin).WithScope(scopes.System), // Allows system admins to delete a tenant
	}
}
