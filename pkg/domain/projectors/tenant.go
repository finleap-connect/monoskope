package projectors

import (
	"context"

	"github.com/google/uuid"
	ed "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/domain/eventdata"
	projectionsApi "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/domain/projections"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/constants/events"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/projections"
	es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing/errors"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type tenantProjector struct {
}

func NewTenantProjector() es.Projector {
	return &tenantProjector{}
}

func (u *tenantProjector) NewProjection(id uuid.UUID) es.Projection {
	return &projections.Tenant{
		Tenant: &projectionsApi.Tenant{
			Id: id.String(),
		},
	}
}

// Project updates the state of the projection according to the given event.
func (u *tenantProjector) Project(ctx context.Context, event es.Event, projection es.Projection) (es.Projection, error) {
	// Get the actual projection type
	i, ok := projection.(*projections.Tenant)
	if !ok {
		return nil, errors.ErrInvalidProjectionType
	}

	userId := event.Metadata()["user_id"]

	// Apply the changes for the event.
	switch event.EventType() {
	case events.TenantCreated:
		data := &ed.TenantCreatedEventData{}
		if err := event.Data().ToProto(data); err != nil {
			return projection, err
		}
		i.Name = data.GetName()
		i.Prefix = data.GetPrefix()
		i.Created = timestamppb.Now()
		i.LastModified = timestamppb.Now()
		i.SetCreatedByID(userId)
		i.SetLastModifiedByID(userId)
	case events.TenantUpdated:
		data := &ed.TenantUpdatedEventData{}
		if err := event.Data().ToProto(data); err != nil {
			return projection, err
		}
		i.Name = data.GetUpdate().GetName().Value
		i.LastModified = timestamppb.Now()
		i.SetLastModifiedByID(userId)
	case events.TenantDeleted:
		i.Deleted = timestamppb.Now()
		i.SetDeletedByID(userId)
	default:
		return nil, errors.ErrInvalidEventType
	}

	i.IncrementVersion()

	return i, nil
}
