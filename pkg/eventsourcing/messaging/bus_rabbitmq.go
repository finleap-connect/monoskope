package messaging

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/streadway/amqp"
	rabbitmq "github.com/wagslane/go-rabbitmq"
	evs "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing/errors"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/logger"
)

// rabbitEventBus implements an EventBus using RabbitMQ.
type rabbitEventBus struct {
	log       logger.Logger
	conf      *RabbitEventBusConfig
	publisher *rabbitmq.Publisher
	consumer  *rabbitmq.Consumer
	returns   <-chan rabbitmq.Return
}

func newRabbitEventBus(conf *RabbitEventBusConfig) (*rabbitEventBus, error) {
	err := conf.Validate()
	if err != nil {
		return nil, err
	}

	return &rabbitEventBus{
		conf: conf,
	}, nil
}

// NewRabbitEventBusPublisher creates a new EventBusPublisher for rabbitmq.
func NewRabbitEventBusPublisher(conf *RabbitEventBusConfig) (evs.EventBusPublisher, error) {
	b, err := newRabbitEventBus(conf)
	if err != nil {
		return nil, err
	}

	publisher, returns, err := rabbitmq.NewPublisher(conf.url, *conf.amqpConfig)
	if err != nil {
		return nil, err
	}
	b.publisher = &publisher
	b.returns = returns
	b.log = logger.WithName("publisher").WithValues("name", conf.name)

	go func() {
		for r := range returns {
			b.log.Info("message returned from server: %s", string(r.Body))
		}
	}()

	return b, nil
}

// NewRabbitEventBusConsumer creates a new EventBusConsumer for rabbitmq.
func NewRabbitEventBusConsumer(conf *RabbitEventBusConfig) (evs.EventBusConsumer, error) {
	b, err := newRabbitEventBus(conf)
	if err != nil {
		return nil, err
	}

	consumer, err := rabbitmq.NewConsumer(conf.url, *conf.amqpConfig)
	if err != nil {
		return nil, err
	}
	b.consumer = &consumer
	b.log = logger.WithName("consumer").WithValues("name", conf.name)

	return b, nil
}

// PublishEvent publishes the event on the bus.
func (b *rabbitEventBus) PublishEvent(ctx context.Context, event evs.Event) error {
	re := &rabbitMessage{
		EventType:        event.EventType(),
		Data:             event.Data(),
		Timestamp:        event.Timestamp(),
		AggregateType:    event.AggregateType(),
		AggregateID:      event.AggregateID(),
		AggregateVersion: event.AggregateVersion(),
		Metadata:         event.Metadata(),
	}

	rabbitMessageBytes, err := json.Marshal(re)
	if err != nil {
		b.log.Error(err, errors.ErrCouldNotMarshalEvent.Error())
		return errors.ErrCouldNotMarshalEvent
	}

	err = b.publisher.Publish(
		rabbitMessageBytes,
		[]string{b.generateRoutingKey(event)},
		rabbitmq.WithPublishOptionsContentType("application/json"),
		rabbitmq.WithPublishOptionsExchange(b.conf.exchangeName),
	)
	if err != nil {
		b.log.Error(err, errors.ErrCouldNotPublishEvent.Error())
		return errors.ErrCouldNotPublishEvent
	}

	return nil
}

// AddHandler adds a handler for event matching the EventFilter.
func (b *rabbitEventBus) AddHandler(ctx context.Context, handler evs.EventHandler, matchers ...evs.EventMatcher) error {
	if matchers == nil {
		b.log.Error(errors.ErrMatcherMustNotBeNil, errors.ErrMatcherMustNotBeNil.Error())
		return errors.ErrMatcherMustNotBeNil
	}
	if handler == nil {
		b.log.Error(errors.ErrHandlerMustNotBeNil, errors.ErrHandlerMustNotBeNil.Error())
		return errors.ErrHandlerMustNotBeNil
	}

	err := b.addHandler(ctx, handler, matchers...)
	if err != nil {
		b.log.Error(err, "Adding handler failed.")
		return err
	}

	b.log.Info("Handler added.")
	return nil
}

