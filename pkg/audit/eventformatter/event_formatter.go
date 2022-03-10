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

package eventformatter

import (
	"context"
	"fmt"
	esApi "github.com/finleap-connect/monoskope/pkg/api/eventsourcing"
	es "github.com/finleap-connect/monoskope/pkg/eventsourcing"
	"github.com/google/uuid"
	"io"
	"strings"
)

const (
	AppendUpdateTo   = "\n- “%s“ to “%s“"
	AppendUpdateFrom = " from “%s“"
)

type EventFormatter interface {
	GetFormattedDetails(context.Context, *esApi.Event) (string, error)
}

type BaseEventFormatter struct {
	EsClient esApi.EventStoreClient
}

// TODO: find a better place, move to domain package?
// TODO: ticket: domain -> snapshots (same idea as e.g. domain -> projections)

func (f *BaseEventFormatter) CreateSnapshot(ctx context.Context, projector es.Projector, eventFilter *esApi.EventFilter) (es.Projection, error) {
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

func (f *BaseEventFormatter) AppendUpdate(field string, update string, old string, strBuilder *strings.Builder) {
	if update != "" {
		strBuilder.WriteString(fmt.Sprintf(AppendUpdateTo, field, update))
		if old != "" {
			strBuilder.WriteString(fmt.Sprintf(AppendUpdateFrom, old))
		}
	}
}
