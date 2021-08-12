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
}

func NewQueryHandlerDomain(ctx context.Context, eventBus eventsourcing.EventBusConsumer, esClient eventsourcingApi.EventStoreClient) (*QueryHandlerDomain, error) {
	d := &QueryHandlerDomain{}

	// Setup repositories
	d.UserRoleBindingRepository = repositories.NewUserRoleBindingRepository(esr.NewInMemoryRepository())
	d.UserRepository = repositories.NewUserRepository(esr.NewInMemoryRepository(), d.UserRoleBindingRepository)
	d.TenantRepository = repositories.NewTenantRepository(esr.NewInMemoryRepository())
	d.TenantUserRepository = repositories.NewTenantUserRepository(d.UserRepository, d.UserRoleBindingRepository)
	d.ClusterRepository = repositories.NewClusterRepository(esr.NewInMemoryRepository())

	// Setup projectors
	userProjector := projectors.NewUserProjector()
	userRoleBindingProjector := projectors.NewUserRoleBindingProjector()
	tenantProjector := projectors.NewTenantProjector()
	clusterProjector := projectors.NewClusterProjector()

	// Setup handler
	userProjectingHandler := eventhandler.NewProjectingEventHandler(userProjector, d.UserRepository)
	tenantProjectingHandler := eventhandler.NewProjectingEventHandler(tenantProjector, d.TenantRepository)
	userRoleBindingProjectingHandler := eventhandler.NewProjectingEventHandler(userRoleBindingProjector, d.UserRoleBindingRepository)
	clusterProjectingHandler := eventhandler.NewProjectingEventHandler(clusterProjector, d.ClusterRepository)

	// Setup middleware
	replayHandler := eventhandler.NewEventStoreReplayMiddleware(esClient)
	//
	userHandlerChain := eventsourcing.UseEventHandlerMiddleware(userProjectingHandler, replayHandler, eventhandler.NewEventStoreRefreshMiddleware(esClient, time.Second*30, aggregates.User))
	userRoleBindingHandlerChain := eventsourcing.UseEventHandlerMiddleware(userRoleBindingProjectingHandler, replayHandler, eventhandler.NewEventStoreRefreshMiddleware(esClient, time.Second*30, aggregates.UserRoleBinding))
	tenantHandlerChain := eventsourcing.UseEventHandlerMiddleware(tenantProjectingHandler, replayHandler, eventhandler.NewEventStoreRefreshMiddleware(esClient, time.Second*30, aggregates.Tenant))
	clusterHandlerChain := eventsourcing.UseEventHandlerMiddleware(clusterProjectingHandler, replayHandler, eventhandler.NewEventStoreRefreshMiddleware(esClient, time.Second*30, aggregates.Cluster))

	// Setup matcher for event bus
	userMatcher := eventBus.Matcher().MatchAggregateType(aggregates.User)
	tenantMatcher := eventBus.Matcher().MatchAggregateType(aggregates.Tenant)
	userRoleBindingMatcher := eventBus.Matcher().MatchAggregateType(aggregates.UserRoleBinding)
	clusterMatcher := eventBus.Matcher().MatchAggregateType(aggregates.Cluster)

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

	return d, nil
}
