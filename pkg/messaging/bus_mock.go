package messaging

import (
	"context"

	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/storage"
)

type MockEventBus struct {
}

func NewMockEventBusPublisher() EventBusPublisher {
	return &MockEventBus{}
}

func (b *MockEventBus) Connect(ctx context.Context) *MessageBusError {
	// mock
	return nil
}

func (b *MockEventBus) PublishEvent(ctx context.Context, event storage.Event) *MessageBusError {
	panic("not implemented")
}

func (b *MockEventBus) Matcher() EventMatcher {
	panic("not implemented")
}

func (b *MockEventBus) Close() error {
	// mock
	return nil
}
