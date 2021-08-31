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

package projectors

import (
	"context"
	"time"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/domain/eventdata"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/constants/aggregates"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/constants/events"
	metadata "gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/metadata"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/projections"
	es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing"
)

var (
	expectedDisplayName         = "the one cluster"
	expectedName                = "one-cluster"
	expectedApiServerAddress    = "one.example.com"
	expectedClusterCACertBundle = []byte("This should be a certificate")
	expectedJWT                 = "thisisnotajwt"
)

var _ = Describe("domain/cluster_repo", func() {
	ctx := context.Background()
	userId := uuid.New()

	mdManager, err := metadata.NewDomainMetadataManager(ctx)
	Expect(err).ToNot(HaveOccurred())
	mdManager.SetUserInformation(&metadata.UserInformation{
		Id:     userId,
		Name:   "admin",
		Email:  "admin@monoskope.io",
		Issuer: "monoskope",
	})
	ctx = mdManager.GetContext()

	It("can handle ClusterCreated events", func() {
		clusterProjector := NewClusterProjector()
		clusterProjection := clusterProjector.NewProjection(uuid.New())
		protoClusterCreatedEventData := &eventdata.ClusterCreated{
			Name:                expectedDisplayName,
			Label:               expectedName,
			ApiServerAddress:    expectedApiServerAddress,
			CaCertificateBundle: expectedClusterCACertBundle,
		}
		clusterCreatedEventData := es.ToEventDataFromProto(protoClusterCreatedEventData)
		event := es.NewEvent(ctx, events.ClusterCreated, clusterCreatedEventData, time.Now().UTC(), aggregates.Cluster, uuid.New(), 1)

		clusterProjection, err := clusterProjector.Project(context.Background(), event, clusterProjection)
		Expect(err).NotTo(HaveOccurred())

		Expect(clusterProjection.Version()).To(Equal(uint64(1)))

		cluster, ok := clusterProjection.(*projections.Cluster)
		Expect(ok).To(BeTrue())

		Expect(cluster.GetDisplayName()).To(Equal(expectedDisplayName))
		Expect(cluster.GetName()).To(Equal(expectedName))
		Expect(cluster.GetApiServerAddress()).To(Equal(expectedApiServerAddress))
		Expect(cluster.GetCaCertBundle()).To(Equal(expectedClusterCACertBundle))

		dp := cluster.DomainProjection
		Expect(dp.Created).ToNot(BeNil())
	})

	It("can handle ClusterCreatedV2 events", func() {
		clusterProjector := NewClusterProjector()
		clusterProjection := clusterProjector.NewProjection(uuid.New())
		protoClusterCreatedEventData := &eventdata.ClusterCreatedV2{
			DisplayName:         expectedDisplayName,
			Name:                expectedName,
			ApiServerAddress:    expectedApiServerAddress,
			CaCertificateBundle: expectedClusterCACertBundle,
		}
		clusterCreatedEventData := es.ToEventDataFromProto(protoClusterCreatedEventData)
		event := es.NewEvent(ctx, events.ClusterCreatedV2, clusterCreatedEventData, time.Now().UTC(), aggregates.Cluster, uuid.New(), 1)

		clusterProjection, err := clusterProjector.Project(context.Background(), event, clusterProjection)
		Expect(err).NotTo(HaveOccurred())

		Expect(clusterProjection.Version()).To(Equal(uint64(1)))

		cluster, ok := clusterProjection.(*projections.Cluster)
		Expect(ok).To(BeTrue())

		Expect(cluster.GetDisplayName()).To(Equal(expectedDisplayName))
		Expect(cluster.GetName()).To(Equal(expectedName))
		Expect(cluster.GetApiServerAddress()).To(Equal(expectedApiServerAddress))
		Expect(cluster.GetCaCertBundle()).To(Equal(expectedClusterCACertBundle))

		dp := cluster.DomainProjection
		Expect(dp.Created).ToNot(BeNil())
	})

	It("can handle ClusterBootstrapTokenCreated events", func() {
		clusterProjector := NewClusterProjector()
		clusterProjection := clusterProjector.NewProjection(uuid.New())
		protoClusterCreatedEventData := &eventdata.ClusterCreatedV2{
			DisplayName:         expectedDisplayName,
			Name:                expectedName,
			ApiServerAddress:    expectedApiServerAddress,
			CaCertificateBundle: expectedClusterCACertBundle,
		}
		clusterCreatedEventData := es.ToEventDataFromProto(protoClusterCreatedEventData)
		clusterProjection, err := clusterProjector.Project(context.Background(), es.NewEvent(ctx, events.ClusterCreatedV2, clusterCreatedEventData, time.Now().UTC(), aggregates.Cluster, uuid.New(), 1), clusterProjection)
		Expect(err).NotTo(HaveOccurred())

		Expect(clusterProjection.Version()).To(Equal(uint64(1)))

		cluster, ok := clusterProjection.(*projections.Cluster)
		Expect(ok).To(BeTrue())
		Expect(cluster.GetDisplayName()).To(Equal(expectedDisplayName))
		Expect(cluster.GetName()).To(Equal(expectedName))
		Expect(cluster.GetApiServerAddress()).To(Equal(expectedApiServerAddress))
		Expect(cluster.GetCaCertBundle()).To(Equal(expectedClusterCACertBundle))

		protoTokenCreatedEventData := &eventdata.ClusterBootstrapTokenCreated{
			Jwt: expectedJWT,
		}
		tokenCreatedEventData := es.ToEventDataFromProto(protoTokenCreatedEventData)
		clusterProjection, err = clusterProjector.Project(context.Background(), es.NewEvent(ctx, events.ClusterBootstrapTokenCreated, tokenCreatedEventData, time.Now().UTC(), aggregates.Cluster, uuid.New(), 1), clusterProjection)
		Expect(err).NotTo(HaveOccurred())
		Expect(clusterProjection.Version()).To(Equal(uint64(2)))
		cluster, ok = clusterProjection.(*projections.Cluster)
		Expect(ok).To(BeTrue())
		Expect(cluster.GetBootstrapToken()).To(Equal(expectedJWT))

		dp := cluster.DomainProjection
		Expect(dp.LastModified).ToNot(BeNil())
	})

})
