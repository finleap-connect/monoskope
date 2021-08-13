package projections

import (
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/domain/projections"
)

type DomainProjection struct {
	projections.LifecycleMetadata
	version uint64
}

func NewDomainProjection() *DomainProjection {
	return &DomainProjection{}
}

// Version implements the Version method of the Projection interface.
func (p *DomainProjection) Version() uint64 {
	return p.version
}

// IncrementVersion implements the IncrementVersion method of the Projection interface.
func (p *DomainProjection) IncrementVersion() {
	p.version++
}
