package event_sourcing

import (
	"context"
)

// CommandHandler is an interface that all handlers of commands must implement.
type CommandHandler interface {
	// HandleCommand handles a command.
	HandleCommand(context.Context, Command) error
}

type commandHandlerChain struct {
	handlers []CommandHandler
}

// HandleCommand handles all commands chained in this chain.
func (c *commandHandlerChain) HandleCommand(ctx context.Context, cmd Command) error {
	for _, handler := range c.handlers {
		if err := handler.HandleCommand(ctx, cmd); err != nil {
			return err
		}
	}

	return nil
}

// ChainCommandHandler wraps a CommandHandler in one or more middleware.
func ChainCommandHandler(handlers ...CommandHandler) CommandHandler {
	return &commandHandlerChain{
		handlers: handlers,
	}
}

// EventHandler is an interface that all handlers of events must implement.
type EventHandler interface {
	// HandleEvent handles an event.
	HandleEvent(context.Context, Event) error
}
