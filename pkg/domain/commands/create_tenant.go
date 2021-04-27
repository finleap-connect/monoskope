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
	es.DefaultCommandRegistry.RegisterCommand(NewCreateTenantCommand)
}

// CreateTenantCommand is a command for creating a tenant.
type CreateTenantCommand struct {
	*es.BaseCommand
	cmdData.CreateTenantCommandData
}

// NewCreateTenantCommand creates a CreateTenantCommand.
func NewCreateTenantCommand(id uuid.UUID) es.Command {
	return &CreateTenantCommand{
		BaseCommand: es.NewBaseCommand(id, aggregates.Tenant, commands.CreateTenant),
	}
}

func (c *CreateTenantCommand) SetData(a *anypb.Any) error {
	if err := a.UnmarshalTo(&c.CreateTenantCommandData); err != nil {
		return err
	}
	return nil
}

// Policies returns the Role/Scope/Resource combination allowed to execute.
func (c *CreateTenantCommand) Policies(ctx context.Context) []es.Policy {
	return []es.Policy{
		es.NewPolicy().WithRole(roles.Admin).WithScope(scopes.System), // Allows system admins to create a tenant
	}
}
