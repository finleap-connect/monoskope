package repositories

import (
	"context"
	"fmt"

	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/projections"
	es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/event_sourcing"
)

type userRepository struct {
	es.Repository
}

// Repository is a repository for reading and writing user projections.
type UserRepository interface {
	ReadOnlyUserRepository
	WriteOnlyUserRepository
}

// ReadOnlyUserRepository is a repository for reading user projections.
type ReadOnlyUserRepository interface {
	es.ReadOnlyRepository

	// ByEmail searches for the a user projection by it's email address.
	ByEmail(context.Context, string) (*projections.User, error)
}

// WriteOnlyUserRepository is a repository for reading user projections.
type WriteOnlyUserRepository interface {
	es.WriteOnlyRepository
}

// NewUserRepository creates a repository for reading and writing user projections.
func NewUserRepository(base es.Repository) UserRepository {
	return &userRepository{
		Repository: base,
	}
}

// ByEmail searches for the a user projection by it's email address.
func (r *userRepository) ByEmail(ctx context.Context, email string) (*projections.User, error) {
	ps, err := r.All(ctx)
	if err != nil {
		return nil, err
	}

	for _, p := range ps {
		if u, ok := p.(*projections.User); ok {
			if u.Email == email {
				return u, nil
			}
		}
	}
	return nil, fmt.Errorf("not found")
}
