package eventhandler

import (
	"context"

	es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing/errors"
)

type projectionRepoEventHandler struct {
	projector  es.Projector
	repository es.Repository
}

// NewProjectionRepositoryEventHandler creates an EventHandler which applies incoming events on a Projector and updates the Repository accordingly.
func NewProjectionRepositoryEventHandler(projector es.Projector, repository es.Repository) es.EventHandler {
	return &projectionRepoEventHandler{
		projector:  projector,
		repository: repository,
	}
}

// HandleEvent implements the HandleEvent method of the es.EventHandler interface.
func (h *projectionRepoEventHandler) HandleEvent(ctx context.Context, event es.Event) error {
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
	if projection.Version() >= event.AggregateVersion() {
		// Ignore old/duplicate events.
		return nil
	}
	if projection.Version()+1 != event.AggregateVersion() {
		// Version of event is not exactly one higher than the projection.
		return errors.ErrProjectionOutdated
	}

	// Apply event on projection.
	projection, err = h.projector.Project(ctx, event, projection)
	if err != nil {
		return err
	}

	// Check version again.
	if projection.Version() != event.AggregateVersion() {
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
