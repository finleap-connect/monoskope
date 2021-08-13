package projections

import (
	"github.com/google/uuid"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/domain/projections"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing"
)

type Cluster struct {
	*DomainProjection
	*projections.Cluster
}

func NewClusterProjection(id uuid.UUID) eventsourcing.Projection {
	dp := NewDomainProjection()
	return &Cluster{
		DomainProjection: dp,
		Cluster: &projections.Cluster{
			Id:       id.String(),
			Metadata: &dp.LifecycleMetadata,
		},
	}
}

// ID implements the ID method of the Projection interface.
func (p *Cluster) ID() uuid.UUID {
	return uuid.MustParse(p.Id)
}

// Proto gets the underlying proto representation.
func (p *Cluster) Proto() *projections.Cluster {
	return p.Cluster
}
