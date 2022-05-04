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
	"io"
	"time"

	"github.com/finleap-connect/monoskope/pkg/domain/constants/aggregates"
	"github.com/finleap-connect/monoskope/pkg/domain/constants/roles"
	"github.com/finleap-connect/monoskope/pkg/domain/constants/scopes"
	grpcUtil "github.com/finleap-connect/monoskope/pkg/grpc"
	"github.com/finleap-connect/monoskope/pkg/jwt"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/finleap-connect/monoskope/internal/gateway/auth"
	domainApi "github.com/finleap-connect/monoskope/pkg/api/domain"
	cmdData "github.com/finleap-connect/monoskope/pkg/api/domain/commanddata"
	esApi "github.com/finleap-connect/monoskope/pkg/api/eventsourcing"
	cmd "github.com/finleap-connect/monoskope/pkg/domain/commands"
	commandTypes "github.com/finleap-connect/monoskope/pkg/domain/constants/commands"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("AuditLog Test", func() {
	ctx := context.Background()
	userEmail := "jane.dou@monoskope.io"
	expectedValidity := time.Hour * 1

	getAdminAuthToken := func() string {
		signer := testEnv.gatewayTestEnv.JwtTestEnv.CreateSigner()
		token := auth.NewAuthToken(&jwt.StandardClaims{Name: testEnv.gatewayTestEnv.AdminUser.Name, Email: testEnv.gatewayTestEnv.AdminUser.Email}, testEnv.gatewayTestEnv.GetApiAddr(), testEnv.gatewayTestEnv.AdminUser.ID().String(), expectedValidity)
		authToken, err := signer.GenerateSignedToken(token)
		Expect(err).ToNot(HaveOccurred())
		return authToken
	}

	commandHandlerClient := func() esApi.CommandHandlerClient {
		chAddr := testEnv.commandHandlerTestEnv.GetApiAddr()
		_, chClient, err := grpcUtil.NewClientWithInsecureAuth(ctx, chAddr, getAdminAuthToken(), esApi.NewCommandHandlerClient)
		Expect(err).ToNot(HaveOccurred())
		return chClient
	}

	auditLogServiceClient := func() domainApi.AuditLogClient {
		_, client, err := grpcUtil.NewClientWithInsecureAuth(ctx, testEnv.queryHandlerTestEnv.GetApiAddr(), getAdminAuthToken(), domainApi.NewAuditLogClient)
		Expect(err).ToNot(HaveOccurred())
		return client
	}

	initEvents := func(commandHandlerClient func() esApi.CommandHandlerClient) time.Time {
		// CreateUser
		command, err := cmd.AddCommandData(
			cmd.CreateCommand(uuid.Nil, commandTypes.CreateUser),
			&cmdData.CreateUserCommandData{Name: "Jane Dou", Email: "jane.dou@monoskope.io"},
		)
		Expect(err).ToNot(HaveOccurred())
		var reply *esApi.CommandReply
		Eventually(func(g Gomega) {
			reply, err = commandHandlerClient().Execute(ctx, command)
			g.Expect(err).ToNot(HaveOccurred())
		}).Should(Succeed())
		userId := uuid.MustParse(reply.AggregateId)

		// CreateUserRoleBinding on system level
		command, err = cmd.AddCommandData(
			cmd.CreateCommand(uuid.Nil, commandTypes.CreateUserRoleBinding),
			&cmdData.CreateUserRoleBindingCommandData{Role: roles.Admin.String(), Scope: scopes.System.String(), UserId: userId.String(), Resource: &wrapperspb.StringValue{Value: uuid.New().String()}},
		)
		Expect(err).ToNot(HaveOccurred())
		Eventually(func(g Gomega) {
			reply, err = commandHandlerClient().Execute(ctx, command)
			g.Expect(err).ToNot(HaveOccurred())
		}).Should(Succeed())
		userRoleBindingId := uuid.MustParse(reply.AggregateId)

		// UpdateUser
		command, err = cmd.AddCommandData(
			cmd.CreateCommand(userId, commandTypes.UpdateUser),
			&cmdData.UpdateUserCommandData{Name: &wrapperspb.StringValue{Value: "Jane New"}},
		)
		Expect(err).ToNot(HaveOccurred())
		Eventually(func(g Gomega) {
			reply, err = commandHandlerClient().Execute(ctx, command)
			g.Expect(err).ToNot(HaveOccurred())
		}).Should(Succeed())

		// CreateTenant
		command, err = cmd.AddCommandData(
			cmd.CreateCommand(uuid.Nil, commandTypes.CreateTenant),
			&cmdData.CreateTenantCommandData{Name: "Tenant Y", Prefix: "ty"},
		)
		Expect(err).ToNot(HaveOccurred())
		Eventually(func(g Gomega) {
			reply, err = commandHandlerClient().Execute(ctx, command)
			g.Expect(err).ToNot(HaveOccurred())
		}).Should(Succeed())
		tenantId := uuid.MustParse(reply.AggregateId)

		// CreateUserRoleBinding on tenant level
		command, err = cmd.AddCommandData(
			cmd.CreateCommand(uuid.Nil, commandTypes.CreateUserRoleBinding),
			&cmdData.CreateUserRoleBindingCommandData{Role: roles.User.String(), Scope: scopes.Tenant.String(), UserId: userId.String(), Resource: &wrapperspb.StringValue{Value: tenantId.String()}},
		)
		Expect(err).ToNot(HaveOccurred())
		Eventually(func(g Gomega) {
			reply, err = commandHandlerClient().Execute(ctx, command)
			g.Expect(err).ToNot(HaveOccurred())
		}).Should(Succeed())
		_ = uuid.MustParse(reply.AggregateId)

		// UpdateTenant
		command, err = cmd.AddCommandData(
			cmd.CreateCommand(tenantId, commandTypes.UpdateTenant),
			&cmdData.UpdateTenantCommandData{Name: &wrapperspb.StringValue{Value: "Tenant Z"}},
		)
		Expect(err).ToNot(HaveOccurred())
		Eventually(func(g Gomega) {
			reply, err = commandHandlerClient().Execute(ctx, command)
			g.Expect(err).ToNot(HaveOccurred())
		}).Should(Succeed())

		// 6 events
		midTime := time.Now().UTC()

		// CreateCluster
		command, err = cmd.AddCommandData(
			cmd.CreateCommand(uuid.Nil, commandTypes.CreateCluster),
			&cmdData.CreateCluster{DisplayName: "Cluster Y", Name: "cluster-y", ApiServerAddress: "y.cluster.com", CaCertBundle: []byte("This should be a certificate")},
		)
		Expect(err).ToNot(HaveOccurred())
		Eventually(func(g Gomega) {
			reply, err = commandHandlerClient().Execute(ctx, command)
			g.Expect(err).ToNot(HaveOccurred())
		}).Should(Succeed())
		clusterId := uuid.MustParse(reply.AggregateId)

		// UpdateCluster
		command, err = cmd.AddCommandData(
			cmd.CreateCommand(clusterId, commandTypes.UpdateCluster),
			&cmdData.UpdateCluster{DisplayName: &wrapperspb.StringValue{Value: "Cluster Z"}, ApiServerAddress: &wrapperspb.StringValue{Value: "z.cluster.com"}, CaCertBundle: []byte("This should be a new certificate")},
		)
		Expect(err).ToNot(HaveOccurred())
		Eventually(func(g Gomega) {
			reply, err = commandHandlerClient().Execute(ctx, command)
			g.Expect(err).ToNot(HaveOccurred())
		}).Should(Succeed())

		// CreateTenantClusterBinding
		command, err = cmd.AddCommandData(
			cmd.CreateCommand(uuid.Nil, commandTypes.CreateTenantClusterBinding),
			&cmdData.CreateTenantClusterBindingCommandData{TenantId: tenantId.String(), ClusterId: clusterId.String()},
		)
		Expect(err).ToNot(HaveOccurred())
		Eventually(func(g Gomega) {
			reply, err = commandHandlerClient().Execute(ctx, command)
			g.Expect(err).ToNot(HaveOccurred())
		}).Should(Succeed())
		tenantClusterBindingId := uuid.MustParse(reply.AggregateId)

		// RequestCertificate
		command, err = cmd.AddCommandData(
			cmd.CreateCommand(uuid.Nil, commandTypes.RequestCertificate),
			&cmdData.RequestCertificate{
				ReferencedAggregateId:   clusterId.String(),
				ReferencedAggregateType: aggregates.Cluster.String(),
				SigningRequest:          []byte("-----BEGIN CERTIFICATE REQUEST-----this is a CSR-----END CERTIFICATE REQUEST-----"),
			},
		)
		Expect(err).ToNot(HaveOccurred())
		Eventually(func(g Gomega) {
			reply, err = commandHandlerClient().Execute(ctx, command)
			g.Expect(err).ToNot(HaveOccurred())
		}).Should(Succeed())

		// DeleteUser
		_, err = commandHandlerClient().Execute(ctx,
			cmd.CreateCommand(userId, commandTypes.DeleteUser))
		Expect(err).ToNot(HaveOccurred())

		// DeleteUserRoleBinding
		_, err = commandHandlerClient().Execute(ctx,
			cmd.CreateCommand(userRoleBindingId, commandTypes.DeleteUserRoleBinding))
		Expect(err).ToNot(HaveOccurred())

		// DeleteTenant
		_, err = commandHandlerClient().Execute(ctx,
			cmd.CreateCommand(tenantId, commandTypes.DeleteTenant))
		Expect(err).ToNot(HaveOccurred())

		// DeleteTenantClusterBinding
		reply, err = commandHandlerClient().Execute(ctx,
			cmd.CreateCommand(tenantClusterBindingId, commandTypes.DeleteTenantClusterBinding))
		Expect(err).ToNot(HaveOccurred())

		// DeleteCluster
		reply, err = commandHandlerClient().Execute(ctx,
			cmd.CreateCommand(clusterId, commandTypes.DeleteCluster))
		Expect(err).ToNot(HaveOccurred())

		return midTime
	}

	It("can provide human-readable events/overviews", func() {
		minTime := time.Now().UTC()
		midTime := initEvents(commandHandlerClient)
		maxTime := time.Now().UTC()

		When("getting by date range", func() {
			By("using a general range")
			dateRange := &domainApi.GetAuditLogByDateRangeRequest{
				MinTimestamp: timestamppb.New(minTime),
				MaxTimestamp: timestamppb.New(maxTime),
			}

			events, err := auditLogServiceClient().GetByDateRange(ctx, dateRange)
			Expect(err).ToNot(HaveOccurred())

			for {
				e, err := events.Recv()
				if err == io.EOF {
					break
				}
				Expect(err).ToNot(HaveOccurred())

				Expect(e.When).ToNot(BeEmpty())
				Expect(e.Issuer).ToNot(BeEmpty())
				Expect(e.IssuerId).ToNot(BeEmpty())
				Expect(e.EventType).ToNot(BeEmpty())
				Expect(e.Details).ToNot(BeEmpty())
			}

			By("using a custom range")
			dateRange.MaxTimestamp = timestamppb.New(midTime)

			events, err = auditLogServiceClient().GetByDateRange(ctx, dateRange)
			Expect(err).ToNot(HaveOccurred())

			counter := 0
			for {
				_, err := events.Recv()
				if err == io.EOF {
					break
				}
				Expect(err).ToNot(HaveOccurred())
				counter++
			}
			Expect(counter).To(Equal(6)) // see midTime definition
		})

		When("getting by user", func() {
			events, err := auditLogServiceClient().GetByUser(ctx, &domainApi.GetByUserRequest{
				Email: wrapperspb.String(userEmail),
				DateRange: &domainApi.GetAuditLogByDateRangeRequest{
					MinTimestamp: timestamppb.New(minTime),
					MaxTimestamp: timestamppb.New(maxTime),
				},
			})
			Expect(err).ToNot(HaveOccurred())

			for {
				e, err := events.Recv()
				if err == io.EOF {
					break
				}
				Expect(err).ToNot(HaveOccurred())

				Expect(e.Details).To(ContainSubstring(userEmail))
			}
		})

		When("getting user actions", func() {
			events, err := auditLogServiceClient().GetUserActions(ctx, &domainApi.GetUserActionsRequest{
				Email: wrapperspb.String(testEnv.gatewayTestEnv.AdminUser.Email),
				DateRange: &domainApi.GetAuditLogByDateRangeRequest{
					MinTimestamp: timestamppb.New(minTime),
					MaxTimestamp: timestamppb.New(maxTime),
				},
			})
			Expect(err).ToNot(HaveOccurred())

			for {
				e, err := events.Recv()
				if err == io.EOF {
					break
				}
				Expect(err).ToNot(HaveOccurred())

				Expect(e.Issuer).To(Equal(testEnv.gatewayTestEnv.AdminUser.Email))
			}
		})

		When("getting users overview", func() {
			overviews, err := auditLogServiceClient().GetUsersOverview(ctx, &domainApi.GetUsersOverviewRequest{
				Timestamp: timestamppb.New(maxTime),
			})
			Expect(err).ToNot(HaveOccurred())

			for {
				o, err := overviews.Recv()
				if err == io.EOF {
					break
				}
				Expect(err).ToNot(HaveOccurred())

				Expect(o.Name).ToNot(BeEmpty())
				Expect(o.Email).ToNot(BeEmpty())
				Expect(o.Details).ToNot(BeEmpty())
			}
		})
	})

	It("can not provide human-readable events", func() {
		When("getting user actions with a range that exceeds one year", func() {
			minTime := time.Date(2021, time.December, 1, 0, 0, 0, 0, time.UTC)
			maxTime := time.Date(2022, time.December, 1, 0, 0, 0, 1, time.UTC)
			request := &domainApi.GetUserActionsRequest{
				Email: wrapperspb.String(testEnv.gatewayTestEnv.AdminUser.Email),
				DateRange: &domainApi.GetAuditLogByDateRangeRequest{
					MinTimestamp: timestamppb.New(minTime),
					MaxTimestamp: timestamppb.New(maxTime),
				},
			}
			events, err := auditLogServiceClient().GetUserActions(ctx, request)
			Expect(err).ToNot(HaveOccurred())

			_, err = events.Recv()
			Expect(err).To(HaveOccurred())
		})
	})
})
