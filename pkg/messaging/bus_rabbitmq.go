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
	shutdown        chan bool
}

// NewRabbitEventBusPublisher creates a new EventBusPublisher for rabbitmq.
func NewRabbitEventBusPublisher(conf *rabbitEventBusConfig) (EventBusPublisher, error) {
	err := conf.Validate()
	if err != nil {
		return nil, err
	}

	b := &rabbitEventBus{
		log:      logger.WithName("publisher").WithValues("name", conf.Name),
		conf:     conf,
		shutdown: make(chan bool),
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
		log:      logger.WithName("consumer").WithValues("name", conf.Name),
		conf:     conf,
		shutdown: make(chan bool),
	}
	return b, nil
}

// Connect starts automatic reconnect with rabbitmq
func (b *rabbitEventBus) Connect(ctx context.Context) *messageBusError {
	go b.handleReconnect(b.conf.Url)
	for {
		select {
		case <-time.After(300 * time.Millisecond):
			if b.isReady {
				return nil
			}
		case <-b.shutdown:
		case <-ctx.Done():
			return &messageBusError{
				Err: ErrMessageNotConnected,
			}
		}
	}
}

// PublishEvent publishes the event on the bus.
func (b *rabbitEventBus) PublishEvent(ctx context.Context, event storage.Event) *messageBusError {
	resendsLeft := b.conf.MaxResends
	for {
		err := b.publishEvent(event)
		if err != nil {
			b.log.Error(err, "Publish failed. Retrying...")
			select {
			case <-b.shutdown:
				b.log.Info("Publish failed because of shutdown.")
				return &messageBusError{
					Err: ErrCouldNotPublishEvent,
				}
			case <-ctx.Done():
				b.log.Info("Publish failed because context deadline exceeded.")
				return &messageBusError{
					Err:     ErrContextDeadlineExceeded,
					BaseErr: ctx.Err(),
				}
			case <-time.After(b.conf.ResendDelay):
				if resendsLeft > 0 {
					resendsLeft--
					continue
				} else {
					b.log.Info("Publish failed.")
					return &messageBusError{
						Err: ErrCouldNotPublishEvent,
					}
				}
			}
		}

		retry := true
		for retry {
			select {
			case confirm := <-b.notifyConfirm:
				if confirm.Ack {
					b.log.Info("Publish confirmed.", "DeliveryTag", confirm.DeliveryTag)
					return nil
				}
			case <-time.After(b.conf.ResendDelay):
				b.log.Info("Publish wasn't confirmed. Retrying...", "resends left", resendsLeft)
				retry = false
			case <-ctx.Done():
				b.log.Info("Publish failed because context deadline exceeded.")
				return &messageBusError{
					Err:     ErrContextDeadlineExceeded,
					BaseErr: ctx.Err(),
				}
			case <-b.shutdown:
				b.log.Info("Publish failed.")
				return &messageBusError{
					Err: ErrCouldNotPublishEvent,
				}
			}
		}

		if resendsLeft > 0 {
			resendsLeft--
			continue
		}

		b.log.Info("Publish failed.")
		return &messageBusError{
			Err: ErrCouldNotPublishEvent,
		}
	}
}

