package messaging

import (
	"context"

	evs "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing"
)

type mockEventBus struct {
}

func NewMockEventBusPublisher() evs.EventBusPublisher {
	return &mockEventBus{}
}

func (b *mockEventBus) Connect(ctx context.Context) error {
	// mock
	return nil
}

func (b *mockEventBus) PublishEvent(ctx context.Context, event evs.Event) error {
	panic("not implemented")
}

func (b *mockEventBus) Matcher() evs.EventMatcher {
	panic("not implemented")
}

func (b *mockEventBus) Close() error {
	// mock
	return nil
}
