package domain

import (
	"context"
	"time"

	eventsourcingApi "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/eventsourcing"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/constants/aggregates"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/handler"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/projectors"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/repositories"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing/eventhandler"
	esr "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing/repositories"
)

type QueryHandlerDomain struct {
	UserRoleBindingRepository repositories.UserRoleBindingRepository
	UserRepository            repositories.UserRepository
	TenantRepository          repositories.TenantRepository
	TenantUserRepository      repositories.ReadOnlyTenantUserRepository
	ClusterRepository         repositories.ClusterRepository
	CertificateRepository     repositories.CertificateRepository
}

func NewQueryHandlerDomain(ctx context.Context, eventBus eventsourcing.EventBusConsumer, esClient eventsourcingApi.EventStoreClient) (*QueryHandlerDomain, error) {
	d := &QueryHandlerDomain{}

	// Setup repositories
	d.UserRoleBindingRepository = repositories.NewUserRoleBindingRepository(esr.NewInMemoryRepository())
	d.UserRepository = repositories.NewUserRepository(esr.NewInMemoryRepository(), d.UserRoleBindingRepository)
	d.TenantRepository = repositories.NewTenantRepository(esr.NewInMemoryRepository())
	d.TenantUserRepository = repositories.NewTenantUserRepository(d.UserRepository, d.UserRoleBindingRepository)
	d.CertificateRepository = repositories.NewCertificateRepository(esr.NewInMemoryRepository())
	d.ClusterRepository = repositories.NewClusterRepository(esr.NewInMemoryRepository())

	// Setup projectors
	userProjector := projectors.NewUserProjector()
	userRoleBindingProjector := projectors.NewUserRoleBindingProjector()
	tenantProjector := projectors.NewTenantProjector()
	clusterProjector := projectors.NewClusterProjector()
	certificateProjector := projectors.NewCertificateProjector()

	// Setup handler
	userProjectingHandler := eventhandler.NewProjectingEventHandler(userProjector, d.UserRepository)
	tenantProjectingHandler := eventhandler.NewProjectingEventHandler(tenantProjector, d.TenantRepository)
	userRoleBindingProjectingHandler := eventhandler.NewProjectingEventHandler(userRoleBindingProjector, d.UserRoleBindingRepository)
	clusterProjectingHandler := eventhandler.NewProjectingEventHandler(clusterProjector, d.ClusterRepository)
	certificateProjectHandler := eventhandler.NewProjectingEventHandler(certificateProjector, d.CertificateRepository)

	// Setup middleware
	refreshDuration := time.Second * 30
	userHandlerChain := eventsourcing.UseEventHandlerMiddleware(userProjectingHandler, eventhandler.NewEventStoreReplayMiddleware(esClient), eventhandler.NewEventStoreRefreshMiddleware(esClient, refreshDuration))
	userRoleBindingHandlerChain := eventsourcing.UseEventHandlerMiddleware(userRoleBindingProjectingHandler, eventhandler.NewEventStoreReplayMiddleware(esClient), eventhandler.NewEventStoreRefreshMiddleware(esClient, refreshDuration))
	tenantHandlerChain := eventsourcing.UseEventHandlerMiddleware(tenantProjectingHandler, eventhandler.NewEventStoreReplayMiddleware(esClient), eventhandler.NewEventStoreRefreshMiddleware(esClient, refreshDuration))
	clusterHandlerChain := eventsourcing.UseEventHandlerMiddleware(clusterProjectingHandler, eventhandler.NewEventStoreReplayMiddleware(esClient), eventhandler.NewEventStoreRefreshMiddleware(esClient, refreshDuration))
	certificateHandlerChain := eventsourcing.UseEventHandlerMiddleware(certificateProjectHandler, eventhandler.NewEventStoreReplayMiddleware(esClient), eventhandler.NewEventStoreRefreshMiddleware(esClient, refreshDuration))

	// Setup matcher for event bus
	userMatcher := eventBus.Matcher().MatchAggregateType(aggregates.User)
	tenantMatcher := eventBus.Matcher().MatchAggregateType(aggregates.Tenant)
	userRoleBindingMatcher := eventBus.Matcher().MatchAggregateType(aggregates.UserRoleBinding)
	clusterMatcher := eventBus.Matcher().MatchAggregateType(aggregates.Cluster)
	certificateMatcher := eventBus.Matcher().MatchAggregateType(aggregates.Certificate)

	// Register event handler with event bus
	if err := eventBus.AddHandler(ctx, userHandlerChain, userMatcher); err != nil {
		return nil, err
	}
	if err := eventBus.AddHandler(ctx, tenantHandlerChain, tenantMatcher); err != nil {
		return nil, err
	}
	if err := eventBus.AddHandler(ctx, userRoleBindingHandlerChain, userRoleBindingMatcher); err != nil {
		return nil, err
	}
	if err := eventBus.AddHandler(ctx, clusterHandlerChain, clusterMatcher); err != nil {
		return nil, err
	}
	if err := eventBus.AddHandler(ctx, certificateHandlerChain, certificateMatcher); err != nil {
		return nil, err
	}

	// Start repo warming
	if err := handler.WarmUp(ctx, esClient, aggregates.User, userHandlerChain); err != nil {
		return nil, err
	}
	if err := handler.WarmUp(ctx, esClient, aggregates.UserRoleBinding, userRoleBindingHandlerChain); err != nil {
		return nil, err
	}
	if err := handler.WarmUp(ctx, esClient, aggregates.Tenant, tenantHandlerChain); err != nil {
		return nil, err
	}
	if err := handler.WarmUp(ctx, esClient, aggregates.Cluster, clusterHandlerChain); err != nil {
		return nil, err
	}
	if err := handler.WarmUp(ctx, esClient, aggregates.Certificate, certificateHandlerChain); err != nil {
		return nil, err
	}

	return d, nil
}
