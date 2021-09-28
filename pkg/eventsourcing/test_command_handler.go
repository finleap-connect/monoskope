// Copyright 2021 Monoskope Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
