// Copyright 2022 Monoskope Authors
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
	"github.com/google/uuid"
)

type userRepository struct {
	DomainRepository[*projections.User]
	roleBindingRepo UserRoleBindingRepository
}

// UserRepository is a repository for reading and writing user projections.
type UserRepository interface {
	DomainRepository[*projections.User]
	// ByUserId searches for the a user projection by it's id.
	ByUserId(context.Context, uuid.UUID) (*projections.User, error)
	// ByEmail searches for the a user projection by it's email address.
	ByEmail(context.Context, string) (*projections.User, error)
	// ByEmail searches for the a user projection by it's email address.
	ByEmailIncludingDeleted(context.Context, string) ([]*projections.User, error)
	// GetCount returns the user count
	GetCount(context.Context, bool) (int, error)
}

// NewUserRepository creates a repository for reading and writing user projections.
func NewUserRepository(repository es.Repository[*projections.User], roleBindingRepo UserRoleBindingRepository) UserRepository {
	return &userRepository{
		DomainRepository: NewDomainRepository(repository),
		roleBindingRepo:  roleBindingRepo,
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

// ByUserId searches for the a user projection by it's id.
func (r *userRepository) ByUserId(ctx context.Context, id uuid.UUID) (*projections.User, error) {
	user, err := r.ById(ctx, id)
	if err != nil {
		return nil, err
	}

	// Find roles of user
	err = r.addRolesToUser(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// ByEmail searches for a user projection by it's email address.
func (r *userRepository) ByEmail(ctx context.Context, email string) (*projections.User, error) {
	users, err := r.AllWith(ctx, false)
	if err != nil {
		return nil, err
	}

	for _, user := range users {
		if email == user.Email {
			// Find roles of user
			err = r.addRolesToUser(ctx, user)
			if err != nil {
				return nil, err
			}

			return user, nil
		}
	}

	return nil, errors.ErrUserNotFound
}

// ByEmail searches for a user projection by it's email address.
func (r *userRepository) ByEmailIncludingDeleted(ctx context.Context, email string) ([]*projections.User, error) {
	users, err := r.AllWith(ctx, true)
	if err != nil {
		return nil, err
	}

	var result []*projections.User
	for _, user := range users {
		if email == user.Email {
			// Find roles of user
			err = r.addRolesToUser(ctx, user)
			if err != nil {
				return nil, err
			}
			result = append(result, user)
		}
	}

	return result, nil
}

// All searches for all user projections.
func (r *userRepository) GetCount(ctx context.Context, includeDeleted bool) (int, error) {
	users, err := r.AllWith(ctx, includeDeleted)
	if err != nil {
		return 0, err
	}
	return len(users), nil
}
