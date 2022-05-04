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
	esErrors "github.com/finleap-connect/monoskope/pkg/eventsourcing/errors"

	projections "github.com/finleap-connect/monoskope/pkg/domain/projections"
	es "github.com/finleap-connect/monoskope/pkg/eventsourcing"
	"github.com/google/uuid"
)

type userRoleBindingRepository struct {
	es.Repository
}

// Repository is a repository for reading and writing UserRoleBinding projections.
type UserRoleBindingRepository interface {
	es.Repository
	ReadOnlyUserRoleBindingRepository
}

// ReadOnlyUserRepository is a repository for reading UserRoleBinding projections.
type ReadOnlyUserRoleBindingRepository interface {
	// GetByUserRoleBindingId searches for the UserRoleBinding projections by its id.
	GetByUserRoleBindingId(ctx context.Context, id string) (*projections.UserRoleBinding, error)
	// ByUserId searches for all UserRoleBinding projection's by the user id.
	ByUserId(context.Context, uuid.UUID) ([]*projections.UserRoleBinding, error)
	// ByUserIdAndScope searches for all UserRoleBinding projection's by the user id and the scope.
	ByUserIdAndScope(context.Context, uuid.UUID, es.Scope) ([]*projections.UserRoleBinding, error)
	// ByScopeAndResource returns all UserRoleBinding projections matching the given scope and resource.
	ByScopeAndResource(context.Context, es.Scope, uuid.UUID) ([]*projections.UserRoleBinding, error)
}

// NewUserRepository creates a repository for reading and writing UserRoleBinding projections.
func NewUserRoleBindingRepository(repository es.Repository) UserRoleBindingRepository {
	return &userRoleBindingRepository{
		Repository: repository,
	}
}

// GetByUserRoleBindingId searches for the UserRoleBinding projections by its id.
func (r *userRoleBindingRepository) GetByUserRoleBindingId(ctx context.Context, id string) (*projections.UserRoleBinding, error) {
	projectionUUID, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	projection, err := r.ById(ctx, projectionUUID)
	if err != nil {
		return nil, err
	}

	userRoleBinding, ok := projection.(*projections.UserRoleBinding)
	if !ok {
		return nil, esErrors.ErrInvalidProjectionType
	}

	return userRoleBinding, nil
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

// ByUserIdAndScope searches for all UserRoleBinding projection's by the a user id and the scope.
func (r *userRoleBindingRepository) ByUserIdAndScope(ctx context.Context, userId uuid.UUID, scope es.Scope) ([]*projections.UserRoleBinding, error) {
	ps, err := r.All(ctx)
	if err != nil {
		return nil, err
	}

	var userRoleBindings []*projections.UserRoleBinding
	for _, projection := range ps {
		if userRoleBinding, ok := projection.(*projections.UserRoleBinding); ok {
			if userId.String() == userRoleBinding.GetUserId() && userRoleBinding.Scope == scope.String() {
				userRoleBindings = append(userRoleBindings, userRoleBinding)
			}
		}
	}
	return userRoleBindings, nil
}

// ByScopeAndResource returns all UserRoleBinding projections matching the given scope and resource.
func (r *userRoleBindingRepository) ByScopeAndResource(ctx context.Context, scope es.Scope, resource uuid.UUID) ([]*projections.UserRoleBinding, error) {
	ps, err := r.All(ctx)
	if err != nil {
		return nil, err
	}

	var userRoleBindings []*projections.UserRoleBinding
	for _, projection := range ps {
		if userRoleBinding, ok := projection.(*projections.UserRoleBinding); ok {
			if scope.String() == userRoleBinding.GetScope() && resource.String() == userRoleBinding.Resource {
				userRoleBindings = append(userRoleBindings, userRoleBinding)
			}
		}
	}
	return userRoleBindings, nil
}
