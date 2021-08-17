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
	es.DefaultCommandRegistry.RegisterCommand(NewUpdateClusterCommand)
}

// UpdateClusterCommand is a command for updating a cluster.
type UpdateClusterCommand struct {
	*es.BaseCommand
	cmdData.UpdateCluster
}

// NewUpdateClusterCommand creates an UpdateClusterCommand.
func NewUpdateClusterCommand(id uuid.UUID) es.Command {
	return &UpdateClusterCommand{
		BaseCommand: es.NewBaseCommand(id, aggregates.Cluster, commands.UpdateCluster),
	}
}

func (c *UpdateClusterCommand) SetData(a *anypb.Any) error {
	return a.UnmarshalTo(&c.UpdateCluster)
}

// Policies returns the Role/Scope/Resource combination allowed to execute.
func (c *UpdateClusterCommand) Policies(ctx context.Context) []es.Policy {
	return []es.Policy{
		es.NewPolicy().WithRole(roles.Admin).WithScope(scopes.System), // Allows system admins to update a cluster
	}
}
