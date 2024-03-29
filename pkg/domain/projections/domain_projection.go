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
	"google.golang.org/protobuf/types/known/timestamppb"
)

type domainProjection struct {
	projections.LifecycleMetadata
	version uint64
}
type DomainProjection interface {
	LifecycleMetadata
	ID() uuid.UUID
	Version() uint64
	IncrementVersion()
	GetLifecycleMetadata() *projections.LifecycleMetadata
}

type LifecycleMetadata interface {
	// GetCreated returns when the projection has been created
	GetCreated() *timestamppb.Timestamp
	// GetCreatedById returns by whom the projection has been created
	GetCreatedById() string
	// GetLastModified returns when the projection has been last modified
	GetLastModified() *timestamppb.Timestamp
	// GetLastModifiedById returns by whom the projection has been last modified
	GetLastModifiedById() string
	// GetDeletedById returns ry whom the projection has been deleted
	GetDeletedById() string
	// GetDeleted returns when the projection has been deleted
	GetDeleted() *timestamppb.Timestamp
	// IsDeleted returns if the projection has been deleted
	IsDeleted() bool
}

func NewDomainProjection() DomainProjection {
	return &domainProjection{}
}

// ID implements the ID method of the Projection interface.
func (p *domainProjection) ID() uuid.UUID {
	panic("not implemented")
}

// Version implements the Version method of the Projection interface.
func (p *domainProjection) Version() uint64 {
	return p.version
}

// IncrementVersion implements the IncrementVersion method of the Projection interface.
func (p *domainProjection) IncrementVersion() {
	p.version++
}

// GetLifecycleMetadata implements the GetLifecycleMetadata method of the Projection interface.
func (p *domainProjection) GetLifecycleMetadata() *projections.LifecycleMetadata {
	return &p.LifecycleMetadata
}

// IsDeleted implements DomainProjection
func (p *domainProjection) IsDeleted() bool {
	return p.Deleted != nil
}
