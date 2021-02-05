package repositories

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	projections "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/queryhandler"
	es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/event_sourcing"
)

type userRepository struct {
	userRepo        es.Repository
	roleBindingRepo UserRoleBindingRepository
}

// Repository is a repository for reading and writing user projections.
type UserRepository interface {
	ReadOnlyUserRepository
	WriteOnlyUserRepository
}

// ReadOnlyUserRepository is a repository for reading user projections.
type ReadOnlyUserRepository interface {
	// ByEmail searches for the a user projection by it's email address.
	ByEmail(context.Context, string) (*projections.User, error)
}

// WriteOnlyUserRepository is a repository for reading user projections.
type WriteOnlyUserRepository interface {
}

// NewUserRepository creates a repository for reading and writing user projections.
func NewUserRepository(userRepo es.Repository, roleBindingRepo UserRoleBindingRepository) UserRepository {
	return &userRepository{
		userRepo:        userRepo,
		roleBindingRepo: roleBindingRepo,
	}
}

// ByEmail searches for the a user projection by it's email address.
func (r *userRepository) ByEmail(ctx context.Context, email string) (*projections.User, error) {
	ps, err := r.userRepo.All(ctx)
	if err != nil {
		return nil, err
	}

	for _, p := range ps {
		if u, ok := p.(*projections.User); ok {
			if email == u.Email {
				roles, err := r.roleBindingRepo.ByUserId(ctx, uuid.MustParse(u.GetId()))
				if err != nil {
					return nil, err
				}

				u.Roles = roles
				return u, nil
			}
		}
	}
	return nil, fmt.Errorf("not found")
}
