package projections

import (
	"github.com/google/uuid"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/domain/projections"
)

type User struct {
	*projections.User
	version uint64
}

// ID implements the ID method of the Aggregate interface.
func (p *User) ID() uuid.UUID {
	return uuid.MustParse(p.GetId())
}

// Version implements the Version method of the Aggregate interface.
func (p *User) Version() uint64 {
	return p.version
}

// IncrementVersion implements the IncrementVersion method of the Projection interface.
func (p *User) IncrementVersion() {
	p.version++
}

// Proto gets the underlying proto representation.
func (p *User) Proto() *projections.User {
	return p.User
}
