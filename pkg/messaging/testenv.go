package messaging

import (
	"gitlab.figo.systems/platform/monoskope/monoskope/internal/test"
)

type messageBusTestEnv struct {
	*test.TestEnv
	amqpURL string
}

func (env *messageBusTestEnv) Shutdown() error {
	return env.TestEnv.Shutdown()
}
