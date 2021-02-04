package projectors

import (
	"context"

	aggregates "gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/constants/aggregates"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/projections"
	es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/event_sourcing"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/event_sourcing/errors"
)

type userProjector struct{}

func NewUserProjector() es.Projector {
	return &userProjector{}
}

// EvenTypes returns the EvenTypes for which events should be projected.
func (u *userProjector) EvenTypes() []es.EventType {
	return []es.EventType{}
}

// AggregateTypes returns the AggregateTypes for which events should be projected.
func (u *userProjector) AggregateTypes() []es.AggregateType {
	return []es.AggregateType{
		aggregates.User,
		aggregates.UserRoleBinding,
	}
}

func (u *userProjector) NewProjection() es.Projection {
	return &projections.User{}
}

func (u *userProjector) ValidateVersion(ctx context.Context, event es.Event, projection es.Projection) error {
	// Check version.
	if projection.AggregateVersion() >= event.AggregateVersion() {
		// Ignore old/duplicate events.
		return nil
	}

	if projection.AggregateVersion()+1 != event.AggregateVersion() {
		// Version of event is not exactly one higher than the projection.
		return errors.ErrIncorrectAggregateVersion
	}

	return nil
}

// Project updates the state of the projection occording to the given event.
func (u *userProjector) Project(ctx context.Context, event es.Event, projection es.Projection) (es.Projection, error) {
	return projection, nil
}
