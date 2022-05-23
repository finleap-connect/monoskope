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

type tenantProjector struct {
	*domainProjector
}

func NewTenantProjector() es.Projector[*projections.Tenant] {
	return &tenantProjector{
		domainProjector: NewDomainProjector(),
	}
}

func (t *tenantProjector) NewProjection(id uuid.UUID) *projections.Tenant {
	return projections.NewTenantProjection(id)
}

// Project updates the state of the projection according to the given event.
func (t *tenantProjector) Project(ctx context.Context, event es.Event, p *projections.Tenant) (*projections.Tenant, error) {
	// Apply the changes for the event.
	switch event.EventType() {
	case events.TenantCreated:
		data := new(eventdata.TenantCreated)
		if err := event.Data().ToProto(data); err != nil {
			return p, err
		}

		p.Name = data.GetName()
		p.Prefix = data.GetPrefix()

		if err := t.projectCreated(event, p.DomainProjection); err != nil {
			return nil, err
		}
	case events.TenantUpdated:
		data := new(eventdata.TenantUpdated)
		if err := event.Data().ToProto(data); err != nil {
			return p, err
		}
		p.Name = data.GetName().GetValue()
	case events.TenantDeleted:
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
