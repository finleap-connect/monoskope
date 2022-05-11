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
	"github.com/finleap-connect/monoskope/internal/gateway/auth"
	"github.com/finleap-connect/monoskope/pkg/domain/projections"
	es "github.com/finleap-connect/monoskope/pkg/eventsourcing"
	"github.com/finleap-connect/monoskope/pkg/logger"
	"github.com/google/uuid"
	timestamp "google.golang.org/protobuf/types/known/timestamppb"
)

type domainProjector struct {
	log logger.Logger
}

// NewDomainProjector returns a new basic domain projector
func NewDomainProjector() *domainProjector {
	return &domainProjector{
		log: logger.WithName("domain-projector"),
	}
}

// getUserIdFromEvent gets the UserID from event metadata
func (*domainProjector) getUserIdFromEvent(event es.Event) (uuid.UUID, error) {
	userId, err := uuid.Parse(event.Metadata()[auth.HeaderAuthId])
	if err != nil {
		return uuid.Nil, err
	}
	return userId, nil
}

// projectModified updates the modified metadata
func (p *domainProjector) projectModified(event es.Event, dp *projections.DomainProjection) error {
	// Get UserID from event metadata
	userId, err := p.getUserIdFromEvent(event)
	if err != nil {
		p.log.Info("Event metadata do not contain user information.", "EventType", event.EventType(), "AggregateType", event.AggregateType(), "AggregateID", event.AggregateID())
		return err
	}

	dp.LastModified = timestamp.New(event.Timestamp())
	dp.LastModifiedById = userId.String()

	return nil
}

// projectCreated updates the created metadata
func (p *domainProjector) projectCreated(event es.Event, dp *projections.DomainProjection) error {
	// Get UserID from event metadata
	userId, err := p.getUserIdFromEvent(event)
	if err != nil {
		p.log.Info("Event metadata do not contain user information.", "EventType", event.EventType(), "AggregateType", event.AggregateType(), "AggregateID", event.AggregateID())
		return err
	}

	dp.Created = timestamp.New(event.Timestamp())
	dp.CreatedById = userId.String()

	return p.projectModified(event, dp)
}

// projectDeleted updates the deleted metadata
func (p *domainProjector) projectDeleted(event es.Event, dp *projections.DomainProjection) error {
	// Get UserID from event metadata
	userId, err := p.getUserIdFromEvent(event)
	if err != nil {
		p.log.Info("Event metadata do not contain user information.", "EventType", event.EventType(), "AggregateType", event.AggregateType(), "AggregateID", event.AggregateID())
		return err
	}

	dp.Deleted = timestamp.New(event.Timestamp())
	dp.DeletedById = userId.String()

	return p.projectModified(event, dp)
}
