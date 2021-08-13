package projections

import (
	"github.com/google/uuid"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/domain/projections"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing"
)

type User struct {
	*DomainProjection
	*projections.User
}

func NewUserProjection(id uuid.UUID) eventsourcing.Projection {
	dp := NewDomainProjection()
	return &User{
		DomainProjection: dp,
		User: &projections.User{
			Id:       id.String(),
			Metadata: &dp.LifecycleMetadata,
		},
	}
}

// ID implements the ID method of the Aggregate interface.
func (p *User) ID() uuid.UUID {
	return uuid.MustParse(p.Id)
}

// Proto gets the underlying proto representation.
func (p *User) Proto() *projections.User {
	return p.User
}
