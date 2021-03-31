package projections

import (
	"github.com/google/uuid"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/domain/projections"
)

type UserRoleBinding struct {
	*DomainProjection
	*projections.UserRoleBinding
}

func NewUserRoleBinding(id uuid.UUID) *UserRoleBinding {
	dp := NewDomainProjection()
	return &UserRoleBinding{
		DomainProjection: dp,
		UserRoleBinding: &projections.UserRoleBinding{
			Id:       id.String(),
			Metadata: &dp.ProjectionMetadata,
		},
	}
}

// ID implements the ID method of the Aggregate interface.
func (p *UserRoleBinding) ID() uuid.UUID {
	return uuid.MustParse(p.Id)
}

// Proto gets the underlying proto representation.
func (p *UserRoleBinding) Proto() *projections.UserRoleBinding {
	return p.UserRoleBinding
}
