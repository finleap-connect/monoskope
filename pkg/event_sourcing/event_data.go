package event_sourcing

import (
	"github.com/golang/protobuf/ptypes"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/runtime/protoiface"
	"google.golang.org/protobuf/types/known/anypb"
)

// EventData is any additional data for an event.
type EventData *anypb.Any

func NewEventData() EventData {
	return &anypb.Any{}
}

func NewEventDataFromProto(m protoiface.MessageV1) (EventData, error) {
	return ptypes.MarshalAny(m)
}

func NewProtoFromEventData(data EventData, m protoiface.MessageV1) error {
	return ptypes.UnmarshalAny(data, m)
}

// MarshalEventData marshals any given object into event data.
func MarshalEventData(v EventData) ([]byte, error) {
	var ed *anypb.Any = v
	return protojson.Marshal(ed)
}

// UnmarshalEventData unmarshals event data to a given type.
func UnmarshalEventData(rawData []byte, v EventData) error {
	var ed *anypb.Any = v
	return protojson.Unmarshal(rawData, ed)
}
