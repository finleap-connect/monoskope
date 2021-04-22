package domain

import (
	"context"

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
	UserRepository                repositories.UserRepository
	TenantRepository              repositories.TenantRepository
	ClusterRegistrationRepository repositories.ClusterRegistrationRepository
}

func NewQueryHandlerDomain(ctx context.Context, messageBusConsumer eventsourcing.EventBusConsumer, esClient eventsourcingApi.EventStoreClient) (*QueryHandlerDomain, error) {
	qhDomain := &QueryHandlerDomain{}

	// Setup repositories
	userRoleBindings := repositories.NewUserRoleBindingRepository(esr.NewInMemoryRepository())
	qhDomain.UserRepository = repositories.NewUserRepository(esr.NewInMemoryRepository(), userRoleBindings)
	qhDomain.TenantRepository = repositories.NewTenantRepository(esr.NewInMemoryRepository(), qhDomain.UserRepository)
	qhDomain.ClusterRegistrationRepository = repositories.NewClusterRegistrationRepository(esr.NewInMemoryRepository(), qhDomain.UserRepository)

	// Setup event handler and middleware
	userEventHandler := eventsourcing.UseEventHandlerMiddleware(
		eventhandler.NewProjectingEventHandler(
			projectors.NewUserProjector(),
			qhDomain.UserRepository,
		),
		eventhandler.NewEventStoreReplayMiddleware(esClient).Middleware,
	)

	// Register event handler with message bus
	err := messageBusConsumer.AddHandler(ctx,
		userEventHandler,
		messageBusConsumer.Matcher().MatchAggregateType(aggregates.User),
	)
	if err != nil {
		return nil, err
	}

	// Setup event handler and middleware
	userRoleBindingEventHandler := eventsourcing.UseEventHandlerMiddleware(
		eventhandler.NewProjectingEventHandler(
			projectors.NewUserRoleBindingProjector(),
			userRoleBindings,
		),
		eventhandler.NewEventStoreReplayMiddleware(esClient).Middleware,
	)

	// Register event handler with message bus
	err = messageBusConsumer.AddHandler(ctx,
		userRoleBindingEventHandler,
		messageBusConsumer.Matcher().MatchAggregateType(aggregates.UserRoleBinding),
	)
	if err != nil {
		return nil, err
	}

	// Setup event handler and middleware
	tenantEventHandler := eventsourcing.UseEventHandlerMiddleware(
		eventhandler.NewProjectingEventHandler(
			projectors.NewTenantProjector(),
			qhDomain.TenantRepository,
		),
		eventhandler.NewEventStoreReplayMiddleware(esClient).Middleware,
	)

	// Register event handler with message bus
	err = messageBusConsumer.AddHandler(ctx,
		tenantEventHandler,
		messageBusConsumer.Matcher().MatchAggregateType(aggregates.Tenant),
	)
	if err != nil {
		return nil, err
	}

	// Start repo warming
	userRepoWarming := handler.NewRepoWarmingMiddleware(esClient, aggregates.User, userEventHandler)
	if err := userRepoWarming.WarmUp(ctx); err != nil {
		return nil, err
	}
	userRoleBindingRepoWarming := handler.NewRepoWarmingMiddleware(esClient, aggregates.UserRoleBinding, userRoleBindingEventHandler)
	if err := userRoleBindingRepoWarming.WarmUp(ctx); err != nil {
		return nil, err
	}
	tenantRepoWarming := handler.NewRepoWarmingMiddleware(esClient, aggregates.Tenant, tenantEventHandler)
	if err := tenantRepoWarming.WarmUp(ctx); err != nil {
		return nil, err
	}

	return qhDomain, nil
}
