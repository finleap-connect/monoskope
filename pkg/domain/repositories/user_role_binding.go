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

	"github.com/google/uuid"
	projections "gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/projections"
	es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing"
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
	// ByUserId searches for all UserRoleBinding projection's by the a user id.
	ByUserId(context.Context, uuid.UUID) ([]*projections.UserRoleBinding, error)
	// ByScopeAndResource returns all UserRoleBinding projections matching the given scope and resource.
	ByScopeAndResource(context.Context, es.Scope, uuid.UUID) ([]*projections.UserRoleBinding, error)
}

// NewUserRepository creates a repository for reading and writing UserRoleBinding projections.
func NewUserRoleBindingRepository(repository es.Repository) UserRoleBindingRepository {
	return &userRoleBindingRepository{
		Repository: repository,
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
