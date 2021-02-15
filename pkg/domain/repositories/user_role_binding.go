package repositories

import (
	"context"

	"github.com/google/uuid"
	projections "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/domain/projections"
	es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing"
)

type userRoleBindingRepository struct {
	es.ProjectionRepository
}

// Repository is a repository for reading and writing UserRoleBinding projections.
type UserRoleBindingRepository interface {
	es.ProjectionRepository
	ReadOnlyUserRoleBindingRepository
	WriteOnlyUserRoleBindingRepository
}

// ReadOnlyUserRepository is a repository for reading UserRoleBinding projections.
type ReadOnlyUserRoleBindingRepository interface {
	// ByUserId searches for all UserRoleBinding projection's by the a user id.
	ByUserId(context.Context, uuid.UUID) ([]*projections.UserRoleBinding, error)
}

// WriteOnlyUserRepository is a repository for reading UserRoleBinding projections.
type WriteOnlyUserRoleBindingRepository interface {
}

// NewUserRepository creates a repository for reading and writing UserRoleBinding projections.
func NewUserRoleBindingRepository(repository es.ProjectionRepository) UserRoleBindingRepository {
	return &userRoleBindingRepository{
		ProjectionRepository: repository,
	}
}

// ByUserId searches for all UserRoleBinding projection's by the a user id.
func (r *userRoleBindingRepository) ByUserId(ctx context.Context, userId uuid.UUID) ([]*projections.UserRoleBinding, error) {
	ps, err := r.All(ctx)
	if err != nil {
		return nil, err
	}

	var userRoleBindings []*projections.UserRoleBinding
	for _, projection := range ps {
		if userRoleBinding, ok := projection.(*projections.UserRoleBinding); ok {
			if userId.String() == userRoleBinding.GetUserId() {
				userRoleBindings = append(userRoleBindings, userRoleBinding)
			}
		}
	}
	return userRoleBindings, nil
}
