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

type userProjector struct {
	*domainProjector
}

func NewUserProjector() es.Projector[*projections.User] {
	return &userProjector{
		domainProjector: NewDomainProjector(),
	}
}

func (u *userProjector) NewProjection(id uuid.UUID) *projections.User {
	return projections.NewUserProjection(id)
}

// Project updates the state of the projection according to the given event.
func (u *userProjector) Project(ctx context.Context, event es.Event, p *projections.User) (*projections.User, error) {
	// Apply the changes for the event.
	switch event.EventType() {
	case events.UserCreated:
		data := &eventdata.UserCreated{}
		if err := event.Data().ToProto(data); err != nil {
			return p, err
		}

		p.Email = data.GetEmail()
		p.Name = data.GetName()
		p.Source = data.GetSource()

		if err := u.projectCreated(event, p.DomainProjection); err != nil {
			return nil, err
		}
	case events.UserUpdated:
		data := &eventdata.UserUpdated{}
		if err := event.Data().ToProto(data); err != nil {
			return p, err
		}

		p.Name = data.GetName()
	case events.UserDeleted:
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
