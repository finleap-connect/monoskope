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

	"github.com/finleap-connect/monoskope/pkg/domain/errors"
	"github.com/finleap-connect/monoskope/pkg/domain/projections"
	es "github.com/finleap-connect/monoskope/pkg/eventsourcing"
	"github.com/google/uuid"
)

type tenantClusterBindingRepository struct {
	DomainRepository[*projections.TenantClusterBinding]
}

// TenantClusterBindingRepository is a repository for reading and writing tenantclusterbinding projections.
type TenantClusterBindingRepository interface {
	DomainRepository[*projections.TenantClusterBinding]
	GetByTenantId(ctx context.Context, tenantId uuid.UUID) ([]*projections.TenantClusterBinding, error)
	GetByClusterId(ctx context.Context, tenantId uuid.UUID) ([]*projections.TenantClusterBinding, error)
	GetByTenantAndClusterId(ctx context.Context, tenantId, clusterId uuid.UUID) (*projections.TenantClusterBinding, error)
}

// NewTenantClusterBindingRepository creates a repository for reading and writing tenantclusterbinding projections.
func NewTenantClusterBindingRepository(repository es.Repository[*projections.TenantClusterBinding]) TenantClusterBindingRepository {
	return &tenantClusterBindingRepository{
		NewDomainRepository(repository),
	}
}

// GetByTenantId searches for the TenantClusterBinding projections by tenant id.
func (r *tenantClusterBindingRepository) GetByTenantId(ctx context.Context, tenantId uuid.UUID) ([]*projections.TenantClusterBinding, error) {
	ps, err := r.AllWith(ctx, false)
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

// GetByClusterId searches for the TenantClusterBinding projections by cluster id.
func (r *tenantClusterBindingRepository) GetByClusterId(ctx context.Context, clusterId uuid.UUID) ([]*projections.TenantClusterBinding, error) {
	ps, err := r.AllWith(ctx, false)
	if err != nil {
		return nil, err
	}

	var bindings []*projections.TenantClusterBinding
	for _, p := range ps {
		if p.ClusterId == clusterId.String() {
			bindings = append(bindings, p)
		}
	}
	return bindings, nil
}

// GetByTenantAndClusterId searches the TenantClusterBinding projection by tenant and cluster id.
func (r *tenantClusterBindingRepository) GetByTenantAndClusterId(ctx context.Context, tenantId, clusterId uuid.UUID) (*projections.TenantClusterBinding, error) {
	ps, err := r.AllWith(ctx, false)
	if err != nil {
		return nil, err
	}

	for _, p := range ps {
		if p.TenantId == tenantId.String() && p.ClusterId == clusterId.String() {
			return p, nil
		}
	}
	return nil, errors.ErrTenantClusterBindingNotFound
}
