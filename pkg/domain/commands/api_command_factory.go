package commands

import (
	"github.com/google/uuid"
	esApi "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/eventsourcing/commands"
	es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/known/anypb"
)

// CreateCommand builds up a new proto command with the given type and data.
func CreateCommand(id uuid.UUID, commandType es.CommandType, commandData protoreflect.ProtoMessage) (*esApi.Command, error) {
	data, err := CreateCommandData(commandData)
	if err != nil {
		return nil, err
	}

	return &esApi.Command{
		Id:   id.String(),
		Type: commandType.String(),
		Data: data,
	}, nil
}

func CreateCommandData(commandData protoreflect.ProtoMessage) (*anypb.Any, error) {
	data := &anypb.Any{}
	if err := data.MarshalFrom(commandData); err != nil {
		return nil, err
	}
	return data, nil
}
