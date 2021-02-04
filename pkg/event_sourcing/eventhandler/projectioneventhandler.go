package eventhandler

import (
	"context"

	es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/event_sourcing"
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
		return err
	}

	// Check that the projection is one version behind the event.

	// Apply event on projection.
	projection, err = h.projector.Project(ctx, event, projection)
	if err != nil {
		return err
	}

	// Update projection in repo.
	err = h.repository.Upsert(ctx, projection)
	if err != nil {
		return err
	}

	return nil
}
