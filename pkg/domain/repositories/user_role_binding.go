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

	projections "github.com/finleap-connect/monoskope/pkg/domain/projections"
	es "github.com/finleap-connect/monoskope/pkg/eventsourcing"
	"github.com/google/uuid"
)

type userRoleBindingRepository struct {
	DomainRepository[*projections.UserRoleBinding]
}

// UserRoleBindingRepository is a repository for reading and writing UserRoleBinding projections.
type UserRoleBindingRepository interface {
	DomainRepository[*projections.UserRoleBinding]
	// ByUserId searches for all UserRoleBinding projection's by the user id.
	ByUserId(context.Context, uuid.UUID) ([]*projections.UserRoleBinding, error)
	// ByUserIdScopeAndResource searches for all UserRoleBinding projection's by the user id and the scope.
	ByUserIdScopeAndResource(context.Context, uuid.UUID, es.Scope, string) ([]*projections.UserRoleBinding, error)
	// ByScopeAndResource returns all UserRoleBinding projections matching the given scope and resource.
	ByScopeAndResource(context.Context, es.Scope, uuid.UUID) ([]*projections.UserRoleBinding, error)
}

// NewUserRoleBindingRepository creates a repository for reading and writing UserRoleBinding projections.
func NewUserRoleBindingRepository(repository es.Repository[*projections.UserRoleBinding]) UserRoleBindingRepository {
	return &userRoleBindingRepository{
		NewDomainRepository(repository),
	}
}

// ByUserId searches for all UserRoleBinding projection's by the a user id.
func (r *userRoleBindingRepository) ByUserId(ctx context.Context, userId uuid.UUID) ([]*projections.UserRoleBinding, error) {
	ps, err := r.AllWith(ctx, false)
	if err != nil {
		return nil, err
	}

	var userRoleBindings []*projections.UserRoleBinding
	for _, userRoleBinding := range ps {
		if userId.String() == userRoleBinding.GetUserId() {
			userRoleBindings = append(userRoleBindings, userRoleBinding)
		}
	}
	return userRoleBindings, nil
}

// ByUserIdScopeAndResource searches for all UserRoleBinding projection's by the a user id and the scope.
func (r *userRoleBindingRepository) ByUserIdScopeAndResource(ctx context.Context, userId uuid.UUID, scope es.Scope, resource string) ([]*projections.UserRoleBinding, error) {
	ps, err := r.AllWith(ctx, false)
	if err != nil {
		return nil, err
	}

	var userRoleBindings []*projections.UserRoleBinding
	for _, userRoleBinding := range ps {
		if userId.String() == userRoleBinding.GetUserId() && userRoleBinding.Scope == string(scope) && userRoleBinding.Resource == resource {
			userRoleBindings = append(userRoleBindings, userRoleBinding)
		}
	}
	return userRoleBindings, nil
}

// ByScopeAndResource returns all UserRoleBinding projections matching the given scope and resource.
func (r *userRoleBindingRepository) ByScopeAndResource(ctx context.Context, scope es.Scope, resource uuid.UUID) ([]*projections.UserRoleBinding, error) {
	ps, err := r.AllWith(ctx, false)
	if err != nil {
		return nil, err
	}

	var userRoleBindings []*projections.UserRoleBinding
	for _, userRoleBinding := range ps {
		if string(scope) == userRoleBinding.GetScope() && resource.String() == userRoleBinding.Resource {
			userRoleBindings = append(userRoleBindings, userRoleBinding)
		}
	}
	return userRoleBindings, nil
}
