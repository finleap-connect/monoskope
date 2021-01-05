package messaging

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/streadway/amqp"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/logger"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/storage"
)

const (
	exchangeName   = "m8_events"     // name of the monoskope exchange
	reconnectDelay = 5 * time.Second // When reconnecting to the server after connection failure
	reInitDelay    = 3 * time.Second // When setting up the channel after a channel exception
	resendDelay    = 3 * time.Second // When resending messages the server didn't confirm
	maxResend      = 5               // How many times resending messages the server didn't confirm
)

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
	addr             string
	connection       *amqp.Connection
	channel          *amqp.Channel
	notifyConnClose  chan *amqp.Error
	notifyChanClose  chan *amqp.Error
	notifyConfirm    chan amqp.Confirmation
	isReady          bool
	routingKeyPrefix string
	name             string
	done             chan bool
	mu               sync.Mutex
}

// changeChannel takes a new channel to the queue,
// and updates the channel listeners to reflect this.
func (b *RabbitEventBus) changeChannel(channel *amqp.Channel) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.channel = channel
	b.notifyChanClose = make(chan *amqp.Error)
	b.notifyConfirm = make(chan amqp.Confirmation, 1)
	b.channel.NotifyClose(b.notifyChanClose)
	b.channel.NotifyPublish(b.notifyConfirm)
}

// init will initialize channel & declare queue
func (b *RabbitEventBus) init(conn *amqp.Connection) error {
	ch, err := conn.Channel()
	if err != nil {
		return err
	}

	// Indicate we only want 1 message to acknowledge at a time.
	if err := ch.Qos(1, 0, false); err != nil {
		return err
	}

	err = ch.Confirm(false)
	if err != nil {
		return err
	}

	err = ch.ExchangeDeclare(
		exchangeName,       // name
		amqp.ExchangeTopic, // type
		true,               // durable
		false,              // auto-deleted
		false,              // internal
		false,              // no-wait
		nil,                // arguments
	)
	if err != nil {
		return err
	}

	b.changeChannel(ch)
	b.isReady = true

	return nil
}

