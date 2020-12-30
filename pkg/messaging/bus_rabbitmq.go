package messaging

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/streadway/amqp"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/storage"
)

const (
	exchangeName = "m8_events"
)

type RabbitMatcher struct {
	topicPrefix   string
	eventType     string
	aggregateType string
}

// MatchAny matches any event.
func (m *RabbitMatcher) MatchAny() EventMatcher {
	m.eventType = "*"
	m.aggregateType = "*"
	return m
}

// MatchEvent matches a specific event type, nil events never match.
func (m *RabbitMatcher) MatchEvent(eventType storage.EventType) EventMatcher {
	m.eventType = string(eventType)
	return m
}

// MatchAggregate matches a specific aggregate type, nil events never match.
func (m *RabbitMatcher) MatchAggregate(aggregateType storage.AggregateType) EventMatcher {
	m.aggregateType = string(aggregateType)
	return m
}

// generateRoutingKey returns the routing key for events
func (m *RabbitMatcher) generateRoutingKey() string {
	return fmt.Sprintf("%s.%s.%s", m.topicPrefix, m.aggregateType, m.eventType)
}

// RabbitEventBus implements an EventBus using RabbitMQ.
type RabbitEventBus struct {
	conn        *amqp.Connection
	channel     *amqp.Channel
	topicPrefix string
}

func (b *RabbitEventBus) getChannel(forceNew bool) (*amqp.Channel, error) {
	if b.channel != nil && !forceNew {
		return b.channel, nil
	}
	ch, err := b.conn.Channel()
	if err != nil {
		return nil, err
	}
	b.channel = ch
	return ch, nil
}

func (b *RabbitEventBus) initExchange() error {
	ch, err := b.getChannel(false)
	if err != nil {
		return err
	}

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
	return fmt.Sprintf("%s.%s.%s", b.topicPrefix, event.AggregateType(), event.EventType())
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

	ch, err := b.getChannel(false)
	if err != nil {
		return err
	}

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
func (b *RabbitEventBus) AddReceiver(matcher EventMatcher, receiver EventReceiver) error {
	if matcher == nil {
		return errors.New("matcher can't be nil")
	}
	rabbitMatcher, ok := matcher.(*RabbitMatcher)
	if !ok {
		return errors.New("matcher must be of type RabbitMatcher")
	}
	if receiver == nil {
		return errors.New("receiver can't be nil")
	}

	ch, err := b.getChannel(false)
	if err != nil {
		return err
	}

	q, err := ch.QueueDeclare(
		"",    // name
		false, // durable
		true,  // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return err
	}

	err = ch.QueueBind(
		q.Name,                             // queue name
		rabbitMatcher.generateRoutingKey(), // routing key
		exchangeName,                       // exchange
		false,
		nil)
	if err != nil {
		return err
	}

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto ack
		false,  // exclusive
		false,  // no local
		false,  // no wait
		nil,    // args
	)

	go func() {
		for d := range msgs {
			log.Printf(" [x] %s", d.Body)
			_ = d.Ack(false)
		}
	}()

	return err
}

func (b *RabbitEventBus) Matcher() EventMatcher {
	return &RabbitMatcher{
		topicPrefix: b.topicPrefix,
	}
}
