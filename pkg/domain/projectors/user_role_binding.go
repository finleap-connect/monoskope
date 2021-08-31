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

	if err := u.projectModified(event, p.DomainProjection); err != nil {
		return nil, err
	}
	p.IncrementVersion()

	return p, nil
}
