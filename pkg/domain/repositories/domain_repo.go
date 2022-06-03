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
)

type domainRepository[T projections.DomainProjection] struct {
	es.Repository[T]
}

// Repository is a repository for reading projections.
type DomainRepository[T projections.DomainProjection] interface {
	es.Repository[T]
	// All returns all projections in the repository.
	AllWith(ctx context.Context, includeDeleted bool) ([]T, error)
}

// NewDomainRepository creates a repository for reading and writing domain projections.
func NewDomainRepository[T projections.DomainProjection](repository es.Repository[T]) DomainRepository[T] {
	return &domainRepository[T]{
		repository,
	}
}

// AllWith returns all projections in the repository.
func (r *domainRepository[T]) AllWith(ctx context.Context, includeDeleted bool) ([]T, error) {
	ps, err := r.Repository.All(ctx)
	if err != nil {
		return nil, err
	}
	var projections []T
	for _, p := range ps {
		if !p.GetDeleted().IsValid() || includeDeleted {
			projections = append(projections, p)
		}
	}
	return projections, nil
}
