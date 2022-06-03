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

package formatters

import (
	"context"
	"io"

	esApi "github.com/finleap-connect/monoskope/pkg/api/eventsourcing"
	es "github.com/finleap-connect/monoskope/pkg/eventsourcing"
	"github.com/google/uuid"
)

// Snapshotter implements basic snapshot creation for audit formatters
type Snapshotter[T es.Projection] struct {
	esClient  esApi.EventStoreClient
	projector es.Projector[T]
}

func NewSnapshotter[T es.Projection](esClient esApi.EventStoreClient, projector es.Projector[T]) *Snapshotter[T] {
	return &Snapshotter[T]{esClient, projector}
}

// CreateSnapshot creates a snapshot based on an event-filter and the corresponding projector for
// the aggregate of which the id is used in the filter.
// This is a temporary implementation until snapshots are fully implemented,
// and it is not meant to be used extensively.
func (s *Snapshotter[T]) CreateSnapshot(ctx context.Context, eventFilter *esApi.EventFilter) (T, error) {
	projection := s.projector.NewProjection(uuid.New())
	aggregateEvents, err := s.esClient.Retrieve(ctx, eventFilter)

	var nilResult T
	if err != nil {
		return nilResult, err
	}

	for {
		e, err := aggregateEvents.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nilResult, err
		}

		event, err := es.NewEventFromProto(e)
		if err != nil {
			return nilResult, err
		}

		projection, err = s.projector.Project(ctx, event, projection)
		if err != nil {
			return nilResult, err
		}
	}

	return projection, nil
}
