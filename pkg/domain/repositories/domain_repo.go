package repositories

import (
	"context"

	"github.com/google/uuid"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/projections"
	es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing"
)

type domainRepository struct {
	es.Repository
	userRepo UserRepository
}

func NewDomainRepository(repository es.Repository, userRepo UserRepository) *domainRepository {
	return &domainRepository{
		Repository: repository,
		userRepo:   userRepo,
	}
}

func (r *domainRepository) addMetadata(ctx context.Context, dp *projections.DomainProjection) error {
	if id := dp.CreatedById; id != uuid.Nil {
		createdBy, err := r.userRepo.ByUserId(ctx, id)
		if err != nil {
			return err
		}
		dp.CreatedBy = createdBy.User
	}

	if id := dp.LastModifiedById; id != uuid.Nil {
		lastModifiedBy, err := r.userRepo.ByUserId(ctx, id)
		if err != nil {
			return err
		}
		dp.LastModifiedBy = lastModifiedBy.User
	}

	if id := dp.DeletedById; id != uuid.Nil {
		deletedBy, err := r.userRepo.ByUserId(ctx, id)
		if err != nil {
			return err
		}
		dp.DeletedBy = deletedBy.User
	}

	return nil
}
