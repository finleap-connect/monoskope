package projectors

import (
	"context"

	"github.com/google/uuid"
	aggregates "gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/constants/aggregates"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/projections"
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
	return projections.NewUser(uuid.Nil, "", "", []*projections.UserRoleBinding{})
}

// Project updates the state of the projection occording to the given event.
func (u *userRoleBindingProjector) Project(ctx context.Context, event es.Event, projection es.Projection) (es.Projection, error) {
	return projection, nil
}
