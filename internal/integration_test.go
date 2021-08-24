package internal

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	ch "gitlab.figo.systems/platform/monoskope/monoskope/internal/commandhandler"
	"gitlab.figo.systems/platform/monoskope/monoskope/internal/eventstore"
	"gitlab.figo.systems/platform/monoskope/monoskope/internal/queryhandler"
	testReactor "gitlab.figo.systems/platform/monoskope/monoskope/internal/test/reactor"
	domainApi "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/domain"
	cmdData "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/domain/commanddata"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/domain/common"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/domain/eventdata"
	esApi "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/eventsourcing"
	cmd "gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/commands"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/constants/aggregates"
	commandTypes "gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/constants/commands"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/constants/events"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/constants/roles"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/constants/scopes"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/errors"
	metadata "gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/metadata"
	es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

var _ = Describe("integration", func() {
	ctx := context.Background()

	expectedClusterDisplayName := "the one cluster"
	expectedClusterName := "one-cluster"
	expectedClusterApiServerAddress := "one.example.com"
	expectedClusterCACertBundle := []byte("This should be a certificate")

	mdManager, err := metadata.NewDomainMetadataManager(ctx)
	Expect(err).ToNot(HaveOccurred())

	mdManager.SetUserInformation(&metadata.UserInformation{
		Name:   "admin",
		Email:  "admin@monoskope.io",
		Issuer: "monoskope",
	})

	commandHandlerClient := func() esApi.CommandHandlerClient {
		chAddr := testEnv.commandHandlerTestEnv.GetApiAddr()
		_, chClient, err := ch.NewServiceClient(ctx, chAddr)
		Expect(err).ToNot(HaveOccurred())
		return chClient
	}

	userServiceClient := func() domainApi.UserClient {
		addr := testEnv.queryHandlerTestEnv.GetApiAddr()
		_, client, err := queryhandler.NewUserClient(ctx, addr)
		Expect(err).ToNot(HaveOccurred())
		return client
	}

	tenantServiceClient := func() domainApi.TenantClient {
		addr := testEnv.queryHandlerTestEnv.GetApiAddr()
		_, client, err := queryhandler.NewTenantClient(ctx, addr)
		Expect(err).ToNot(HaveOccurred())
		return client
	}

	clusterServiceClient := func() domainApi.ClusterClient {
		addr := testEnv.queryHandlerTestEnv.GetApiAddr()
		_, client, err := queryhandler.NewClusterClient(ctx, addr)
		Expect(err).ToNot(HaveOccurred())
		return client
	}

	certificateServiceClient := func() domainApi.CertificateClient {
		addr := testEnv.queryHandlerTestEnv.GetApiAddr()
		_, client, err := queryhandler.NewCertificateClient(ctx, addr)
		Expect(err).ToNot(HaveOccurred())
		return client
	}

	eventStoreClient := func() esApi.EventStoreClient {
		addr := testEnv.eventStoreTestEnv.GetApiAddr()
		_, client, err := eventstore.NewEventStoreClient(ctx, addr)
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

			reply, err := commandHandlerClient().Execute(mdManager.GetOutgoingGrpcContext(), command)
			Expect(err).ToNot(HaveOccurred())
			Expect(uuid.Nil).ToNot(Equal(reply.AggregateId))

			// update userId, as the "create" command will have changed it.
			userId := uuid.MustParse(reply.AggregateId)

			// Wait to propagate
			time.Sleep(1000 * time.Millisecond)

			user, err := userServiceClient().GetByEmail(ctx, wrapperspb.String("jane.doe@monoskope.io"))
			Expect(err).ToNot(HaveOccurred())
			Expect(user).ToNot(BeNil())
			Expect(user.GetEmail()).To(Equal("jane.doe@monoskope.io"))
			Expect(user.Id).To(Equal(userId.String()))

			userRoleBindingId := uuid.New()
			command, err = cmd.AddCommandData(
				cmd.CreateCommand(userRoleBindingId, commandTypes.CreateUserRoleBinding),
				&cmdData.CreateUserRoleBindingCommandData{Role: roles.Admin.String(), Scope: scopes.System.String(), UserId: userId.String()},
			)
			Expect(err).ToNot(HaveOccurred())

			reply, err = commandHandlerClient().Execute(mdManager.GetOutgoingGrpcContext(), command)
			Expect(err).ToNot(HaveOccurred())
			Expect(userRoleBindingId.String()).ToNot(Equal(reply.AggregateId))

			// update userRolebBindingId, as the "create" command will have changed it.
			userRoleBindingId = uuid.MustParse(reply.AggregateId)

			// Wait to propagate
			time.Sleep(1000 * time.Millisecond)

			// Creating the same rolebinding again should fail
			command.Id = uuid.New().String()
			_, err = commandHandlerClient().Execute(mdManager.GetOutgoingGrpcContext(), command)
			Expect(err).To(HaveOccurred())

			user, err = userServiceClient().GetByEmail(ctx, wrapperspb.String("jane.doe@monoskope.io"))
			Expect(err).ToNot(HaveOccurred())
			Expect(user).ToNot(BeNil())
			Expect(user.Roles[0].Role).To(Equal(roles.Admin.String()))
			Expect(user.Roles[0].Scope).To(Equal(scopes.System.String()))

			_, err = commandHandlerClient().Execute(mdManager.GetOutgoingGrpcContext(), cmd.CreateCommand(userRoleBindingId, commandTypes.DeleteUserRoleBinding))
			Expect(err).ToNot(HaveOccurred())

			// Wait to propagate
			time.Sleep(1000 * time.Millisecond)

			user, err = userServiceClient().GetByEmail(ctx, wrapperspb.String("jane.doe@monoskope.io"))
			Expect(err).ToNot(HaveOccurred())
			Expect(user).ToNot(BeNil())
		})
		It("can accept Nil as an Id when creating a user", func() {
			command, err := cmd.AddCommandData(
				cmd.CreateCommand(uuid.Nil, commandTypes.CreateUser),
				&cmdData.CreateUserCommandData{Name: "Jane Doe", Email: "jane.doe2@monoskope.io"},
			)
			Expect(err).ToNot(HaveOccurred())

			reply, err := commandHandlerClient().Execute(mdManager.GetOutgoingGrpcContext(), command)
			Expect(err).ToNot(HaveOccurred())
			Expect(reply.AggregateId).ToNot(Equal(uuid.Nil.String()))
		})
		It("fail to create a user which already exists", func() {
			command, err := cmd.AddCommandData(
				cmd.CreateCommand(uuid.New(), commandTypes.CreateUser),
				&cmdData.CreateUserCommandData{Name: "admin", Email: "admin@monoskope.io"},
			)
			Expect(err).ToNot(HaveOccurred())

			_, err = commandHandlerClient().Execute(mdManager.GetOutgoingGrpcContext(), command)
			Expect(err).To(HaveOccurred())
			Expect(errors.TranslateFromGrpcError(err)).To(Equal(errors.ErrUserAlreadyExists))
		})
	})

	Context("tenant management", func() {
		It("can manage a tenant", func() {
			user, err := userServiceClient().GetByEmail(ctx, wrapperspb.String("admin@monoskope.io"))
			Expect(err).ToNot(HaveOccurred())

			tenantId := uuid.New()
			command, err := cmd.AddCommandData(
				cmd.CreateCommand(tenantId, commandTypes.CreateTenant),
				&cmdData.CreateTenantCommandData{Name: "Tenant X", Prefix: "tx"},
			)
			Expect(err).ToNot(HaveOccurred())

			reply, err := commandHandlerClient().Execute(mdManager.GetOutgoingGrpcContext(), command)
			Expect(err).ToNot(HaveOccurred())
			Expect(tenantId.String()).ToNot(Equal(reply.AggregateId))

			// update tenantId, as the "create" command will have changed it.
			tenantId = uuid.MustParse(reply.AggregateId)

			// Wait to propagate
			time.Sleep(1000 * time.Millisecond)

			tenant, err := tenantServiceClient().GetByName(ctx, wrapperspb.String("Tenant X"))
			Expect(err).ToNot(HaveOccurred())
			Expect(tenant).ToNot(BeNil())
			Expect(tenant.GetName()).To(Equal("Tenant X"))
			Expect(tenant.GetPrefix()).To(Equal("tx"))
			Expect(tenant.Id).To(Equal(tenantId.String()))

			command, err = cmd.AddCommandData(
				cmd.CreateCommand(tenantId, commandTypes.UpdateTenant),
				&cmdData.UpdateTenantCommandData{Name: &wrapperspb.StringValue{Value: "DIIIETER"}},
			)
			Expect(err).ToNot(HaveOccurred())

			_, err = commandHandlerClient().Execute(mdManager.GetOutgoingGrpcContext(), command)
			Expect(err).ToNot(HaveOccurred())

			// Wait to propagate
			time.Sleep(1000 * time.Millisecond)

			tenant, err = tenantServiceClient().GetByName(ctx, wrapperspb.String("DIIIETER"))
			Expect(err).ToNot(HaveOccurred())
			Expect(tenant).ToNot(BeNil())
			Expect(tenant.Metadata.GetLastModifiedById()).To(Equal(user.Id))

			_, err = commandHandlerClient().Execute(mdManager.GetOutgoingGrpcContext(), cmd.CreateCommand(tenantId, commandTypes.DeleteTenant))
			Expect(err).ToNot(HaveOccurred())

			// Wait to propagate
			time.Sleep(1000 * time.Millisecond)

			tenant, err = tenantServiceClient().GetByName(ctx, wrapperspb.String("DIIIETER"))
			Expect(err).ToNot(HaveOccurred())
			Expect(tenant).ToNot(BeNil())
			Expect(tenant.Metadata.GetDeletedById()).To(Equal(user.GetId()))

			Expect(tenant.Metadata.Created).NotTo(BeNil())

		})
		It("can accept Nil as ID when creating a tenant", func() {
			command, err := cmd.AddCommandData(
				cmd.CreateCommand(uuid.Nil, commandTypes.CreateTenant),
				&cmdData.CreateTenantCommandData{Name: "Tenant X", Prefix: "tx"},
			)
			Expect(err).ToNot(HaveOccurred())

			reply, err := commandHandlerClient().Execute(mdManager.GetOutgoingGrpcContext(), command)
			Expect(err).ToNot(HaveOccurred())
			Expect(reply.AggregateId).ToNot(Equal(uuid.Nil.String()))
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
			testReactor, err := testReactor.NewTestReactor()
			Expect(err).ToNot(HaveOccurred())
			err = testReactor.Setup(ctx, testEnv.eventStoreTestEnv, eventStoreClient())
			Expect(err).ToNot(HaveOccurred())

			reply, err := commandHandlerClient().Execute(mdManager.GetOutgoingGrpcContext(), command)
			Expect(err).ToNot(HaveOccurred())
			Expect(uuid.Nil).ToNot(Equal(reply.AggregateId))

			// update clusterId, as the "create" command will have changed it.
			clusterId := uuid.MustParse(reply.AggregateId)

			// Wait to propagate
			time.Sleep(1000 * time.Millisecond)

			cluster, err := clusterServiceClient().GetByName(ctx, wrapperspb.String(expectedClusterName))
			Expect(err).ToNot(HaveOccurred())
			Expect(cluster).ToNot(BeNil())
			Expect(cluster.GetDisplayName()).To(Equal(expectedClusterDisplayName))
			Expect(cluster.GetName()).To(Equal(expectedClusterName))
			Expect(cluster.GetApiServerAddress()).To(Equal(expectedClusterApiServerAddress))
			Expect(cluster.GetCaCertBundle()).To(Equal(expectedClusterCACertBundle))

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

			time.Sleep(1 * time.Second)

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

			time.Sleep(3 * time.Second)

			tokenValue, err := clusterServiceClient().GetBootstrapToken(ctx, wrapperspb.String(clusterId.String()))
			Expect(err).ToNot(HaveOccurred())
			Expect(tokenValue.GetValue()).To(Equal("this is a valid JWT, honest!"))

		})
	})

	Context("cert management", func() {
		It("can create and query a certificate", func() {

			testReactor, err := testReactor.NewTestReactor()
			Expect(err).ToNot(HaveOccurred())
			err = testReactor.Setup(ctx, testEnv.eventStoreTestEnv, eventStoreClient())
			Expect(err).ToNot(HaveOccurred())

			clusterInfo, err := clusterServiceClient().GetByName(ctx, &wrapperspb.StringValue{Value: expectedClusterName})
			Expect(err).ToNot(HaveOccurred())
			Expect(clusterInfo).ToNot(BeNil())

			command, err := cmd.AddCommandData(
				cmd.CreateCommand(uuid.Nil, commandTypes.RequestCertificate),
				&cmdData.RequestCertificate{
					ReferencedAggregateId:   clusterInfo.Id,
					ReferencedAggregateType: aggregates.Cluster.String(),
					SigningRequest:          []byte("this is a CSR"),
				},
			)
			Expect(err).ToNot(HaveOccurred())

			certRequestCmdReply, err := commandHandlerClient().Execute(mdManager.GetOutgoingGrpcContext(), command)
			Expect(err).ToNot(HaveOccurred())

			// Wait to propagate
			time.Sleep(1000 * time.Millisecond)

			observed := testReactor.GetObservedEvents()
			Expect(len(observed)).ToNot(Equal(0))
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

			// Wait to propagate
			time.Sleep(1000 * time.Millisecond)

			certificate, err := certificateServiceClient().GetCertificate(ctx,
				&domainApi.GetCertificateRequest{
					AggregateId:   clusterInfo.GetId(),
					AggregateType: aggregates.Cluster.String(),
				})
			Expect(err).ToNot(HaveOccurred())
			Expect(certificate.GetCertificate()).ToNot(BeNil())
			Expect(certificate.GetCaCertBundle()).ToNot(BeNil())
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
