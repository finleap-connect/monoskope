package projectors

import (
	"context"

	"github.com/google/uuid"
	ed "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/domain/eventdata"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/constants/events"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/projections"
	es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing/errors"
)

type tenantProjector struct {
	*DomainProjector
}

func NewTenantProjector() es.Projector {
	return &tenantProjector{
		DomainProjector: NewDomainProjector(),
	}
}

func (u *tenantProjector) NewProjection(id uuid.UUID) es.Projection {
	return projections.NewTenantProjection(id)
}

// Project updates the state of the projection according to the given event.
func (u *tenantProjector) Project(ctx context.Context, event es.Event, projection es.Projection) (es.Projection, error) {
	// Get the actual projection type
	p, ok := projection.(*projections.Tenant)
	if !ok {
		return nil, errors.ErrInvalidProjectionType
	}

	// Apply the changes for the event.
	switch event.EventType() {
	case events.TenantCreated:
		data := &ed.TenantCreatedEventData{}
		if err := event.Data().ToProto(data); err != nil {
			return projection, err
		}
		p.Name = data.GetName()
		p.Prefix = data.GetPrefix()

		if err := u.ProjectCreated(event, p.DomainProjection); err != nil {
			return nil, err
		}
	case events.TenantUpdated:
		data := &ed.TenantUpdatedEventData{}
		if err := event.Data().ToProto(data); err != nil {
			return projection, err
		}
		p.Name = data.GetName().GetValue()
		if err := u.ProjectModified(event, p.DomainProjection); err != nil {
			return nil, err
		}
	case events.TenantDeleted:
		if err := u.ProjectDeleted(event, p.DomainProjection); err != nil {
			return nil, err
		}
	default:
		return nil, errors.ErrInvalidEventType
	}

	p.IncrementVersion()

	return p, nil
}
