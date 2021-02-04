package event_sourcing

import (
	"context"
	"fmt"

	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/eventdata/test"
)

type testEventHandler struct {
}

func (c *testEventHandler) HandleEvent(ctx context.Context, event Event) error {
	switch event.EventType() {
	case testEventType:
		proto := &test.TestEventData{}
		err := event.Data().ToProto(proto)
		if err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("couldn't handle event")
}
