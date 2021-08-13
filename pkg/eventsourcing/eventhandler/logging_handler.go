package eventhandler

import (
	"context"

	es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/logger"
)

type loggingEventHandler struct {
	log logger.Logger
}

// NewLoggingEventHandler creates an EventHandler which automates storing Events in the EventStore when a Logging has emitted any.
func NewLoggingEventHandler() *loggingEventHandler {
	return &loggingEventHandler{
		log: logger.WithName("loggingEventHandler"),
	}
}

// HandleEvent implements the HandleEvent method of the es.EventHandler interface.
func (m *loggingEventHandler) HandleEvent(ctx context.Context, event es.Event) error {
	m.log.Info("Handling event.", "event", event.String())
	return nil
}
