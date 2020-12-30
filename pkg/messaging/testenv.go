package messaging

import (
	"github.com/streadway/amqp"
	"gitlab.figo.systems/platform/monoskope/monoskope/internal/test"
)

type MessageBusTestEnv struct {
	*test.TestEnv
	RabbitConn *amqp.Connection
	Publisher  EventBusPublisher
	Consumer   EventBusConsumer
}

func (env *MessageBusTestEnv) Shutdown() error {
	if env.RabbitConn != nil {
		defer env.RabbitConn.Close()
	}
	return env.TestEnv.Shutdown()
}
