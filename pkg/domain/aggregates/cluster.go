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

// ClusterAggregate is an aggregate for K8s Clusters.
type ClusterAggregate struct {
	DomainAggregateBase
	aggregateManager es.AggregateStore
	name             string
	label            string
	apiServerAddr    string
	caCertBundle     []byte
	bootstrapToken   string
}

// ClusterAggregate creates a new ClusterAggregate
func NewClusterAggregate(id uuid.UUID, aggregateManager es.AggregateStore) es.Aggregate {
	return &ClusterAggregate{
		DomainAggregateBase: DomainAggregateBase{
			BaseAggregate: es.NewBaseAggregate(aggregates.Cluster, id),
		},
		aggregateManager: aggregateManager,
	}
}

// HandleCommand implements the HandleCommand method of the Aggregate interface.
func (a *ClusterAggregate) HandleCommand(ctx context.Context, cmd es.Command) error {
	if err := a.Authorize(ctx, cmd); err != nil {
		return err
	}
	if err := a.validate(ctx, cmd); err != nil {
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
	default:
		return fmt.Errorf("couldn't handle command of type '%s'", cmd.CommandType())
	}
}

// validate validates the current state of the aggregate and if a specific command is valid in the current state
func (a *ClusterAggregate) validate(ctx context.Context, cmd es.Command) error {
	switch cmd := cmd.(type) {
	case *commands.CreateClusterCommand:
		if a.Exists() {
			return domainErrors.ErrClusterAlreadyExists
		}

		// Get all aggregates of same type
		aggregates, err := a.aggregateManager.All(ctx, a.Type())
		if err != nil {
			return err
		}

		if containsCluster(aggregates, cmd.GetName()) {
			return domainErrors.ErrClusterAlreadyExists
		}
		return nil
	default:
		return a.Validate(ctx, cmd)
	}
}

func containsCluster(values []es.Aggregate, name string) bool {
	for _, value := range values {
		d, ok := value.(*ClusterAggregate)
		if ok {
			if d.name == name {
				return true
			}
		}
	}
	return false
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
		err := event.Data().ToProto(data)
		if err != nil {
			return err
		}
		a.bootstrapToken = data.GetJWT()
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
