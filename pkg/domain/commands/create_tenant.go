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

// CreateTenantCommand is a command for creating a tenant.
type CreateTenantCommand struct {
	*es.BaseCommand
	cmdData.CreateTenantCommandData
}

// NewCreateTenantCommand creates a CreateTenantCommand.
func NewCreateTenantCommand() *CreateTenantCommand {
	return &CreateTenantCommand{
		BaseCommand: es.NewBaseCommand(aggregates.Tenant, commands.CreateTenant),
	}
}

func (c *CreateTenantCommand) SetData(a *anypb.Any) error {
	return a.UnmarshalTo(&c.CreateTenantCommandData)
}

// Policies returns the Role/Scope/Resource combination allowed to execute.
func (c *CreateTenantCommand) Policies(ctx context.Context) []es.Policy {
	return []es.Policy{
		es.NewPolicy().WithRole(roles.Admin).WithScope(scopes.System), // Allows system admins to create a tenant
	}
}
