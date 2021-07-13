package repositories

import (
	"context"

	domApi "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/domain"
	projections "gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/projections"
	es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing"
	esErrors "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing/errors"
)

type certificateRepository struct {
	*domainRepository
}

// TenantRepository is a repository for reading and writing tenant projections.
type CertificateRepository interface {
	es.Repository
	ReadOnlyCertificateRepository
	WriteOnlyCertificateRepository
}

// ReadOnlyTenantRepository is a repository for reading tenant projections.
type ReadOnlyCertificateRepository interface {
	// ById searches for the a tenant projection by it's id.
	GetCertificate(context.Context, *domApi.GetCertificateRequest) (*projections.Certificate, error)
}

// WriteOnlyTenantRepository is a repository for writing tenant projections.
type WriteOnlyCertificateRepository interface {
}

// NewCertificateRepository creates a repository for reading and writing certificate projections.
func NewCertificateRepository(repository es.Repository, userRepo UserRepository) CertificateRepository {
	return &certificateRepository{
		domainRepository: NewDomainRepository(repository, userRepo),
	}
}

// Retrieve certificates for a specified aggregate ID and type.
func (r *certificateRepository) GetCertificate(ctx context.Context, req *domApi.GetCertificateRequest) (*projections.Certificate, error) {
	projection, err := r.ById(ctx, uuid)
	if err != nil {
		return nil, err
	}

	tenant, ok := projection.(*projections.Tenant)
	if !ok {
		return nil, esErrors.ErrInvalidProjectionType
	}

	err = r.addMetadata(ctx, tenant.DomainProjection)
	if err != nil {
		return nil, err
	}

	return tenant, nil
}
