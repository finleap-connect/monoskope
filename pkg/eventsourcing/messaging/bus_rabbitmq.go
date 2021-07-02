package messaging

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/streadway/amqp"
	evs "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing/errors"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/logger"
)

// rabbitEventBus implements an EventBus using RabbitMQ.
type rabbitEventBus struct {
	log             logger.Logger
	conf            *RabbitEventBusConfig
	connection      *amqp.Connection
	channel         *amqp.Channel
	notifyConnClose chan *amqp.Error
	notifyChanClose chan *amqp.Error
	notifyPublish   chan amqp.Confirmation
	isPublisher     bool
	isConnected     bool
	ctx             context.Context
	cancel          context.CancelFunc
}

func newRabbitEventBus(conf *RabbitEventBusConfig) *rabbitEventBus {
	ctx, cancel := context.WithCancel(context.Background())

	return &rabbitEventBus{
		conf:   conf,
		ctx:    ctx,
		cancel: cancel,
	}
}

// NewRabbitEventBusPublisher creates a new EventBusPublisher for rabbitmq.
func NewRabbitEventBusPublisher(conf *RabbitEventBusConfig) (evs.EventBusPublisher, error) {
	err := conf.Validate()
	if err != nil {
		return nil, err
	}

	b := newRabbitEventBus(conf)
	b.log = logger.WithName("publisher").WithValues("name", conf.name)
	b.isPublisher = true

	return b, nil
}

// NewRabbitEventBusConsumer creates a new EventBusConsumer for rabbitmq.
func NewRabbitEventBusConsumer(conf *RabbitEventBusConfig) (evs.EventBusConsumer, error) {
	err := conf.Validate()
	if err != nil {
		return nil, err
	}

	b := newRabbitEventBus(conf)
	b.log = logger.WithName("consumer").WithValues("name", conf.name)
	return b, nil
}

// Open starts automatic reconnect with rabbitmq
func (b *rabbitEventBus) Open(ctx context.Context) error {
	go b.handleReconnect(b.conf.url)
	for !b.isConnected {
		select {
		case <-b.ctx.Done():
			b.log.Info("Connection aborted because of shutdown.")
			return errors.ErrContextDeadlineExceeded
		case <-ctx.Done():
			b.log.Info("Connection aborted because context deadline exceeded.")
			return errors.ErrMessageNotConnected
		case <-time.After(300 * time.Millisecond):
		}
	}
	return nil
}

// PublishEvent publishes the event on the bus.
func (b *rabbitEventBus) PublishEvent(ctx context.Context, event evs.Event) error {
	resendsLeft := b.conf.maxResends
	for resendsLeft > 0 {
		resendsLeft--

		err := b.publishEvent(event)
		if err != nil {
			select {
			case <-b.ctx.Done():
				b.log.Info("Publish failed because of shutdown.")
				return errors.ErrCouldNotPublishEvent
			case <-ctx.Done():
				b.log.Info("Publish failed because context deadline exceeded.")
				return errors.ErrContextDeadlineExceeded
			case <-time.After(b.conf.resendDelay):
				b.log.Error(err, "Publish failed. Retrying...")
				continue
			}
		}

		select {
		case confirmed := <-b.notifyPublish:
			if confirmed.Ack {
				b.log.V(logger.DebugLevel).Info("Publish confirmed.", "DeliveryTag", confirmed.DeliveryTag)
				return nil
			} else {
				b.log.Info("Publish wasn't confirmed. Retrying...", "resends left", resendsLeft, "DeliveryTag", confirmed.DeliveryTag)
			}
		case <-ctx.Done():
			b.log.Info("Publish failed because context deadline exceeded.")
			return errors.ErrContextDeadlineExceeded
		case <-b.ctx.Done():
			b.log.Info("Publish failed because of shutdown.")
			return errors.ErrCouldNotPublishEvent
		case <-time.After(b.conf.resendDelay):
			b.log.Info("Publish failed. Wasn't confirmed within timeout.")
			return errors.ErrCouldNotPublishEvent
		}
	}

	b.log.Info("Publish failed.")
	return errors.ErrCouldNotPublishEvent
}

