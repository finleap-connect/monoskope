package projections

import (
	"github.com/google/uuid"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/domain/projections"
)

type Tenant struct {
	*projections.Tenant
	userRefs userReferenceIDs
	version  uint64
}

type userReferenceIDs struct {
	CreatedBy      string
	LastModifiedBy string
	DeletedBy      string
}

// ID implements the ID method of the Aggregate interface.
func (p *Tenant) ID() uuid.UUID {
	return uuid.MustParse(p.GetId())
}

// Version implements the Version method of the Aggregate interface.
func (p *Tenant) Version() uint64 {
	return p.version
}

// IncrementVersion implements the IncrementVersion method of the Projection interface.
func (p *Tenant) IncrementVersion() {
	p.version++
}

// Proto gets the underlying proto representation.
func (p *Tenant) Proto() *projections.Tenant {
	return p.Tenant
}

func (p *Tenant) GetCreatedByID() string {
	return p.userRefs.CreatedBy
}

func (p *Tenant) SetCreatedByID(id string) {
	p.userRefs.CreatedBy = id
}

func (p *Tenant) GetLastModifiedByID() string {
	return p.userRefs.LastModifiedBy
}

func (p *Tenant) SetLastModifiedByID(id string) {
	p.userRefs.LastModifiedBy = id
}

func (p *Tenant) GetDeletedByID() string {
	return p.userRefs.DeletedBy
}

func (p *Tenant) SetDeletedByID(id string) {
	p.userRefs.DeletedBy = id
}
