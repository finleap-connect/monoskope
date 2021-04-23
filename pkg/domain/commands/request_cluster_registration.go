package commands

import (
	"context"

	"github.com/google/uuid"
	cmdData "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/domain/commanddata"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/constants/aggregates"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/constants/commands"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/constants/roles"
	es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing"
	"google.golang.org/protobuf/types/known/anypb"
)

func init() {
	es.DefaultCommandRegistry.RegisterCommand(NewRequestClusterRegistrationCommand)
}

// RequestClusterRegistration is a command for registering a new cluster with the m8 Control Plane.
type RequestClusterRegistration struct {
	*es.BaseCommand
	cmdData.RequestClusterRegistration
}

func NewRequestClusterRegistrationCommand(id uuid.UUID) es.Command {
	return &RequestClusterRegistration{
		BaseCommand: es.NewBaseCommand(id, aggregates.ClusterRegistration, commands.RequestClusterRegistration),
	}
}

func (c *RequestClusterRegistration) SetData(a *anypb.Any) error {
	return a.UnmarshalTo(&c.RequestClusterRegistration)
}

func (c *RequestClusterRegistration) Policies(ctx context.Context) []es.Policy {
	return []es.Policy{
		es.NewPolicy().WithRole(roles.K8sOperator), // Allows for K8sOperators
	}
}
