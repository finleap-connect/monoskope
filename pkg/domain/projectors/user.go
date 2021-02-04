package projectors

import (
	"context"

	"github.com/google/uuid"
	aggregates "gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/constants/aggregates"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/projections"
	es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/event_sourcing"
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
	return projections.NewUser(uuid.Nil, "", "", []*projections.UserRoleBinding{})
}

// Project updates the state of the projection occording to the given event.
func (u *userProjector) Project(ctx context.Context, event es.Event, projection es.Projection) (es.Projection, error) {
	return projection, nil
}
