package messaging

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/streadway/amqp"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/logger"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/storage"
)

// rabbitEventBus implements an EventBus using RabbitMQ.
type rabbitEventBus struct {
	log             logger.Logger
	conf            *rabbitEventBusConfig
	connection      *amqp.Connection
	channel         *amqp.Channel
	notifyConnClose chan *amqp.Error
	notifyChanClose chan *amqp.Error
	notifyConfirm   chan amqp.Confirmation
	isReady         bool
	done            chan bool
	mu              sync.Mutex
}

// NewRabbitEventBusPublisher creates a new EventBusPublisher for rabbitmq.
func NewRabbitEventBusPublisher(conf *rabbitEventBusConfig) (EventBusPublisher, error) {
	err := conf.Validate()
	if err != nil {
		return nil, err
	}

	b := &rabbitEventBus{
		log:  logger.WithName("publisher").WithValues("name", conf.Name),
		conf: conf,
		done: make(chan bool),
	}
	return b, nil
}

// NewRabbitEventBusConsumer creates a new EventBusConsumer for rabbitmq.
func NewRabbitEventBusConsumer(conf *rabbitEventBusConfig) (EventBusConsumer, error) {
	err := conf.Validate()
	if err != nil {
		return nil, err
	}

	b := &rabbitEventBus{
		log:  logger.WithName("consumer").WithValues("name", conf.Name),
		conf: conf,
		done: make(chan bool),
	}
	return b, nil
}

// Connect starts automatic reconnect with rabbitmq
func (b *rabbitEventBus) Connect(ctx context.Context) *MessageBusError {
	go b.handleReconnect(b.conf.Url)
	for {
		select {
		case <-time.After(300 * time.Millisecond):
			if b.isReady {
				return nil
			}
		case <-b.done:
		case <-ctx.Done():
			return &MessageBusError{
				Err: ErrMessageNotConnected,
			}
		}
	}
}

// PublishEvent publishes the event on the bus.
func (b *rabbitEventBus) PublishEvent(ctx context.Context, event storage.Event) *MessageBusError {
	if !b.isReady {
		return &MessageBusError{
			Err: ErrMessageNotConnected,
		}
	}

	resendsLeft := b.conf.MaxResends
	for {
		err := b.publishEvent(ctx, event)
		if err != nil {
			b.log.Error(err, "Publish failed. Retrying...")
			select {
			case <-b.done:
				return &MessageBusError{
					Err: ErrCouldNotPublishEvent,
				}
			case <-time.After(b.conf.ResendDelay):
				if resendsLeft > 0 {
					resendsLeft--
					continue
				} else {
					b.log.Info("Publish failed.")
					return &MessageBusError{
						Err: ErrCouldNotPublishEvent,
					}
				}
			}
		}

		select {
		case confirm := <-b.notifyConfirm:
			if confirm.Ack {
				b.log.Info("Publish confirmed.")
				return nil
			}
		case <-time.After(b.conf.ResendDelay):
			b.log.Info("Publish wasn't confirmed. Retrying...", "resends left", resendsLeft)

			if resendsLeft > 0 {
				resendsLeft--
				continue
			}
		case <-b.done:
		}

		b.log.Info("Publish failed.")
		return &MessageBusError{
			Err: ErrCouldNotPublishEvent,
		}
	}
}

