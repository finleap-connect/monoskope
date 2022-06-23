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
	cmdData "github.com/finleap-connect/monoskope/pkg/api/domain/commanddata"
	cmd "github.com/finleap-connect/monoskope/pkg/domain/commands"
	"github.com/finleap-connect/monoskope/pkg/domain/constants/aggregates"
	commandTypes "github.com/finleap-connect/monoskope/pkg/domain/constants/commands"
	"github.com/finleap-connect/monoskope/pkg/domain/constants/roles"
	"github.com/finleap-connect/monoskope/pkg/domain/constants/scopes"
	"github.com/finleap-connect/monoskope/pkg/domain/projections"
	"github.com/google/uuid"
	"io"
	"time"

	grpcUtil "github.com/finleap-connect/monoskope/pkg/grpc"
	"github.com/finleap-connect/monoskope/pkg/jwt"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/finleap-connect/monoskope/internal/gateway/auth"
	domainApi "github.com/finleap-connect/monoskope/pkg/api/domain"
	esApi "github.com/finleap-connect/monoskope/pkg/api/eventsourcing"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("AuditLog Test", func() {
	ctx := context.Background()
	expectedValidity := time.Hour * 1
	expectedNumUsers := 2                    // SUPER_USERS
	expectedNumEventsDoneOnAdmin := 2        // creation and admin role
	expectedNumEventsDoneByAdmin := 0        // to be counted see initEvents
	expectedNumEventsDoneByAdminMidTime := 0 // to be counted see initEvents

	getAdminAuthToken := func() string {
		signer := testEnv.gatewayTestEnv.JwtTestEnv.CreateSigner()
		token := auth.NewAuthToken(&jwt.StandardClaims{Name: testEnv.gatewayTestEnv.AdminUser.Name, Email: testEnv.gatewayTestEnv.AdminUser.Email}, testEnv.gatewayTestEnv.GetApiAddr(), testEnv.gatewayTestEnv.AdminUser.Id, expectedValidity)
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

	userServiceClient := func() domainApi.UserClient {
		_, client, err := grpcUtil.NewClientWithInsecureAuth(ctx, testEnv.queryHandlerTestEnv.GetApiAddr(), getAdminAuthToken(), domainApi.NewUserClient)
		Expect(err).ToNot(HaveOccurred())
		return client
	}

	// see PR #172
	adminWorkaround := func() {
		// the admin user created by commandhandler (SUPER_USERS) and known by queryHandler
		adminUser, err := userServiceClient().GetByEmail(ctx, wrapperspb.String(testEnv.gatewayTestEnv.AdminUser.Email))
		Expect(err).ToNot(HaveOccurred())

		// clean up gateway repos to avoid side effects
		err = testEnv.gatewayTestEnv.UserRepo.Remove(ctx, testEnv.gatewayTestEnv.AdminUser.ID())
		Expect(err).ToNot(HaveOccurred())
		gatewayAdminRoleBindings, err := testEnv.gatewayTestEnv.UserRoleBindingRepo.ByUserId(ctx, testEnv.gatewayTestEnv.AdminUser.ID())
		Expect(err).ToNot(HaveOccurred())
		err = testEnv.gatewayTestEnv.UserRoleBindingRepo.Remove(ctx, gatewayAdminRoleBindings[0].ID())
		Expect(err).ToNot(HaveOccurred())

		// replace the gateway admin user with queryHandler one
		testEnv.gatewayTestEnv.AdminUser = projections.NewUserProjection(uuid.MustParse(adminUser.Id))
		testEnv.gatewayTestEnv.AdminUser.Email = adminUser.Email
		testEnv.gatewayTestEnv.AdminUser.Name = adminUser.Name
		testEnv.gatewayTestEnv.AdminUser.Metadata = adminUser.Metadata
		testEnv.gatewayTestEnv.AdminUser.Source = adminUser.Source
		err = testEnv.gatewayTestEnv.UserRepo.Upsert(ctx, testEnv.gatewayTestEnv.AdminUser)
		Expect(err).ToNot(HaveOccurred())
		adminRoleBinding := projections.NewUserRoleBinding(uuid.New())
		adminRoleBinding.UserId = adminUser.Id
		adminRoleBinding.Role = string(roles.Admin)
		adminRoleBinding.Scope = string(scopes.System)
		err = testEnv.gatewayTestEnv.UserRoleBindingRepo.Upsert(ctx, adminRoleBinding)
		Expect(err).ToNot(HaveOccurred())
	}

	initEvents := func(commandHandlerClient func() esApi.CommandHandlerClient) time.Time {
		adminWorkaround() // remove when issue #182 is resolved

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
		expectedNumEventsDoneByAdmin++
		expectedNumUsers++

		// CreateUserRoleBinding on system level
		command, err = cmd.AddCommandData(
			cmd.CreateCommand(uuid.Nil, commandTypes.CreateUserRoleBinding),
			&cmdData.CreateUserRoleBindingCommandData{Role: string(roles.Admin), Scope: string(scopes.System), UserId: userId.String(), Resource: &wrapperspb.StringValue{Value: uuid.New().String()}},
		)
		Expect(err).ToNot(HaveOccurred())
		Eventually(func(g Gomega) {
			reply, err = commandHandlerClient().Execute(ctx, command)
			g.Expect(err).ToNot(HaveOccurred())
		}).Should(Succeed())
		userRoleBindingId := uuid.MustParse(reply.AggregateId)
		expectedNumEventsDoneByAdmin++

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
		expectedNumEventsDoneByAdmin++

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
		expectedNumEventsDoneByAdmin++

		// CreateUserRoleBinding on tenant level
		command, err = cmd.AddCommandData(
			cmd.CreateCommand(uuid.Nil, commandTypes.CreateUserRoleBinding),
			&cmdData.CreateUserRoleBindingCommandData{Role: string(roles.User), Scope: string(scopes.Tenant), UserId: userId.String(), Resource: &wrapperspb.StringValue{Value: tenantId.String()}},
		)
		Expect(err).ToNot(HaveOccurred())
		Eventually(func(g Gomega) {
			reply, err = commandHandlerClient().Execute(ctx, command)
			g.Expect(err).ToNot(HaveOccurred())
		}).Should(Succeed())
		_ = uuid.MustParse(reply.AggregateId)
		expectedNumEventsDoneByAdmin++

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
		expectedNumEventsDoneByAdmin++

		midTime := time.Now().UTC()
		expectedNumEventsDoneByAdminMidTime = expectedNumEventsDoneByAdmin

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
		expectedNumEventsDoneByAdmin++

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
		expectedNumEventsDoneByAdmin++

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
		expectedNumEventsDoneByAdmin++

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
		expectedNumEventsDoneByAdmin++

		// DeleteUser
		_, err = commandHandlerClient().Execute(ctx,
			cmd.CreateCommand(userId, commandTypes.DeleteUser))
		Expect(err).ToNot(HaveOccurred())
		expectedNumEventsDoneByAdmin++

		// DeleteUserRoleBinding
		_, err = commandHandlerClient().Execute(ctx,
			cmd.CreateCommand(userRoleBindingId, commandTypes.DeleteUserRoleBinding))
		Expect(err).ToNot(HaveOccurred())
		expectedNumEventsDoneByAdmin++

		// DeleteTenant
		_, err = commandHandlerClient().Execute(ctx,
			cmd.CreateCommand(tenantId, commandTypes.DeleteTenant))
		Expect(err).ToNot(HaveOccurred())
		expectedNumEventsDoneByAdmin++

		// DeleteTenantClusterBinding
		reply, err = commandHandlerClient().Execute(ctx,
			cmd.CreateCommand(tenantClusterBindingId, commandTypes.DeleteTenantClusterBinding))
		Expect(err).ToNot(HaveOccurred())
		expectedNumEventsDoneByAdmin++

		// DeleteCluster
		reply, err = commandHandlerClient().Execute(ctx,
			cmd.CreateCommand(clusterId, commandTypes.DeleteCluster))
		Expect(err).ToNot(HaveOccurred())
		expectedNumEventsDoneByAdmin++

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

			counter := 0
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
				counter++
			}
			Expect(counter).To(Equal(expectedNumEventsDoneByAdmin))

			By("using a custom range")
			dateRange.MaxTimestamp = timestamppb.New(midTime)

			events, err = auditLogServiceClient().GetByDateRange(ctx, dateRange)
			Expect(err).ToNot(HaveOccurred())

			counter = 0
			for {
				_, err := events.Recv()
				if err == io.EOF {
					break
				}
				Expect(err).ToNot(HaveOccurred())
				counter++
			}
			Expect(counter).To(Equal(expectedNumEventsDoneByAdminMidTime))
		})

		When("getting by user", func() {
			events, err := auditLogServiceClient().GetByUser(ctx, &domainApi.GetByUserRequest{
				Email: wrapperspb.String(testEnv.gatewayTestEnv.AdminUser.Email),
				DateRange: &domainApi.GetAuditLogByDateRangeRequest{
					MinTimestamp: timestamppb.New(time.Now().Add(time.Hour * -1).UTC()),
					MaxTimestamp: timestamppb.New(maxTime),
				},
			})
			Expect(err).ToNot(HaveOccurred())

			counter := 0
			for {
				e, err := events.Recv()
				if err == io.EOF {
					break
				}
				Expect(err).ToNot(HaveOccurred())

				Expect(e.IssuerId).ToNot(Equal(testEnv.gatewayTestEnv.AdminUser.Id))
				Expect(e.Details).To(ContainSubstring(testEnv.gatewayTestEnv.AdminUser.Email))
				counter++
			}
			Expect(counter).To(Equal(expectedNumEventsDoneOnAdmin))
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

			counter := 0
			for {
				e, err := events.Recv()
				if err == io.EOF {
					break
				}
				Expect(err).ToNot(HaveOccurred())

				Expect(e.Issuer).To(Equal(testEnv.gatewayTestEnv.AdminUser.Email))
				Expect(e.IssuerId).To(Equal(testEnv.gatewayTestEnv.AdminUser.Id))
				counter++
			}
			Expect(counter).To(Equal(expectedNumEventsDoneByAdmin))
		})

		When("getting users overview", func() {
			overviews, err := auditLogServiceClient().GetUsersOverview(ctx, &domainApi.GetUsersOverviewRequest{
				Timestamp: timestamppb.New(maxTime),
			})
			Expect(err).ToNot(HaveOccurred())

			counter := 0
			for {
				o, err := overviews.Recv()
				if err == io.EOF {
					break
				}
				Expect(err).ToNot(HaveOccurred())

				Expect(o.Name).ToNot(BeEmpty())
				Expect(o.Email).ToNot(BeEmpty())
				Expect(o.Details).ToNot(BeEmpty())
				counter++
			}
			Expect(counter).To(Equal(expectedNumUsers))
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
