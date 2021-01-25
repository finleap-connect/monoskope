package messaging

import (
	"context"

	evs "gitlab.figo.systems/platform/monoskope/monoskope/pkg/event_sourcing"
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

func (b *mockEventBus) PublishEvent(ctx context.Context, event evs.Event) *messageBusError {
	panic("not implemented")
}

func (b *mockEventBus) Matcher() EventMatcher {
	panic("not implemented")
}

func (b *mockEventBus) Close() error {
	// mock
	return nil
}
