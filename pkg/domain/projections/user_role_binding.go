package projections

import (
	"github.com/google/uuid"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/domain/projections"
)

type UserRoleBinding struct {
	projections.UserRoleBinding
	version uint64
}

// ID implements the ID method of the Aggregate interface.
func (p *UserRoleBinding) ID() uuid.UUID {
	return uuid.MustParse(p.GetId())
}

// Version implements the Version method of the Aggregate interface.
func (p *UserRoleBinding) Version() uint64 {
	return p.version
}

// IncrementVersion implements the IncrementVersion method of the Projection interface.
func (p *UserRoleBinding) IncrementVersion() {
	p.version++
}

// Proto gets the underlying proto representation.
func (p *UserRoleBinding) Proto() *projections.UserRoleBinding {
	return &p.UserRoleBinding
}
