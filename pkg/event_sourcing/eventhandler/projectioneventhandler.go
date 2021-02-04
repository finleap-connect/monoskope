package eventhandler

import (
	"context"

	es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/event_sourcing"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/event_sourcing/errors"
)

type ProjectionEventHandler struct {
	projector  es.Projector
	repository es.Repository
}

func NewProjectionEventHandler(projector es.Projector, repository es.Repository) es.EventHandler {
	return &ProjectionEventHandler{
		projector:  projector,
		repository: repository,
	}
}

// HandleEvent implements the HandleEvent method of the es.EventHandler interface.
func (h *ProjectionEventHandler) HandleEvent(ctx context.Context, event es.Event) error {
	projection, err := h.repository.ById(ctx, event.AggregateID())

	// If error is not found create new projection.
	if err != nil {
		if err == errors.ErrProjectionNotFound {
			projection = h.projector.NewProjection()
		} else {
			return err
		}
	}

	// Check version.
	if projection.AggregateVersion() >= event.AggregateVersion() {
		// Ignore old/duplicate events.
		return nil
	}
	if projection.AggregateVersion()+1 != event.AggregateVersion() {
		// Version of event is not exactly one higher than the projection.
		return errors.ErrIncorrectAggregateVersion
	}

	// Apply event on projection.
	projection, err = h.projector.Project(ctx, event, projection)
	if err != nil {
		return err
	}

	// Check version again.
	if projection.AggregateVersion() != event.AggregateVersion() {
		// Project version and Event version do not match after projection.
		return errors.ErrIncorrectAggregateVersion
	}

	if projection == nil {
		// Remove projection from repo.
		return h.repository.Remove(ctx, event.AggregateID())
	} else {
		// Upsert projection in repo.
		return h.repository.Upsert(ctx, projection)
	}
}
