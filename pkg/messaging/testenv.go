package messaging

import (
	"gitlab.figo.systems/platform/monoskope/monoskope/internal/test"
)

type MessageBusTestEnv struct {
	*test.TestEnv
	amqpURL string
}

func (env *MessageBusTestEnv) Shutdown() error {
	return env.TestEnv.Shutdown()
}
