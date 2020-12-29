package messaging

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/streadway/amqp"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/storage"
)

const (
	exchangeName = "m8_events"
)

// RabbitEventBus implements an EventBus using RabbitMQ.
type RabbitEventBus struct {
	conn        *amqp.Connection
	topicPrefix string
}

func (b *RabbitEventBus) initExchange() error {
	ch, err := b.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	err = ch.ExchangeDeclare(
		exchangeName, // name
		"topic",      // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
	return err
}

func (b *RabbitEventBus) generateRoutingKey(event storage.Event) string {
	return fmt.Sprintf("%s.%s", b.topicPrefix, event.AggregateType())
}

func NewRabbitEventBusPublisher(conn *amqp.Connection, topicPrefix string) (EventBusPublisher, error) {
	s := &RabbitEventBus{
		conn:        conn,
		topicPrefix: topicPrefix,
	}
	err := s.initExchange()
	if err != nil {
		return nil, err
	}
	return s, nil
}

func NewRabbitEventBusConsumer(conn *amqp.Connection, topicPrefix string) (EventBusConsumer, error) {
	s := &RabbitEventBus{
		conn:        conn,
		topicPrefix: topicPrefix,
	}
	return s, nil
}

// PublishEvent publishes the event on the bus.
func (b *RabbitEventBus) PublishEvent(ctx context.Context, event storage.Event) error {
	bytes, err := json.Marshal(event)
	if err != nil {
		return err
	}

	ch, err := b.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	err = ch.Publish(
		exchangeName,                // exchange
		b.generateRoutingKey(event), // routingKey
		false,                       // mandatory
		false,                       // immediate
		amqp.Publishing{
			ContentType: "text/json",
			Body:        bytes,
		})

	return err
}

// AddReceiver adds a receiver for event matching the EventFilter.
func (b *RabbitEventBus) AddReceiver(EventMatcher, EventReceiver) {

}
