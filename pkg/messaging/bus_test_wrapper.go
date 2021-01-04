package messaging

import (
	"context"

	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/logger"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/storage"
)

// TestEventBus implements an EventBus using RabbitMQ and adds debug logging.
type TestEventBus struct {
	log       logger.Logger
	publisher EventBusPublisher
	consumer  EventBusConsumer
}

func NewTestEventBusPublisher(log logger.Logger, publisher EventBusPublisher) EventBusPublisher {
	return &TestEventBus{
		log:       log,
		publisher: publisher,
	}
}

func NewTestEventBusConsumer(log logger.Logger, consumer EventBusConsumer) EventBusConsumer {
	return &TestEventBus{
		log:      log,
		consumer: consumer,
	}
}

func (b *TestEventBus) PublishEvent(ctx context.Context, event storage.Event) *MessageBusError {
	b.log.Info("Publishing event...", "event", event.String())
	err := b.publisher.PublishEvent(ctx, event)
	if err != nil {
		b.log.Error(err, "Error publishing event.", "event", event.String())
	} else {
		b.log.Info("Published event.", "event", event.String())
	}
	return err
}

func (b *TestEventBus) AddReceiver(receiver EventReceiver, matchers ...EventMatcher) error {
	b.log.Info("Adding receiver...")
	err := b.consumer.AddReceiver(func(e storage.Event) error {
		err := receiver(e)
		if err != nil {
			b.log.Error(err, "Error calling receiver.", "event", e.String())
		} else {
			b.log.Info("Received event", "event", e.String())
		}
		return err
	}, matchers...)
	if err != nil {
		b.log.Error(err, "Error adding receiver.")
	} else {
		b.log.Info("Receiver added.")
	}
	return err
}

func (b *TestEventBus) Matcher() EventMatcher {
	return b.consumer.Matcher()
}

func (b *TestEventBus) AddErrorHandler(eh ErrorHandler) {
	b.consumer.AddErrorHandler(eh)
}

// Close frees all disposable resources
func (b *TestEventBus) Close() error {
	b.log.Info("Closing event bus...")
	if b.publisher != nil {
		if err := b.publisher.Close(); err != nil {
			return err
		}
	}
	if b.consumer != nil {
		if err := b.consumer.Close(); err != nil {
			return err
		}
	}
	b.log.Info("Closed event bus.")
	return nil
}
