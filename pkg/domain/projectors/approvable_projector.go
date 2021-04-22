package projectors

import (
	"github.com/google/uuid"
	"gitlab.figo.systems/platform/monoskope/monoskope/internal/gateway"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/projections"
	es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing"
	timestamp "google.golang.org/protobuf/types/known/timestamppb"
)

type approvableProjector struct {
	*domainProjector
}

// NewApprovableProjector returns a new basic approvable projector
func NewApprovableProjector() *approvableProjector {
	return &approvableProjector{
		domainProjector: NewDomainProjector(),
	}
}

// getUserIdFromEvent gets the UserID from event metadata
func (*approvableProjector) getUserIdFromEvent(event es.Event) (uuid.UUID, error) {
	userId, err := uuid.Parse(event.Metadata()[gateway.HeaderAuthId])
	if err != nil {
		return uuid.Nil, err
	}
	return userId, nil
}

// projectApproved updates the modified metadata
func (p *approvableProjector) projectApproved(event es.Event, dp *projections.ApprovableProjection) error {
	// Get UserID from event metadata
	userId, err := p.getUserIdFromEvent(event)
	if err != nil {
		p.log.Info("Event metadata do not contain user information.", "EventType", event.EventType(), "AggregateType", event.AggregateType(), "AggregateID", event.AggregateID())
		return nil
	}

	dp.Approved = timestamp.New(event.Timestamp())
	dp.ApprovedById = userId

	return p.projectModified(event, dp.DomainProjection)
}

// projectDenied updates the modified metadata
func (p *approvableProjector) projectDenied(event es.Event, dp *projections.ApprovableProjection) error {
	// Get UserID from event metadata
	userId, err := p.getUserIdFromEvent(event)
	if err != nil {
		p.log.Info("Event metadata do not contain user information.", "EventType", event.EventType(), "AggregateType", event.AggregateType(), "AggregateID", event.AggregateID())
		return nil
	}

	dp.Denied = timestamp.New(event.Timestamp())
	dp.DeniedById = userId

	return p.projectModified(event, dp.DomainProjection)
}
