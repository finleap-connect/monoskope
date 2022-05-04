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
	esApi "github.com/finleap-connect/monoskope/pkg/api/eventsourcing"
	es "github.com/finleap-connect/monoskope/pkg/eventsourcing"
	"github.com/google/uuid"
	"io"
)

// FormatterBase is the base implementation for all audit formatters
type FormatterBase struct {
	EsClient esApi.EventStoreClient
}

// CreateSnapshot creates a snapshot based on an event-filter and the corresponding projector for
// the aggregate of which the id is used in the filter.
// This is a temporary implementation until snapshots are fully implemented,
// and it is not meant to be used extensively.
func (f *FormatterBase) CreateSnapshot(ctx context.Context, projector es.Projector, eventFilter *esApi.EventFilter) (es.Projection, error) {
	projection := projector.NewProjection(uuid.New())
	aggregateEvents, err := f.EsClient.Retrieve(ctx, eventFilter)
	if err != nil {
		return nil, err
	}

	for {
		e, err := aggregateEvents.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		event, err := es.NewEventFromProto(e)
		if err != nil {
			return nil, err
		}

		projection, err = projector.Project(ctx, event, projection)
		if err != nil {
			return nil, err
		}

	}

	return projection, nil
}
