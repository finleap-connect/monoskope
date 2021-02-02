package user

import (
	"context"
	"fmt"

	es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/event_sourcing"
)

type userRepository struct {
	es.Repository
}

// Repository is a repository for reading and writing user projections.
type UserRepository interface {
	es.Repository

	// ByEmail searches for the a user projection by it's email address.
	ByEmail(context.Context, string) (*User, error)
}

// NewUserRepository creates a repository for reading and writing user projections.
func NewUserRepository(base es.Repository) UserRepository {
	return &userRepository{
		Repository: base,
	}
}

// ByEmail searches for the a user projection by it's email address.
func (r *userRepository) ByEmail(ctx context.Context, email string) (*User, error) {
	ps, err := r.All(ctx)
	if err != nil {
		return nil, err
	}

	for _, p := range ps {
		if u, ok := p.(*User); ok {
			if u.Email == email {
				return u, nil
			}
		}
	}
	return nil, fmt.Errorf("not found")
}
