package commandhandler

import (
	"context"

	es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing"
)

type storingAggregateHandler struct {
	aggregateManager es.AggregateStore
}

// NewAggregateHandler creates a new CommandHandler which handles aggregates.
func NewAggregateHandler(aggregateManager es.AggregateStore) es.CommandHandler {
	return &storingAggregateHandler{
		aggregateManager: aggregateManager,
	}
}

// HandleCommand implements the CommandHandler interface
func (h *storingAggregateHandler) HandleCommand(ctx context.Context, cmd es.Command) (*es.CommandReply, error) {
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

	return reply, err
}
