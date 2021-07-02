package messagebus

import (
	"context"
	"os"
	"time"

	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing"
	esMessaging "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing/messaging"
)

func getMsgBusUrl() string {
	var msgbusUrl string
	if v := os.Getenv("BUS_URL"); v != "" {
		msgbusUrl = v
	}
	return msgbusUrl
}

func getRabbitConf(name, msgbusPrefix string, useTLS bool) (*esMessaging.RabbitEventBusConfig, error) {
	rabbitConf := esMessaging.NewRabbitEventBusConfig(name, getMsgBusUrl(), msgbusPrefix)

	if useTLS {
		err := rabbitConf.ConfigureTLS()
		if err != nil {
			return nil, err
		}
	}

	return rabbitConf, nil
}

func NewEventBusConsumer(name, msgbusPrefix string) (eventsourcing.EventBusConsumer, error) {
	rabbitConf, err := getRabbitConf(name, msgbusPrefix, true)
	if err != nil {
		return nil, err
	}
	return NewEventBusConsumerFromConfig(rabbitConf)
}

func NewEventBusConsumerFromConfig(rabbitConf *esMessaging.RabbitEventBusConfig) (eventsourcing.EventBusConsumer, error) {
	consumer, err := esMessaging.NewRabbitEventBusConsumer(rabbitConf)
	if err != nil {
		return nil, err
	}

	ctx, cancelFunc := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancelFunc()
	err = consumer.Open(ctx)
	if err != nil {
		return nil, err
	}

	return consumer, nil
}

func NewEventBusPublisher(name, msgbusPrefix string) (eventsourcing.EventBusPublisher, error) {
	rabbitConf, err := getRabbitConf(name, msgbusPrefix, true)
	if err != nil {
		return nil, err
	}

	publisher, err := esMessaging.NewRabbitEventBusPublisher(rabbitConf)
	if err != nil {
		return nil, err
	}

	ctx, cancelFunc := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancelFunc()
	err = publisher.Open(ctx)
	if err != nil {
		return nil, err
	}

	return publisher, nil
}
