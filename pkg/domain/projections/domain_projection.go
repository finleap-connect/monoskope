package projections

import (
	"github.com/google/uuid"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/domain/projections"
)

type DomainProjection struct {
	projections.LifecycleMetadata
	id               uuid.UUID
	version          uint64
	CreatedById      uuid.UUID
	LastModifiedById uuid.UUID
	DeletedById      uuid.UUID
}

func NewDomainProjection() *DomainProjection {
	return &DomainProjection{}
}

// ID implements the ID method of the Aggregate interface.
func (p *DomainProjection) ID() uuid.UUID {
	return p.id
}

// Version implements the Version method of the Projection interface.
func (p *DomainProjection) Version() uint64 {
	return p.version
}

// IncrementVersion implements the IncrementVersion method of the Projection interface.
func (p *DomainProjection) IncrementVersion() {
	p.version++
}
