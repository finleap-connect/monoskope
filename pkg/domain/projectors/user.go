package projectors

import (
	"context"

	ed "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/eventdata/user"
	projections "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/queryhandler"
	aggregates "gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/constants/aggregates"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/constants/events"
	es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing/errors"
)

type userProjector struct {
}

func NewUserProjector() es.Projector {
	return &userProjector{}
}

// AggregateType returns the AggregateType for which events should be projected.
func (u *userProjector) AggregateType() es.AggregateType {
	return aggregates.User
}

func (u *userProjector) NewProjection() es.Projection {
	return &projections.User{}
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

		i.Email = data.GetEmail()
		i.Name = data.GetName()
	default:
		return nil, errors.ErrInvalidEventType
	}

	i.AggregateVersion++
	return i, nil
}
