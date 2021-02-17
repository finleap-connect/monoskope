package eventsourcing

import (
	"context"
	"fmt"
)

type testCommandHandler struct {
	val int32
}

func (c *testCommandHandler) HandleCommand(ctx context.Context, cmd Command) error {
	switch cmd := cmd.(type) {
	case *testCommand:
		cmd.TestCount++
		cmd.Test = fmt.Sprintf("%s%v", cmd.Test, c.val)
		return nil
	}
	return fmt.Errorf("couldn't handle command")
}

func newTestCommandHandler() *testCommandHandler {
	return &testCommandHandler{}
}
