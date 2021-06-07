package projectors

import (
	"context"

	"github.com/google/uuid"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/domain/eventdata"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/constants/events"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/projections"
	es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing/errors"
)

type clusterProjector struct {
	*domainProjector
}

func NewClusterProjector() es.Projector {
	return &clusterProjector{
		domainProjector: NewDomainProjector(),
	}
}

func (u *clusterProjector) NewProjection(id uuid.UUID) es.Projection {
	return projections.NewCluster(id)
}

// Project updates the state of the projection according to the given event.
func (c *clusterProjector) Project(ctx context.Context, event es.Event, projection es.Projection) (es.Projection, error) {
	// Get the actual projection type
	p, ok := projection.(*projections.Cluster)
	if !ok {
		return nil, errors.ErrInvalidProjectionType
	}

	// Apply the changes for the event.
	switch event.EventType() {
	case events.ClusterBootstrapTokenCreated:
		data := &eventdata.ClusterBootstrapTokenCreated{}
		if err := event.Data().ToProto(data); err != nil {
			return projection, err
		}
		p.BootstrapToken = data.GetJWT()
		if err := c.projectModified(event, p.DomainProjection); err != nil {
			return nil, err
		}
	case events.ClusterCreated:
		data := &eventdata.ClusterCreated{}
		if err := event.Data().ToProto(data); err != nil {
			return projection, err
		}
		p.Name = data.GetName()
		p.Label = data.GetLabel()
		p.ApiServerAddress = data.GetApiServerAddress()
		p.ClusterCACertBundle = data.GetCaCertificateBundle()
		if err := c.projectCreated(event, p.DomainProjection); err != nil {
			return nil, err
		}
	default:
		return nil, errors.ErrInvalidEventType
	}

	p.IncrementVersion()

	return p, nil
}
