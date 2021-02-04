package messaging

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/streadway/amqp"
	evs "gitlab.figo.systems/platform/monoskope/monoskope/pkg/event_sourcing"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/logger"
)

// rabbitEventBus implements an EventBus using RabbitMQ.
type rabbitEventBus struct {
	log             logger.Logger
	conf            *rabbitEventBusConfig
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

func newRabbitEventBus(conf *rabbitEventBusConfig) *rabbitEventBus {
	ctx, cancel := context.WithCancel(context.Background())

	return &rabbitEventBus{
		conf:   conf,
		ctx:    ctx,
		cancel: cancel,
	}
}

// NewRabbitEventBusPublisher creates a new EventBusPublisher for rabbitmq.
func NewRabbitEventBusPublisher(conf *rabbitEventBusConfig) (evs.EventBusPublisher, error) {
	err := conf.Validate()
	if err != nil {
		return nil, err
	}

	b := newRabbitEventBus(conf)
	b.log = logger.WithName("publisher").WithValues("name", conf.Name)
	b.isPublisher = true

	return b, nil
}

// NewRabbitEventBusConsumer creates a new EventBusConsumer for rabbitmq.
func NewRabbitEventBusConsumer(conf *rabbitEventBusConfig) (evs.EventBusConsumer, error) {
	err := conf.Validate()
	if err != nil {
		return nil, err
	}

	b := newRabbitEventBus(conf)
	b.log = logger.WithName("consumer").WithValues("name", conf.Name)
	return b, nil
}

// Connect starts automatic reconnect with rabbitmq
func (b *rabbitEventBus) Connect(ctx context.Context) error {
	go b.handleReconnect(b.conf.Url)
	for !b.isConnected {
		select {
		case <-b.ctx.Done():
			b.log.Info("Connection aborted because of shutdown.")
			return evs.ErrContextDeadlineExceeded
		case <-ctx.Done():
			b.log.Info("Connection aborted because context deadline exceeded.")
			return evs.ErrMessageNotConnected
		case <-time.After(300 * time.Millisecond):
		}
	}
	return nil
}

// PublishEvent publishes the event on the bus.
func (b *rabbitEventBus) PublishEvent(ctx context.Context, event evs.Event) error {
	resendsLeft := b.conf.MaxResends
	for resendsLeft > 0 {
		resendsLeft--

		err := b.publishEvent(event)
		if err != nil {
			select {
			case <-b.ctx.Done():
				b.log.Info("Publish failed because of shutdown.")
				return evs.ErrCouldNotPublishEvent
			case <-ctx.Done():
				b.log.Info("Publish failed because context deadline exceeded.")
				return evs.ErrContextDeadlineExceeded
			case <-time.After(b.conf.ResendDelay):
				b.log.Error(err, "Publish failed. Retrying...")
				continue
			}
		}

		select {
		case confirmed := <-b.notifyPublish:
			if confirmed.Ack {
				b.log.Info("Publish confirmed.", "DeliveryTag", confirmed.DeliveryTag)
				return nil
			} else {
				b.log.Info("Publish wasn't confirmed. Retrying...", "resends left", resendsLeft, "DeliveryTag", confirmed.DeliveryTag)
			}
		case <-ctx.Done():
			b.log.Info("Publish failed because context deadline exceeded.")
			return evs.ErrContextDeadlineExceeded
		case <-b.ctx.Done():
			b.log.Info("Publish failed because of shutdown.")
			return evs.ErrCouldNotPublishEvent
		case <-time.After(b.conf.ResendDelay):
			b.log.Info("Publish failed. Wasn't confirmed within timeout.")
			return evs.ErrCouldNotPublishEvent
		}
	}

	b.log.Info("Publish failed.")
	return evs.ErrCouldNotPublishEvent
}

// publishEvent will push to the queue without checking for
// confirmation. It returns an error if it fails to connect.
// No guarantees are provided for whether the server will
// recieve the message.
func (b *rabbitEventBus) publishEvent(event evs.Event) error {
	if !b.isConnected {
		return evs.ErrMessageNotConnected
	}

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
		b.log.Error(err, evs.ErrCouldNotMarshalEvent.Error())
		return evs.ErrCouldNotMarshalEvent
	}

	err = b.channel.Publish(
		b.conf.ExchangeName,         // exchange
		b.generateRoutingKey(event), // routingKey
		false,                       // mandatory
		false,                       // immediate
		amqp.Publishing{
			ContentType:  "text/json",
			Body:         bytes,
			DeliveryMode: amqp.Transient,
		})

	if err != nil {
		b.log.Error(err, evs.ErrCouldNotPublishEvent.Error())
		return evs.ErrCouldNotPublishEvent
	}
	return nil
}

