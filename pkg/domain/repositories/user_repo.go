package repositories

import (
	"context"

	"github.com/google/uuid"
	projectionsApi "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/domain/projections"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/errors"
	projections "gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/projections"
	es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing"
	esErrors "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing/errors"
)

type userRepository struct {
	es.Repository
	roleBindingRepo UserRoleBindingRepository
}

// Repository is a repository for reading and writing user projections.
type UserRepository interface {
	es.Repository
	ReadOnlyUserRepository
	WriteOnlyUserRepository
}

// ReadOnlyUserRepository is a repository for reading user projections.
type ReadOnlyUserRepository interface {
	// ById searches for the a user projection by it's id.
	ByUserId(context.Context, uuid.UUID) (*projections.User, error)
	// ByEmail searches for the a user projection by it's email address.
	ByEmail(context.Context, string) (*projections.User, error)
}

// WriteOnlyUserRepository is a repository for writing user projections.
type WriteOnlyUserRepository interface {
}

// NewUserRepository creates a repository for reading and writing user projections.
func NewUserRepository(repository es.Repository, roleBindingRepo UserRoleBindingRepository) UserRepository {
	return &userRepository{
		Repository:      repository,
		roleBindingRepo: roleBindingRepo,
	}
}

func (r *userRepository) addRolesToUser(ctx context.Context, user *projections.User) error {
	// Find roles of user
	roles, err := r.roleBindingRepo.ByUserId(ctx, user.ID())
	if err != nil {
		return err
	}
	user.Roles = toProtoRoles(roles)

	return nil
}

func toProtoRoles(roles []*projections.UserRoleBinding) []*projectionsApi.UserRoleBinding {
	var mapped []*projectionsApi.UserRoleBinding
	for _, role := range roles {
		mapped = append(mapped, role.Proto())
	}
	return mapped
}

// ById searches for the a user projection by it's id.
func (r *userRepository) ByUserId(ctx context.Context, id uuid.UUID) (*projections.User, error) {
	projection, err := r.ById(ctx, id)
	if err != nil {
		return nil, err
	}

	if user, ok := projection.(*projections.User); !ok {
		return nil, esErrors.ErrInvalidProjectionType
	} else {
		// Find roles of user
		err = r.addRolesToUser(ctx, user)
		if err != nil {
			return nil, err
		}

		return user, nil
	}
}

// ByEmail searches for the a user projection by it's email address.
func (r *userRepository) ByEmail(ctx context.Context, email string) (*projections.User, error) {
	ps, err := r.All(ctx)
	if err != nil {
		return nil, err
	}

	var user *projections.User
	for _, p := range ps {
		if u, ok := p.(*projections.User); ok {
			if email == u.Email {
				// User found
				user = u
			}
		}
	}

	if user != nil {
		// Find roles of user
		err = r.addRolesToUser(ctx, user)
		if err != nil {
			return nil, err
		}
		return user, nil
	}

	return nil, errors.ErrUserNotFound
}
