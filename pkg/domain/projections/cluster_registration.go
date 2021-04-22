package projections

import (
	"github.com/google/uuid"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/domain/projections"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing"
)

type ClusterRegistration struct {
	*DomainProjection
	*projections.ClusterRegistration
}

func NewClusterRegistrationProjection(id uuid.UUID) eventsourcing.Projection {
	dp := NewDomainProjection()
	return &ClusterRegistration{
		DomainProjection: dp,
		ClusterRegistration: &projections.ClusterRegistration{
			Id:       id.String(),
			Metadata: &dp.ProjectionMetadata,
		},
	}
}

// ID implements the ID method of the Aggregate interface.
func (p *ClusterRegistration) ID() uuid.UUID {
	return uuid.MustParse(p.Id)
}

// Proto gets the underlying proto representation.
func (p *ClusterRegistration) Proto() *projections.ClusterRegistration {
	return p.ClusterRegistration
}
