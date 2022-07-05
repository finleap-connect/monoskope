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
	"regexp"
	"strings"
	"time"

	"github.com/finleap-connect/monoskope/internal/gateway/auth"
	domainApi "github.com/finleap-connect/monoskope/pkg/api/domain"
	cmdData "github.com/finleap-connect/monoskope/pkg/api/domain/commanddata"
	esApi "github.com/finleap-connect/monoskope/pkg/api/eventsourcing"
	cmd "github.com/finleap-connect/monoskope/pkg/domain/commands"
	"github.com/finleap-connect/monoskope/pkg/domain/constants/aggregates"
	commandTypes "github.com/finleap-connect/monoskope/pkg/domain/constants/commands"
	fConsts "github.com/finleap-connect/monoskope/pkg/domain/constants/formatters"
	"github.com/finleap-connect/monoskope/pkg/domain/constants/roles"
	"github.com/finleap-connect/monoskope/pkg/domain/constants/scopes"
	"github.com/finleap-connect/monoskope/pkg/domain/constants/users"
	"github.com/finleap-connect/monoskope/pkg/domain/projections"
	grpcUtil "github.com/finleap-connect/monoskope/pkg/grpc"
	"github.com/finleap-connect/monoskope/pkg/jwt"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

var _ = Describe("AuditLog Test", func() {
	ctx := context.Background()

	var (
		userId           = uuid.New()
		userEmail        = "jane.dou@monoskope.io"
		expectedValidity = time.Hour * 1
		// see initEvents
		expectedNumUsers                    = 0
		expectedNumEventsDoneOnUser         = 0
		expectedNumEventsDoneByAdmin        = 0
		expectedNumEventsDoneByAdminMidTime = 0
		expectedDetailMsgs                  []string
		expectedUserOverviewDetailMsgs      []string
		expectedUserOverviewRoleMsgs        []string
		expectedUserOverviewTenantMsgs      []string
	)

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

		// SUPER_USERS
		for _, su := range testEnv.superUsers {
			expectedUserOverviewDetailMsgs = append(expectedUserOverviewDetailMsgs, strings.ReplaceAll(fConsts.UserCreatedOverviewDetailsFormat.Sprint(su, "system@"+users.BASE_DOMAIN, "x"), fConsts.Quote("x"), ""))
			expectedUserOverviewRoleMsgs = append(expectedUserOverviewRoleMsgs, fConsts.UserRoleBindingOverviewDetailsFormat.Sprint(scopes.System, roles.Admin))
			expectedUserOverviewTenantMsgs = append(expectedUserOverviewTenantMsgs, "")
			expectedNumUsers++
		}

		// CreateUser
		command, err := cmd.AddCommandData(
			cmd.CreateCommand(uuid.Nil, commandTypes.CreateUser),
			&cmdData.CreateUserCommandData{Name: "Jane Dou", Email: userEmail},
		)
		Expect(err).ToNot(HaveOccurred())
		var reply *esApi.CommandReply
		Eventually(func(g Gomega) {
			reply, err = commandHandlerClient().Execute(ctx, command)
			g.Expect(err).ToNot(HaveOccurred())
		}).Should(Succeed())
		userId = uuid.MustParse(reply.AggregateId)
		expectedNumEventsDoneByAdmin++
		expectedNumEventsDoneOnUser++
		expectedNumUsers++
		expectedDetailMsgs = append(expectedDetailMsgs, fConsts.UserCreatedDetailsFormat.Sprint(testEnv.gatewayTestEnv.AdminUser.Email, userEmail))

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
		expectedNumEventsDoneOnUser++
		expectedDetailMsgs = append(expectedDetailMsgs, fConsts.UserRoleAddedDetailsFormat.Sprint(testEnv.gatewayTestEnv.AdminUser.Email, roles.Admin, scopes.System, userEmail))

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
		expectedNumEventsDoneOnUser++
		expectedDetailMsgs = append(expectedDetailMsgs, fConsts.UserUpdatedDetailsFormat.Sprint(testEnv.gatewayTestEnv.AdminUser.Email))

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
		expectedDetailMsgs = append(expectedDetailMsgs, fConsts.TenantCreatedDetailsFormat.Sprint(testEnv.gatewayTestEnv.AdminUser.Email, "Tenant Y", "ty"))

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
		expectedNumEventsDoneOnUser++
		expectedDetailMsgs = append(expectedDetailMsgs, fConsts.UserRoleAddedDetailsFormat.Sprint(testEnv.gatewayTestEnv.AdminUser.Email, roles.User, scopes.Tenant, userEmail))
		expectedUserOverviewRoleMsgs = append(expectedUserOverviewRoleMsgs, fConsts.UserRoleBindingOverviewDetailsFormat.Sprint(scopes.Tenant, roles.User))
		expectedUserOverviewTenantMsgs = append(expectedUserOverviewTenantMsgs, fConsts.TenantUserRoleBindingOverviewDetailsFormat.Sprint("Tenant Z", roles.User))

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
		expectedDetailMsgs = append(expectedDetailMsgs, fConsts.TenantUpdatedDetailsFormat.Sprint(testEnv.gatewayTestEnv.AdminUser.Email))

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
		expectedDetailMsgs = append(expectedDetailMsgs, fConsts.ClusterCreatedDetailsFormat.Sprint(testEnv.gatewayTestEnv.AdminUser.Email, "cluster-y"))

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
		expectedDetailMsgs = append(expectedDetailMsgs, fConsts.ClusterUpdatedDetailsFormat.Sprint(testEnv.gatewayTestEnv.AdminUser.Email))

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
		expectedDetailMsgs = append(expectedDetailMsgs, fConsts.TenantClusterBindingCreatedDetailsFormat.Sprint(testEnv.gatewayTestEnv.AdminUser.Email, "Tenant Z", "Cluster Z"))

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
		expectedDetailMsgs = append(expectedDetailMsgs, fConsts.CertificateRequestedDetailsFormat.Sprint(testEnv.gatewayTestEnv.AdminUser.Email))

		// DeleteUser
		_, err = commandHandlerClient().Execute(ctx,
			cmd.CreateCommand(userId, commandTypes.DeleteUser))
		Expect(err).ToNot(HaveOccurred())
		expectedNumEventsDoneByAdmin++
		expectedNumEventsDoneOnUser++
		expectedDetailMsgs = append(expectedDetailMsgs, fConsts.UserDeletedDetailsFormat.Sprint(testEnv.gatewayTestEnv.AdminUser.Email, userEmail))
		expectedUserOverviewDetailMsgs = append(expectedUserOverviewDetailMsgs, strings.ReplaceAll(fConsts.UserDeletedOverviewDetailsFormat.Sprint(testEnv.gatewayTestEnv.AdminUser.Email, "x"), fConsts.Quote("x"), ""))

		// DeleteUserRoleBinding
		_, err = commandHandlerClient().Execute(ctx,
			cmd.CreateCommand(userRoleBindingId, commandTypes.DeleteUserRoleBinding))
		Expect(err).ToNot(HaveOccurred())
		expectedNumEventsDoneByAdmin++
		expectedDetailMsgs = append(expectedDetailMsgs, fConsts.UserRoleBindingDeletedDetailsFormat.Sprint(testEnv.gatewayTestEnv.AdminUser.Email, roles.Admin, scopes.System, userEmail))

		// DeleteTenant
		_, err = commandHandlerClient().Execute(ctx,
			cmd.CreateCommand(tenantId, commandTypes.DeleteTenant))
		Expect(err).ToNot(HaveOccurred())
		expectedNumEventsDoneByAdmin++
		expectedDetailMsgs = append(expectedDetailMsgs, fConsts.TenantDeletedDetailsFormat.Sprint(testEnv.gatewayTestEnv.AdminUser.Email, "Tenant Z"))

		// DeleteTenantClusterBinding
		reply, err = commandHandlerClient().Execute(ctx,
			cmd.CreateCommand(tenantClusterBindingId, commandTypes.DeleteTenantClusterBinding))
		Expect(err).ToNot(HaveOccurred())
		expectedNumEventsDoneByAdmin++
		expectedDetailMsgs = append(expectedDetailMsgs, fConsts.TenantClusterBindingDeletedDetailsFormat.Sprint(testEnv.gatewayTestEnv.AdminUser.Email, "Cluster Z", "Tenant Z"))

		// DeleteCluster
		reply, err = commandHandlerClient().Execute(ctx,
			cmd.CreateCommand(clusterId, commandTypes.DeleteCluster))
		Expect(err).ToNot(HaveOccurred())
		expectedNumEventsDoneByAdmin++
		expectedDetailMsgs = append(expectedDetailMsgs, fConsts.ClusterDeletedDetailsFormat.Sprint(testEnv.gatewayTestEnv.AdminUser.Email, "Cluster Z"))

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

				Expect(e.Timestamp).ToNot(BeNil())
				Expect(e.Issuer).ToNot(BeEmpty())
				Expect(e.IssuerId).ToNot(BeEmpty())
				Expect(e.EventType).ToNot(BeEmpty())
				Expect(regexp.MatchString(`.*`+regexp.QuoteMeta(strings.TrimSpace(expectedDetailMsgs[counter]))+`.*`, e.Details)).To(BeTrue())
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
				Email: wrapperspb.String(userEmail),
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

				Expect(e.IssuerId).ToNot(Equal(userId))
				Expect(e.Details).To(ContainSubstring(userEmail))
				counter++
			}
			Expect(counter).To(Equal(expectedNumEventsDoneOnUser))
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

			// "shared" testEnv workaround
			// ginkgo v2 should solve this by utilizing baforeAll/afterAll?
			knownUsersSet := make(map[string]struct{}, len(testEnv.superUsers)+1)
			for _, s := range testEnv.superUsers {
				knownUsersSet[s] = struct{}{}
			}
			knownUsersSet[userEmail] = struct{}{}

			counter := 0
			for {
				o, err := overviews.Recv()
				if err == io.EOF {
					break
				}
				Expect(err).ToNot(HaveOccurred())

				Expect(o.Name).ToNot(BeEmpty())
				Expect(o.Email).ToNot(BeEmpty())
				if _, known := knownUsersSet[o.Email]; known {
					Expect(regexp.MatchString(`.*`+regexp.QuoteMeta(strings.TrimSpace(expectedUserOverviewRoleMsgs[counter]))+`.*`, o.Roles)).To(BeTrue())
					Expect(regexp.MatchString(`.*`+regexp.QuoteMeta(strings.TrimSpace(expectedUserOverviewTenantMsgs[counter]))+`.*`, o.Tenants)).To(BeTrue())
					Expect(regexp.MatchString(`.*`+regexp.QuoteMeta(strings.TrimSpace(expectedUserOverviewDetailMsgs[counter]))+`.*`, o.Details)).To(BeTrue())
					counter++
				}
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
