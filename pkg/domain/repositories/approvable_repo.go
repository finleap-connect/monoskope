package repositories

import (
	"context"

	"github.com/google/uuid"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/projections"
	es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing"
)

type approvableRepository struct {
	*domainRepository
}

func NewApprovableRepository(repository es.Repository, userRepo UserRepository) *approvableRepository {
	return &approvableRepository{
		domainRepository: NewDomainRepository(repository, userRepo),
	}
}

func (r *approvableRepository) addMetadata(ctx context.Context, dp *projections.ApprovableProjection) error {
	if id := dp.ApprovedById; id != uuid.Nil {
		user, err := r.userRepo.ByUserId(ctx, id)
		if err != nil {
			return err
		}
		dp.ApprovedBy = user.User
	}

	if id := dp.DeniedById; id != uuid.Nil {
		user, err := r.userRepo.ByUserId(ctx, id)
		if err != nil {
			return err
		}
		dp.DeniedBy = user.User
	}

	return r.domainRepository.addMetadata(ctx, dp.DomainProjection)
}
