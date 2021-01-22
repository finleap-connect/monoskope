package event_sourcing

import (
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/known/anypb"
)

// EventData is any additional data for an event.
type EventData []byte

func ToAny(d EventData) (*anypb.Any, error) {
	a := &anypb.Any{}
	err := protojson.Unmarshal(d, a)
	if err != nil {
		return nil, err
	}
	return a, nil
}

func ToEventDataFromAny(a *anypb.Any) (EventData, error) {
	bytes, err := protojson.Marshal(a)
	if err != nil {
		return bytes, err
	}
	return bytes, nil
}

func ToEventDataFromProto(m protoreflect.ProtoMessage) (EventData, error) {
	a := &anypb.Any{}
	err := a.MarshalFrom(m)
	if err != nil {
		return EventData{}, err
	}
	return ToEventDataFromAny(a)
}

func ToType(d EventData, m protoreflect.ProtoMessage) error {
	return protojson.Unmarshal(d, m)
}
