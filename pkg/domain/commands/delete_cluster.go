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
	es.DefaultCommandRegistry.RegisterCommand(NewDeleteClusterCommand)
}

// DeleteClusterCommand is a command for deleting a cluster.
type DeleteClusterCommand struct {
	*es.BaseCommand
}

// NewDeleteClusterCommand creates a DeleteClusterCommand.
func NewDeleteClusterCommand(id uuid.UUID) es.Command {
	return &DeleteClusterCommand{
		BaseCommand: es.NewBaseCommand(id, aggregates.Cluster, commands.DeleteCluster),
	}
}

func (c *DeleteClusterCommand) SetData(a *anypb.Any) error {
	return nil
}

// Policies returns the Role/Scope/Resource combination allowed to execute.
func (c *DeleteClusterCommand) Policies(ctx context.Context) []es.Policy {
	return []es.Policy{
		es.NewPolicy().WithRole(roles.Admin).WithScope(scopes.System), // Allows system admins
	}
}
