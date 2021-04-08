package eventhandler

import (
	"context"
	"fmt"

	es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing"
	esErrors "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing/errors"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/logger"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/util"
)

type ProjectionOutdatedError struct {
	ProjectionVersion uint64
}

func (e *ProjectionOutdatedError) Error() string {
	return fmt.Sprintf("projection version %v is outdated", e.ProjectionVersion)
}

type projectionRepoEventHandler struct {
	log        logger.Logger
	projector  es.Projector
	repository es.Repository
}

// NewProjectingEventHandler creates an EventHandler which applies incoming events on a Projector and updates the Repository accordingly.
func NewProjectingEventHandler(projector es.Projector, repository es.Repository) es.EventHandler {
	return &projectionRepoEventHandler{
		log:        logger.WithName("projection-repo-middleware"),
		projector:  projector,
		repository: repository,
	}
}

// HandleEvent implements the HandleEvent method of the es.EventHandler interface.
func (h *projectionRepoEventHandler) HandleEvent(ctx context.Context, event es.Event) error {
	if util.GetOperationMode() == util.DEVELOPMENT {
		h.log.Info("Projecting event...", "EventType", event.EventType(), "AggregateType", event.AggregateType())
	}

	projection, err := h.repository.ById(ctx, event.AggregateID())

	// If error is not found create new projection.
	if err != nil {
		if err == esErrors.ErrProjectionNotFound {
			projection = h.projector.NewProjection(event.AggregateID())
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
		return &ProjectionOutdatedError{ProjectionVersion: projection.Version()}
	}

	// Apply event on projection.
	projection, err = h.projector.Project(ctx, event, projection)
	if err != nil {
		return err
	}

	// Check version again.
	if projection.Version() != event.AggregateVersion() {
		// Project version and Event version do not match after projection.
		return esErrors.ErrIncorrectAggregateVersion
	}

	// Upsert projection in repo.
	return h.repository.Upsert(ctx, projection)
}
