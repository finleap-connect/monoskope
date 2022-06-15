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
	projections "github.com/finleap-connect/monoskope/pkg/domain/projections"
	es "github.com/finleap-connect/monoskope/pkg/eventsourcing"
)

type tenantRepository struct {
	DomainRepository[*projections.Tenant]
}

// TenantRepository is a repository for reading and writing tenant projections.
type TenantRepository interface {
	DomainRepository[*projections.Tenant]
	// ByName searches for the a tenant projection by it's name
	ByName(context.Context, string) (*projections.Tenant, error)
}

// NewTenantRepository creates a repository for reading and writing tenant projections.
func NewTenantRepository(repository es.Repository[*projections.Tenant]) TenantRepository {
	return &tenantRepository{
		NewDomainRepository(repository),
	}
}

// ByName searches for a tenant projection by its name.
func (r *tenantRepository) ByName(ctx context.Context, name string) (*projections.Tenant, error) {
	ps, err := r.AllWith(ctx, true)
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
