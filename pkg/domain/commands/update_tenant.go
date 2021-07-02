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

func init() {
	es.DefaultCommandRegistry.RegisterCommand(NewUpdateTenantCommand)
}

// UpdateTenantCommand is a command for updating a tenant.
type UpdateTenantCommand struct {
	*es.BaseCommand
	cmdData.UpdateTenantCommandData
}

// NewUpdateTenantCommand creates an UpdateTenantCommand.
func NewUpdateTenantCommand(id uuid.UUID) es.Command {
	return &UpdateTenantCommand{
		BaseCommand: es.NewBaseCommand(id, aggregates.Tenant, commands.UpdateTenant),
	}
}

func (c *UpdateTenantCommand) SetData(a *anypb.Any) error {
	return a.UnmarshalTo(&c.UpdateTenantCommandData)
}

// Policies returns the Role/Scope/Resource combination allowed to execute.
func (c *UpdateTenantCommand) Policies(ctx context.Context) []es.Policy {
	return []es.Policy{
		es.NewPolicy().WithRole(roles.Admin).WithScope(scopes.System), // Allows system admins to update a tenant
	}
}
