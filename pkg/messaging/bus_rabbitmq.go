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

// MatchEventType matches a specific event type, nil events never match.
func (m *RabbitMatcher) MatchEventType(eventType storage.EventType) EventMatcher {
	m.eventType = string(eventType)
	return m
}

// MatchAggregateType matches a specific aggregate type, nil events never match.
func (m *RabbitMatcher) MatchAggregateType(aggregateType storage.AggregateType) EventMatcher {
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
	log              logger.Logger
	conn             *amqp.Connection
	channel          *amqp.Channel
	queues           []*amqp.Queue
	routingKeyPrefix string
	name             string
	errHandler       ErrorHandler
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
	return fmt.Sprintf("%s.%s.%s", b.routingKeyPrefix, event.AggregateType(), event.EventType())
}

/*
NewRabbitEventBusPublisher creates a new EventBusPublisher for rabbitmq.

- routingKeyPrefix defaults to "m8"
*/
func NewRabbitEventBusPublisher(log logger.Logger, conn *amqp.Connection, routingKeyPrefix string) (EventBusPublisher, error) {
	if routingKeyPrefix == "" {
		routingKeyPrefix = "m8"
	}
	s := &RabbitEventBus{
		conn:             conn,
		routingKeyPrefix: routingKeyPrefix,
		log:              log,
	}
	err := s.initExchange()
	if err != nil {
		return nil, err
	}
	return s, nil
}

/*
NewRabbitEventBusConsumer creates a new EventBusConsumer for rabbitmq.

- routingKeyPrefix defaults to "m8"
*/
func NewRabbitEventBusConsumer(log logger.Logger, conn *amqp.Connection, consumerName, routingKeyPrefix string) (EventBusConsumer, error) {
	if routingKeyPrefix == "" {
		routingKeyPrefix = "m8"
	}
	b := &RabbitEventBus{
		conn:             conn,
		routingKeyPrefix: routingKeyPrefix,
		log:              log,
		queues:           make([]*amqp.Queue, 0),
		name:             consumerName,
	}
	return b, nil
}

// PublishEvent publishes the event on the bus.
func (b *RabbitEventBus) PublishEvent(ctx context.Context, event storage.Event) *MessageBusError {
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
		return &MessageBusError{
			Err:     ErrCouldNotMarshalEvent,
			BaseErr: err,
		}
	}

	ch, err := b.getChannel(false)
	if err != nil {
		b.log.Error(err, ErrCouldNotPublishEvent.Error())
		return &MessageBusError{
			Err:     ErrCouldNotPublishEvent,
			BaseErr: err,
		}
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
		return &MessageBusError{
			Err:     ErrCouldNotPublishEvent,
			BaseErr: err,
		}
	}
	return nil
}

// AddErrorHandler adds a handler function to call if any error occurs
func (b *RabbitEventBus) AddErrorHandler(eh ErrorHandler) {
	b.errHandler = eh
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

	errchan := make(chan *amqp.Error)
	go b.handleErr(errchan)
	ch.NotifyClose(errchan)

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
	b.log.Info(fmt.Sprintf("Registering new handler with queue '%s'...", q.Name))

	for _, matcher := range matchers {
		rabbitMatcher, ok := matcher.(*RabbitMatcher)
		if !ok {
			b.log.Error(ErrMatcherMustNotBeNil, ErrMatcherMustNotBeNil.Error())
			return ErrMatcherMustNotBeNil
		}

		routingKey := rabbitMatcher.generateRoutingKey()
		err = ch.QueueBind(
			q.Name,       // queue name
			routingKey,   // routing key
			exchangeName, // exchange
			false,
			nil)
		if err != nil {
			return err
		}
		b.log.Info(fmt.Sprintf("Routing key '%s' bound for queue '%s'...", routingKey, q.Name))
	}

	msgs, err := ch.Consume(
		q.Name, // queue
		b.name, // consumer
		false,  // auto ack
		true,   // exclusive
		false,  // no local
		false,  // no wait
		nil,    // args
	)
	if err != nil {
		return err
	}
	go b.handle(q.Name, msgs, receiver)

	return nil
}

func (b *RabbitEventBus) handleErr(msgs <-chan *amqp.Error) {
	for d := range msgs {
		if b.errHandler != nil {
			b.errHandler(MessageBusError{
				Err:     ErrMessageBusConnection,
				BaseErr: d,
			})
		}
	}
}

func (b *RabbitEventBus) handle(qName string, msgs <-chan amqp.Delivery, receiver EventReceiver) {
	b.log.Info(fmt.Sprintf("Handler for queue '%s' started.", qName))
	for d := range msgs {
		re := &rabbitEvent{}
		err := json.Unmarshal(d.Body, re)
		if err != nil {
			b.log.Error(err, "Failed to unmarshal event.", "event", d.Body)
			_ = d.Nack(false, false)
		}
		err = receiver(storage.NewEvent(re.EventType, re.Data, re.Timestamp, re.AggregateType, re.AggregateID, re.AggregateVersion))
		if err != nil {
			_ = d.Nack(false, false)
		} else {
			_ = d.Ack(false)
		}
	}
	b.log.Info(fmt.Sprintf("Handler for queue '%s' stopped.", qName))
}

// Matcher returns a new EventMatcher of type RabbitMatcher
func (b *RabbitEventBus) Matcher() EventMatcher {
	matcher := &RabbitMatcher{
		topicPrefix: b.routingKeyPrefix,
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
		b.log.Info(fmt.Sprintf("Deleting queue '%s'...", q.Name))
		_, err := ch.QueueDelete(q.Name, false, true, true)
		if err != nil {
			return err
		}
	}

	return nil
}
