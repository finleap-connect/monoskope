package messaging

import (
	"github.com/streadway/amqp"
	"gitlab.figo.systems/platform/monoskope/monoskope/internal/test"
)

type MessageBusTestEnv struct {
	*test.TestEnv
	RabbitConn []*amqp.Connection
	Publisher  EventBusPublisher
	Consumer   EventBusConsumer
}

func (env *MessageBusTestEnv) Shutdown() error {
	if env.RabbitConn != nil {
		for _, conn := range env.RabbitConn {
			err := conn.Close()
			if err != nil {
				return err
			}
		}
	}
	return env.TestEnv.Shutdown()
}
