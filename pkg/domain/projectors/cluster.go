package projectors

import (
	"bytes"
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
	return projections.NewClusterProjection(id)
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
	case events.ClusterCreated:
		data := new(eventdata.ClusterCreated)
		if err := event.Data().ToProto(data); err != nil {
			return nil, err
		}

		p.DisplayName = data.GetName()
		p.Name = data.GetLabel()
		p.ApiServerAddress = data.GetApiServerAddress()
		p.CaCertBundle = data.GetCaCertificateBundle()

		if err := c.projectCreated(event, p.DomainProjection); err != nil {
			return nil, err
		}
	case events.ClusterCreatedV2:
		data := new(eventdata.ClusterCreatedV2)
		if err := event.Data().ToProto(data); err != nil {
			return nil, err
		}

		p.DisplayName = data.GetDisplayName()
		p.Name = data.GetName()
		p.ApiServerAddress = data.GetApiServerAddress()
		p.CaCertBundle = data.GetCaCertificateBundle()

		if err := c.projectCreated(event, p.DomainProjection); err != nil {
			return nil, err
		}
	case events.ClusterUpdated:
		data := new(eventdata.ClusterUpdated)
		if err := event.Data().ToProto(data); err != nil {
			return nil, err
		}
		if len(data.GetDisplayName()) > 0 && p.DisplayName != data.GetDisplayName() {
			p.DisplayName = data.GetDisplayName()
		}
		if len(data.GetApiServerAddress()) > 0 && p.ApiServerAddress != data.GetApiServerAddress() {
			p.ApiServerAddress = data.GetApiServerAddress()
		}
		if len(data.GetCaCertificateBundle()) > 0 && bytes.Equal(p.CaCertBundle, data.GetCaCertificateBundle()) {
			p.CaCertBundle = data.GetCaCertificateBundle()
		}
	case events.ClusterBootstrapTokenCreated:
		data := new(eventdata.ClusterBootstrapTokenCreated)
		if err := event.Data().ToProto(data); err != nil {
			return nil, err
		}
		p.BootstrapToken = data.GetJwt()
	case events.ClusterDeleted:
		if err := c.projectDeleted(event, p.DomainProjection); err != nil {
			return nil, err
		}
	default:
		return nil, errors.ErrInvalidEventType
	}

	if err := c.projectModified(event, p.DomainProjection); err != nil {
		return nil, err
	}
	p.IncrementVersion()

	return p, nil
}
