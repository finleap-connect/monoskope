package repositories

import (
	"context"

	"github.com/google/uuid"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/errors"
	projections "gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/projections"
	es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing"
	esErrors "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing/errors"
)

type tenantRepository struct {
	*domainRepository
}

// TenantRepository is a repository for reading and writing tenant projections.
type TenantRepository interface {
	es.Repository
	ReadOnlyTenantRepository
	WriteOnlyTenantRepository
}

// ReadOnlyTenantRepository is a repository for reading tenant projections.
type ReadOnlyTenantRepository interface {
	// ById searches for the a tenant projection by it's id.
	ByTenantId(context.Context, string) (*projections.Tenant, error)
	// ByName searches for the a tenant projection by it's name
	ByName(context.Context, string) (*projections.Tenant, error)
	// GetAll searches for all tenant projections.
	GetAll(context.Context, bool) ([]*projections.Tenant, error)
}

// WriteOnlyTenantRepository is a repository for writing tenant projections.
type WriteOnlyTenantRepository interface {
}

// NewTenantRepository creates a repository for reading and writing tenant projections.
func NewTenantRepository(repository es.Repository, userRepo UserRepository) TenantRepository {
	return &tenantRepository{
		domainRepository: NewDomainRepository(repository, userRepo),
	}
}

// ByTenantId searches for a tenant projection by its id.
func (r *tenantRepository) ByTenantId(ctx context.Context, id string) (*projections.Tenant, error) {
	uuid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

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

// ByTenantName searches for a tenant projection by its name.
func (r *tenantRepository) ByName(ctx context.Context, name string) (*projections.Tenant, error) {
	ps, err := r.GetAll(ctx, true)
	if err != nil {
		return nil, err
	}

	for _, t := range ps {
		if name == t.Name {
			return t, nil
		}
	}

	return nil, errors.ErrTenantNotFound
}

// All searches for the a tenant projections.
func (r *tenantRepository) GetAll(ctx context.Context, includeDeleted bool) ([]*projections.Tenant, error) {
	ps, err := r.All(ctx)
	if err != nil {
		return nil, err
	}

	var tenants []*projections.Tenant
	for _, p := range ps {
		if t, ok := p.(*projections.Tenant); ok {
			if !t.GetDeleted().IsValid() || includeDeleted {
				// Add metadata
				err = r.addMetadata(ctx, t.DomainProjection)
				if err != nil {
					return nil, err
				}
				tenants = append(tenants, t)
			}
		} else {
			return nil, esErrors.ErrInvalidProjectionType
		}
	}
	return tenants, nil
}
