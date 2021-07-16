package eventsourcing

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

type testCommandHandler struct {
	val int32
}

func (c *testCommandHandler) HandleCommand(ctx context.Context, cmd Command) (*CommandReply, error) {
	switch cmd := cmd.(type) {
	case *testCommand:
		cmd.Test = fmt.Sprintf("%s%v", cmd.Test, c.val)
		reply := &CommandReply{
			// This simulates a create command. Set a new ID.
			Id:      uuid.New(),
			Version: uint64(cmd.TestCount),
		}
		return reply, nil
	}
	return nil, fmt.Errorf("couldn't handle command")
}

func newTestCommandHandler() *testCommandHandler {
	return &testCommandHandler{}
}
