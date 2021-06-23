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
	es.DefaultCommandRegistry.RegisterCommand(NewRequestClusterCertificateCommand)
}

// RequestClusterCertificateCommand is a command for creating a cluster.
type RequestClusterCertificateCommand struct {
	*es.BaseCommand
	cmdData.RequestClusterOperatorCertificate
}

// NewRequestClusterCertificateCommand creates a RequestClusterCertificateCommand.
func NewRequestClusterCertificateCommand(id uuid.UUID) es.Command {
	return &RequestClusterCertificateCommand{
		BaseCommand: es.NewBaseCommand(id, aggregates.Cluster, commands.RequestClusterCertificate),
	}
}

func (c *RequestClusterCertificateCommand) SetData(a *anypb.Any) error {
	return a.UnmarshalTo(&c.RequestClusterOperatorCertificate)
}

// Policies returns the Role/Scope/Resource combination allowed to execute.
func (c *RequestClusterCertificateCommand) Policies(ctx context.Context) []es.Policy {
	return []es.Policy{
		es.NewPolicy().WithRole(roles.Admin).WithScope(scopes.System),       // Allows system admins
		es.NewPolicy().WithRole(roles.K8sOperator).WithScope(scopes.System), // Allows k8s operators
	}
}
