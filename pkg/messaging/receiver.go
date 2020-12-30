package messaging

import (
	"context"

	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/storage"
)

type EventReceiver func(context.Context, storage.Event) error
