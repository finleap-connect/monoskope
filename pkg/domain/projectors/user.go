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

type userProjector struct {
}

func NewUserProjector() es.Projector {
	return &userProjector{}
}

func (u *userProjector) NewProjection(id uuid.UUID) es.Projection {
	return &projections.User{
		User: &projectionsApi.User{
			Id: id.String(),
		},
	}
}

// Project updates the state of the projection occording to the given event.
func (u *userProjector) Project(ctx context.Context, event es.Event, projection es.Projection) (es.Projection, error) {
	// Get the actual projection type
	i, ok := projection.(*projections.User)
	if !ok {
		return nil, errors.ErrInvalidProjectionType
	}

	// Apply the changes for the event.
	switch event.EventType() {
	case events.UserCreated:
		data := &ed.UserCreatedEventData{}
		if err := event.Data().ToProto(data); err != nil {
			return projection, err
		}

		i.Id = event.AggregateID().String()
		i.Email = data.GetEmail()
		i.Name = data.GetName()
	default:
		return nil, errors.ErrInvalidEventType
	}

	i.IncrementVersion()

	return i, nil
}
