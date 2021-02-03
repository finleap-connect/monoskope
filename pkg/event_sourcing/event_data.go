package event_sourcing

import (
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/known/anypb"
)

// EventData is any additional data for an event. Internally this is represented by protojson.
type EventData []byte

// toEventDataFromAny marshalls a given any to protojson
func toEventDataFromAny(a *anypb.Any) (EventData, error) {
	bytes, err := protojson.Marshal(a)
	if err != nil {
		return bytes, err
	}
	return bytes, nil
}

// ToEventDataFromProto marshalls m into EventData.
func ToEventDataFromProto(m protoreflect.ProtoMessage) (EventData, error) {
	a := &anypb.Any{}
	err := a.MarshalFrom(m)
	if err != nil {
		return EventData{}, err
	}
	return toEventDataFromAny(a)
}

// toAny unmarshalls protojson to an any
func (d EventData) toAny() (*anypb.Any, error) {
	a := &anypb.Any{}
	err := protojson.Unmarshal(d, a)
	if err != nil {
		return nil, err
	}
	return a, nil
}

// ToProto the contents the EventData into m.
func (d EventData) ToProto(m protoreflect.ProtoMessage) error {
	a, err := d.toAny()
	if err != nil {
		return err
	}
	return a.UnmarshalTo(m)
}
