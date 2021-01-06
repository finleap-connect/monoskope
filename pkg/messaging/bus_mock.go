package messaging

import (
	"context"

	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/storage"
)

type mockEventBus struct {
}

func NewMockEventBusPublisher() EventBusPublisher {
	return &mockEventBus{}
}

func (b *mockEventBus) Connect(ctx context.Context) *messageBusError {
	// mock
	return nil
}

func (b *mockEventBus) PublishEvent(ctx context.Context, event storage.Event) *messageBusError {
	panic("not implemented")
}

func (b *mockEventBus) Matcher() EventMatcher {
	panic("not implemented")
}

func (b *mockEventBus) Close() error {
	// mock
	return nil
}
