package repositories

import (
	"context"

	"github.com/google/uuid"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/projections"
	es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/event_sourcing"
)

type userRoleBindingRepository struct {
	es.Repository
}

// Repository is a repository for reading and writing user projections.
type UserRoleBindingRepository interface {
	ReadOnlyUserRoleBindingRepository
	WriteOnlyUserRoleBindingRepository
}

// ReadOnlyUserRepository is a repository for reading user projections.
type ReadOnlyUserRoleBindingRepository interface {
	es.ReadOnlyRepository

	// ByEmail searches for the a user projection by it's email address.
	ByUserId(context.Context, uuid.UUID) ([]*projections.UserRoleBinding, error)
}

// WriteOnlyUserRepository is a repository for reading user projections.
type WriteOnlyUserRoleBindingRepository interface {
	es.WriteOnlyRepository
}

// NewUserRepository creates a repository for reading and writing user projections.
func NewUserRoleBindingRepository(base es.Repository) UserRoleBindingRepository {
	return &userRoleBindingRepository{
		Repository: base,
	}
}

// ByEmail searches for the a user projection by it's email address.
func (r *userRoleBindingRepository) ByUserId(ctx context.Context, userId uuid.UUID) ([]*projections.UserRoleBinding, error) {
	ps, err := r.All(ctx)
	if err != nil {
		return nil, err
	}

	var bindings []*projections.UserRoleBinding
	for _, p := range ps {
		if b, ok := p.(*projections.UserRoleBinding); ok {
			if b.UserId() == userId {
				bindings = append(bindings, b)
			}
		}
	}
	return bindings, nil
}
