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

package projections

import (
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/domain/projections"
)

type DomainProjection struct {
	projections.LifecycleMetadata
	version uint64
}

func NewDomainProjection() *DomainProjection {
	return &DomainProjection{}
}

// Version implements the Version method of the Projection interface.
func (p *DomainProjection) Version() uint64 {
	return p.version
}

// IncrementVersion implements the IncrementVersion method of the Projection interface.
func (p *DomainProjection) IncrementVersion() {
	p.version++
}
