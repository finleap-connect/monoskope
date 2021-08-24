package projections

import (
	"github.com/google/uuid"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/domain/projections"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing"
)

type Certificate struct {
	*DomainProjection
	*projections.Certificate
	SigningRequest []byte
}

func NewCertificateProjection(id uuid.UUID) eventsourcing.Projection {
	dp := NewDomainProjection()
	return &Certificate{
		DomainProjection: dp,
		Certificate: &projections.Certificate{
			Id:       id.String(),
			Metadata: &dp.LifecycleMetadata,
		},
	}
}

// ID implements the ID method of the Aggregate interface.
func (p *Certificate) ID() uuid.UUID {
	return uuid.MustParse(p.ReferencedAggregateId)
}

// Proto gets the underlying proto representation.
func (p *Certificate) Proto() *projections.Certificate {
	return p.Certificate
}