// AddReceiver adds a receiver for event matching the EventFilter.
func (b *rabbitEventBus) AddReceiver(ctx context.Context, handler evs.EventHandler, matchers ...evs.EventMatcher) error {
	if matchers == nil {
		b.log.Error(evs.ErrMatcherMustNotBeNil, evs.ErrMatcherMustNotBeNil.Error())
		return evs.ErrMatcherMustNotBeNil
	}
	if handler == nil {
		b.log.Error(evs.ErrReceiverMustNotBeNil, evs.ErrReceiverMustNotBeNil.Error())
		return evs.ErrReceiverMustNotBeNil
	}

	resendsLeft := b.conf.MaxResends
	for {
		err := b.addReceiver(ctx, handler, matchers...)
		if err != nil {
			b.log.Info("Adding receiver failed. Retrying...", "error", err.Error())
			select {
			case <-b.ctx.Done():
				return evs.ErrCouldNotAddReceiver
			case <-ctx.Done():
				return evs.ErrContextDeadlineExceeded
			case <-time.After(b.conf.ResendDelay):
				if resendsLeft > 0 {
					resendsLeft--
					continue
				} else {
					b.log.Info("Adding receiver failed.")
					return evs.ErrCouldNotPublishEvent
				}
			}
		}

		b.log.Info("Receiver added.")
		return nil
	}
}

// addReceiver creates a queue along with bindings for the given matchers
func (b *rabbitEventBus) addReceiver(ctx context.Context, handler evs.EventHandler, matchers ...evs.EventMatcher) error {
	if !b.isConnected {
		return evs.ErrMessageNotConnected
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
		return evs.ErrMessageBusConnection
	}
	b.log.Info(fmt.Sprintf("Queue declared '%s'.", q.Name))

	for _, matcher := range matchers {
		rabbitMatcher, ok := matcher.(*rabbitMatcher)
		if !ok {
			b.log.Error(evs.ErrMatcherMustNotBeNil, evs.ErrMatcherMustNotBeNil.Error())
			return evs.ErrMatcherMustNotBeNil
		}

		routingKey := rabbitMatcher.generateRoutingKey()
		err = b.channel.QueueBind(
			q.Name,              // queue name
			routingKey,          // routing key
			b.conf.ExchangeName, // exchange
			false,               // no wait
			nil)
		if err != nil {
			return evs.ErrMessageBusConnection
		}
		b.log.Info(fmt.Sprintf("Routing key '%s' bound for queue '%s'...", routingKey, q.Name))
	}

	msgs, err := b.channel.Consume(
		q.Name, // queue
		fmt.Sprintf("%s-%s", q.Name, b.conf.Name), // consumer
		false, // auto ack
		true,  // exclusive
		false, // no local
		false, // no wait
		nil,   // args
	)
	if err != nil {
		return evs.ErrMessageBusConnection
	}
	go b.handleIncomingMessages(ctx, q.Name, msgs, handler)

	return nil
}

// Matcher returns a new EventMatcher of type RabbitMatcher
func (b *rabbitEventBus) Matcher() evs.EventMatcher {
	matcher := &rabbitMatcher{
		routingKeyPrefix: b.conf.RoutingKeyPrefix,
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
			b.conf.ExchangeName, // name
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
			case <-time.After(b.conf.ReconnectDelay):
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
				case <-time.After(b.conf.ReInitDelay):
					continue
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
	return fmt.Sprintf("%s.%s.%s", b.conf.RoutingKeyPrefix, event.AggregateType(), event.EventType())
}

// handleIncomingMessages handles the routing of the received messages and ack/nack based on receiver result
func (b *rabbitEventBus) handleIncomingMessages(ctx context.Context, qName string, msgs <-chan amqp.Delivery, handler evs.EventHandler) {
	b.log.Info(fmt.Sprintf("Handler for queue '%s' started.", qName))
	for d := range msgs {
		b.log.Info(fmt.Sprintf("Handler received event from queue '%s'.", qName))

		re := &rabbitEvent{}
		err := json.Unmarshal(d.Body, re)
		if err != nil {
			b.log.Error(err, "Failed to unmarshal event.", "event", d.Body)
			if err := d.Nack(false, false); err != nil {
				b.log.Error(err, "Failed to NACK event.")
			}
		}

		err = handler.HandleEvent(ctx, evs.NewEvent(re.EventType, re.Data, re.Timestamp, re.AggregateType, re.AggregateID, re.AggregateVersion))
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

// rabbitEvent implements the message body transfered via rabbitmq
type rabbitEvent struct {
	EventType        evs.EventType
	Data             evs.EventData
	Timestamp        time.Time
	AggregateType    evs.AggregateType
	AggregateID      uuid.UUID
	AggregateVersion uint64
}
