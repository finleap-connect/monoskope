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

package internal

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/finleap-connect/monoskope/internal/gateway/auth"
	mock_reactor "github.com/finleap-connect/monoskope/internal/test/reactor"
	domainApi "github.com/finleap-connect/monoskope/pkg/api/domain"
	cmdData "github.com/finleap-connect/monoskope/pkg/api/domain/commanddata"
	"github.com/finleap-connect/monoskope/pkg/api/domain/common"
	"github.com/finleap-connect/monoskope/pkg/api/domain/eventdata"
	"github.com/finleap-connect/monoskope/pkg/api/domain/projections"
	esApi "github.com/finleap-connect/monoskope/pkg/api/eventsourcing"
	cmd "github.com/finleap-connect/monoskope/pkg/domain/commands"
	"github.com/finleap-connect/monoskope/pkg/domain/constants/aggregates"
	commandTypes "github.com/finleap-connect/monoskope/pkg/domain/constants/commands"
	"github.com/finleap-connect/monoskope/pkg/domain/constants/events"
	"github.com/finleap-connect/monoskope/pkg/domain/constants/roles"
	"github.com/finleap-connect/monoskope/pkg/domain/constants/scopes"
	"github.com/finleap-connect/monoskope/pkg/domain/errors"
	es "github.com/finleap-connect/monoskope/pkg/eventsourcing"
	grpcUtil "github.com/finleap-connect/monoskope/pkg/grpc"
	"github.com/finleap-connect/monoskope/pkg/jwt"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

