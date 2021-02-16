package commandhandler

import (
	"context"

	es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing/errors"
)

type storingAggregateHandler struct {
	aggregateManager es.AggregateManager
}

// NewAggregateHandler creates a new CommandHandler which handles aggregates.
func NewAggregateHandler(aggregateStore es.AggregateManager) es.CommandHandler {
	return &storingAggregateHandler{
		aggregateManager: aggregateStore,
	}
}

// HandleCommand implements the CommandHandler interface
func (h *storingAggregateHandler) HandleCommand(ctx context.Context, cmd es.Command) error {
	var aggregate es.Aggregate

	// Load the aggregate from the store
	if aggregate, err := h.aggregateManager.Get(ctx, cmd.AggregateType(), cmd.AggregateID()); err != nil {
		return err
	} else if aggregate == nil {
		return errors.ErrAggregateNotFound
	}

	// Apply the command to the aggregate
	if err := aggregate.HandleCommand(ctx, cmd); err != nil {
		return err
	}

	// Store any emitted events
	return h.aggregateManager.Update(ctx, aggregate)
}
