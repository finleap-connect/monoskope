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
	"sort"

	"github.com/finleap-connect/monoskope/pkg/domain/errors"
	projections "github.com/finleap-connect/monoskope/pkg/domain/projections"
	es "github.com/finleap-connect/monoskope/pkg/eventsourcing"
	esErrors "github.com/finleap-connect/monoskope/pkg/eventsourcing/errors"
	"github.com/google/uuid"
)

type tenantRepository struct {
	es.Repository
}

// TenantRepository is a repository for reading and writing tenant projections.
type TenantRepository interface {
	es.Repository
	// ById searches for the a tenant projection by it's id.
	ByTenantId(context.Context, string) (*projections.Tenant, error)
	// ByName searches for the a tenant projection by it's name
	ByName(context.Context, string) (*projections.Tenant, error)
	// GetAll searches for all tenant projections.
	GetAll(context.Context, bool) ([]*projections.Tenant, error)
}

// NewTenantRepository creates a repository for reading and writing tenant projections.
func NewTenantRepository(repository es.Repository) TenantRepository {
	return &tenantRepository{
		Repository: repository,
	}
}

// ByTenantId searches for a tenant projection by its id.
func (r *tenantRepository) ByTenantId(ctx context.Context, id string) (*projections.Tenant, error) {
	uuid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	projection, err := r.ById(ctx, uuid)
	if err != nil {
		return nil, err
	}

	tenant, ok := projection.(*projections.Tenant)
	if !ok {
		return nil, esErrors.ErrInvalidProjectionType
	}

	return tenant, nil
}

// ByTenantName searches for a tenant projection by its name.
func (r *tenantRepository) ByName(ctx context.Context, name string) (*projections.Tenant, error) {
	ps, err := r.GetAll(ctx, true)
	if err != nil {
		return nil, err
	}

	for _, t := range ps {
		if name == t.Name {
			return t, nil
		}
	}

	return nil, errors.ErrTenantNotFound
}

// All searches for the a tenant projections.
func (r *tenantRepository) GetAll(ctx context.Context, includeDeleted bool) ([]*projections.Tenant, error) {
	ps, err := r.All(ctx)
	if err != nil {
		return nil, err
	}

	var tenants []*projections.Tenant
	for _, p := range ps {
		if t, ok := p.(*projections.Tenant); ok {
			if !t.GetDeleted().IsValid() || includeDeleted {
				tenants = append(tenants, t)
			}
		} else {
			return nil, esErrors.ErrInvalidProjectionType
		}
	}

	sort.Slice(tenants, func(i, j int) bool {
		return tenants[i].Name > tenants[j].Name
	})

	return tenants, nil
}
