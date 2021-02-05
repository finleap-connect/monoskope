package eventsourcing

import (
	"context"
)

// CommandHandler is an interface that all handlers of commands must implement.
type CommandHandler interface {
	// HandleCommand handles a command.
	HandleCommand(context.Context, Command) error
}

// EventHandler is an interface that all handlers of events must implement.
type EventHandler interface {
	// HandleEvent handles an event.
	HandleEvent(context.Context, Event) error
}

// CommandHandlerMiddleware is a function that middlewares can implement to be
// able to chain.
type CommandHandlerMiddleware func(CommandHandler) CommandHandler

// UseCommandHandlerMiddleware wraps a CommandHandler in one or more middleware.
func UseCommandHandlerMiddleware(h CommandHandler, middleware ...CommandHandlerMiddleware) CommandHandler {
	// Apply in reverse order.
	for i := len(middleware) - 1; i >= 0; i-- {
		m := middleware[i]
		h = m(h)
	}
	return h
}

// EventHandlerMiddleware is a function that middlewares can implement to be
// able to chain.
type EventHandlerMiddleware func(EventHandler) EventHandler

// UseEventHandlerMiddleware wraps a EventHandler in one or more middleware.
func UseEventHandlerMiddleware(h EventHandler, middleware ...EventHandlerMiddleware) EventHandler {
	// Apply in reverse order.
	for i := len(middleware) - 1; i >= 0; i-- {
		m := middleware[i]
		h = m(h)
	}
	return h
}
