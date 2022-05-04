// Copyright 2022 Monoskope Authors
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
	cmdApi "github.com/finleap-connect/monoskope/pkg/api/eventsourcing/commands"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/anypb"
)

const (
	testCommandType CommandType = "TestCommandType"
)

// testCommand is a command for tests.
type testCommand struct {
	aggregateId uuid.UUID
	cmdApi.TestCommandData
}

func (c *testCommand) AggregateID() uuid.UUID { return c.aggregateId }
func (c *testCommand) AggregateType() AggregateType {
	return testAggregateType
}
func (c *testCommand) CommandType() CommandType { return testCommandType }
func (c *testCommand) SetData(a *anypb.Any) error {
	return a.UnmarshalTo(&c.TestCommandData)
}