// publishEvent will push to the queue without checking for
// confirmation. It returns an error if it fails to connect.
// No guarantees are provided for whether the server will
// receive the message.
func (b *rabbitEventBus) publishEvent(event evs.Event) error {
	if !b.isConnected {
		return errors.ErrMessageNotConnected
	}

	re := &rabbitMessage{
		EventType:        event.EventType(),
		Data:             event.Data(),
		Timestamp:        event.Timestamp(),
		AggregateType:    event.AggregateType(),
		AggregateID:      event.AggregateID(),
		AggregateVersion: event.AggregateVersion(),
		Metadata:         event.Metadata(),
	}

	bytes, err := json.Marshal(re)
	if err != nil {
		b.log.Error(err, errors.ErrCouldNotMarshalEvent.Error())
		return errors.ErrCouldNotMarshalEvent
	}

	err = b.channel.Publish(
		b.conf.exchangeName,         // exchange
		b.generateRoutingKey(event), // routingKey
		false,                       // mandatory
		false,                       // immediate
		amqp.Publishing{
			ContentType:  "text/json",
			Body:         bytes,
			DeliveryMode: amqp.Transient,
		})

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

	resendsLeft := b.conf.maxResends
	for {
		err := b.addHandler(ctx, handler, matchers...)
		if err != nil {
			b.log.Info("Adding handler failed. Retrying...", "error", err.Error())
			select {
			case <-b.ctx.Done():
				return errors.ErrCouldNotAddHandler
			case <-ctx.Done():
				return errors.ErrContextDeadlineExceeded
			case <-time.After(b.conf.resendDelay):
				if resendsLeft > 0 {
					resendsLeft--
					continue
				} else {
					b.log.Info("Adding handler failed.")
					return errors.ErrCouldNotAddHandler
				}
			}
		}

		b.log.Info("Handler added.")
		return nil
	}
}

// addHandler creates a queue along with bindings for the given matchers
func (b *rabbitEventBus) addHandler(ctx context.Context, handler evs.EventHandler, matchers ...evs.EventMatcher) error {
	if !b.isConnected {
		return errors.ErrMessageNotConnected
	}

	q, err := b.channel.QueueDeclare(
		"",    // queue name autogenerated
		false, // durable
		true,  // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return errors.ErrMessageBusConnection
	}
	b.log.Info(fmt.Sprintf("Queue declared '%s'.", q.Name))

	for _, matcher := range matchers {
		rabbitMatcher, ok := matcher.(*rabbitMatcher)
		if !ok {
			b.log.Error(errors.ErrMatcherMustNotBeNil, errors.ErrMatcherMustNotBeNil.Error())
			return errors.ErrMatcherMustNotBeNil
		}

		routingKey := rabbitMatcher.generateRoutingKey()
		err = b.channel.QueueBind(
			q.Name,              // queue name
			routingKey,          // routing key
			b.conf.exchangeName, // exchange
			false,               // no wait
			nil)
		if err != nil {
			return errors.ErrMessageBusConnection
		}
		b.log.Info(fmt.Sprintf("Routing key '%s' bound for queue '%s'...", routingKey, q.Name))
	}

	msgs, err := b.channel.Consume(
		q.Name, // queue
		fmt.Sprintf("%s-%s", q.Name, b.conf.name), // consumer
		false, // auto ack
		true,  // exclusive
		false, // no local
		false, // no wait
		nil,   // args
	)
	if err != nil {
		return errors.ErrMessageBusConnection
	}
	go b.handleIncomingMessages(ctx, q.Name, msgs, handler)

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

	b.cancel()
	go b.flushConfirms()
	err := b.connection.Close()
	if err != nil {
		return err
	}

	b.log.Info("Shutdown complete.")

	return nil
}

func (b *rabbitEventBus) channelClosed() {
	b.isConnected = false
	if b.channel != nil {
		b.log.Info("Closing previous channel...")
		_ = b.channel.Close()
	}
	b.channel = nil
}

func (b *rabbitEventBus) connectionClosed() {
	b.isConnected = false
	if b.connection != nil {
		b.log.Info("Closing previous connection...")
		_ = b.connection.Close()
	}
	b.connection = nil
}

// changeChannel takes a new channel to the queue,
// and updates the channel listeners to reflect this.
func (b *rabbitEventBus) changeChannel(channel *amqp.Channel) {
	b.channel = channel
	b.notifyChanClose = b.channel.NotifyClose(make(chan *amqp.Error))

	if b.isPublisher {
		b.notifyPublish = b.channel.NotifyPublish(make(chan amqp.Confirmation, 1))
	}

	b.isConnected = true
}

