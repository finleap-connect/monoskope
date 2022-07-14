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

package integration

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/finleap-connect/monoskope/internal/gateway/auth"
	mock_reactor "github.com/finleap-connect/monoskope/internal/test/reactor"
	domainApi "github.com/finleap-connect/monoskope/pkg/api/domain"
	cmdData "github.com/finleap-connect/monoskope/pkg/api/domain/commanddata"
	"github.com/finleap-connect/monoskope/pkg/api/domain/projections"
	esApi "github.com/finleap-connect/monoskope/pkg/api/eventsourcing"
	cmd "github.com/finleap-connect/monoskope/pkg/domain/commands"
	commandTypes "github.com/finleap-connect/monoskope/pkg/domain/constants/commands"
	"github.com/finleap-connect/monoskope/pkg/domain/constants/roles"
	"github.com/finleap-connect/monoskope/pkg/domain/constants/scopes"
	"github.com/finleap-connect/monoskope/pkg/domain/errors"
	"github.com/finleap-connect/monoskope/pkg/domain/mock"
	grpcUtil "github.com/finleap-connect/monoskope/pkg/grpc"
	"github.com/finleap-connect/monoskope/pkg/jwt"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

var _ = Describe("internal/integration_test", func() {
	ctx := context.Background()

	expectedUserName := "Jane Doe"
	expectedUserEmail := "jane.doe@monoskope.io"

	expectedTenantName := "Tenant X"
	expectedTenantNameUpdated := "tenantx"
	expectedTenantPrefix := "tx"

	expectedClusterDisplayName := "the one cluster"
	expectedClusterName := "one-cluster"
	expectedClusterApiServerAddress := "one.example.com"
	expectedClusterCACertBundle := []byte("This should be a certificate")

	getAdminAuthToken := func() string {
		signer := testEnv.GatewayTestEnv.JwtTestEnv.CreateSigner()
		token := auth.NewAuthToken(&jwt.StandardClaims{Name: mock.TestAdminUser.Name, Email: mock.TestAdminUser.Email}, testEnv.GatewayTestEnv.GetApiAddr(), mock.TestAdminUser.Id, time.Minute*10)
		authToken, err := signer.GenerateSignedToken(token)
		Expect(err).ToNot(HaveOccurred())
		return authToken
	}

	commandHandlerClient := func() esApi.CommandHandlerClient {
		_, chClient, err := grpcUtil.NewClientWithInsecureAuth(ctx, testEnv.CommandHandlerTestEnv.GetApiAddr(), getAdminAuthToken(), esApi.NewCommandHandlerClient)
		Expect(err).ToNot(HaveOccurred())
		return chClient
	}

	userServiceClient := func() domainApi.UserClient {
		_, client, err := grpcUtil.NewClientWithInsecureAuth(ctx, testEnv.QueryHandlerTestEnv.GetApiAddr(), getAdminAuthToken(), domainApi.NewUserClient)
		Expect(err).ToNot(HaveOccurred())
		return client
	}

	tenantServiceClient := func() domainApi.TenantClient {
		_, client, err := grpcUtil.NewClientWithInsecureAuth(ctx, testEnv.QueryHandlerTestEnv.GetApiAddr(), getAdminAuthToken(), domainApi.NewTenantClient)
		Expect(err).ToNot(HaveOccurred())
		return client
	}

	clusterServiceClient := func() domainApi.ClusterClient {
		_, client, err := grpcUtil.NewClientWithInsecureAuth(ctx, testEnv.QueryHandlerTestEnv.GetApiAddr(), getAdminAuthToken(), domainApi.NewClusterClient)
		Expect(err).ToNot(HaveOccurred())
		return client
	}

	clusterAccessClient := func() domainApi.ClusterAccessClient {
		_, client, err := grpcUtil.NewClientWithInsecureAuth(ctx, testEnv.QueryHandlerTestEnv.GetApiAddr(), getAdminAuthToken(), domainApi.NewClusterAccessClient)
		Expect(err).ToNot(HaveOccurred())
		return client
	}

	eventStoreClient := func() esApi.EventStoreClient {
		_, client, err := grpcUtil.NewClientWithInsecure(ctx, testEnv.EventStoreTestEnv.GetApiAddr(), esApi.NewEventStoreClient)
		Expect(err).ToNot(HaveOccurred())
		return client
	}

	Context("user management", func() {
		It("can manage a user", func() {
			By("creating the user")
			command, err := cmd.AddCommandData(
				cmd.CreateCommand(uuid.Nil, commandTypes.CreateUser),
				&cmdData.CreateUserCommandData{Name: expectedUserName, Email: expectedUserEmail},
			)
			Expect(err).ToNot(HaveOccurred())

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
				user, err = userServiceClient().GetByEmail(ctx, wrapperspb.String(expectedUserEmail))
				g.Expect(err).ToNot(HaveOccurred())
				g.Expect(user).ToNot(BeNil())
				g.Expect(user.GetEmail()).To(Equal(expectedUserEmail))
				g.Expect(user.Id).To(Equal(userId.String()))
			}).Should(Succeed())

			By("ensuring the same user can't be created again")
			command, err = cmd.AddCommandData(
				cmd.CreateCommand(uuid.Nil, commandTypes.CreateUser),
				&cmdData.CreateUserCommandData{Name: expectedUserName,
					Email: strings.ToUpper(expectedUserEmail)}, // regardless of the case
			)
			_, err = commandHandlerClient().Execute(ctx, command)
			Expect(err).To(HaveOccurred())

			By("giving the user system admin role")
			command, err = cmd.AddCommandData(
				cmd.CreateCommand(uuid.Nil, commandTypes.CreateUserRoleBinding),
				&cmdData.CreateUserRoleBindingCommandData{Role: string(roles.Admin), Scope: string(scopes.System), UserId: userId.String()},
			)
			Expect(err).ToNot(HaveOccurred())

			reply, err = commandHandlerClient().Execute(ctx, command)
			Expect(err).ToNot(HaveOccurred())

			// update userRoleBindingId, as the "create" command will have changed it.
			userRoleBindingId := uuid.MustParse(reply.AggregateId)

			Eventually(func(g Gomega) {
				user, err = userServiceClient().GetByEmail(ctx, wrapperspb.String(expectedUserEmail))
				g.Expect(err).ToNot(HaveOccurred())
				g.Expect(user).ToNot(BeNil())
				g.Expect(len(user.Roles)).To(BeNumerically(">=", 1))
				g.Expect(user.Roles[0].Role).To(Equal(string(roles.Admin)))
				g.Expect(user.Roles[0].Scope).To(Equal(string(scopes.System)))
			}).Should(Succeed())

			By("ensuring the same role (system admin) can't be given again")
			command.Id = uuid.New().String()
			_, err = commandHandlerClient().Execute(ctx, command)
			Expect(err).To(HaveOccurred())

			By("removing/deleting the users system admin role")
			_, err = commandHandlerClient().Execute(ctx, cmd.CreateCommand(userRoleBindingId, commandTypes.DeleteUserRoleBinding))
			Expect(err).ToNot(HaveOccurred())

			Eventually(func(g Gomega) {
				user, err = userServiceClient().GetByEmail(ctx, wrapperspb.String(expectedUserEmail))
				g.Expect(err).ToNot(HaveOccurred())
				g.Expect(user).ToNot(BeNil())
				g.Expect(user.Roles).To(BeEmpty())
			}).Should(Succeed())

			By("giving the user the same role (system admin) again (after deleting the old one)")
			Eventually(func(g Gomega) {
				command.Id = uuid.New().String()
				_, err = commandHandlerClient().Execute(ctx, command)
				g.Expect(err).ToNot(HaveOccurred())
			}).Should(Succeed())

			Eventually(func(g Gomega) {
				user, err = userServiceClient().GetByEmail(ctx, wrapperspb.String(expectedUserEmail))
				g.Expect(err).ToNot(HaveOccurred())
				g.Expect(user).ToNot(BeNil())
				g.Expect(len(user.Roles)).To(BeNumerically(">=", 1))
				g.Expect(user.Roles[0].Role).To(Equal(string(roles.Admin)))
				g.Expect(user.Roles[0].Scope).To(Equal(string(scopes.System)))
			}).Should(Succeed())

			By("deleting the user")
			_, err = commandHandlerClient().Execute(ctx,
				cmd.CreateCommand(userId, commandTypes.DeleteUser))
			Expect(err).ToNot(HaveOccurred())

			Eventually(func(g Gomega) {
				user, err = userServiceClient().GetByEmail(ctx, wrapperspb.String(expectedUserEmail))
				g.Expect(err).To(HaveOccurred())
				g.Expect(user).To(BeNil())
			}).Should(Succeed())

			By("recreating the user after deletion")
			command, err = cmd.AddCommandData(
				cmd.CreateCommand(uuid.New(), commandTypes.CreateUser),
				&cmdData.CreateUserCommandData{Name: expectedUserName, Email: expectedUserEmail},
			)
			Expect(err).ToNot(HaveOccurred())

			reply, err = commandHandlerClient().Execute(ctx, command)
			Expect(err).ToNot(HaveOccurred())
			userIdNew := uuid.MustParse(reply.AggregateId)

			Eventually(func(g Gomega) {
				user, err = userServiceClient().GetByEmail(ctx, wrapperspb.String(expectedUserEmail))
				g.Expect(err).ToNot(HaveOccurred())
				g.Expect(user).ToNot(BeNil())
				g.Expect(user.GetEmail()).To(Equal(expectedUserEmail))
				g.Expect(user.Id).To(Equal(userIdNew.String()))
				g.Expect(user.Id).ToNot(Equal(userId.String()))
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
				&cmdData.CreateUserCommandData{Name: mock.TestAdminUser.Name, Email: mock.TestAdminUser.Email},
			)
			Expect(err).ToNot(HaveOccurred())

			_, err = commandHandlerClient().Execute(ctx, command)
			Expect(err).To(HaveOccurred())
			Expect(errors.TranslateFromGrpcError(err)).To(Equal(errors.ErrUserAlreadyExists))
		})
	})
	Context("tenant management", func() {
		It("can manage a tenant", func() {
			By("creating the tenant")
			tenantId := uuid.New()
			command, err := cmd.AddCommandData(
				cmd.CreateCommand(tenantId, commandTypes.CreateTenant),
				&cmdData.CreateTenantCommandData{Name: expectedTenantName, Prefix: expectedTenantPrefix},
			)
			Expect(err).ToNot(HaveOccurred())

			reply, err := commandHandlerClient().Execute(ctx, command)
			Expect(err).ToNot(HaveOccurred())
			Expect(tenantId.String()).ToNot(Equal(reply.AggregateId))

			// update tenantId, as the "create" command will have changed it.
			tenantId = uuid.MustParse(reply.AggregateId)

			var tenant *projections.Tenant
			Eventually(func(g Gomega) {
				tenant, err = tenantServiceClient().GetByName(ctx, wrapperspb.String(expectedTenantName))
				g.Expect(err).ToNot(HaveOccurred())
				g.Expect(tenant).ToNot(BeNil())
				g.Expect(tenant.GetName()).To(Equal(expectedTenantName))
				g.Expect(tenant.GetPrefix()).To(Equal(expectedTenantPrefix))
				g.Expect(tenant.Id).To(Equal(tenantId.String()))
			}).Should(Succeed())

			By("ensuring the same tenant can't be created again")
			command, err = cmd.AddCommandData(
				cmd.CreateCommand(tenantId, commandTypes.CreateTenant),
				&cmdData.CreateTenantCommandData{Name: strings.ToUpper(expectedTenantName), // regardless of the case
					Prefix: expectedTenantPrefix},
			)
			_, err = commandHandlerClient().Execute(ctx, command)
			Expect(err).To(HaveOccurred())

			By("updating the tenant")
			command, err = cmd.AddCommandData(
				cmd.CreateCommand(tenantId, commandTypes.UpdateTenant),
				&cmdData.UpdateTenantCommandData{Name: &wrapperspb.StringValue{Value: expectedTenantNameUpdated}},
			)
			Expect(err).ToNot(HaveOccurred())

			_, err = commandHandlerClient().Execute(ctx, command)
			Expect(err).ToNot(HaveOccurred())

			Eventually(func(g Gomega) {
				tenant, err = tenantServiceClient().GetByName(ctx, wrapperspb.String(expectedTenantNameUpdated))
				g.Expect(err).ToNot(HaveOccurred())
				g.Expect(tenant).ToNot(BeNil())
			}).Should(Succeed())

			By("deleting the tenant")
			_, err = commandHandlerClient().Execute(ctx, cmd.CreateCommand(tenantId, commandTypes.DeleteTenant))
			Expect(err).ToNot(HaveOccurred())

			Eventually(func(g Gomega) {
				tenant, err = tenantServiceClient().GetByName(ctx, wrapperspb.String(expectedTenantNameUpdated))
				g.Expect(err).ToNot(HaveOccurred())
				g.Expect(tenant).ToNot(BeNil())
				g.Expect(tenant.Metadata.Deleted).NotTo(BeNil())
			}).Should(Succeed())

			By("recreating the tenant after deletion")
			command, err = cmd.AddCommandData(
				cmd.CreateCommand(uuid.New(), commandTypes.CreateTenant),
				&cmdData.CreateTenantCommandData{Name: expectedTenantNameUpdated, Prefix: expectedTenantPrefix},
			)
			Expect(err).ToNot(HaveOccurred())

			reply, err = commandHandlerClient().Execute(ctx, command)
			Expect(err).ToNot(HaveOccurred())

			tenantIdNew := uuid.MustParse(reply.AggregateId)

			Eventually(func(g Gomega) {
				tenant, err = tenantServiceClient().GetByName(ctx, wrapperspb.String(expectedTenantNameUpdated))
				g.Expect(err).ToNot(HaveOccurred())
				g.Expect(tenant).ToNot(BeNil())
				g.Expect(tenant.Id).ToNot(Equal(tenantId.String()))
				g.Expect(tenant.Id).To(Equal(tenantIdNew.String()))
				g.Expect(tenant.Metadata.Deleted).To(BeNil())
			}, "5s").Should(Succeed())
		})
		It("can accept Nil as ID when creating a tenant", func() {
			command, err := cmd.AddCommandData(
				cmd.CreateCommand(uuid.Nil, commandTypes.CreateTenant),
				&cmdData.CreateTenantCommandData{Name: "Tenant K", Prefix: "tk"},
			)
			Expect(err).ToNot(HaveOccurred())

			reply, err := commandHandlerClient().Execute(ctx, command)
			Expect(err).ToNot(HaveOccurred())
			Expect(reply.AggregateId).ToNot(Equal(uuid.Nil.String()))
			Expect(int(reply.Version)).To(BeNumerically("==", 1))
		})
	})
	Context("cluster management", func() {
		It("can manage a cluster", func() {
			By("creating the cluster")
			command, err := cmd.AddCommandData(
				cmd.CreateCommand(uuid.Nil, commandTypes.CreateCluster),
				&cmdData.CreateCluster{DisplayName: expectedClusterDisplayName, Name: expectedClusterName, ApiServerAddress: expectedClusterApiServerAddress, CaCertBundle: expectedClusterCACertBundle},
			)
			Expect(err).ToNot(HaveOccurred())

			// set up reactor for checking JWTs later
			testReactor := mock_reactor.NewTestReactor()
			defer testReactor.Close()

			err = testReactor.Setup(ctx, testEnv.EventStoreTestEnv, eventStoreClient())
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
				g.Expect(cluster.Id).To(Equal(clusterId.String()))
				g.Expect(cluster.GetDisplayName()).To(Equal(expectedClusterDisplayName))
				g.Expect(cluster.GetName()).To(Equal(expectedClusterName))
				g.Expect(cluster.GetApiServerAddress()).To(Equal(expectedClusterApiServerAddress))
				g.Expect(cluster.GetCaCertBundle()).To(Equal(expectedClusterCACertBundle))
			}).Should(Succeed())

			By("ensuring the same cluster can't be created again")
			command, err = cmd.AddCommandData(
				cmd.CreateCommand(uuid.Nil, commandTypes.CreateCluster),
				&cmdData.CreateCluster{DisplayName: expectedClusterDisplayName,
					Name:             strings.ToUpper(expectedClusterName), // regardless of the case and white spaces
					ApiServerAddress: expectedClusterApiServerAddress, CaCertBundle: expectedClusterCACertBundle},
			)
			_, err = commandHandlerClient().Execute(ctx, command)
			Expect(err).To(HaveOccurred())

			By("deleting the cluster")
			_, err = commandHandlerClient().Execute(ctx, cmd.CreateCommand(clusterId, commandTypes.DeleteCluster))
			Expect(err).ToNot(HaveOccurred())

			Eventually(func(g Gomega) {
				cluster, err = clusterServiceClient().GetByName(ctx, wrapperspb.String(expectedClusterName))
				g.Expect(err).ToNot(HaveOccurred())
				g.Expect(cluster).ToNot(BeNil())
				g.Expect(cluster.Metadata.Deleted).NotTo(BeNil())
			}).Should(Succeed())

			By("recreating the cluster after deletion")
			command, err = cmd.AddCommandData(
				cmd.CreateCommand(uuid.Nil, commandTypes.CreateCluster),
				&cmdData.CreateCluster{DisplayName: expectedClusterDisplayName, Name: expectedClusterName, ApiServerAddress: expectedClusterApiServerAddress, CaCertBundle: expectedClusterCACertBundle},
			)
			Expect(err).ToNot(HaveOccurred())

			reply, err = commandHandlerClient().Execute(ctx, command)
			Expect(err).ToNot(HaveOccurred())
			Expect(uuid.Nil).ToNot(Equal(reply.AggregateId))

			// update clusterId, as the "create" command will have changed it.
			clusterIdNew := uuid.MustParse(reply.AggregateId)

			Eventually(func(g Gomega) {
				cluster, err = clusterServiceClient().GetByName(ctx, wrapperspb.String(expectedClusterName))
				g.Expect(err).ToNot(HaveOccurred())
				g.Expect(cluster).ToNot(BeNil())
				g.Expect(cluster.Id).ToNot(Equal(clusterId.String()))
				g.Expect(cluster.Id).To(Equal(clusterIdNew.String()))
			}).Should(Succeed())
		})
		It("can grant a tenant access to a cluster", func() {
			// create the tenant
			command, err := cmd.AddCommandData(
				cmd.CreateCommand(uuid.Nil, commandTypes.CreateTenant),
				&cmdData.CreateTenantCommandData{Name: "Tenant Z", Prefix: "tz"},
			)
			Expect(err).ToNot(HaveOccurred())

			reply, err := commandHandlerClient().Execute(ctx, command)
			Expect(err).ToNot(HaveOccurred())

			tenantId := uuid.MustParse(reply.AggregateId)

			var tenant *projections.Tenant
			Eventually(func(g Gomega) {
				tenant, err = tenantServiceClient().GetByName(ctx, wrapperspb.String("Tenant Z"))
				g.Expect(err).ToNot(HaveOccurred())
				g.Expect(tenant).ToNot(BeNil())
				g.Expect(tenant.Id).To(Equal(tenantId.String()))
			}).Should(Succeed())

			// create the cluster
			command, err = cmd.AddCommandData(
				cmd.CreateCommand(uuid.Nil, commandTypes.CreateCluster),
				&cmdData.CreateCluster{DisplayName: "Cluster Z", Name: "cluster-z", ApiServerAddress: "z.cluster.com", CaCertBundle: []byte("cluster z certificate")},
			)
			Expect(err).ToNot(HaveOccurred())

			reply, err = commandHandlerClient().Execute(ctx, command)
			Expect(err).ToNot(HaveOccurred())

			clusterId := uuid.MustParse(reply.AggregateId)

			var cluster *projections.Cluster
			Eventually(func(g Gomega) {
				cluster, err = clusterServiceClient().GetByName(ctx, wrapperspb.String("cluster-z"))
				g.Expect(err).ToNot(HaveOccurred())
				g.Expect(cluster).ToNot(BeNil())
				g.Expect(cluster.Id).To(Equal(clusterId.String()))
			}).Should(Succeed())

			By("granting the tenant access to the cluster")
			command, err = cmd.AddCommandData(
				cmd.CreateCommand(uuid.Nil, commandTypes.CreateTenantClusterBinding),
				&cmdData.CreateTenantClusterBindingCommandData{TenantId: tenantId.String(), ClusterId: clusterId.String()},
			)
			Expect(err).ToNot(HaveOccurred())
			reply, err = commandHandlerClient().Execute(ctx, command)
			Expect(err).ToNot(HaveOccurred())

			tenantClusterBindingId := uuid.MustParse(reply.AggregateId)

			Eventually(func(g Gomega) {
				tenantClusterBinding, err := clusterAccessClient().GetTenantClusterMappingByTenantAndClusterId(ctx, &domainApi.GetClusterMappingRequest{ClusterId: clusterId.String(), TenantId: tenantId.String()})
				Expect(err).ToNot(HaveOccurred())
				Expect(tenantClusterBinding).ToNot(BeNil())
				Expect(tenantClusterBinding.Id).To(Equal(tenantClusterBindingId.String()))
			}, "5s").Should(Succeed())

			By("ensuring the same access can't be granted again")
			command.Id = uuid.New().String()
			_, err = commandHandlerClient().Execute(ctx, command)
			Expect(err).To(HaveOccurred())

			By("revoking the tenant access to the cluster")
			_, err = commandHandlerClient().Execute(ctx, cmd.CreateCommand(tenantClusterBindingId, commandTypes.DeleteTenantClusterBinding))
			Expect(err).ToNot(HaveOccurred())

			Eventually(func(g Gomega) {
				_, err := clusterAccessClient().GetTenantClusterMappingByTenantAndClusterId(ctx, &domainApi.GetClusterMappingRequest{ClusterId: clusterId.String(), TenantId: tenantId.String()})
				Expect(err).To(HaveOccurred())
			}).Should(Succeed())

			By("granting the tenant access to the cluster again (after revoking the old one)")
			command, err = cmd.AddCommandData(
				cmd.CreateCommand(uuid.Nil, commandTypes.CreateTenantClusterBinding),
				&cmdData.CreateTenantClusterBindingCommandData{TenantId: tenantId.String(), ClusterId: clusterId.String()},
			)
			Expect(err).ToNot(HaveOccurred())
			reply, err = commandHandlerClient().Execute(ctx, command)
			Expect(err).ToNot(HaveOccurred())

			tenantClusterBindingIdNew := uuid.MustParse(reply.AggregateId)

			Eventually(func(g Gomega) {
				tenantClusterBinding, err := clusterAccessClient().GetTenantClusterMappingByTenantAndClusterId(ctx, &domainApi.GetClusterMappingRequest{ClusterId: clusterId.String(), TenantId: tenantId.String()})
				Expect(err).ToNot(HaveOccurred())
				Expect(tenantClusterBinding).ToNot(BeNil())
				Expect(tenantClusterBinding.Id).ToNot(Equal(tenantClusterBindingId.String()))
				Expect(tenantClusterBinding.Id).To(Equal(tenantClusterBindingIdNew.String()))
			}).Should(Succeed())
		})
	})
	It("can scrape event store metrics", func() {
		res, err := http.Get(fmt.Sprintf("http://%s/metrics", testEnv.EventStoreTestEnv.MetricsListener.Addr()))
		Expect(err).ToNot(HaveOccurred())
		defer res.Body.Close()
		Expect(res.StatusCode).To(Equal(200))
	})
})
