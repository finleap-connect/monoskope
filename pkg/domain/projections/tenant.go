package projections

import (
	"github.com/google/uuid"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/domain/projections"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing"
)

type Tenant struct {
	*DomainProjection
	*projections.Tenant
}

func NewTenantProjection(id uuid.UUID) eventsourcing.Projection {
	dp := NewDomainProjection()
	return &Tenant{
		DomainProjection: dp,
		Tenant: &projections.Tenant{
			Id:       id.String(),
			Metadata: &dp.LifecycleMetadata,
		},
	}
}

// ID implements the ID method of the Aggregate interface.
func (p *Tenant) ID() uuid.UUID {
	return uuid.MustParse(p.Id)
}

// Proto gets the underlying proto representation.
func (p *Tenant) Proto() *projections.Tenant {
	return p.Tenant
}
