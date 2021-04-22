package projectors

import (
	"context"

	"github.com/google/uuid"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/domain/eventdata"
	projectionsApi "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/domain/projections"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/constants/events"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/projections"
	es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing/errors"
)

type clusterregistrationProjector struct {
	*approvableProjector
}

func NewClusterRegistrationProjector() es.Projector {
	return &clusterregistrationProjector{
		approvableProjector: NewApprovableProjector(),
	}
}

func (t *clusterregistrationProjector) NewProjection(id uuid.UUID) es.Projection {
	return projections.NewClusterRegistrationProjection(id)
}

// Project updates the state of the projection according to the given event.
func (t *clusterregistrationProjector) Project(ctx context.Context, event es.Event, projection es.Projection) (es.Projection, error) {
	// Get the actual projection type
	p, ok := projection.(*projections.ClusterRegistration)
	if !ok {
		return nil, errors.ErrInvalidProjectionType
	}

	// Apply the changes for the event.
	switch event.EventType() {
	case events.ClusterRegistrationRequested:
		data := &eventdata.ClusterRegistered{}
		if err := event.Data().ToProto(data); err != nil {
			return projection, err
		}
		p.Name = data.GetName()
		p.ApiServerAddress = data.GetApiServerAddress()
		p.ClusterCACert = data.GetCaCertificate()
		p.Status = projectionsApi.ClusterRegistration_REQUESTED

		if err := t.projectCreated(event, p.DomainProjection); err != nil {
			return nil, err
		}
	case events.ClusterRegistrationApproved:
		p.Status = projectionsApi.ClusterRegistration_APPROVED
		if err := t.projectApproved(event, p.ApprovableProjection); err != nil {
			return nil, err
		}
	case events.ClusterRegistrationDenied:
		p.Status = projectionsApi.ClusterRegistration_DENIED
		if err := t.projectDenied(event, p.ApprovableProjection); err != nil {
			return nil, err
		}
	default:
		return nil, errors.ErrInvalidEventType
	}

	p.IncrementVersion()

	return p, nil
}
