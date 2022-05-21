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

package projectors

import (
	"context"
	"time"

	"github.com/finleap-connect/monoskope/pkg/api/domain/eventdata"
	"github.com/finleap-connect/monoskope/pkg/domain/constants/aggregates"
	"github.com/finleap-connect/monoskope/pkg/domain/constants/events"
	metadata "github.com/finleap-connect/monoskope/pkg/domain/metadata"
	es "github.com/finleap-connect/monoskope/pkg/eventsourcing"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	expectedDisplayName         = "the one cluster"
	expectedName                = "one-cluster"
	expectedApiServerAddress    = "one.example.com"
	expectedClusterCACertBundle = []byte("This should be a certificate")
	expectedJWT                 = "thisisnotajwt"
)

var _ = Describe("domain/projectors/cluster", func() {
	ctx := context.Background()
	userId := uuid.New()

	mdManager, err := metadata.NewDomainMetadataManager(ctx)
	Expect(err).ToNot(HaveOccurred())
	mdManager.SetUserInformation(&metadata.UserInformation{
		Id:    userId,
		Name:  "admin",
		Email: "admin@monoskope.io",
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
		Expect(clusterProjection.GetDisplayName()).To(Equal(expectedDisplayName))
		Expect(clusterProjection.GetName()).To(Equal(expectedName))
		Expect(clusterProjection.GetApiServerAddress()).To(Equal(expectedApiServerAddress))
		Expect(clusterProjection.GetCaCertBundle()).To(Equal(expectedClusterCACertBundle))

		dp := clusterProjection.DomainProjection
		Expect(dp.GetCreated()).ToNot(BeNil())
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

		Expect(clusterProjection.GetDisplayName()).To(Equal(expectedDisplayName))
		Expect(clusterProjection.GetName()).To(Equal(expectedName))
		Expect(clusterProjection.GetApiServerAddress()).To(Equal(expectedApiServerAddress))
		Expect(clusterProjection.GetCaCertBundle()).To(Equal(expectedClusterCACertBundle))

		dp := clusterProjection.DomainProjection
		Expect(dp.GetCreated()).ToNot(BeNil())
	})

	It("can handle ClusterBootstrapTokenCreated events", func() {
		clusterProjector := NewClusterProjector()
		clusterProjection := clusterProjector.NewProjection(uuid.New())

		clusterProjection, err := clusterProjector.Project(
			context.Background(),
			es.NewEvent(ctx,
				events.ClusterCreatedV2,
				es.ToEventDataFromProto(&eventdata.ClusterCreatedV2{
					DisplayName:         expectedDisplayName,
					Name:                expectedName,
					ApiServerAddress:    expectedApiServerAddress,
					CaCertificateBundle: expectedClusterCACertBundle,
				}),
				time.Now().UTC(),
				aggregates.Cluster,
				uuid.New(),
				1),
			clusterProjection,
		)
		Expect(err).NotTo(HaveOccurred())
		Expect(clusterProjection.Version()).To(Equal(uint64(1)))

		Expect(clusterProjection.GetDisplayName()).To(Equal(expectedDisplayName))
		Expect(clusterProjection.GetName()).To(Equal(expectedName))
		Expect(clusterProjection.GetApiServerAddress()).To(Equal(expectedApiServerAddress))
		Expect(clusterProjection.GetCaCertBundle()).To(Equal(expectedClusterCACertBundle))

		protoTokenCreatedEventData := &eventdata.ClusterBootstrapTokenCreated{
			Jwt: expectedJWT,
		}
		tokenCreatedEventData := es.ToEventDataFromProto(protoTokenCreatedEventData)
		clusterProjection, err = clusterProjector.Project(context.Background(), es.NewEvent(ctx, events.ClusterBootstrapTokenCreated, tokenCreatedEventData, time.Now().UTC(), aggregates.Cluster, uuid.New(), 1), clusterProjection)
		Expect(err).NotTo(HaveOccurred())
		Expect(clusterProjection.Version()).To(Equal(uint64(2)))
		Expect(clusterProjection.GetBootstrapToken()).To(Equal(expectedJWT))

		dp := clusterProjection.DomainProjection
		Expect(dp.GetLastModified()).ToNot(BeNil())
	})

	It("can handle ClusterUpdated events", func() {
		clusterProjector := NewClusterProjector()
		clusterProjection := clusterProjector.NewProjection(uuid.New())

		newDisplayName := "new-display-name"
		newApiServerAddress := "https://new-api-server-address"
		newClusterCaCertificate := []byte("new-ca-cert")

		clusterProjection, err := clusterProjector.Project(
			context.Background(),
			es.NewEvent(ctx,
				events.ClusterCreatedV2,
				es.ToEventDataFromProto(&eventdata.ClusterCreatedV2{
					DisplayName:         expectedDisplayName,
					Name:                expectedName,
					ApiServerAddress:    expectedApiServerAddress,
					CaCertificateBundle: expectedClusterCACertBundle,
				}),
				time.Now().UTC(),
				aggregates.Cluster,
				uuid.New(),
				1),
			clusterProjection,
		)
		Expect(err).NotTo(HaveOccurred())
		Expect(clusterProjection.Version()).To(Equal(uint64(1)))

		clusterProjection, err = clusterProjector.Project(
			context.Background(),
			es.NewEvent(ctx,
				events.ClusterUpdated,
				es.ToEventDataFromProto(&eventdata.ClusterUpdated{
					DisplayName:         newDisplayName,
					ApiServerAddress:    newApiServerAddress,
					CaCertificateBundle: newClusterCaCertificate,
				}),
				time.Now().UTC(),
				aggregates.Cluster,
				uuid.New(),
				2),
			clusterProjection,
		)
		Expect(err).NotTo(HaveOccurred())
		Expect(clusterProjection.Version()).To(Equal(uint64(2)))

		Expect(clusterProjection.GetDisplayName()).To(Equal(newDisplayName))
		Expect(clusterProjection.GetApiServerAddress()).To(Equal(newApiServerAddress))
		Expect(clusterProjection.GetCaCertBundle()).To(Equal(newClusterCaCertificate))
	})

})