// addHandler creates a queue along with bindings for the given matchers
func (b *rabbitEventBus) addHandler(ctx context.Context, handler evs.EventHandler, matchers ...evs.EventMatcher) error {
	var routingKeys []string
	for _, matcher := range matchers {
		rabbitMatcher, ok := matcher.(*rabbitMatcher)
		if !ok {
			b.log.Error(errors.ErrMatcherMustNotBeNil, errors.ErrMatcherMustNotBeNil.Error())
			return errors.ErrMatcherMustNotBeNil
		}
		routingKeys = append(routingKeys, rabbitMatcher.generateRoutingKey())
	}

	err := b.consumer.StartConsuming(
		func(d rabbitmq.Delivery) bool {
			return b.handleIncomingMessages(ctx, d, handler)
		},
		b.conf.name,
		routingKeys,
		rabbitmq.WithConsumeOptionsBindingExchangeName(b.conf.exchangeName),
		rabbitmq.WithConsumeOptionsBindingExchangeKind(amqp.ExchangeTopic),
		rabbitmq.WithConsumeOptionsBindingExchangeAutoDelete,
		rabbitmq.WithConsumeOptionsQueueExclusive,
	)
	if err != nil {
		return errors.ErrMessageBusConnection
	}

	return nil
}

// Matcher returns a new EventMatcher of type RabbitMatcher
func (b *rabbitEventBus) Matcher() evs.EventMatcher {
	matcher := &rabbitMatcher{
		routingKeyPrefix: b.conf.routingKeyPrefix,
	}
	return matcher.Any()
}

// Close will cleanly shutdown the channel and connection.
func (b *rabbitEventBus) Close() error {
	b.log.Info("Shutting down...")

	if b.publisher != nil {
		b.publisher.StopPublishing()
	}
	if b.consumer != nil {
		b.consumer.Disconnect()
	}

	b.log.Info("Shutdown complete.")

	return nil
}

// generateRoutingKey generates the routing key for an event.
func (b *rabbitEventBus) generateRoutingKey(event evs.Event) string {
	return fmt.Sprintf("%s.%s.%s", b.conf.routingKeyPrefix, event.AggregateType(), event.EventType())
}

// handleIncomingMessages handles the routing of the received messages and ack/nack based on handler result
func (b *rabbitEventBus) handleIncomingMessages(ctx context.Context, d rabbitmq.Delivery, handler evs.EventHandler) bool {
	re := &rabbitEvent{}
	err := json.Unmarshal(d.Body, re)
	if err != nil {
		b.log.Error(err, "Failed to unmarshal event.", "event", d.Body)
		return false
	}

	err = handler.HandleEvent(ctx, re)
	if err != nil {
		b.log.Error(err, "Handling event failed.")
		return false
	}

	return true
}

// rabbitEvent implements the message body transferred via rabbitmq
type rabbitMessage struct {
	EventType        evs.EventType
	Data             evs.EventData
	Timestamp        time.Time
	AggregateType    evs.AggregateType
	AggregateID      uuid.UUID
	AggregateVersion uint64
	Metadata         map[string]string
}

// rabbitEvent is the private implementation of the Event interface for a rabbitmq message bus.
type rabbitEvent struct {
	rabbitMessage
}

// EventType implements the EventType method of the Event interface.
func (e rabbitEvent) EventType() evs.EventType {
	return e.rabbitMessage.EventType
}

// Data implements the Data method of the Event interface.
func (e rabbitEvent) Data() evs.EventData {
	return e.rabbitMessage.Data
}

// Timestamp implements the Timestamp method of the Event interface.
func (e rabbitEvent) Timestamp() time.Time {
	return e.rabbitMessage.Timestamp
}

// AggregateType implements the AggregateType method of the Event interface.
func (e rabbitEvent) AggregateType() evs.AggregateType {
	return e.rabbitMessage.AggregateType
}

// AggregateID implements the AggregateID method of the Event interface.
func (e rabbitEvent) AggregateID() uuid.UUID {
	return e.rabbitMessage.AggregateID
}

// AggregateVersion implements the AggregateVersion method of the Event interface.
func (e rabbitEvent) AggregateVersion() uint64 {
	return e.rabbitMessage.AggregateVersion
}

// AggregateVersion implements the AggregateVersion method of the Event interface.
func (e rabbitEvent) Metadata() map[string]string {
	return e.rabbitMessage.Metadata
}

// String implements the String method of the Event interface.
func (e rabbitEvent) String() string {
	return fmt.Sprintf("%s@%d", e.rabbitMessage.EventType, e.rabbitMessage.AggregateVersion)
}
