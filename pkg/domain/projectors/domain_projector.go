package projectors

import (
	"github.com/google/uuid"
	"gitlab.figo.systems/platform/monoskope/monoskope/internal/gateway"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/projections"
	es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing"
	timestamp "google.golang.org/protobuf/types/known/timestamppb"
)

type DomainProjector struct {
}

func NewDomainProjector() *DomainProjector {
	return &DomainProjector{}
}

// GetUserIdFromEvent gets the UserID from event metadata
func (*DomainProjector) GetUserIdFromEvent(event es.Event) (uuid.UUID, error) {
	userId, err := uuid.Parse(event.Metadata()[gateway.HeaderAuthId])
	if err != nil {
		return uuid.Nil, err
	}
	return userId, nil
}

func (p *DomainProjector) ProjectModified(event es.Event, dp *projections.DomainProjection) error {
	// Get UserID from event metadata
	userId, err := p.GetUserIdFromEvent(event)
	if err != nil {
		return err
	}

	dp.LastModified = timestamp.New(event.Timestamp())
	dp.LastModifiedById = userId

	return nil
}

func (p *DomainProjector) ProjectCreated(event es.Event, dp *projections.DomainProjection) error {
	// Get UserID from event metadata
	userId, err := p.GetUserIdFromEvent(event)
	if err != nil {
		return err
	}

	dp.Created = timestamp.New(event.Timestamp())
	dp.CreatedById = userId

	return p.ProjectModified(event, dp)
}

func (p *DomainProjector) ProjectDeleted(event es.Event, dp *projections.DomainProjection) error {
	// Get UserID from event metadata
	userId, err := p.GetUserIdFromEvent(event)
	if err != nil {
		return err
	}

	dp.Deleted = timestamp.New(event.Timestamp())
	dp.DeletedById = userId

	return p.ProjectModified(event, dp)
}
