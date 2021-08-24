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

type tenantProjector struct {
	*domainProjector
}

func NewTenantProjector() es.Projector {
	return &tenantProjector{
		domainProjector: NewDomainProjector(),
	}
}

func (t *tenantProjector) NewProjection(id uuid.UUID) es.Projection {
	return projections.NewTenantProjection(id)
}

// Project updates the state of the projection according to the given event.
func (t *tenantProjector) Project(ctx context.Context, event es.Event, projection es.Projection) (es.Projection, error) {
	// Get the actual projection type
	p, ok := projection.(*projections.Tenant)
	if !ok {
		return nil, errors.ErrInvalidProjectionType
	}

	// Apply the changes for the event.
	switch event.EventType() {
	case events.TenantCreated:
		data := new(eventdata.TenantCreated)
		if err := event.Data().ToProto(data); err != nil {
			return projection, err
		}

		p.Name = data.GetName()
		p.Prefix = data.GetPrefix()

		if err := t.projectCreated(event, p.DomainProjection); err != nil {
			return nil, err
		}
	case events.TenantUpdated:
		data := new(eventdata.TenantUpdated)
		if err := event.Data().ToProto(data); err != nil {
			return projection, err
		}
		p.Name = data.GetName().GetValue()
	case events.TenantDeleted:
		if err := t.projectDeleted(event, p.DomainProjection); err != nil {
			return nil, err
		}
	default:
		return nil, errors.ErrInvalidEventType
	}

	if err := t.projectModified(event, p.DomainProjection); err != nil {
		return nil, err
	}
	p.IncrementVersion()

	return p, nil
}
