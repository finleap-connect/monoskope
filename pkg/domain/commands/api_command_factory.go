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

package commands

import (
	esApi "github.com/finleap-connect/monoskope/pkg/api/eventsourcing/commands"
	es "github.com/finleap-connect/monoskope/pkg/eventsourcing"
	"github.com/google/uuid"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/known/anypb"
)

// CreateCommand builds up a new proto command with the given type and data.
func CreateCommand(aggregateId uuid.UUID, commandType es.CommandType) *esApi.Command {
	return &esApi.Command{
		Id:   aggregateId.String(),
		Type: commandType.String(),
	}
}

func AddCommandData(command *esApi.Command, commandData protoreflect.ProtoMessage) (*esApi.Command, error) {
	data, err := CreateCommandData(commandData)
	if err != nil {
		return nil, err
	}
	command.Data = data
	return command, nil
}

func CreateCommandData(commandData protoreflect.ProtoMessage) (*anypb.Any, error) {
	data := &anypb.Any{}
	if err := data.MarshalFrom(commandData); err != nil {
		return nil, err
	}
	return data, nil
}
