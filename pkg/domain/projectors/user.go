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

type userProjector struct {
	*domainProjector
}

func NewUserProjector() es.Projector {
	return &userProjector{
		domainProjector: NewDomainProjector(),
	}
}

func (u *userProjector) NewProjection(id uuid.UUID) es.Projection {
	return projections.NewUserProjection(id)
}

// Project updates the state of the projection according to the given event.
func (u *userProjector) Project(ctx context.Context, event es.Event, projection es.Projection) (es.Projection, error) {
	// Get the actual projection type
	p, ok := projection.(*projections.User)
	if !ok {
		return nil, errors.ErrInvalidProjectionType
	}

	// Apply the changes for the event.
	switch event.EventType() {
	case events.UserCreated:
		data := &eventdata.UserCreated{}
		if err := event.Data().ToProto(data); err != nil {
			return projection, err
		}

		p.Email = data.GetEmail()
		p.Name = data.GetName()

		if err := u.projectCreated(event, p.DomainProjection); err != nil {
			return nil, err
		}
	default:
		return nil, errors.ErrInvalidEventType
	}

	if err := u.projectModified(event, p.DomainProjection); err != nil {
		return nil, err
	}
	p.IncrementVersion()

	return p, nil
}