// handleReInit will wait for a channel error
// and then continuously attempt to re-initialize both channels
func (b *RabbitEventBus) handleReInit(conn *amqp.Connection) bool {
	for {
		b.isReady = false

		err := b.init(conn)

		if err != nil {
			b.log.Info("Failed to initialize channel. Retrying...", "error", err.Error())

			select {
			case <-b.done:
				b.log.Info("Automatic reinit stopped.")
				return false
			case <-time.After(reInitDelay):
				continue
			}
		}

		select {
		case errConnClose := <-b.notifyConnClose:
			if _, ok := <-b.done; !ok {
				b.log.Info("Connection closed.")
				return false
			}
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
func (b *RabbitEventBus) changeConnection(connection *amqp.Connection) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.connection = connection
	b.notifyConnClose = make(chan *amqp.Error)
	b.connection.NotifyClose(b.notifyConnClose)
}

// connect will create a new AMQP connection
func (b *RabbitEventBus) connect(addr string) (*amqp.Connection, error) {
	b.log.Info("Attempting to connect...")
	conn, err := amqp.Dial(addr)

	if err != nil {
		return nil, err
	}

	b.changeConnection(conn)
	b.log.Info("Connection established.")

	return conn, nil
}

// handleReconnect will wait for a connection error on
// notifyConnClose, and then continuously attempt to reconnect.
func (b *RabbitEventBus) handleReconnect(addr string) {
	for {
		b.isReady = false

		conn, err := b.connect(addr)
		if err != nil {
			b.log.Info("Failed to connect. Retrying...", "error", err.Error())

			select {
			case <-b.done:
				return
			case <-time.After(reconnectDelay):
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
func (b *RabbitEventBus) generateRoutingKey(event storage.Event) string {
	return fmt.Sprintf("%s.%s.%s", b.routingKeyPrefix, event.AggregateType(), event.EventType())
}

// handle handles the routing of the received messages and ack/nack based on receiver result
func (b *RabbitEventBus) handle(qName string, msgs <-chan amqp.Delivery, receiver EventReceiver) {
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

/*
NewRabbitEventBusPublisher creates a new EventBusPublisher for rabbitmq.

- routingKeyPrefix defaults to "m8"
*/
func NewRabbitEventBusPublisher(addr, name, routingKeyPrefix string) (EventBusPublisher, error) {
	if routingKeyPrefix == "" {
		routingKeyPrefix = "m8"
	}
	b := &RabbitEventBus{
		routingKeyPrefix: routingKeyPrefix,
		log:              logger.WithName("publisher").WithValues("name", name),
		addr:             addr,
		done:             make(chan bool),
	}
	return b, nil
}

/*
NewRabbitEventBusConsumer creates a new EventBusConsumer for rabbitmq.

- routingKeyPrefix defaults to "m8"
*/
func NewRabbitEventBusConsumer(addr, name, routingKeyPrefix string) (EventBusConsumer, error) {
	if routingKeyPrefix == "" {
		routingKeyPrefix = "m8"
	}
	b := &RabbitEventBus{
		routingKeyPrefix: routingKeyPrefix,
		log:              logger.WithName("consumer").WithValues("name", name),
		name:             name,
		addr:             addr,
		done:             make(chan bool),
	}
	return b, nil
}

// Connect starts automatic reconnect with rabbitmq
func (b *RabbitEventBus) Connect(ctx context.Context) *MessageBusError {
	go b.handleReconnect(b.addr)
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
func (b *RabbitEventBus) PublishEvent(ctx context.Context, event storage.Event) *MessageBusError {
	if !b.isReady {
		return &MessageBusError{
			Err: ErrMessageNotConnected,
		}
	}

	resendsLeft := maxResend
	for {
		err := b.publishEvent(ctx, event)
		if err != nil {
			b.log.Error(err, "Publish failed. Retrying...")
			select {
			case <-b.done:
				return &MessageBusError{
					Err: ErrCouldNotPublishEvent,
				}
			case <-time.After(resendDelay):
				continue
			}
		}

		select {
		case confirm := <-b.notifyConfirm:
			if confirm.Ack {
				b.log.Info("Publish confirmed.")
				return nil
			}
		case <-b.done:
			return &MessageBusError{
				Err: ErrCouldNotPublishEvent,
			}
		case <-time.After(resendDelay):
			b.log.Info("Publish wasn't confirmed. Retrying...", "resends left", resendsLeft)
		}

		if resendsLeft > 0 {
			resendsLeft--
			continue
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
func (b *RabbitEventBus) publishEvent(ctx context.Context, event storage.Event) *MessageBusError {
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

// AddReceiver adds a receiver for event matching the EventFilter.
func (b *RabbitEventBus) AddReceiver(receiver EventReceiver, matchers ...EventMatcher) *MessageBusError {
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
			q.Name,       // queue name
			routingKey,   // routing key
			exchangeName, // exchange
			false,        // no wait
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
		q.Name,                               // queue
		fmt.Sprintf("%s-%s", q.Name, b.name), // consumer
		false,                                // auto ack
		true,                                 // exclusive
		false,                                // no local
		false,                                // no wait
		nil,                                  // args
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
func (b *RabbitEventBus) Matcher() EventMatcher {
	matcher := &RabbitMatcher{
		topicPrefix: b.routingKeyPrefix,
	}
	return matcher.Any()
}

// Close will cleanly shutdown the channel and connection.
func (b *RabbitEventBus) Close() error {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.log.Info("Closing channel and connection...")

	close(b.done)
	b.isReady = false

	err := b.connection.Close()
	if err != nil {
		return err
	}
	b.log.Info("Shutdown complete.")

	return nil
}

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
