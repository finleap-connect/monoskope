package aggregates

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	ed "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/domain/eventdata"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/commands"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/constants/aggregates"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/constants/events"
	es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing"
)

type ClusterState string

const (
	Requested ClusterState = "Requested"
	Approved  ClusterState = "Approved"
	Denied    ClusterState = "Denied"
)

// ClusterAggregate is an aggregate for K8s Clusters.
type ClusterAggregate struct {
	DomainAggregateBase
	name          string
	apiServerAddr string
	caCertificate []byte
	state         ClusterState
}

// ClusterAggregate creates a new ClusterAggregate
func NewClusterAggregate(id uuid.UUID) es.Aggregate {
	return &ClusterAggregate{
		DomainAggregateBase: DomainAggregateBase{
			BaseAggregate: es.NewBaseAggregate(aggregates.Cluster, id),
		},
	}
}

// HandleCommand implements the HandleCommand method of the Aggregate interface.
func (a *ClusterAggregate) HandleCommand(ctx context.Context, cmd es.Command) error {
	if err := a.Authorize(ctx, cmd); err != nil {
		return err
	}
	if err := a.Validate(ctx, cmd); err != nil {
		return err
	}

	switch cmd := cmd.(type) {
	case *commands.RequestClusterRegistration:
		_ = a.AppendEvent(ctx, events.ClusterRegistrationRequested, nil)
		return nil
	case *commands.DeleteClusterCommand:
		_ = a.AppendEvent(ctx, events.ClusterDeleted, nil)
		return nil
	default:
		return fmt.Errorf("couldn't handle command of type '%s'", cmd.CommandType())
	}
}

// ApplyEvent implements the ApplyEvent method of the Aggregate interface.
func (a *ClusterAggregate) ApplyEvent(event es.Event) error {
	switch event.EventType() {
	case events.ClusterRegistrationRequested:
		data := &ed.ClusterRegisteredEventData{}
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
	case events.ClusterDeleted:
		a.SetDeleted(true)
	default:
		return fmt.Errorf("couldn't handle event of type '%s'", event.EventType())
	}
	return nil
}
