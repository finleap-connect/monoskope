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

package aggregates

import (
	"time"

	"github.com/finleap-connect/monoskope/pkg/api/domain/eventdata"
	"github.com/finleap-connect/monoskope/pkg/domain/constants/events"
	es "github.com/finleap-connect/monoskope/pkg/eventsourcing"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

var _ = Describe("Unit Test for Cluster Aggregate", func() {
	It("should set the data from a command to the resultant event", func() {
		ctx := createSysAdminCtx()
		agg := NewClusterAggregate(NewTestAggregateManager())

		reply, err := createCluster(ctx, agg)
		Expect(err).NotTo(HaveOccurred())
		Expect(reply.Id).ToNot(Equal(uuid.Nil))
		Expect(reply.Version).To(Equal(uint64(0)))

		event := agg.UncommittedEvents()[0]

		Expect(event.EventType()).To(Equal(events.ClusterCreatedV3))

		data := new(eventdata.ClusterCreatedV2)
		err = event.Data().ToProto(data)
		Expect(err).NotTo(HaveOccurred())

		Expect(data.Name).To(Equal(expectedClusterName))
		Expect(data.ApiServerAddress).To(Equal(expectedClusterApiServerAddress))
		Expect(data.CaCertificateBundle).To(Equal(expectedClusterCACertBundle))

	})

	It("should apply the data from an event to the aggregate", func() {

		ctx := createSysAdminCtx()
		agg := NewClusterAggregate(NewTestAggregateManager())

		ed := es.ToEventDataFromProto(&eventdata.ClusterCreatedV2{
			DisplayName:         expectedClusterDisplayName,
			Name:                expectedClusterName,
			ApiServerAddress:    expectedClusterApiServerAddress,
			CaCertificateBundle: expectedClusterCACertBundle,
		})
		esEvent := es.NewEvent(ctx, events.ClusterCreatedV2, ed, time.Now().UTC(),
			agg.Type(), agg.ID(), agg.Version())

		err := agg.ApplyEvent(esEvent)
		Expect(err).NotTo(HaveOccurred())

		Expect(agg.(*ClusterAggregate).name).To(Equal(expectedClusterName))
		Expect(agg.(*ClusterAggregate).apiServerAddr).To(Equal(expectedClusterApiServerAddress))
		Expect(agg.(*ClusterAggregate).caCertBundle).To(Equal(expectedClusterCACertBundle))
	})

	Context("cluster update", func() {
		It("should update the DisplayName", func() {
			ctx := createSysAdminCtx()
			agg := NewClusterAggregate(NewTestAggregateManager())

			expectedNewName := "the-new-name"
			ed := es.ToEventDataFromProto(&eventdata.ClusterUpdatedV2{
				Name: wrapperspb.String(expectedNewName),
			})
			esEvent := es.NewEvent(ctx, events.ClusterUpdatedV2, ed, time.Now().UTC(),
				agg.Type(), agg.ID(), agg.Version()+1)

			err := agg.ApplyEvent(esEvent)
			Expect(err).NotTo(HaveOccurred())

			Expect(agg.(*ClusterAggregate).name).To(Equal(expectedNewName))
		})
		It("should update the ApiServerAddress", func() {
			ctx := createSysAdminCtx()
			agg := NewClusterAggregate(NewTestAggregateManager())

			expectedValue := "https://some-new-address.io"
			ed := es.ToEventDataFromProto(&eventdata.ClusterUpdated{
				ApiServerAddress: expectedValue,
			})
			esEvent := es.NewEvent(ctx, events.ClusterUpdated, ed, time.Now().UTC(),
				agg.Type(), agg.ID(), agg.Version()+1)

			err := agg.ApplyEvent(esEvent)
			Expect(err).NotTo(HaveOccurred())

			Expect(agg.(*ClusterAggregate).apiServerAddr).To(Equal(expectedValue))
		})
		It("should update the CaCertificateBundle", func() {
			ctx := createSysAdminCtx()
			agg := NewClusterAggregate(NewTestAggregateManager())

			expectedValue := []byte("some-shiny-new-ca-cert")
			ed := es.ToEventDataFromProto(&eventdata.ClusterUpdated{
				CaCertificateBundle: expectedValue,
			})
			esEvent := es.NewEvent(ctx, events.ClusterUpdated, ed, time.Now().UTC(),
				agg.Type(), agg.ID(), agg.Version()+1)

			err := agg.ApplyEvent(esEvent)
			Expect(err).NotTo(HaveOccurred())

			Expect(agg.(*ClusterAggregate).caCertBundle).To(Equal(expectedValue))
		})
	})
})
