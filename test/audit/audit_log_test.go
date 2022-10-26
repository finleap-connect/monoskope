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

package audit

import (
	"context"
	"fmt"
	"io"
	"regexp"
	"strings"
	"time"

	"github.com/finleap-connect/monoskope/internal/gateway/auth"
	domainApi "github.com/finleap-connect/monoskope/pkg/api/domain"
	cmdData "github.com/finleap-connect/monoskope/pkg/api/domain/commanddata"
	esApi "github.com/finleap-connect/monoskope/pkg/api/eventsourcing"
	cmd "github.com/finleap-connect/monoskope/pkg/domain/commands"
	commandTypes "github.com/finleap-connect/monoskope/pkg/domain/constants/commands"
	fConsts "github.com/finleap-connect/monoskope/pkg/domain/constants/formatters"
	"github.com/finleap-connect/monoskope/pkg/domain/constants/roles"
	"github.com/finleap-connect/monoskope/pkg/domain/constants/scopes"
	"github.com/finleap-connect/monoskope/pkg/domain/mock"
	grpcUtil "github.com/finleap-connect/monoskope/pkg/grpc"
	"github.com/finleap-connect/monoskope/pkg/jwt"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

var _ = Describe("test/audit/audit_log_test", func() {
	var (
		ctx              = context.Background()
		userId           = uuid.New()
		userEmail        = "userxyz@monoskope.io"
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
		expectedClusterCACertBundle         = []byte(`-----BEGIN CERTIFICATE-----
		MIICnTCCAkSgAwIBAgIQMo7x823NtJ/Xyy1Wl+8+yzAKBggqhkjOPQQDAjAnMSUw
		IwYDVQQDExxyb290Lm1vbm9za29wZS5jbHVzdGVyLmxvY2FsMB4XDTIxMDYwMjAy
		MTAxNVoXDTIxMDYwNDAyMTAxNVowMDESMBAGA1UEChMJTW9ub3Nrb3BlMRowGAYD
		VQQDExFtOC1hdXRoZW50aWNhdGlvbjCCASIwDQYJKoZIhvcNAQEBBQADggEPADCC
		AQoCggEBANCKZWW0el3OzPw7914TC1Ld2At/xIh/3zoiawQcbS8mrjnVMO2oSomY
		mks6sEaWp4p80PwJkzSplpgoJmEOYqps+YXo+1NLp66bFPkAbMEZDsZ4QmrQQ7X3
		iv5IaDFW4vSGJFSkTQnUmedlhrWguasOD3vL0Pek89L8kQ09+YlDk/fpBZUXFADU
		+ef4GjTkWJzkg32dSOudJDYD4wUPczTFlRO097MBBlaMb4LKYfDfjuUKRCOAL3LD
		7kKAatHKeoADuBptUv/lQLExGNzlhRteaLocTHHab2hs+NCFYABv2Px5Tcnbw8g+
		/r/97gwKkpFeF5p4WhdVgbDYd2MGUlMCAwEAAaN+MHwwHQYDVR0lBBYwFAYIKwYB
		BQUHAwIGCCsGAQUFBwMBMAwGA1UdEwEB/wQCMAAwHwYDVR0jBBgwFoAUI+RyVqj0
		J9qH8l3pbY9KUkoTgHQwLAYDVR0RBCUwI4IJbG9jYWxob3N0hwR/AAABhxAAAAAA
		AAAAAAAAAAAAAAAAMAoGCCqGSM49BAMCA0cAMEQCIEPbvMo2YvqlYQtdkQwlhJci
		mTlsDv6VmO4WfCjrQdwLAiA+N0eeiL/yLPC5ReaPYQ7PeoXbc9+EPR2FBDrkiBbA
		8w==
		-----END CERTIFICATE-----`)
	)

	getAdminAuthToken := func() string {
		signer := testEnv.GatewayTestEnv.JwtTestEnv.CreateSigner()
		token := auth.NewAuthToken(&jwt.StandardClaims{Name: mock.TestAdminUser.Name, Email: mock.TestAdminUser.Email}, testEnv.GatewayTestEnv.GetApiAddr(), mock.TestAdminUser.Id, expectedValidity)
		authToken, err := signer.GenerateSignedToken(token)
		Expect(err).ToNot(HaveOccurred())
		return authToken
	}

	commandHandlerClient := func() esApi.CommandHandlerClient {
		chAddr := testEnv.CommandHandlerTestEnv.GetApiAddr()
		_, chClient, err := grpcUtil.NewClientWithInsecureAuth(ctx, chAddr, getAdminAuthToken(), esApi.NewCommandHandlerClient)
		Expect(err).ToNot(HaveOccurred())
		return chClient
	}

	auditLogServiceClient := func() domainApi.AuditLogClient {
		_, client, err := grpcUtil.NewClientWithInsecureAuth(ctx, testEnv.QueryHandlerTestEnv.GetApiAddr(), getAdminAuthToken(), domainApi.NewAuditLogClient)
		Expect(err).ToNot(HaveOccurred())
		return client
	}

	initEvents := func(commandHandlerClient func() esApi.CommandHandlerClient) time.Time {
		// CreateUser
		command := cmd.NewCommandWithData(
			uuid.Nil, commandTypes.CreateUser,
			&cmdData.CreateUserCommandData{Name: "XYZ", Email: userEmail},
		)
		var reply *esApi.CommandReply
		Eventually(func(g Gomega) {
			var err error
			reply, err = commandHandlerClient().Execute(ctx, command)
			g.Expect(err).ToNot(HaveOccurred())
		}).Should(Succeed())
		userId = uuid.MustParse(reply.AggregateId)
		expectedNumEventsDoneByAdmin++
		expectedNumEventsDoneOnUser++
		expectedNumUsers++
		expectedDetailMsgs = append(expectedDetailMsgs, fConsts.UserCreatedDetailsFormat.Sprint(mock.TestAdminUser.Email, userEmail))

		// CreateUserRoleBinding on system level
		command = cmd.NewCommandWithData(
			uuid.Nil, commandTypes.CreateUserRoleBinding,
			&cmdData.CreateUserRoleBindingCommandData{Role: string(roles.Admin), Scope: string(scopes.System), UserId: userId.String(), Resource: &wrapperspb.StringValue{Value: uuid.New().String()}},
		)
		Eventually(func(g Gomega) {
			var err error
			reply, err = commandHandlerClient().Execute(ctx, command)
			g.Expect(err).ToNot(HaveOccurred())
		}).Should(Succeed())
		userRoleBindingId := uuid.MustParse(reply.AggregateId)
		expectedNumEventsDoneByAdmin++
		expectedNumEventsDoneOnUser++
		expectedDetailMsgs = append(expectedDetailMsgs, fConsts.UserRoleAddedDetailsFormat.Sprint(mock.TestAdminUser.Email, roles.Admin, scopes.System, userEmail))

		// UpdateUser
		command = cmd.NewCommandWithData(
			userId, commandTypes.UpdateUser,
			&cmdData.UpdateUserCommandData{Name: &wrapperspb.StringValue{Value: "XYZ New"}},
		)
		Eventually(func(g Gomega) {
			var err error
			reply, err = commandHandlerClient().Execute(ctx, command)
			g.Expect(err).ToNot(HaveOccurred())
		}).Should(Succeed())
		expectedNumEventsDoneByAdmin++
		expectedNumEventsDoneOnUser++
		expectedDetailMsgs = append(expectedDetailMsgs, fConsts.UserUpdatedDetailsFormat.Sprint(mock.TestAdminUser.Email))

		// CreateTenant
		command = cmd.NewCommandWithData(
			uuid.Nil, commandTypes.CreateTenant,
			&cmdData.CreateTenantCommandData{Name: "Tenant Y", Prefix: "ty"},
		)
		Eventually(func(g Gomega) {
			var err error
			reply, err = commandHandlerClient().Execute(ctx, command)
			g.Expect(err).ToNot(HaveOccurred())
		}).Should(Succeed())
		tenantId := uuid.MustParse(reply.AggregateId)
		expectedNumEventsDoneByAdmin++
		expectedDetailMsgs = append(expectedDetailMsgs, fConsts.TenantCreatedDetailsFormat.Sprint(mock.TestAdminUser.Email, "Tenant Y", "ty"))

		// CreateUserRoleBinding on tenant level
		command = cmd.NewCommandWithData(
			uuid.Nil, commandTypes.CreateUserRoleBinding,
			&cmdData.CreateUserRoleBindingCommandData{Role: string(roles.User), Scope: string(scopes.Tenant), UserId: userId.String(), Resource: &wrapperspb.StringValue{Value: tenantId.String()}},
		)
		Eventually(func(g Gomega) {
			var err error
			reply, err = commandHandlerClient().Execute(ctx, command)
			g.Expect(err).ToNot(HaveOccurred())
		}).Should(Succeed())
		_ = uuid.MustParse(reply.AggregateId)
		expectedNumEventsDoneByAdmin++
		expectedNumEventsDoneOnUser++
		expectedDetailMsgs = append(expectedDetailMsgs, fConsts.UserRoleAddedDetailsFormat.Sprint(mock.TestAdminUser.Email, roles.User, scopes.Tenant, userEmail))
		expectedUserOverviewRoleMsgs = append(expectedUserOverviewRoleMsgs, fConsts.UserRoleBindingOverviewDetailsFormat.Sprint(scopes.Tenant, roles.User))
		expectedUserOverviewTenantMsgs = append(expectedUserOverviewTenantMsgs, fConsts.TenantUserRoleBindingOverviewDetailsFormat.Sprint("Tenant Z", roles.User))

		// UpdateTenant
		command = cmd.NewCommandWithData(
			tenantId, commandTypes.UpdateTenant,
			&cmdData.UpdateTenantCommandData{Name: &wrapperspb.StringValue{Value: "Tenant Z"}},
		)
		Eventually(func(g Gomega) {
			var err error
			reply, err = commandHandlerClient().Execute(ctx, command)
			g.Expect(err).ToNot(HaveOccurred())
		}).Should(Succeed())
		expectedNumEventsDoneByAdmin++
		expectedDetailMsgs = append(expectedDetailMsgs, fConsts.TenantUpdatedDetailsFormat.Sprint(mock.TestAdminUser.Email))

		midTime := time.Now().UTC()
		expectedNumEventsDoneByAdminMidTime = expectedNumEventsDoneByAdmin

		// CreateCluster
		command = cmd.NewCommandWithData(
			uuid.Nil, commandTypes.CreateCluster,
			&cmdData.CreateCluster{Name: "cluster-y", ApiServerAddress: "y.cluster.com", CaCertBundle: expectedClusterCACertBundle},
		)
		Eventually(func(g Gomega) {
			var err error
			reply, err = commandHandlerClient().Execute(ctx, command)
			g.Expect(err).ToNot(HaveOccurred())
		}).Should(Succeed())
		clusterId := uuid.MustParse(reply.AggregateId)
		expectedNumEventsDoneByAdmin++
		expectedDetailMsgs = append(expectedDetailMsgs, fConsts.ClusterCreatedDetailsFormat.Sprint(mock.TestAdminUser.Email, "cluster-y"))

		// UpdateCluster
		command = cmd.NewCommandWithData(
			clusterId, commandTypes.UpdateCluster,
			&cmdData.UpdateCluster{Name: &wrapperspb.StringValue{Value: "cluster-z"}, ApiServerAddress: &wrapperspb.StringValue{Value: "z.cluster.com"}, CaCertBundle: expectedClusterCACertBundle},
		)
		Eventually(func(g Gomega) {
			var err error
			reply, err = commandHandlerClient().Execute(ctx, command)
			g.Expect(err).ToNot(HaveOccurred())
		}).Should(Succeed())
		expectedNumEventsDoneByAdmin++
		expectedDetailMsgs = append(expectedDetailMsgs, fConsts.ClusterUpdatedDetailsFormat.Sprint(mock.TestAdminUser.Email))

		// CreateTenantClusterBinding
		command = cmd.NewCommandWithData(
			uuid.Nil, commandTypes.CreateTenantClusterBinding,
			&cmdData.CreateTenantClusterBindingCommandData{TenantId: tenantId.String(), ClusterId: clusterId.String()},
		)
		Eventually(func(g Gomega) {
			var err error
			reply, err = commandHandlerClient().Execute(ctx, command)
			g.Expect(err).ToNot(HaveOccurred())
		}).Should(Succeed())
		tenantClusterBindingId := uuid.MustParse(reply.AggregateId)
		expectedNumEventsDoneByAdmin++
		expectedDetailMsgs = append(expectedDetailMsgs, fConsts.TenantClusterBindingCreatedDetailsFormat.Sprint(mock.TestAdminUser.Email, "Tenant Z", "Cluster Z"))

		// DeleteUser
		var err error
		reply, err = commandHandlerClient().Execute(ctx, cmd.NewCommand(userId, commandTypes.DeleteUser))
		Expect(err).ToNot(HaveOccurred())
		expectedNumEventsDoneByAdmin++
		expectedNumEventsDoneOnUser++
		expectedDetailMsgs = append(expectedDetailMsgs, fConsts.UserDeletedDetailsFormat.Sprint(mock.TestAdminUser.Email, userEmail))
		expectedUserOverviewDetailMsgs = append(expectedUserOverviewDetailMsgs, strings.ReplaceAll(fConsts.UserDeletedOverviewDetailsFormat.Sprint(mock.TestAdminUser.Email, "x"), fConsts.Quote("x"), ""))

		// DeleteUserRoleBinding
		_, err = commandHandlerClient().Execute(ctx,
			cmd.NewCommand(userRoleBindingId, commandTypes.DeleteUserRoleBinding))
		Expect(err).ToNot(HaveOccurred())
		expectedNumEventsDoneByAdmin++
		expectedDetailMsgs = append(expectedDetailMsgs, fConsts.UserRoleBindingDeletedDetailsFormat.Sprint(mock.TestAdminUser.Email, roles.Admin, scopes.System, userEmail))

		// DeleteTenant
		_, err = commandHandlerClient().Execute(ctx,
			cmd.NewCommand(tenantId, commandTypes.DeleteTenant))
		Expect(err).ToNot(HaveOccurred())
		expectedNumEventsDoneByAdmin++
		expectedDetailMsgs = append(expectedDetailMsgs, fConsts.TenantDeletedDetailsFormat.Sprint(mock.TestAdminUser.Email, "Tenant Z"))

		// DeleteTenantClusterBinding
		_, err = commandHandlerClient().Execute(ctx, cmd.NewCommand(tenantClusterBindingId, commandTypes.DeleteTenantClusterBinding))
		Expect(err).ToNot(HaveOccurred())
		expectedNumEventsDoneByAdmin++
		expectedDetailMsgs = append(expectedDetailMsgs, fConsts.TenantClusterBindingDeletedDetailsFormat.Sprint(mock.TestAdminUser.Email, "Cluster Z", "Tenant Z"))

		// DeleteCluster
		_, err = commandHandlerClient().Execute(ctx, cmd.NewCommand(clusterId, commandTypes.DeleteCluster))
		Expect(err).ToNot(HaveOccurred())
		expectedNumEventsDoneByAdmin++
		expectedDetailMsgs = append(expectedDetailMsgs, fConsts.ClusterDeletedDetailsFormat.Sprint(mock.TestAdminUser.Email, "Cluster Z"))

		return midTime
	}

	It("can provide human-readable events/overviews", func() {
		fromTime := time.Now().UTC()
		betweenTime := initEvents(commandHandlerClient)
		toTime := time.Now().UTC()

		When("getting by date range", func() {
			By("using a general range")
			dateRange := &domainApi.GetAuditLogByDateRangeRequest{
				MinTimestamp: timestamppb.New(fromTime),
				MaxTimestamp: timestamppb.New(toTime),
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
			dateRange.MaxTimestamp = timestamppb.New(betweenTime)

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
					MinTimestamp: timestamppb.New(fromTime),
					MaxTimestamp: timestamppb.New(toTime),
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
				Email: wrapperspb.String(mock.TestAdminUser.Email),
				DateRange: &domainApi.GetAuditLogByDateRangeRequest{
					MinTimestamp: timestamppb.New(fromTime),
					MaxTimestamp: timestamppb.New(toTime),
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

				Expect(e.Issuer).To(Equal(mock.TestAdminUser.Email))
				Expect(e.IssuerId).To(Equal(mock.TestAdminUser.Id))
				counter++
			}
			Expect(counter).To(Equal(expectedNumEventsDoneByAdmin))
		})

		When("getting users overview", func() {
			overviews, err := auditLogServiceClient().GetUsersOverview(ctx, &domainApi.GetUsersOverviewRequest{
				Timestamp: timestamppb.New(toTime),
			})
			Expect(err).ToNot(HaveOccurred())

			// "shared" testEnv workaround
			// ginkgo v2 should solve this by utilizing baforeAll/afterAll?
			knownUsersSet := make(map[string]struct{})
			for _, s := range mock.TestMockUsers {
				knownUsersSet[s.Email] = struct{}{}
			}

			counter := 0
			for {
				o, err := overviews.Recv()
				if err == io.EOF {
					break
				}
				Expect(err).ToNot(HaveOccurred())

				if _, known := knownUsersSet[o.Email]; known {
					continue
				}

				fmt.Printf("%v", o)
				Expect(o.Name).ToNot(BeEmpty())
				Expect(o.Email).ToNot(BeEmpty())
				Expect(regexp.MatchString(`.*`+regexp.QuoteMeta(strings.TrimSpace(expectedUserOverviewRoleMsgs[counter]))+`.*`, o.Roles)).To(BeTrue())
				Expect(regexp.MatchString(`.*`+regexp.QuoteMeta(strings.TrimSpace(expectedUserOverviewTenantMsgs[counter]))+`.*`, o.Tenants)).To(BeTrue())
				Expect(regexp.MatchString(`.*`+regexp.QuoteMeta(strings.TrimSpace(expectedUserOverviewDetailMsgs[counter]))+`.*`, o.Details)).To(BeTrue())
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
				Email: wrapperspb.String(mock.TestAdminUser.Email),
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