// publishEvent will push to the queue without checking for
// confirmation. It returns an error if it fails to connect.
// No guarantees are provided for whether the server will
// recieve the message.
func (b *rabbitEventBus) publishEvent(event storage.Event) *messageBusError {
	if !b.isReady {
		return &messageBusError{
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
		return &messageBusError{
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
			ContentType:  "text/json",
			Body:         bytes,
			DeliveryMode: amqp.Transient,
		})

	if err != nil {
		b.log.Error(err, ErrCouldNotPublishEvent.Error())
		return &messageBusError{
			Err:     ErrCouldNotPublishEvent,
			BaseErr: err,
		}
	}
	return nil
}

// AddReceiver adds a receiver for event matching the EventFilter.
func (b *rabbitEventBus) AddReceiver(ctx context.Context, receiver EventReceiver, matchers ...EventMatcher) *messageBusError {
	if matchers == nil {
		b.log.Error(ErrMatcherMustNotBeNil, ErrMatcherMustNotBeNil.Error())
		return &messageBusError{
			Err: ErrMatcherMustNotBeNil,
		}
	}
	if receiver == nil {
		b.log.Error(ErrReceiverMustNotBeNil, ErrReceiverMustNotBeNil.Error())
		return &messageBusError{
			Err: ErrReceiverMustNotBeNil,
		}
	}

	resendsLeft := b.conf.MaxResends
	for {
		err := b.addReceiver(receiver, matchers...)
		if err != nil {
			b.log.Info("Adding receiver failed. Retrying...", "error", err.Cause().Error())
			select {
			case <-b.shutdown:
				return &messageBusError{
					Err: ErrCouldNotAddReceiver,
				}
			case <-ctx.Done():
				return &messageBusError{
					Err: ErrContextDeadlineExceeded,
				}
			case <-time.After(b.conf.ResendDelay):
				if resendsLeft > 0 {
					resendsLeft--
					continue
				} else {
					b.log.Info("Adding receiver failed.")
					return &messageBusError{
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
func (b *rabbitEventBus) addReceiver(receiver EventReceiver, matchers ...EventMatcher) *messageBusError {
	if !b.isReady {
		return &messageBusError{
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
		return &messageBusError{
			Err:     ErrMessageBusConnection,
			BaseErr: err,
		}
	}
	b.log.Info(fmt.Sprintf("Queue declared '%s'.", q.Name))

	for _, matcher := range matchers {
		rabbitMatcher, ok := matcher.(*rabbitMatcher)
		if !ok {
			b.log.Error(ErrMatcherMustNotBeNil, ErrMatcherMustNotBeNil.Error())
			return &messageBusError{
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
			return &messageBusError{
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
		return &messageBusError{
			Err:     ErrMessageBusConnection,
			BaseErr: err,
		}
	}
	go b.handle(q.Name, msgs, receiver)

	return nil
}

// Matcher returns a new EventMatcher of type RabbitMatcher
func (b *rabbitEventBus) Matcher() EventMatcher {
	matcher := &rabbitMatcher{
		routingKeyPrefix: b.conf.RoutingKeyPrefix,
	}
	return matcher.Any()
}

// Close will cleanly shutdown the channel and connection.
func (b *rabbitEventBus) Close() error {
	b.log.Info("Shutting down...")

	b.isReady = false
	close(b.shutdown)

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
	b.channel = channel
	b.notifyChanClose = b.channel.NotifyClose(make(chan *amqp.Error))
	b.notifyConfirm = channel.NotifyPublish(make(chan amqp.Confirmation, 1))
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
			case <-b.shutdown:
				b.log.Info("Aborting init. Shutting down...")
				return false
			case <-time.After(b.conf.ReInitDelay):
				continue
			}
		}

		select {
		case <-b.shutdown:
			return false
		case errConnClose := <-b.notifyConnClose:
			b.isReady = false
			if errConnClose != nil {
				b.log.Info("Connection closed. Reconnecting...", "error", errConnClose.Error())
			} else {
				b.log.Info("Connection closed. Reconnecting...")
			}
			return true
		case errChanClose := <-b.notifyChanClose:
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
	b.connection = connection
	b.notifyConnClose = make(chan *amqp.Error)
	b.connection.NotifyClose(b.notifyConnClose)
}

// connect will create a new AMQP connection
func (b *rabbitEventBus) connect(addr string) (*amqp.Connection, error) {
	b.log.Info("Attempting to connect...")
	conn, err := amqp.DialConfig(addr, *b.conf.amqpConfig)
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
			case <-b.shutdown:
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
			b.log.Info(fmt.Sprintf("Handler received event from queue '%s'.", qName))

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
		case <-b.shutdown:
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
