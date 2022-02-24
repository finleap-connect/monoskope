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
	cmdApi "github.com/finleap-connect/monoskope/pkg/api/eventsourcing/commands"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"google.golang.org/protobuf/types/known/anypb"
)

var _ = Describe("EventData", func() {
	getProto := func() *cmdApi.TestCommandData {
		return &cmdApi.TestCommandData{Test: "Hello world!"}
	}
	eventDataFromProto := func() EventData {
		eventData := ToEventDataFromProto(getProto())
		Expect(eventData).ToNot(BeNil())
		return eventData
	}

	It("can create from proto", func() {
		_ = eventDataFromProto()
	})
	It("can create from any", func() {
		proto := getProto()
		any := anypb.Any{}
		err := any.MarshalFrom(proto)
		Expect(err).ToNot(HaveOccurred())

		eventData := toEventDataFromAny(&any)
		Expect(eventData).ToNot(BeNil())
	})
	It("can unmarshal to any", func() {
		eventData := eventDataFromProto()
		any, err := eventData.toAny()
		Expect(err).ToNot(HaveOccurred())
		Expect(any).ToNot(BeNil())
	})
	It("can unmarshal to proto", func() {
		eventData := eventDataFromProto()
		proto := &cmdApi.TestCommandData{}
		err := eventData.ToProto(proto)
		Expect(err).ToNot(HaveOccurred())
		Expect(proto.GetTest()).To(Equal(getProto().GetTest()))
	})
	It("can resolve message type on unmarshal", func() {
		eventData := eventDataFromProto()
		proto, err := eventData.Unmarshal()
		Expect(err).ToNot(HaveOccurred())
		Expect(proto).ToNot(BeNil())
		_, ok := proto.(*cmdApi.TestCommandData)
		Expect(ok).To(BeTrue())
	})
})
