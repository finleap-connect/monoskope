package messaging

import (
	"context"

	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/logger"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/storage"
)

// TestWrapperEventBus wraps an actual event bus and adds logging for debug purposes.
type TestWrapperEventBus struct {
	log       logger.Logger
	publisher EventBusPublisher
	consumer  EventBusConsumer
}

func NewTestEventBusPublisher(log logger.Logger, publisher EventBusPublisher) EventBusPublisher {
	return &TestWrapperEventBus{
		log:       log,
		publisher: publisher,
	}
}

func NewTestEventBusConsumer(log logger.Logger, consumer EventBusConsumer) EventBusConsumer {
	return &TestWrapperEventBus{
		log:      log,
		consumer: consumer,
	}
}

func (b *TestWrapperEventBus) Connect(ctx context.Context) *MessageBusError {
	var err *MessageBusError

	if b.publisher != nil {
		err = b.publisher.Connect(ctx)
	}

	if b.consumer != nil {
		err = b.consumer.Connect(ctx)
	}

	return err
}

func (b *TestWrapperEventBus) PublishEvent(ctx context.Context, event storage.Event) *MessageBusError {
	b.log.Info("Publishing event...", "event", event.String())
	err := b.publisher.PublishEvent(ctx, event)
	if err != nil {
		b.log.Error(err, "Error publishing event.", "event", event.String())
	} else {
		b.log.Info("Published event.", "event", event.String())
	}
	return err
}

func (b *TestWrapperEventBus) AddReceiver(receiver EventReceiver, matchers ...EventMatcher) *MessageBusError {
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

func (b *TestWrapperEventBus) Matcher() EventMatcher {
	return b.consumer.Matcher()
}

// Close frees all disposable resources
func (b *TestWrapperEventBus) Close() error {
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
