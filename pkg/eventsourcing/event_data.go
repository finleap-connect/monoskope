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
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

// EventData is any additional data for an event. Internally this is represented by protojson.
type EventData []byte

// toEventDataFromAny marshalls a given any to protojson
func toEventDataFromAny(a *anypb.Any) EventData {
	bytes, err := protojson.Marshal(a)
	if err != nil {
		panic(err)
	}
	return bytes
}

// ToEventDataFromProto marshalls m into EventData.
func ToEventDataFromProto(m proto.Message) EventData {
	a := &anypb.Any{}
	if err := a.MarshalFrom(m); err != nil {
		panic(err)
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

// ToProto unmarshals the contents the EventData into m.
func (d EventData) ToProto(m proto.Message) error {
	a, err := d.toAny()
	if err != nil {
		return err
	}
	return a.UnmarshalTo(m)
}

// Unmarshal deserializes the EventData into a proto message using a type resolved from
// the type URL.
func (d EventData) Unmarshal() (proto.Message, error) {
	a, err := d.toAny()
	if err != nil {
		return nil, err
	}
	return a.UnmarshalNew()
}
