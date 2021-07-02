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
func (h *storingAggregateHandler) HandleCommand(ctx context.Context, cmd es.Command) error {
	var err error
	var aggregate es.Aggregate

	// Load the aggregate from the store
	if aggregate, err = h.aggregateManager.Get(ctx, cmd.AggregateType(), cmd.AggregateID()); err != nil {
		return err
	}

	// Apply the command to the aggregate
	if err := aggregate.HandleCommand(ctx, cmd); err != nil {
		return err
	}

	// Store any emitted events
	return h.aggregateManager.Update(ctx, aggregate)
}
