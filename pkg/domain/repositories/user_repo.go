// Copyright 2021 Monoskope Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package repositories

import (
	"context"

	projectionsApi "github.com/finleap-connect/monoskope/pkg/api/domain/projections"
	"github.com/finleap-connect/monoskope/pkg/domain/errors"
	projections "github.com/finleap-connect/monoskope/pkg/domain/projections"
	es "github.com/finleap-connect/monoskope/pkg/eventsourcing"
	esErrors "github.com/finleap-connect/monoskope/pkg/eventsourcing/errors"
	"github.com/google/uuid"
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
	// GetAll searches for all user projection.
	GetAll(context.Context, bool) ([]*projections.User, error)
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
	ps, err := r.GetAll(ctx, true)
	if err != nil {
		return nil, err
	}

	for _, u := range ps {
		if email == u.Email {
			return u, nil
		}
	}

	return nil, errors.ErrUserNotFound
}

// All searches for all user projections.
func (r *userRepository) GetAll(ctx context.Context, includeDeleted bool) ([]*projections.User, error) {
	ps, err := r.All(ctx)
	if err != nil {
		return nil, err
	}

	var users []*projections.User
	for _, p := range ps {
		if u, ok := p.(*projections.User); ok {
			// Find roles of user
			err = r.addRolesToUser(ctx, u)
			if err != nil {
				return nil, err
			}
			users = append(users, u)
		}
	}
	return users, nil
}
