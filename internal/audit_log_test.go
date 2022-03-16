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

package internal

import (
	"context"
	"github.com/finleap-connect/monoskope/pkg/domain/constants/aggregates"
	"github.com/finleap-connect/monoskope/pkg/domain/constants/roles"
	"github.com/finleap-connect/monoskope/pkg/domain/constants/scopes"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"io"
	"time"

	ch "github.com/finleap-connect/monoskope/internal/commandhandler"
	"github.com/finleap-connect/monoskope/internal/queryhandler"
	domainApi "github.com/finleap-connect/monoskope/pkg/api/domain"
	cmdData "github.com/finleap-connect/monoskope/pkg/api/domain/commanddata"
	esApi "github.com/finleap-connect/monoskope/pkg/api/eventsourcing"
	cmd "github.com/finleap-connect/monoskope/pkg/domain/commands"
	commandTypes "github.com/finleap-connect/monoskope/pkg/domain/constants/commands"
	"github.com/finleap-connect/monoskope/pkg/domain/metadata"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("AuditLog Test", func() {
	ctx := context.Background()

	mdManager, err := metadata.NewDomainMetadataManager(ctx)
	Expect(err).ToNot(HaveOccurred())

	mdManager.SetUserInformation(&metadata.UserInformation{
		Name:  "admin",
		Email: "admin@monoskope.io",
	})

	commandHandlerClient := func() esApi.CommandHandlerClient {
		chAddr := testEnv.commandHandlerTestEnv.GetApiAddr()
		_, chClient, err := ch.NewServiceClient(ctx, chAddr)
		Expect(err).ToNot(HaveOccurred())
		return chClient
	}

	auditLogServiceClient := func() domainApi.AuditLogClient {
		addr := testEnv.queryHandlerTestEnv.GetApiAddr()
		_, client, err := queryhandler.NewAuditLogClient(ctx, addr)
		Expect(err).ToNot(HaveOccurred())
		return client
	}

	It("can provide events by date range", func() {
		minTimestamp := time.Now().UTC()
		midTimestamp := initEvents(commandHandlerClient, mdManager)
		maxTimestamp := time.Now().UTC()

		By("using a general range")

		dateRange := &domainApi.GetAuditLogByDateRangeRequest{
			MinTimestamp: timestamppb.New(minTimestamp),
			MaxTimestamp: timestamppb.New(maxTimestamp),
		}

		Eventually(func(g Gomega) {
			events, err := auditLogServiceClient().GetByDateRange(ctx, dateRange)
			g.Expect(err).ToNot(HaveOccurred())

			for {
				e, err := events.Recv()
				if err == io.EOF {
					break
				}
				g.Expect(err).ToNot(HaveOccurred())

				g.Expect(e.When).ToNot(BeEmpty())
				g.Expect(e.Issuer).ToNot(BeEmpty())
				g.Expect(e.IssuerId).ToNot(BeEmpty())
				g.Expect(e.EventType).ToNot(BeEmpty())
				g.Expect(e.Details).ToNot(BeEmpty())
			}
		}).Should(Succeed())

		By("using a custom range")

		dateRange.MaxTimestamp = timestamppb.New(midTimestamp)

		Eventually(func(g Gomega) {
			events, err := auditLogServiceClient().GetByDateRange(ctx, dateRange)
			g.Expect(err).ToNot(HaveOccurred())

			counter := 0
			for {
				_, err := events.Recv()
				if err == io.EOF {
					break
				}
				g.Expect(err).ToNot(HaveOccurred())
				counter++
			}
			g.Expect(counter).To(Equal(4))
		}).Should(Succeed())
	})
})

func initEvents(commandHandlerClient func() esApi.CommandHandlerClient, mdManager *metadata.DomainMetadataManager) time.Time {
	// CreateUser
	command, err := cmd.AddCommandData(
		cmd.CreateCommand(uuid.Nil, commandTypes.CreateUser),
		&cmdData.CreateUserCommandData{Name: "Jane Dou", Email: "jane.dou@monoskope.io"},
	)
	Expect(err).ToNot(HaveOccurred())
	var reply *esApi.CommandReply
	Eventually(func(g Gomega) {
		reply, err = commandHandlerClient().Execute(mdManager.GetOutgoingGrpcContext(), command)
		g.Expect(err).ToNot(HaveOccurred())
	}).Should(Succeed())
	userId := uuid.MustParse(reply.AggregateId)

	// CreateUserRoleBinding
	command, err = cmd.AddCommandData(
		cmd.CreateCommand(uuid.Nil, commandTypes.CreateUserRoleBinding),
		&cmdData.CreateUserRoleBindingCommandData{Role: roles.Admin.String(), Scope: scopes.System.String(), UserId: userId.String(), Resource: &wrapperspb.StringValue{Value: uuid.New().String()}},
	)
	Expect(err).ToNot(HaveOccurred())
	Eventually(func(g Gomega) {
		reply, err = commandHandlerClient().Execute(mdManager.GetOutgoingGrpcContext(), command)
		g.Expect(err).ToNot(HaveOccurred())
	}).Should(Succeed())
	userRoleBindingId := uuid.MustParse(reply.AggregateId)

	// CreateTenant
	command, err = cmd.AddCommandData(
		cmd.CreateCommand(uuid.Nil, commandTypes.CreateTenant),
		&cmdData.CreateTenantCommandData{Name: "Tenant Y", Prefix: "ty"},
	)
	Expect(err).ToNot(HaveOccurred())
	Eventually(func(g Gomega) {
		reply, err = commandHandlerClient().Execute(mdManager.GetOutgoingGrpcContext(), command)
		g.Expect(err).ToNot(HaveOccurred())
	}).Should(Succeed())
	tenantId := uuid.MustParse(reply.AggregateId)

	// UpdateTenant
	command, err = cmd.AddCommandData(
		cmd.CreateCommand(tenantId, commandTypes.UpdateTenant),
		&cmdData.UpdateTenantCommandData{Name: &wrapperspb.StringValue{Value: "Tenant Z"}},
	)
	Expect(err).ToNot(HaveOccurred())
	Eventually(func(g Gomega) {
		reply, err = commandHandlerClient().Execute(mdManager.GetOutgoingGrpcContext(), command)
		g.Expect(err).ToNot(HaveOccurred())
	}).Should(Succeed())

	// 4 events
	midTimeStamp := time.Now().UTC()

	// CreateCluster
	command, err = cmd.AddCommandData(
		cmd.CreateCommand(uuid.Nil, commandTypes.CreateCluster),
		&cmdData.CreateCluster{DisplayName: "Cluster Y", Name: "cluster-y", ApiServerAddress: "y.cluster.com", CaCertBundle: []byte("This should be a certificate")},
	)
	Expect(err).ToNot(HaveOccurred())
	Eventually(func(g Gomega) {
		reply, err = commandHandlerClient().Execute(mdManager.GetOutgoingGrpcContext(), command)
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
		reply, err = commandHandlerClient().Execute(mdManager.GetOutgoingGrpcContext(), command)
		g.Expect(err).ToNot(HaveOccurred())
	}).Should(Succeed())

	// CreateTenantClusterBinding
	command, err = cmd.AddCommandData(
		cmd.CreateCommand(uuid.Nil, commandTypes.CreateTenantClusterBinding),
		&cmdData.CreateTenantClusterBindingCommandData{TenantId: tenantId.String(), ClusterId: clusterId.String()},
	)
	Expect(err).ToNot(HaveOccurred())
	Eventually(func(g Gomega) {
		reply, err = commandHandlerClient().Execute(mdManager.GetOutgoingGrpcContext(), command)
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
		reply, err = commandHandlerClient().Execute(mdManager.GetOutgoingGrpcContext(), command)
		g.Expect(err).ToNot(HaveOccurred())
	}).Should(Succeed())

	// DeleteUser
	_, err = commandHandlerClient().Execute(mdManager.GetOutgoingGrpcContext(),
		cmd.CreateCommand(userId, commandTypes.DeleteUser))
	Expect(err).ToNot(HaveOccurred())

	// DeleteUserRoleBinding
	_, err = commandHandlerClient().Execute(mdManager.GetOutgoingGrpcContext(),
		cmd.CreateCommand(userRoleBindingId, commandTypes.DeleteUserRoleBinding))
	Expect(err).ToNot(HaveOccurred())

	// DeleteTenant
	_, err = commandHandlerClient().Execute(mdManager.GetOutgoingGrpcContext(),
		cmd.CreateCommand(tenantId, commandTypes.DeleteTenant))
	Expect(err).ToNot(HaveOccurred())

	// DeleteTenantClusterBinding
	reply, err = commandHandlerClient().Execute(mdManager.GetOutgoingGrpcContext(),
		cmd.CreateCommand(tenantClusterBindingId, commandTypes.DeleteTenantClusterBinding))
	Expect(err).ToNot(HaveOccurred())

	// DeleteCluster
	reply, err = commandHandlerClient().Execute(mdManager.GetOutgoingGrpcContext(),
		cmd.CreateCommand(clusterId, commandTypes.DeleteCluster))
	Expect(err).ToNot(HaveOccurred())

	return midTimeStamp
}
