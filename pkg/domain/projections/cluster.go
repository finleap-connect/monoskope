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
	"github.com/google/uuid"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/domain/projections"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing"
)

type Cluster struct {
	*DomainProjection
	*projections.Cluster
}

func NewClusterProjection(id uuid.UUID) eventsourcing.Projection {
	dp := NewDomainProjection()
	return &Cluster{
		DomainProjection: dp,
		Cluster: &projections.Cluster{
			Id:       id.String(),
			Metadata: &dp.LifecycleMetadata,
		},
	}
}

// ID implements the ID method of the Projection interface.
func (p *Cluster) ID() uuid.UUID {
	return uuid.MustParse(p.Id)
}

// Proto gets the underlying proto representation.
func (p *Cluster) Proto() *projections.Cluster {
	return p.Cluster
}