// publishEvent will push to the queue without checking for
// confirmation. It returns an error if it fails to connect.
// No guarantees are provided for whether the server will
// recieve the message.
func (b *rabbitEventBus) publishEvent(ctx context.Context, event storage.Event) *MessageBusError {
	b.mu.Lock()
	defer b.mu.Unlock()
	if !b.isReady {
		return &MessageBusError{
			Err: ErrMessageNotConnected,
		}
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
		b.log.Error(err, ErrCouldNotMarshalEvent.Error())
		return &MessageBusError{
			Err:     ErrCouldNotMarshalEvent,
			BaseErr: err,
		}
	}

	err = b.channel.Publish(
		b.conf.ExchangeName,         // exchange
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

// AddReceiver adds a receiver for event matching the EventFilter.
func (b *rabbitEventBus) AddReceiver(receiver EventReceiver, matchers ...EventMatcher) *MessageBusError {
	if !b.isReady {
		return &MessageBusError{
			Err: ErrMessageNotConnected,
		}
	}

	if matchers == nil {
		b.log.Error(ErrMatcherMustNotBeNil, ErrMatcherMustNotBeNil.Error())
		return &MessageBusError{
			Err: ErrMatcherMustNotBeNil,
		}
	}
	if receiver == nil {
		b.log.Error(ErrReceiverMustNotBeNil, ErrReceiverMustNotBeNil.Error())
		return &MessageBusError{
			Err: ErrReceiverMustNotBeNil,
		}
	}

	resendsLeft := b.conf.MaxResends
	for {
		err := b.addReceiver(receiver, matchers...)
		if err != nil {
			b.log.Info("Adding receiver failed. Retrying...", "error", err.Cause().Error())
			select {
			case <-b.done:
				return &MessageBusError{
					Err: ErrCouldNotAddReceiver,
				}
			case <-time.After(b.conf.ResendDelay):
				if resendsLeft > 0 {
					resendsLeft--
					continue
				} else {
					b.log.Info("Adding receiver failed.")
					return &MessageBusError{
						Err: ErrCouldNotPublishEvent,
					}
				}
			}
		}

		b.log.Info("Receiver added.")
		return nil
	}
}

// addReceiver creates a queue along with bindings for the given matchers
func (b *rabbitEventBus) addReceiver(receiver EventReceiver, matchers ...EventMatcher) *MessageBusError {
	b.mu.Lock()
	defer b.mu.Unlock()
	if !b.isReady {
		return &MessageBusError{
			Err: ErrMessageNotConnected,
		}
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
		return &MessageBusError{
			Err:     ErrMessageBusConnection,
			BaseErr: err,
		}
	}
	b.log.Info(fmt.Sprintf("Queue declared '%s'.", q.Name))

	for _, matcher := range matchers {
		rabbitMatcher, ok := matcher.(*RabbitMatcher)
		if !ok {
			b.log.Error(ErrMatcherMustNotBeNil, ErrMatcherMustNotBeNil.Error())
			return &MessageBusError{
				Err: ErrMatcherMustNotBeNil,
			}
		}

		routingKey := rabbitMatcher.generateRoutingKey()
		err = b.channel.QueueBind(
			q.Name,              // queue name
			routingKey,          // routing key
			b.conf.ExchangeName, // exchange
			false,               // no wait
			nil)
		if err != nil {
			return &MessageBusError{
				Err:     ErrMessageBusConnection,
				BaseErr: err,
			}
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
		return &MessageBusError{
			Err:     ErrMessageBusConnection,
			BaseErr: err,
		}
	}
	go b.handle(q.Name, msgs, receiver)

	return nil
}

// Matcher returns a new EventMatcher of type RabbitMatcher
func (b *rabbitEventBus) Matcher() EventMatcher {
	matcher := &RabbitMatcher{
		routingKeyPrefix: b.conf.RoutingKeyPrefix,
	}
	return matcher.Any()
}

// Close will cleanly shutdown the channel and connection.
func (b *rabbitEventBus) Close() error {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.log.Info("Shutting down...")

	close(b.done)
	b.isReady = false

	err := b.connection.Close()
	if err != nil {
		return err
	}
	b.log.Info("Shutdown complete.")

	return nil
}

// changeChannel takes a new channel to the queue,
// and updates the channel listeners to reflect this.
func (b *rabbitEventBus) changeChannel(channel *amqp.Channel) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.channel = channel
	b.notifyChanClose = make(chan *amqp.Error)
	b.notifyConfirm = make(chan amqp.Confirmation, 1)
	b.channel.NotifyClose(b.notifyChanClose)
	b.channel.NotifyPublish(b.notifyConfirm)
	b.isReady = true
}

// init will initialize channel & declare queue
func (b *rabbitEventBus) init(conn *amqp.Connection) error {
	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	// Indicate we want confirmation of the publish.
	err = ch.Confirm(false)
	if err != nil {
		return err
	}
	// Indicate we only want 1 message to acknowledge at a time.
	if err := ch.Qos(1, 0, false); err != nil {
		return err
	}

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

	b.changeChannel(ch)

	return nil
}

// handleReInit will wait for a channel error
// and then continuously attempt to re-initialize both channels
func (b *rabbitEventBus) handleReInit(conn *amqp.Connection) bool {
	for {
		err := b.init(conn)
		if err != nil {
			b.log.Info("Failed to initialize channel. Retrying...", "error", err.Error())

			select {
			case <-b.done:
				b.log.Info("Aborting init. Shutting down...")
				return false
			case <-time.After(b.conf.ReInitDelay):
				continue
			}
		}

		select {
		case errConnClose := <-b.notifyConnClose:
			if _, ok := <-b.done; !ok {
				b.log.Info("Connection closed caused by shutdown.")
				return false
			}

			b.mu.Lock()
			defer b.mu.Unlock()
			b.isReady = false
			if errConnClose != nil {
				b.log.Info("Connection closed. Reconnecting...", "error", errConnClose.Error())
			} else {
				b.log.Info("Connection closed. Reconnecting...")
			}
			return true
		case errChanClose := <-b.notifyChanClose:
			if _, ok := <-b.done; !ok {
				b.log.Info("Channel closed.")
				return false
			}

			b.mu.Lock()
			defer b.mu.Unlock()
			b.isReady = false
			if errChanClose != nil {
				b.log.Info("Channel closed. Re-running init...", "error", errChanClose.Error())
			} else {
				b.log.Info("Channel closed. Re-running init...")
			}
		}
	}
}

// changeConnection takes a new connection to the queue,
// and updates the close listener to reflect this.
func (b *rabbitEventBus) changeConnection(connection *amqp.Connection) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.connection = connection
	b.notifyConnClose = make(chan *amqp.Error)
	b.connection.NotifyClose(b.notifyConnClose)
}

// connect will create a new AMQP connection
func (b *rabbitEventBus) connect(addr string) (*amqp.Connection, error) {
	b.log.Info("Attempting to connect...")

	conf := amqp.Config{
		Dial: func(network, addr string) (net.Conn, error) {
			return net.DialTimeout(network, addr, 10*time.Second)
		},
		Heartbeat: 10 * time.Second,
	}
	conn, err := amqp.DialConfig(addr, conf)

	if err != nil {
		return nil, err
	}

	b.changeConnection(conn)
	b.log.Info("Connection established.")

	return conn, nil
}

// handleReconnect will wait for a connection error on
// notifyConnClose, and then continuously attempt to reconnect.
func (b *rabbitEventBus) handleReconnect(addr string) {
	for {
		conn, err := b.connect(addr)
		if err != nil {
			b.log.Info("Failed to connect. Retrying...", "error", err.Error())

			select {
			case <-b.done:
				return
			case <-time.After(b.conf.ReconnectDelay):
			}
			continue
		}

		if !b.handleReInit(conn) {
			break
		}
	}
	b.log.Info("Automatic reconnect stopped.")
}

// generateRoutingKey generates the routing key for an event.
func (b *rabbitEventBus) generateRoutingKey(event storage.Event) string {
	return fmt.Sprintf("%s.%s.%s", b.conf.RoutingKeyPrefix, event.AggregateType(), event.EventType())
}

// handle handles the routing of the received messages and ack/nack based on receiver result
func (b *rabbitEventBus) handle(qName string, msgs <-chan amqp.Delivery, receiver EventReceiver) {
	b.log.Info(fmt.Sprintf("Handler for queue '%s' started.", qName))
	for {
		select {
		case d := <-msgs:
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
		case <-b.done:
			b.log.Info(fmt.Sprintf("Handler for queue '%s' stopped.", qName))
			return
		}
	}
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
