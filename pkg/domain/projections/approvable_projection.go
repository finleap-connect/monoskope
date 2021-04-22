package projections

import (
	"github.com/google/uuid"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/domain/projections"
)

type ApprovableProjection struct {
	*DomainProjection
	*projections.ApprovalMetadata
	ApprovedById uuid.UUID
	DeniedById   uuid.UUID
}

func NewApprovableProjection() *ApprovableProjection {
	dp := NewDomainProjection()
	return &ApprovableProjection{
		DomainProjection: dp,
		ApprovalMetadata: &projections.ApprovalMetadata{},
	}
}
