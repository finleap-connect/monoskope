package messaging

import (
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/storage"
)

type EventReceiver func(storage.Event) error
