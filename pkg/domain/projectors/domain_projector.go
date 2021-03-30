package projectors

import (
	"github.com/google/uuid"
	"gitlab.figo.systems/platform/monoskope/monoskope/internal/gateway"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/projections"
	es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing"
	timestamp "google.golang.org/protobuf/types/known/timestamppb"
)

type domainProjector struct {
}

// NewDomainProjector returns a new basic domain projector
func NewDomainProjector() *domainProjector {
	return &domainProjector{}
}

// getUserIdFromEvent gets the UserID from event metadata
func (*domainProjector) getUserIdFromEvent(event es.Event) (uuid.UUID, error) {
	userId, err := uuid.Parse(event.Metadata()[gateway.HeaderAuthId])
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
		return err
	}

	dp.LastModified = timestamp.New(event.Timestamp())
	dp.LastModifiedById = userId

	return nil
}

// projectCreated updates the created metadata
func (p *domainProjector) projectCreated(event es.Event, dp *projections.DomainProjection) error {
	// Get UserID from event metadata
	userId, err := p.getUserIdFromEvent(event)
	if err != nil {
		return err
	}

	dp.Created = timestamp.New(event.Timestamp())
	dp.CreatedById = userId

	return p.projectModified(event, dp)
}

// projectDeleted updates the deleted metadata
func (p *domainProjector) projectDeleted(event es.Event, dp *projections.DomainProjection) error {
	// Get UserID from event metadata
	userId, err := p.getUserIdFromEvent(event)
	if err != nil {
		return err
	}

	dp.Deleted = timestamp.New(event.Timestamp())
	dp.DeletedById = userId

	return p.projectModified(event, dp)
}
