package aggregates

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/domain/eventdata"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/commands"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/constants/aggregates"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/constants/events"
	domainErrors "gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/errors"
	es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing"
)

type ClusterRegistrationState string

const (
	Requested ClusterRegistrationState = "Requested"
	Approved  ClusterRegistrationState = "Approved"
	Denied    ClusterRegistrationState = "Denied"
)

// ClusterRegistrationAggregate is an aggregate for the registration flow of a K8s Clusters.
type ClusterRegistrationAggregate struct {
	DomainAggregateBase
	name          string
	apiServerAddr string
	caCertificate []byte
	state         ClusterRegistrationState
}

// NewClusterRegistrationAggregate creates a new ClusterRegistrationAggregate
func NewClusterRegistrationAggregate(id uuid.UUID) es.Aggregate {
	return &ClusterAggregate{
		DomainAggregateBase: DomainAggregateBase{
			BaseAggregate: es.NewBaseAggregate(aggregates.ClusterRegistration, id),
		},
	}
}

// validate validates the current state of the aggregate and if a specific command is valid in the current state
func (a *ClusterRegistrationAggregate) validate(ctx context.Context, cmd es.Command) error {
	switch cmd := cmd.(type) {
	case *commands.RequestClusterRegistration:
		if a.Exists() {
			return domainErrors.ErrAggregateAlreadyExists
		}
		return nil
	case *commands.ApproveClusterRegistration:
		if a.state != Requested {
			return domainErrors.ErrFailedPrecondition(fmt.Sprintf("The request is in state '%s' and can't be approved", a.state))
		}
		return nil
	case *commands.DenyClusterRegistration:
		if a.state != Requested {
			return domainErrors.ErrFailedPrecondition(fmt.Sprintf("The request is in state '%s' and can't be denied", a.state))
		}
		return nil
	default:
		return a.Validate(ctx, cmd)
	}
}

// HandleCommand implements the HandleCommand method of the Aggregate interface.
func (a *ClusterRegistrationAggregate) HandleCommand(ctx context.Context, cmd es.Command) error {
	if err := a.Authorize(ctx, cmd); err != nil {
		return err
	}
	if err := a.validate(ctx, cmd); err != nil {
		return err
	}

	switch cmd := cmd.(type) {
	case *commands.RequestClusterRegistration:
		ed := es.ToEventDataFromProto(&eventdata.ClusterRegistered{
			Name:             cmd.GetName(),
			ApiServerAddress: cmd.GetApiServerAddress(),
			CaCertificate:    cmd.GetClusterCACert(),
		})
		_ = a.AppendEvent(ctx, events.ClusterRegistrationRequested, ed)
		return nil
	case *commands.ApproveClusterRegistration:
		ed := es.ToEventDataFromProto(&eventdata.ClusterRegistered{
			Name:             a.name,
			ApiServerAddress: a.apiServerAddr,
			CaCertificate:    a.caCertificate,
		})
		_ = a.AppendEvent(ctx, events.ClusterRegistrationApproved, ed)
		return nil
	case *commands.DenyClusterRegistration:
		_ = a.AppendEvent(ctx, events.ClusterRegistrationDenied, nil)
		return nil
	default:
		return fmt.Errorf("couldn't handle command of type '%s'", cmd.CommandType())
	}
}

// ApplyEvent implements the ApplyEvent method of the Aggregate interface.
func (a *ClusterRegistrationAggregate) ApplyEvent(event es.Event) error {
	switch event.EventType() {
	case events.ClusterRegistrationRequested:
		data := &eventdata.ClusterRegistered{}
		err := event.Data().ToProto(data)
		if err != nil {
			return err
		}
		a.state = Requested
		a.name = data.GetName()
		a.apiServerAddr = data.GetApiServerAddress()
		a.caCertificate = data.GetCaCertificate()
	case events.ClusterRegistrationApproved:
		a.state = Approved
	case events.ClusterRegistrationDenied:
		a.state = Denied
	default:
		return fmt.Errorf("couldn't handle event of type '%s'", event.EventType())
	}
	return nil
}
