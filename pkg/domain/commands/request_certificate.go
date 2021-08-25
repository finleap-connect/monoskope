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
	es.DefaultCommandRegistry.RegisterCommand(NewRequestCertificateCommand)
}

// RequestCertificateCommand is a command for requesting a certificate for a given aggregate.
type RequestCertificateCommand struct {
	*es.BaseCommand
	cmdData.RequestCertificate
}

// NewRequestCertificateCommand creates a RequestCertificateCommand.
func NewRequestCertificateCommand(id uuid.UUID) es.Command {
	return &RequestCertificateCommand{
		BaseCommand: es.NewBaseCommand(id, aggregates.Certificate, commands.RequestCertificate),
	}
}

func (c *RequestCertificateCommand) SetData(a *anypb.Any) error {
	return a.UnmarshalTo(&c.RequestCertificate)
}

// Policies returns the Role/Scope/Resource combination allowed to execute.
func (c *RequestCertificateCommand) Policies(ctx context.Context) []es.Policy {
	return []es.Policy{
		es.NewPolicy().WithRole(roles.Admin).WithScope(scopes.System),       // Allows system admins
		es.NewPolicy().WithRole(roles.K8sOperator).WithScope(scopes.System), // Allows k8s operators
	}
}