// init will initialize channel & declare the exchange
func (b *rabbitEventBus) init(conn *amqp.Connection) error {
	ch, err := conn.Channel()
	if err != nil {
		return err
	}

	if b.isPublisher {
		// Indicate we want confirmation of the publish.
		if err = ch.Confirm(false); err != nil {
			return err
		}

		// Declare the exchange, if it exists this doesn't do anything.
		// If the exchange exists but is setup differently we will get an error here.
		err = ch.ExchangeDeclare(
			b.conf.exchangeName, // name
			amqp.ExchangeTopic,  // type
			true,                // durable
			false,               // auto-deleted
			false,               // internal
			false,               // no-wait
			nil,                 // arguments
		)
		if err != nil {
			return err
		}
	} else {
		// Indicate we only want 1 message to acknowledge at a time.
		if err := ch.Qos(1, 0, true); err != nil {
			return err
		}
	}

	b.changeChannel(ch)

	return nil
}

// connect will create a new AMQP connection
func (b *rabbitEventBus) connect(addr string) (*amqp.Connection, error) {
	b.log.Info("Attempting to connect...")

	conn, err := amqp.DialConfig(addr, *b.conf.amqpConfig)
	if err != nil {
		return nil, err
	}
	b.connection = conn

	b.notifyConnClose = conn.NotifyClose(make(chan *amqp.Error))

	b.log.Info("Connection established.")

	return conn, nil
}

// handleReconnect will wait for a connection error on
// notifyConnClose, and then continuously attempt to reconnect.
func (b *rabbitEventBus) handleReconnect(addr string) {
	for {
		b.isConnected = false
		conn, err := b.connect(addr)
		if err != nil {
			b.log.Info("Failed to connect. Retrying...", "error", err.Error())

			select {
			case <-b.ctx.Done():
				return
			case <-time.After(b.conf.reconnectDelay):
				continue
			}
		}

		for {
			b.isConnected = false
			if err = b.init(conn); err != nil {
				b.log.Info("Failed to initialize channel. Retrying...", "error", err.Error())

				select {
				case <-b.ctx.Done():
					b.log.Info("Aborting init. Shutting down...")
					return
				case <-time.After(b.conf.reInitDelay):
					b.channelClosed()
					go b.handleReconnect(addr)
					return
				}
			}

			select {
			case <-b.notifyConnClose:
				b.log.Info("Connection closed. Reconnecting...")
				b.channelClosed()
				go b.handleReconnect(addr)
				return
			case <-b.notifyChanClose:
				b.log.Info("Channel closed. Re-running init...")
				b.connectionClosed()
				continue
			case <-b.ctx.Done():
				return
			}
		}
	}
}

func (b *rabbitEventBus) flushConfirms() {
	for range b.notifyPublish {
		// read channel until closed
	}
}

// generateRoutingKey generates the routing key for an event.
func (b *rabbitEventBus) generateRoutingKey(event evs.Event) string {
	return fmt.Sprintf("%s.%s.%s", b.conf.routingKeyPrefix, event.AggregateType(), event.EventType())
}

// handleIncomingMessages handles the routing of the received messages and ack/nack based on handler result
func (b *rabbitEventBus) handleIncomingMessages(ctx context.Context, qName string, msgs <-chan amqp.Delivery, handler evs.EventHandler) {
	b.log.Info(fmt.Sprintf("Handler for queue '%s' started.", qName))
	for d := range msgs {
		b.log.V(logger.DebugLevel).Info(fmt.Sprintf("Handler received event from queue '%s'.", qName))

		re := &rabbitEvent{}
		err := json.Unmarshal(d.Body, re)
		if err != nil {
			b.log.Error(err, "Failed to unmarshal event.", "event", d.Body)
			if err := d.Nack(false, false); err != nil {
				b.log.Error(err, "Failed to NACK event.")
			}
		}

		err = handler.HandleEvent(ctx, re)
		if err != nil {
			if err := d.Nack(false, false); err != nil {
				b.log.Error(err, "Failed to NACK event.")
			}
		} else {
			if err := d.Ack(false); err != nil {
				b.log.Error(err, "Failed to ACK event.")
			}
		}
	}
	b.log.Info(fmt.Sprintf("Handler for queue '%s' stopped.", qName))
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
