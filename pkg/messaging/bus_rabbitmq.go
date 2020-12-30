package messaging

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/streadway/amqp"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/logger"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/storage"
)

const (
	exchangeName = "m8_events" // name of the monoskope exchange
)

// RabbitMatcher implements the EventMatcher interface for rabbitmq
type RabbitMatcher struct {
	topicPrefix   string
	eventType     string
	aggregateType string
}

// Any matches any event.
func (m *RabbitMatcher) Any() EventMatcher {
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

// rabbitEvent implements the message body transfered via rabbitmq
type rabbitEvent struct {
	EventType        storage.EventType
	Data             storage.EventData
	Timestamp        time.Time
	AggregateType    storage.AggregateType
	AggregateID      uuid.UUID
	AggregateVersion uint64
}

// RabbitEventBus implements an EventBus using RabbitMQ.
type RabbitEventBus struct {
	log         logger.Logger
	conn        *amqp.Connection
	channel     *amqp.Channel
	queues      []*amqp.Queue
	topicPrefix string
	name        string
}

// createChannel creates a new channel for the current connection
func (b *RabbitEventBus) createChannel() (*amqp.Channel, error) {
	ch, err := b.conn.Channel()
	if err != nil {
		return nil, err
	}
	b.channel = ch
	return ch, nil
}

// getChannel returns the existing channel or creates a new one if there is none yet or is forced to
func (b *RabbitEventBus) getChannel(forceNew bool) (*amqp.Channel, error) {
	if forceNew && b.channel != nil {
		err := b.channel.Close()
		b.channel = nil
		if err != nil {
			return nil, err
		}
	}
	if b.channel != nil {
		return b.channel, nil
	}
	return b.createChannel()
}

// initExchange creates the exchange for rabbitmq.
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

// generateRoutingKey generates the routing key for an event.
func (b *RabbitEventBus) generateRoutingKey(event storage.Event) string {
	return fmt.Sprintf("%s.%s.%s", b.topicPrefix, event.AggregateType(), event.EventType())
}

/*
NewRabbitEventBusPublisher creates a new EventBusPublisher for rabbitmq.

- topicPrefix defaults to "*"
*/
func NewRabbitEventBusPublisher(log logger.Logger, conn *amqp.Connection, topicPrefix string) (EventBusPublisher, error) {
	if topicPrefix == "" {
		topicPrefix = "*"
	}
	s := &RabbitEventBus{
		conn:        conn,
		topicPrefix: topicPrefix,
		log:         log,
	}
	err := s.initExchange()
	if err != nil {
		return nil, err
	}
	return s, nil
}

/*
NewRabbitEventBusConsumer creates a new EventBusConsumer for rabbitmq.

- topicPrefix defaults to "*"
*/
func NewRabbitEventBusConsumer(log logger.Logger, conn *amqp.Connection, consumerName, topicPrefix string) (EventBusConsumer, error) {
	if topicPrefix == "" {
		topicPrefix = "*"
	}
	b := &RabbitEventBus{
		conn:        conn,
		topicPrefix: topicPrefix,
		log:         log,
		queues:      make([]*amqp.Queue, 0),
		name:        consumerName,
	}
	return b, nil
}

// PublishEvent publishes the event on the bus.
func (b *RabbitEventBus) PublishEvent(ctx context.Context, event storage.Event) error {
	re := &rabbitEvent{
		EventType:        event.EventType(),
		Data:             event.Data(),
		Timestamp:        event.Timestamp(),
		AggregateType:    event.AggregateType(),
		AggregateID:      event.AggregateID(),
		AggregateVersion: event.AggregateVersion(),
	}

	bytes, err := json.Marshal(re)
	if err != nil {
		b.log.Error(err, ErrCouldNotMarshalEvent.Error())
		return ErrCouldNotMarshalEvent
	}

	ch, err := b.getChannel(false)
	if err != nil {
		b.log.Error(err, ErrCouldNotPublishEvent.Error())
		return ErrCouldNotPublishEvent
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

	if err != nil {
		b.log.Error(err, ErrCouldNotPublishEvent.Error())
		return ErrCouldNotPublishEvent
	}
	return nil
}

// AddReceiver adds a receiver for event matching the EventFilter.
func (b *RabbitEventBus) AddReceiver(receiver EventReceiver, matchers ...EventMatcher) error {
	if matchers == nil {
		b.log.Error(ErrMatcherMustNotBeNil, ErrMatcherMustNotBeNil.Error())
		return ErrMatcherMustNotBeNil
	}
	if receiver == nil {
		b.log.Error(ErrReceiverMustNotBeNil, ErrReceiverMustNotBeNil.Error())
		return ErrReceiverMustNotBeNil
	}

	ch, err := b.createChannel()
	if err != nil {
		return err
	}

	q, err := ch.QueueDeclare(
		"",    // queue name autogenerated
		false, // durable
		true,  // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return err
	}
	b.queues = append(b.queues, &q)

	for _, matcher := range matchers {
		rabbitMatcher, ok := matcher.(*RabbitMatcher)
		if !ok {
			b.log.Error(ErrMatcherMustNotBeNil, ErrMatcherMustNotBeNil.Error())
			return ErrMatcherMustNotBeNil
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
	}

	msgs, err := ch.Consume(
		q.Name, // queue
		b.name, // consumer
		false,  // auto ack
		false,  // exclusive
		false,  // no local
		false,  // no wait
		nil,    // args
	)
	if err != nil {
		return err
	}
	go b.handle(msgs, receiver)

	return err
}

func (b *RabbitEventBus) handle(msgs <-chan amqp.Delivery, receiver EventReceiver) {
	for d := range msgs {
		re := &rabbitEvent{}
		err := json.Unmarshal(d.Body, re)
		if err != nil {
			_ = d.Nack(false, false)
		}
		err = receiver(storage.NewEvent(re.EventType, re.Data, re.Timestamp, re.AggregateType, re.AggregateID, re.AggregateVersion))
		if err != nil {
			_ = d.Nack(false, false)
		} else {
			_ = d.Ack(false)
		}
	}
}

// Matcher returns a new EventMatcher of type RabbitMatcher
func (b *RabbitEventBus) Matcher() EventMatcher {
	matcher := &RabbitMatcher{
		topicPrefix: b.topicPrefix,
	}
	return matcher.Any()
}

// Close frees all disposable resources
func (b *RabbitEventBus) Close() error {
	ch, err := b.getChannel(false)
	if err != nil {
		return err
	}
	defer ch.Close()

	for _, q := range b.queues {
		_, err := ch.QueueDelete(q.Name, true, true, true)
		if err != nil {
			return err
		}
	}

	return nil
}
