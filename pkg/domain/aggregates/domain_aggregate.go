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

package aggregates

import (
	"context"

	domainErrors "github.com/finleap-connect/monoskope/pkg/domain/errors"
	es "github.com/finleap-connect/monoskope/pkg/eventsourcing"
)

type DomainAggregateBase struct {
	*es.BaseAggregate
}

// Validate validates if the aggregate exists and is not deleted
func (a *DomainAggregateBase) Validate(ctx context.Context, cmd es.Command) error {
	if !a.Exists() {
		return domainErrors.ErrNotFound
	}
	if a.Deleted() {
		return domainErrors.ErrDeleted
	}
	return nil
}

// DefaultReply creates a default command reply containing the ID and current version of the aggregate.
func (a *DomainAggregateBase) DefaultReply() *es.CommandReply {
	return &es.CommandReply{
		Id:      a.ID(),
		Version: a.Version(),
	}
}
