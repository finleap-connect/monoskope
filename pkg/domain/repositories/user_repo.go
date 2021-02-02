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

type UserRepository interface {
	ByEmail(context.Context, string) (*projections.User, error)
}

func NewUserRepository(base es.Repository) UserRepository {
	return &userRepository{
		Repository: base,
	}
}

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
