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

type userRoleBindingProjector struct {
	*domainProjector
}

func NewUserRoleBindingProjector() es.Projector {
	return &userRoleBindingProjector{
		domainProjector: NewDomainProjector(),
	}
}

func (u *userRoleBindingProjector) NewProjection(id uuid.UUID) es.Projection {
	return projections.NewUserRoleBinding(id)
}

// Project updates the state of the projection according to the given event.
func (u *userRoleBindingProjector) Project(ctx context.Context, event es.Event, projection es.Projection) (es.Projection, error) {
	// Get the actual projection type
	p, ok := projection.(*projections.UserRoleBinding)
	if !ok {
		return nil, errors.ErrInvalidProjectionType
	}

	// Apply the changes for the event.
	switch event.EventType() {
	case events.UserRoleBindingCreated:
		data := &eventdata.UserRoleAdded{}
		if err := event.Data().ToProto(data); err != nil {
			return projection, err
		}

		p.UserId = data.GetUserId()
		p.Role = data.GetRole()
		p.Scope = data.GetScope()
		p.Resource = data.GetResource()

		if err := u.projectCreated(event, p.DomainProjection); err != nil {
			return nil, err
		}
	case events.UserRoleBindingDeleted:
		if err := u.projectDeleted(event, p.DomainProjection); err != nil {
			return nil, err
		}
	default:
		return nil, errors.ErrInvalidEventType
	}

	p.IncrementVersion()

	return p, nil
}
