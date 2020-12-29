package messaging

import (
	"context"

	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/storage"
)

type EventReceiver interface {
	// HandleEvent handles an event.
	HandleEvent(context.Context, storage.Event) error
}