var _ = Describe("integration", func() {
	ctx := context.Background()

	expectedClusterDisplayName := "the one cluster"
	expectedClusterName := "one-cluster"
	expectedClusterApiServerAddress := "one.example.com"
	expectedClusterCACertBundle := []byte("This should be a certificate")

	expectedTenantName := "tenantx"

	getAdminAuthToken := func() string {
		signer := testEnv.gatewayTestEnv.JwtTestEnv.CreateSigner()
		token := auth.NewAuthToken(&jwt.StandardClaims{Name: testEnv.gatewayTestEnv.AdminUser.Name, Email: testEnv.gatewayTestEnv.AdminUser.Email}, testEnv.gatewayTestEnv.GetApiAddr(), testEnv.gatewayTestEnv.AdminUser.ID().String(), time.Minute*10)
		authToken, err := signer.GenerateSignedToken(token)
		Expect(err).ToNot(HaveOccurred())
		return authToken
	}

	commandHandlerClient := func() esApi.CommandHandlerClient {
		_, chClient, err := grpcUtil.NewClientWithInsecureAuth(ctx, testEnv.commandHandlerTestEnv.GetApiAddr(), getAdminAuthToken(), esApi.NewCommandHandlerClient)
		Expect(err).ToNot(HaveOccurred())
		return chClient
	}

	userServiceClient := func() domainApi.UserClient {
		_, client, err := grpcUtil.NewClientWithInsecureAuth(ctx, testEnv.queryHandlerTestEnv.GetApiAddr(), getAdminAuthToken(), domainApi.NewUserClient)
		Expect(err).ToNot(HaveOccurred())
		return client
	}

	tenantServiceClient := func() domainApi.TenantClient {
		_, client, err := grpcUtil.NewClientWithInsecureAuth(ctx, testEnv.queryHandlerTestEnv.GetApiAddr(), getAdminAuthToken(), domainApi.NewTenantClient)
		Expect(err).ToNot(HaveOccurred())
		return client
	}

	clusterServiceClient := func() domainApi.ClusterClient {
		_, client, err := grpcUtil.NewClientWithInsecureAuth(ctx, testEnv.queryHandlerTestEnv.GetApiAddr(), getAdminAuthToken(), domainApi.NewClusterClient)
		Expect(err).ToNot(HaveOccurred())
		return client
	}

	certificateServiceClient := func() domainApi.CertificateClient {
		_, client, err := grpcUtil.NewClientWithInsecureAuth(ctx, testEnv.queryHandlerTestEnv.GetApiAddr(), getAdminAuthToken(), domainApi.NewCertificateClient)
		Expect(err).ToNot(HaveOccurred())
		return client
	}

	eventStoreClient := func() esApi.EventStoreClient {
		_, client, err := grpcUtil.NewClientWithInsecure(ctx, testEnv.eventStoreTestEnv.GetApiAddr(), esApi.NewEventStoreClient)
		Expect(err).ToNot(HaveOccurred())
		return client
	}

	Context("user management", func() {

		It("can manage a user", func() {
			command, err := cmd.AddCommandData(
				cmd.CreateCommand(uuid.Nil, commandTypes.CreateUser),
				&cmdData.CreateUserCommandData{Name: "Jane Doe", Email: "jane.doe@monoskope.io"},
			)
			Expect(err).ToNot(HaveOccurred())

			// handel admin user propagation
			var reply *esApi.CommandReply
			Eventually(func(g Gomega) {
				reply, err = commandHandlerClient().Execute(ctx, command)
				g.Expect(err).ToNot(HaveOccurred())
				g.Expect(uuid.Nil).ToNot(Equal(reply.AggregateId))
			}).Should(Succeed())

			// update userId, as the "create" command will have changed it.
			userId := uuid.MustParse(reply.AggregateId)

			var user *projections.User
			Eventually(func(g Gomega) {
				user, err = userServiceClient().GetByEmail(ctx, wrapperspb.String("jane.doe@monoskope.io"))
				g.Expect(err).ToNot(HaveOccurred())
				g.Expect(user).ToNot(BeNil())
				g.Expect(user.GetEmail()).To(Equal("jane.doe@monoskope.io"))
				g.Expect(user.Id).To(Equal(userId.String()))
			}).Should(Succeed())

			command, err = cmd.AddCommandData(
				cmd.CreateCommand(uuid.Nil, commandTypes.CreateUserRoleBinding),
				&cmdData.CreateUserRoleBindingCommandData{Role: roles.Admin.String(), Scope: scopes.System.String(), UserId: userId.String()},
			)
			Expect(err).ToNot(HaveOccurred())

			reply, err = commandHandlerClient().Execute(ctx, command)
			Expect(err).ToNot(HaveOccurred())

			// update userRolebBindingId, as the "create" command will have changed it.
			userRoleBindingId := uuid.MustParse(reply.AggregateId)

			// Wait to propagate
			time.Sleep(500 * time.Millisecond)

			// Creating the same rolebinding again should fail
			Eventually(func(g Gomega) {
				command.Id = uuid.New().String()
				_, err = commandHandlerClient().Execute(ctx, command)
				g.Expect(err).To(HaveOccurred())
			}).Should(Succeed())

			user, err = userServiceClient().GetByEmail(ctx, wrapperspb.String("jane.doe@monoskope.io"))
			Expect(err).ToNot(HaveOccurred())
			Expect(user).ToNot(BeNil())
			Expect(user.Roles[0].Role).To(Equal(roles.Admin.String()))
			Expect(user.Roles[0].Scope).To(Equal(scopes.System.String()))

			_, err = commandHandlerClient().Execute(ctx, cmd.CreateCommand(userRoleBindingId, commandTypes.DeleteUserRoleBinding))
			Expect(err).ToNot(HaveOccurred())

			Eventually(func(g Gomega) {
				user, err = userServiceClient().GetByEmail(ctx, wrapperspb.String("jane.doe@monoskope.io"))
				g.Expect(err).ToNot(HaveOccurred())
				g.Expect(user).ToNot(BeNil())
			}).Should(Succeed())
		})
		It("can accept Nil as an Id when creating a user", func() {
			command, err := cmd.AddCommandData(
				cmd.CreateCommand(uuid.Nil, commandTypes.CreateUser),
				&cmdData.CreateUserCommandData{Name: "Jane Doe", Email: "jane.doe2@monoskope.io"},
			)
			Expect(err).ToNot(HaveOccurred())

			reply, err := commandHandlerClient().Execute(ctx, command)
			Expect(err).ToNot(HaveOccurred())
			Expect(reply.AggregateId).ToNot(Equal(uuid.Nil.String()))
		})
		It("fail to create a user which already exists", func() {
			command, err := cmd.AddCommandData(
				cmd.CreateCommand(uuid.New(), commandTypes.CreateUser),
				&cmdData.CreateUserCommandData{Name: "admin", Email: "admin@monoskope.io"},
			)
			Expect(err).ToNot(HaveOccurred())

			_, err = commandHandlerClient().Execute(ctx, command)
			Expect(err).To(HaveOccurred())
			Expect(errors.TranslateFromGrpcError(err)).To(Equal(errors.ErrUserAlreadyExists))
		})
		It("can delete a user", func() {
			createCommand, err := cmd.AddCommandData(
				cmd.CreateCommand(uuid.Nil, commandTypes.CreateUser),
				&cmdData.CreateUserCommandData{Name: "John Doe", Email: "john.doe@monoskope.io"},
			)
			Expect(err).ToNot(HaveOccurred())

			reply, err := commandHandlerClient().Execute(ctx, createCommand)
			Expect(err).ToNot(HaveOccurred())
			Expect(uuid.Nil).ToNot(Equal(reply.AggregateId))

			// update userId, as the "create" command will have changed it.
			userId := uuid.MustParse(reply.AggregateId)

			_, err = commandHandlerClient().Execute(ctx,
				cmd.CreateCommand(userId, commandTypes.DeleteUser))
			Expect(err).ToNot(HaveOccurred())

			var user *projections.User
			Eventually(func(g Gomega) {
				user, err = userServiceClient().GetByEmail(ctx, wrapperspb.String("john.doe@monoskope.io"))
				g.Expect(err).ToNot(HaveOccurred())
				g.Expect(user).ToNot(BeNil())
				g.Expect(user.GetEmail()).To(Equal("john.doe@monoskope.io"))
				g.Expect(user.Id).To(Equal(userId.String()))
				g.Expect(user.GetMetadata().GetDeleted()).ToNot(BeNil())
			}).Should(Succeed())
		})
	})
	Context("tenant management", func() {
		It("can manage a tenant", func() {
			tenantId := uuid.New()
			command, err := cmd.AddCommandData(
				cmd.CreateCommand(tenantId, commandTypes.CreateTenant),
				&cmdData.CreateTenantCommandData{Name: "Tenant X", Prefix: "tx"},
			)
			Expect(err).ToNot(HaveOccurred())

			reply, err := commandHandlerClient().Execute(ctx, command)
			Expect(err).ToNot(HaveOccurred())
			Expect(tenantId.String()).ToNot(Equal(reply.AggregateId))

			// update tenantId, as the "create" command will have changed it.
			tenantId = uuid.MustParse(reply.AggregateId)

			var tenant *projections.Tenant
			Eventually(func(g Gomega) {
				tenant, err = tenantServiceClient().GetByName(ctx, wrapperspb.String("Tenant X"))
				g.Expect(err).ToNot(HaveOccurred())
				g.Expect(tenant).ToNot(BeNil())
				g.Expect(tenant.GetName()).To(Equal("Tenant X"))
				g.Expect(tenant.GetPrefix()).To(Equal("tx"))
				g.Expect(tenant.Id).To(Equal(tenantId.String()))
			}).Should(Succeed())

			command, err = cmd.AddCommandData(
				cmd.CreateCommand(tenantId, commandTypes.UpdateTenant),
				&cmdData.UpdateTenantCommandData{Name: &wrapperspb.StringValue{Value: expectedTenantName}},
			)
			Expect(err).ToNot(HaveOccurred())

			_, err = commandHandlerClient().Execute(ctx, command)
			Expect(err).ToNot(HaveOccurred())

			Eventually(func(g Gomega) {
				tenant, err = tenantServiceClient().GetByName(ctx, wrapperspb.String(expectedTenantName))
				g.Expect(err).ToNot(HaveOccurred())
				g.Expect(tenant).ToNot(BeNil())
			}).Should(Succeed())

			_, err = commandHandlerClient().Execute(ctx, cmd.CreateCommand(tenantId, commandTypes.DeleteTenant))
			Expect(err).ToNot(HaveOccurred())

			Eventually(func(g Gomega) {
				tenant, err = tenantServiceClient().GetByName(ctx, wrapperspb.String(expectedTenantName))
				g.Expect(err).ToNot(HaveOccurred())
				g.Expect(tenant).ToNot(BeNil())
				g.Expect(tenant.Metadata.Created).NotTo(BeNil())
			}).Should(Succeed())
		})
		It("can accept Nil as ID when creating a tenant", func() {
			command, err := cmd.AddCommandData(
				cmd.CreateCommand(uuid.Nil, commandTypes.CreateTenant),
				&cmdData.CreateTenantCommandData{Name: "Tenant X", Prefix: "tx"},
			)
			Expect(err).ToNot(HaveOccurred())

			reply, err := commandHandlerClient().Execute(ctx, command)
			Expect(err).ToNot(HaveOccurred())
			Expect(reply.AggregateId).ToNot(Equal(uuid.Nil.String()))
			Expect(int(reply.Version)).To(BeNumerically("==", 1))
		})
	})

	Context("cluster management", func() {
		It("manage a cluster", func() {
			command, err := cmd.AddCommandData(
				cmd.CreateCommand(uuid.Nil, commandTypes.CreateCluster),
				&cmdData.CreateCluster{DisplayName: expectedClusterDisplayName, Name: expectedClusterName, ApiServerAddress: expectedClusterApiServerAddress, CaCertBundle: expectedClusterCACertBundle},
			)
			Expect(err).ToNot(HaveOccurred())

			// set up reactor for checking JWTs later
			testReactor := mock_reactor.NewTestReactor()
			defer testReactor.Close()

			err = testReactor.Setup(ctx, testEnv.eventStoreTestEnv, eventStoreClient())
			Expect(err).ToNot(HaveOccurred())

			reply, err := commandHandlerClient().Execute(ctx, command)
			Expect(err).ToNot(HaveOccurred())
			Expect(uuid.Nil).ToNot(Equal(reply.AggregateId))

			// update clusterId, as the "create" command will have changed it.
			clusterId := uuid.MustParse(reply.AggregateId)

			var cluster *projections.Cluster
			Eventually(func(g Gomega) {
				cluster, err = clusterServiceClient().GetByName(ctx, wrapperspb.String(expectedClusterName))
				g.Expect(err).ToNot(HaveOccurred())
				g.Expect(cluster).ToNot(BeNil())
				g.Expect(cluster.GetDisplayName()).To(Equal(expectedClusterDisplayName))
				g.Expect(cluster.GetName()).To(Equal(expectedClusterName))
				g.Expect(cluster.GetApiServerAddress()).To(Equal(expectedClusterApiServerAddress))
				g.Expect(cluster.GetCaCertBundle()).To(Equal(expectedClusterCACertBundle))
			}).Should(Succeed())

			By("getting all existing clusters")

			clusterStream, err := clusterServiceClient().GetAll(ctx, &domainApi.GetAllRequest{
				IncludeDeleted: true,
			})
			Expect(err).ToNot(HaveOccurred())
			// Read next
			firstCluster, err := clusterStream.Recv()
			Expect(err).ToNot(HaveOccurred())

			Expect(firstCluster).ToNot(BeNil())
			Expect(firstCluster.GetDisplayName()).To(Equal(expectedClusterDisplayName))
			Expect(firstCluster.GetName()).To(Equal(expectedClusterName))
			Expect(firstCluster.GetApiServerAddress()).To(Equal(expectedClusterApiServerAddress))
			Expect(firstCluster.GetCaCertBundle()).To(Equal(expectedClusterCACertBundle))

			By("by retrieving the bootstrap token")
			observed := testReactor.GetObservedEvents()
			Expect(len(observed)).ToNot(Equal(0))
			Expect(observed[0].AggregateID()).To(Equal(clusterId))

			eventMD := observed[0].Metadata()
			event := es.NewEventWithMetadata(events.ClusterBootstrapTokenCreated,
				es.ToEventDataFromProto(&eventdata.ClusterBootstrapTokenCreated{
					Jwt: "this is a valid JWT, honest!",
				}), time.Now().UTC(),
				observed[0].AggregateType(), observed[0].AggregateID(),
				observed[0].AggregateVersion()+1,
				eventMD)

			err = testReactor.Emit(ctx, event)
			Expect(err).ToNot(HaveOccurred())

			Eventually(func(g Gomega) {
				tokenValue, err := clusterServiceClient().GetBootstrapToken(ctx, wrapperspb.String(clusterId.String()))
				g.Expect(err).ToNot(HaveOccurred())
				g.Expect(tokenValue.GetValue()).To(Equal("this is a valid JWT, honest!"))
			}).Should(Succeed())
		})
	})

	Context("cert management", func() {
		It("can create and query a certificate", func() {
			testReactor := mock_reactor.NewTestReactor()
			defer testReactor.Close()

			err := testReactor.Setup(ctx, testEnv.eventStoreTestEnv, eventStoreClient())
			Expect(err).ToNot(HaveOccurred())

			clusterInfo, err := clusterServiceClient().GetByName(ctx, &wrapperspb.StringValue{Value: expectedClusterName})
			Expect(err).ToNot(HaveOccurred())
			Expect(clusterInfo).ToNot(BeNil())

			command, err := cmd.AddCommandData(
				cmd.CreateCommand(uuid.Nil, commandTypes.RequestCertificate),
				&cmdData.RequestCertificate{
					ReferencedAggregateId:   clusterInfo.Id,
					ReferencedAggregateType: aggregates.Cluster.String(),
					SigningRequest:          []byte("-----BEGIN CERTIFICATE REQUEST-----this is a CSR-----END CERTIFICATE REQUEST-----"),
				},
			)
			Expect(err).ToNot(HaveOccurred())

			certRequestCmdReply, err := commandHandlerClient().Execute(ctx, command)
			Expect(err).ToNot(HaveOccurred())

			var observed []es.Event
			Eventually(func(g Gomega) {
				observed = testReactor.GetObservedEvents()
				g.Expect(len(observed)).To(Equal(1))
			}).Should(Succeed())
			certRequestedEvent := observed[0]
			Expect(certRequestedEvent.AggregateID().String()).To(Equal(certRequestCmdReply.AggregateId))

			err = testReactor.Emit(ctx, es.NewEvent(
				ctx,
				events.CertificateIssued,
				es.ToEventDataFromProto(&eventdata.CertificateIssued{
					Certificate: &common.CertificateChain{
						Ca:          expectedClusterCACertBundle,
						Certificate: []byte("this is a cert"),
					},
				}),
				time.Now().UTC(),
				certRequestedEvent.AggregateType(),
				certRequestedEvent.AggregateID(),
				certRequestedEvent.AggregateVersion()+1),
			)
			Expect(err).ToNot(HaveOccurred())

			Eventually(func(g Gomega) {
				certificate, err := certificateServiceClient().GetCertificate(ctx,
					&domainApi.GetCertificateRequest{
						AggregateId:   clusterInfo.GetId(),
						AggregateType: aggregates.Cluster.String(),
					})
				g.Expect(err).ToNot(HaveOccurred())
				g.Expect(certificate.GetCertificate()).ToNot(BeNil())
				g.Expect(certificate.GetCaCertBundle()).ToNot(BeNil())
			}).Should(Succeed())
		})
	})
})

var _ = Describe("PrometheusMetrics", func() {
	It("can scrape event store metrics", func() {
		res, err := http.Get(fmt.Sprintf("http://%s/metrics", testEnv.eventStoreTestEnv.MetricsListener.Addr()))
		Expect(err).ToNot(HaveOccurred())
		defer res.Body.Close()
		Expect(res.StatusCode).To(Equal(200))
	})
})
