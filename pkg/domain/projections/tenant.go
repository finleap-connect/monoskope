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

package projections

import (
	"github.com/finleap-connect/monoskope/pkg/api/domain/projections"
	"github.com/google/uuid"
)

type Tenant struct {
	DomainProjection
	*projections.Tenant
}

func NewTenantProjection(id uuid.UUID) *Tenant {
	dp := NewDomainProjection()
	return &Tenant{
		DomainProjection: dp,
		Tenant: &projections.Tenant{
			Id:       id.String(),
			Metadata: dp.GetLifecycleMetadata(),
		},
	}
}

// ID implements the ID method of the Aggregate interface.
func (p *Tenant) ID() uuid.UUID {
	return uuid.MustParse(p.Id)
}

// Proto gets the underlying proto representation.
func (p *Tenant) Proto() *projections.Tenant {
	return p.Tenant
}
