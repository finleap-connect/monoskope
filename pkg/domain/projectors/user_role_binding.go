package projectors

import (
	"context"
	"errors"

	ed "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/eventdata/user"
	projections "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/queryhandler"
	aggregates "gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/constants/aggregates"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/constants/events"
	es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/event_sourcing"
)

type userRoleBindingProjector struct {
}

func NewUserRoleBindingProjector() es.Projector {
	return &userRoleBindingProjector{}
}

// AggregateType returns the AggregateType for which events should be projected.
func (u *userRoleBindingProjector) AggregateType() es.AggregateType {
	return aggregates.UserRoleBinding
}

func (u *userRoleBindingProjector) NewProjection() es.Projection {
	return &projections.UserRoleBinding{}
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
		case events.UserCreated:
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

	i.AggregateVersion++
	return i, nil
}
