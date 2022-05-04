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
	"context"
	"time"

	esApi "github.com/finleap-connect/monoskope/pkg/api/eventsourcing"
	testEd "github.com/finleap-connect/monoskope/pkg/api/eventsourcing/eventdata"
	"github.com/finleap-connect/monoskope/pkg/eventsourcing/errors"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

var _ = Describe("EventData", func() {
	var (
		testEventType     EventType     = "TestEventType"
		testAggregateType AggregateType = "TestAggregateType"
	)

	checkProtoStorageEventEquality := func(pe *esApi.Event, se Event) {
		Expect(pe).ToNot(BeNil())
		Expect(se).ToNot(BeNil())
		Expect(pe.Type).To(Equal(se.EventType().String()))
		Expect(pe.Timestamp.AsTime()).To(Equal(se.Timestamp()))
		Expect(pe.AggregateId).To(Equal(se.AggregateID().String()))
		Expect(pe.AggregateType).To(Equal(se.AggregateType().String()))
		Expect(pe.AggregateVersion.GetValue()).To(Equal(se.AggregateVersion()))
		Expect(se.Data()).To(Equal(EventData(pe.Data)))
	}

	It("can convert to storage event from proto", func() {
		proto := &testEd.TestEventData{Hello: "world"}
		ed := ToEventDataFromProto(proto)

		timestamp := time.Now().UTC()
		pe := &esApi.Event{
			Type:             testEventType.String(),
			Timestamp:        timestamppb.New(timestamp),
			AggregateId:      uuid.New().String(),
			AggregateType:    testAggregateType.String(),
			AggregateVersion: wrapperspb.UInt64(0),
			Data:             ed,
		}

		se, err := NewEventFromProto(pe)
		Expect(err).ToNot(HaveOccurred())

		checkProtoStorageEventEquality(pe, se)
	})
	It("can convert to proto event from storage", func() {
		timestamp := time.Now().UTC()
		aggregateId := uuid.New()

		ed := ToEventDataFromProto(&testEd.TestEventData{Hello: "world"})
		se := NewEvent(
			context.Background(),
			EventType("TestType"),
			ed,
			timestamp,
			AggregateType("TestAggregateType"),
			aggregateId,
			0)
		pe := NewProtoFromEvent(se)

		checkProtoStorageEventEquality(pe, se)
	})
	It("fails to convert to storage query from proto filter for invalid aggregate id", func() {
		proto := &testEd.TestEventData{Hello: "world"}
		ed := ToEventDataFromProto(proto)
		pe := &esApi.Event{
			Type:             testEventType.String(),
			Timestamp:        timestamppb.New(time.Now().UTC()),
			AggregateId:      "", // invalid id
			AggregateType:    testAggregateType.String(),
			AggregateVersion: wrapperspb.UInt64(0),
			Data:             ed,
		}

		se, err := NewEventFromProto(pe)
		Expect(err).To(HaveOccurred())
		Expect(se).To(BeNil())
		Expect(err).To(Equal(errors.ErrCouldNotParseAggregateId))
	})
})
