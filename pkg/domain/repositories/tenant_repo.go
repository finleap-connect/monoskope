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
	es.Repository
	userRepo UserRepository
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
}

// WriteOnlyTenantRepository is a repository for writing tenant projections.
type WriteOnlyTenantRepository interface {
}

// NewTenantRepository creates a repository for reading and writing tenant projections.
func NewTenantRepository(repository es.Repository, userRepo UserRepository) TenantRepository {
	return &tenantRepository{
		Repository: repository,
		userRepo:   userRepo,
	}
}

func (r *tenantRepository) addUsersToTenant(ctx context.Context, tenant *projections.Tenant) error {
	createdBy, err := r.userRepo.ByUserId(ctx, tenant.CreatedById)
	if err != nil {
		return err
	}
	tenant.CreatedBy = createdBy.User
	if id := tenant.LastModifiedById; id != uuid.Nil {
		lastModifiedBy, err := r.userRepo.ByUserId(ctx, id)
		if err != nil {
			return err
		}
		tenant.LastModifiedBy = lastModifiedBy.User
	}
	if id := tenant.DeletedById; id != uuid.Nil {
		deletedBy, err := r.userRepo.ByUserId(ctx, id)
		if err != nil {
			return err
		}
		tenant.DeletedBy = deletedBy.User
	}
	return nil
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

	if tenant, ok := projection.(*projections.Tenant); !ok {
		return nil, esErrors.ErrInvalidProjectionType
	} else {
		err = r.addUsersToTenant(ctx, tenant)
		if err != nil {
			return nil, err
		}
		return tenant, nil
	}
}

// ByTenantName searches for a tenant projection by its name.
func (r *tenantRepository) ByName(ctx context.Context, name string) (*projections.Tenant, error) {
	ps, err := r.All(ctx)
	if err != nil {
		return nil, err
	}

	var tenant *projections.Tenant
	for _, p := range ps {
		if t, ok := p.(*projections.Tenant); ok {
			if name == t.Name {
				// User found
				tenant = t
			}
		}
	}
	if tenant != nil {
		// Find users that created, modified or deleted tenant
		err = r.addUsersToTenant(ctx, tenant)
		if err != nil {
			return nil, err
		}
		return tenant, nil
	}

	return nil, errors.ErrTenantNotFound
}
