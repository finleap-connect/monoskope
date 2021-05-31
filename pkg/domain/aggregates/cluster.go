package aggregates

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/domain/eventdata"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/commands"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/constants/aggregates"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/constants/events"
	es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing"
)

// ClusterAggregate is an aggregate for K8s Clusters.
type ClusterAggregate struct {
	DomainAggregateBase
	name                      string
	label                     string
	apiServerAddr             string
	jwt                       string
	caCertBundle              []byte
	certificateSigningRequest []byte
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
	case *commands.CreateClusterCommand:
		ed := es.ToEventDataFromProto(&eventdata.ClusterCreated{
			Name:                cmd.GetName(),
			Label:               cmd.GetLabel(),
			ApiServerAddress:    cmd.GetApiServerAddress(),
			CaCertificateBundle: cmd.GetClusterCACertBundle(),
		})
		_ = a.AppendEvent(ctx, events.ClusterCreated, ed)
		return nil
	case *commands.DeleteClusterCommand:
		_ = a.AppendEvent(ctx, events.ClusterDeleted, nil)
		return nil
	case *commands.RequestClusterCertificateCommand:
		_ = a.AppendEvent(ctx, events.ClusterCertificateRequested, nil)
		return nil
	default:
		return fmt.Errorf("couldn't handle command of type '%s'", cmd.CommandType())
	}
}

// ApplyEvent implements the ApplyEvent method of the Aggregate interface.
func (a *ClusterAggregate) ApplyEvent(event es.Event) error {
	switch event.EventType() {
	case events.ClusterCreated:
		data := &eventdata.ClusterCreated{}
		err := event.Data().ToProto(data)
		if err != nil {
			return err
		}
		a.name = data.GetName()
		a.label = data.GetLabel()
		a.apiServerAddr = data.GetApiServerAddress()
		a.caCertBundle = data.GetCaCertificateBundle()
	case events.ClusterBootstrapTokenCreated:
		data := &eventdata.ClusterBootstrapTokenCreated{}
		a.jwt = data.GetJWT()
	case events.ClusterCertificateRequested:
		data := &eventdata.ClusterCertificateRequested{}
		a.certificateSigningRequest = data.GetCertificateSigningRequest()
	case events.ClusterDeleted:
		a.SetDeleted(true)
	default:
		return fmt.Errorf("couldn't handle event of type '%s'", event.EventType())
	}
	return nil
}

func (a *ClusterAggregate) GetName() string {
	return a.name
}
