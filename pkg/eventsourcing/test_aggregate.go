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

	testEd "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/eventsourcing/eventdata"
)

const (
	testAggregateType AggregateType = "TestAggregateType"
	testEventType     EventType     = "TestEventType"
)

type testAggregate struct {
	*BaseAggregate
	Test string
}

func newTestAggregate() *testAggregate {
	return &testAggregate{
		BaseAggregate: NewBaseAggregate(testAggregateType),
	}
}

func (a *testAggregate) HandleCommand(ctx context.Context, cmd Command) (*CommandReply, error) {
	switch cmd := cmd.(type) {
	case *testCommand:
		ed := ToEventDataFromProto(&testEd.TestEventData{
			Hello: cmd.TestCommandData.GetTest(),
		})
		agg := a.AppendEvent(ctx, testEventType, ed)
		ret := &CommandReply{
			Id:      agg.AggregateID(),
			Version: agg.AggregateVersion(),
		}
		return ret, nil
	}
	return nil, fmt.Errorf("couldn't handle command")
}

func (a *testAggregate) ApplyEvent(ev Event) error {
	return nil
}
