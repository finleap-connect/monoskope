package projectors

import (
	"context"
	"errors"

	"github.com/google/uuid"
	ed "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/domain/eventdata"
	projectionsApi "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/domain/projections"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/constants/events"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/projections"
	es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing"
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
	i, ok := projection.(*projections.UserRoleBinding)
	if !ok {
		return nil, errors.New("model is of incorrect type")
	}

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
	default:
		return nil, errors.New("could not handle event: " + event.String())
	}

	i.IncrementVersion()

	return i, nil
}
