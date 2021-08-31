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

package usecases

import (
	"time"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	esApi "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/eventsourcing"
	es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

var _ = Describe("Converters", func() {
	It("can convert to storage query from proto filter", func() {
		aggregateId := uuid.New()
		aggregateType := es.AggregateType("TestAggregateType")
		maxTimestamp := time.Now().UTC()
		minTimestamp := maxTimestamp.Add(-1 * time.Hour)

		pf := &esApi.EventFilter{
			AggregateId:   wrapperspb.String(aggregateId.String()),
			AggregateType: wrapperspb.String(aggregateType.String()),
			MinVersion:    wrapperspb.UInt64(1),
			MaxVersion:    wrapperspb.UInt64(4),
			MinTimestamp:  timestamppb.New(minTimestamp),
			MaxTimestamp:  timestamppb.New(maxTimestamp),
		}

		q, err := NewStoreQueryFromProto(pf)
		Expect(err).ToNot(HaveOccurred())
		Expect(q).ToNot(BeNil())
		Expect(*q.AggregateId).To(Equal(aggregateId))
		Expect(*q.AggregateType).To(Equal(aggregateType))
		Expect(*q.MinVersion).To(Equal(pf.MinVersion.GetValue()))
		Expect(*q.MaxVersion).To(Equal(pf.MaxVersion.GetValue()))
		Expect(q.MinTimestamp).To(Equal(&minTimestamp))
		Expect(q.MaxTimestamp).To(Equal(&maxTimestamp))
	})
})
