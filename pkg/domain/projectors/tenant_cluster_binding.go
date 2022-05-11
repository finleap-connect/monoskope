// Copyright 2022 Monoskope Authors
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

package projectors

import (
	"context"

	"github.com/finleap-connect/monoskope/pkg/api/domain/eventdata"
	"github.com/finleap-connect/monoskope/pkg/domain/constants/events"
	"github.com/finleap-connect/monoskope/pkg/domain/projections"
	es "github.com/finleap-connect/monoskope/pkg/eventsourcing"
	"github.com/finleap-connect/monoskope/pkg/eventsourcing/errors"
	"github.com/google/uuid"
)

type tenantclusterbindingProjector struct {
	*domainProjector
}

func NewTenantClusterBindingProjector() es.Projector {
	return &tenantclusterbindingProjector{
		domainProjector: NewDomainProjector(),
	}
}

func (t *tenantclusterbindingProjector) NewProjection(id uuid.UUID) es.Projection {
	return projections.NewTenantClusterBindingProjection(id)
}

// Project updates the state of the projection according to the given event.
func (t *tenantclusterbindingProjector) Project(ctx context.Context, event es.Event, projection es.Projection) (es.Projection, error) {
	// Get the actual projection type
	p, ok := projection.(*projections.TenantClusterBinding)
	if !ok {
		return nil, errors.ErrInvalidProjectionType
	}

	// Apply the changes for the event.
	switch event.EventType() {
	case events.TenantClusterBindingCreated:
		data := new(eventdata.TenantClusterBindingCreated)
		if err := event.Data().ToProto(data); err != nil {
			return projection, err
		}

		p.TenantId = data.GetTenantId()
		p.ClusterId = data.GetClusterId()

		if err := t.projectCreated(event, p.DomainProjection); err != nil {
			return nil, err
		}
	case events.TenantClusterBindingDeleted:
		if err := t.projectDeleted(event, p.DomainProjection); err != nil {
			return nil, err
		}
	default:
		return nil, errors.ErrInvalidEventType
	}

	if err := t.projectModified(event, p.DomainProjection); err != nil {
		return nil, err
	}
	p.IncrementVersion()

	return p, nil
}
