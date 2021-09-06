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

package commandhandler

import (
	"context"
	"sync"

	es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing"
)

type storingAggregateHandler struct {
	aggregateManager es.AggregateStore
	mutex            sync.Mutex
}

// NewAggregateHandler creates a new CommandHandler which handles aggregates.
func NewAggregateHandler(aggregateManager es.AggregateStore) es.CommandHandler {
	return &storingAggregateHandler{
		aggregateManager: aggregateManager,
	}
}

// HandleCommand implements the CommandHandler interface
func (h *storingAggregateHandler) HandleCommand(ctx context.Context, cmd es.Command) (*es.CommandReply, error) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	var err error
	var aggregate es.Aggregate

	// Load the aggregate from the store
	if aggregate, err = h.aggregateManager.Get(ctx, cmd.AggregateType(), cmd.AggregateID()); err != nil {
		return nil, err
	}

	// Apply the command to the aggregate
	reply, err := aggregate.HandleCommand(ctx, cmd)
	if err != nil {
		return nil, err
	}

	// Store any emitted events
	err = h.aggregateManager.Update(ctx, aggregate)
	if err != nil {
		return nil, err
	}

	// Set the version the aggregate now has after handling the command.
	reply.Version = aggregate.Version()

	return reply, err
}
