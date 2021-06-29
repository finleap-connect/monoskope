package repositories

import (
	"context"

	"github.com/google/uuid"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/constants/users"
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

func findSystemUser(id uuid.UUID) *projections.User {
	sysUser, ok := users.AvailableSystemUsers[id]
	if !ok {
		return nil
	}
	user := projections.NewUserProjection(sysUser.ID).(*projections.User)
	user.Name = sysUser.Name
	user.Email = sysUser.Email
	return user
}

func (r *domainRepository) addMetadata(ctx context.Context, dp *projections.DomainProjection) error {
	if id := dp.CreatedById; id != uuid.Nil {
		user, err := r.userRepo.ByUserId(ctx, id)
		if err != nil {
			user = findSystemUser(id)
			if user == nil {
				return err
			}
		}
		dp.CreatedBy = user.User
	}

	if id := dp.LastModifiedById; id != uuid.Nil {
		user, err := r.userRepo.ByUserId(ctx, id)
		if err != nil {
			user = findSystemUser(id)
			if user == nil {
				return err
			}
		}
		dp.LastModifiedBy = user.User
	}

	if id := dp.DeletedById; id != uuid.Nil {
		user, err := r.userRepo.ByUserId(ctx, id)
		if err != nil {
			user = findSystemUser(id)
			if user == nil {
				return err
			}
		}
		dp.DeletedBy = user.User
	}

	return nil
}
