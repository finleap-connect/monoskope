package commandhandler

import (
	"context"

	es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing"
)

type aggregateHandler struct {
	aggregateType es.AggregateType
}

// NewAggregateHandler creates a new CommandHandler which handles aggregates.
func NewAggregateHandler(aggregateType es.AggregateType) es.CommandHandler {
	return &aggregateHandler{
		aggregateType: aggregateType,
	}
}

// HandleCommand implements the CommandHandler interface
func (h *aggregateHandler) HandleCommand(ctx context.Context, cmd es.Command) error {
	return nil
}
