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
)

type userRoleBindingProjector struct {
}

func NewUserRoleBindingProjector() es.Projector {
	return &userRoleBindingProjector{}
}

func (u *userRoleBindingProjector) NewProjection(id uuid.UUID) es.Projection {
	return &projections.UserRoleBinding{
		UserRoleBinding: projectionsApi.UserRoleBinding{
			Id: id.String(),
		},
	}
}

// Project updates the state of the projection occording to the given event.
func (u *userRoleBindingProjector) Project(ctx context.Context, event es.Event, projection es.Projection) (es.Projection, error) {
	// Get the actual projection type
	i, ok := projection.(*projections.UserRoleBinding)
	if !ok {
		return nil, errors.ErrInvalidProjectionType
	}

	// Get UserID from event metadata
	// userId := event.Metadata()[gateway.HeaderAuthId]

	// Apply the changes for the event.
	switch event.EventType() {
	case events.UserRoleBindingCreated:
		switch event.EventType() {
		case events.UserRoleBindingCreated:
			data := &ed.UserRoleAddedEventData{}
			if err := event.Data().ToProto(data); err != nil {
				return projection, err
			}

			i.UserId = data.GetUserId()
			i.Role = data.GetRole()
			i.Scope = data.GetScope()
			i.Resource = data.GetResource()
		}
	case events.UserRoleBindingDeleted:
		// i.Deleted = timestamp.Now()
		// i.SetDeletedByID(userId)
	default:
		return nil, errors.ErrInvalidEventType
	}

	i.IncrementVersion()

	return i, nil
}
