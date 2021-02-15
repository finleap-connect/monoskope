package commandhandler

import (
	"context"

	es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing/errors"
)

type storingAggregateHandler struct {
	aggregateType  es.AggregateType
	aggregateStore es.AggregateRepository
}

// NewStoringAggregateHandler creates a new CommandHandler which handles aggregates.
func NewStoringAggregateHandler(aggregateType es.AggregateType, aggregateStore es.AggregateRepository) es.CommandHandler {
	return &storingAggregateHandler{
		aggregateType:  aggregateType,
		aggregateStore: aggregateStore,
	}
}

// HandleCommand implements the CommandHandler interface
func (h *storingAggregateHandler) HandleCommand(ctx context.Context, cmd es.Command) error {
	var aggregate es.Aggregate

	// Load the aggregate from the store
	if aggregate, err := h.aggregateStore.Get(ctx, cmd.AggregateType(), cmd.AggregateID()); err != nil {
		return err
	} else if aggregate == nil {
		return errors.ErrAggregateNotFound
	}

	// Apply the command to the aggregate
	if err := aggregate.HandleCommand(ctx, cmd); err != nil {
		return err
	}

	// Store any emitted events
	return h.aggregateStore.Update(ctx, aggregate)
}
