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

// ApproveClusterRegistration is a command for approval of new cluster to register with the m8 Control Plane.
type ApproveClusterRegistration struct {
	*es.BaseCommand
	cmdData.RequestClusterRegistration
}

func NewApproveClusterRegistration(id uuid.UUID) es.Command {
	return &ApproveClusterRegistration{
		BaseCommand: es.NewBaseCommand(id, aggregates.User, commands.ApproveClusterRegistration),
	}
}

func (c *ApproveClusterRegistration) SetData(a *anypb.Any) error {
	return a.UnmarshalTo(&c.RequestClusterRegistration)
}

func (c *ApproveClusterRegistration) Policies(ctx context.Context) []es.Policy {
	return []es.Policy{
		es.NewPolicy().WithRole(roles.Admin).WithScope(scopes.System), // Allows system admins
	}
}
