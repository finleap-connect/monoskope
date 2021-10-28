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

	projections "github.com/finleap-connect/monoskope/pkg/domain/projections"
	es "github.com/finleap-connect/monoskope/pkg/eventsourcing"
	esErrors "github.com/finleap-connect/monoskope/pkg/eventsourcing/errors"
	"github.com/google/uuid"
)

type tenantClusterBindingRepository struct {
	es.Repository
}

// TenantClusterBindingRepository is a repository for reading and writing tenantclusterbinding projections.
type TenantClusterBindingRepository interface {
	es.Repository
	ReadOnlyTenantClusterBindingRepository
	WriteOnlyTenantClusterBindingRepository
}

// ReadOnlyTenantClusterBindingRepository is a repository for reading tenantclusterbinding projections.
type ReadOnlyTenantClusterBindingRepository interface {
	// GetAll searches for the a TenantClusterBinding projections.
	GetAll(context.Context, bool) ([]*projections.TenantClusterBinding, error)
	GetByTenantId(context.Context, uuid.UUID) ([]*projections.TenantClusterBinding, error)
}

// WriteOnlyTenantClusterBindingRepository is a repository for writing tenantclusterbinding projections.
type WriteOnlyTenantClusterBindingRepository interface {
}

// NewTenantClusterBindingRepository creates a repository for reading and writing tenantclusterbinding projections.
func NewTenantClusterBindingRepository(repository es.Repository) TenantClusterBindingRepository {
	return &tenantClusterBindingRepository{
		Repository: repository,
	}
}

// GetAll searches for the a TenantClusterBinding projections.
func (r *tenantClusterBindingRepository) GetAll(ctx context.Context, includeDeleted bool) ([]*projections.TenantClusterBinding, error) {
	ps, err := r.All(ctx)
	if err != nil {
		return nil, err
	}

	var bindings []*projections.TenantClusterBinding
	for _, p := range ps {
		if t, ok := p.(*projections.TenantClusterBinding); ok {
			if !t.GetDeleted().IsValid() || includeDeleted {
				bindings = append(bindings, t)
			}
		} else {
			return nil, esErrors.ErrInvalidProjectionType
		}
	}
	return bindings, nil
}

// GetAll searches for the a TenantClusterBinding projections.
func (r *tenantClusterBindingRepository) GetByTenantId(ctx context.Context, tenantId uuid.UUID) ([]*projections.TenantClusterBinding, error) {
	ps, err := r.GetAll(ctx, false)
	if err != nil {
		return nil, err
	}

	var bindings []*projections.TenantClusterBinding
	for _, p := range ps {
		if p.TenantId == tenantId.String() {
			bindings = append(bindings, p)
		}
	}
	return bindings, nil
}
