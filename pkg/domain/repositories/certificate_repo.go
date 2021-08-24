package repositories

import (
	"context"

	"github.com/google/uuid"
	domApi "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/domain"
	projections "gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/projections"
	es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing"
	esErrors "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing/errors"
)

type certificateRepository struct {
	es.Repository
}

// CertificateRepository is a repository for reading and writing certificate projections.
type CertificateRepository interface {
	es.Repository
	ReadOnlyCertificateRepository
	WriteOnlyCertificateRepository
}

// ReadOnlyCertificateRepository is a repository for reading certificate projections.
type ReadOnlyCertificateRepository interface {
	// GetCertificate retrieves certificates by aggregate type and id
	GetCertificate(context.Context, *domApi.GetCertificateRequest) (*projections.Certificate, error)
}

// WriteOnlyCertificateRepository is a repository for writing certificate projections.
type WriteOnlyCertificateRepository interface {
}

// NewCertificateRepository creates a repository for reading and writing certificate projections.
func NewCertificateRepository(repository es.Repository) CertificateRepository {
	return &certificateRepository{
		Repository: repository,
	}
}

// Retrieve certificates for a specified aggregate ID and type.
func (r *certificateRepository) GetCertificate(ctx context.Context, req *domApi.GetCertificateRequest) (*projections.Certificate, error) {
	id, err := uuid.Parse(req.AggregateId)
	if err != nil {
		return nil, err
	}

	projection, err := r.ById(ctx, id)
	if err != nil {
		return nil, err
	}

	certificate, ok := projection.(*projections.Certificate)
	if !ok {
		return nil, esErrors.ErrInvalidProjectionType
	}

	return certificate, nil
}
