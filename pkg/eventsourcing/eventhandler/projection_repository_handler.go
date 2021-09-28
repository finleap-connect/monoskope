// Copyright 2021 Monoskope Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package eventhandler

import (
	"context"
	"fmt"
	"sync"

	es "github.com/finleap-connect/monoskope/pkg/eventsourcing"
	esErrors "github.com/finleap-connect/monoskope/pkg/eventsourcing/errors"
	"github.com/finleap-connect/monoskope/pkg/logger"
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
	mutex      sync.Mutex
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
	h.mutex.Lock()
	defer h.mutex.Unlock()
	h.log.Info("Projecting event...", "event", event.String())

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
		h.log.Error(err, "Projecting event failed.", "event", event.String())
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
