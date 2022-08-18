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

package domain

import (
	"context"
	"time"

	eventsourcingApi "github.com/finleap-connect/monoskope/pkg/api/eventsourcing"
	"github.com/finleap-connect/monoskope/pkg/domain/constants/aggregates"
	"github.com/finleap-connect/monoskope/pkg/domain/handler"
	"github.com/finleap-connect/monoskope/pkg/domain/projections"
	"github.com/finleap-connect/monoskope/pkg/domain/projectors"
	"github.com/finleap-connect/monoskope/pkg/domain/repositories"
	"github.com/finleap-connect/monoskope/pkg/eventsourcing"
	"github.com/finleap-connect/monoskope/pkg/eventsourcing/eventhandler"
	esr "github.com/finleap-connect/monoskope/pkg/eventsourcing/repositories"
)

type GatewayDomain struct {
	UserRoleBindingRepository      repositories.UserRoleBindingRepository
	UserRepository                 repositories.UserRepository
	TenantRepository               repositories.TenantRepository
	TenantUserRepository           repositories.TenantUserRepository
	ClusterRepository              repositories.ClusterRepository
	TenantClusterBindingRepository repositories.TenantClusterBindingRepository
	ClusterAccessRepo              repositories.ClusterAccessRepository
}

func NewGatewayDomain(ctx context.Context, eventBus eventsourcing.EventBusConsumer, esClient eventsourcingApi.EventStoreClient) (*GatewayDomain, error) {
	d := new(GatewayDomain)

	// Setup repositories
	d.UserRoleBindingRepository = repositories.NewUserRoleBindingRepository(esr.NewInMemoryRepository[*projections.UserRoleBinding]())
	d.UserRepository = repositories.NewUserRepository(esr.NewInMemoryRepository[*projections.User](), d.UserRoleBindingRepository)
	d.TenantRepository = repositories.NewTenantRepository(esr.NewInMemoryRepository[*projections.Tenant]())
	d.ClusterRepository = repositories.NewClusterRepository(esr.NewInMemoryRepository[*projections.Cluster]())
	d.TenantClusterBindingRepository = repositories.NewTenantClusterBindingRepository(esr.NewInMemoryRepository[*projections.TenantClusterBinding]())
	d.ClusterAccessRepo = repositories.NewClusterAccessRepository(d.TenantClusterBindingRepository, d.ClusterRepository, d.UserRoleBindingRepository, d.TenantRepository)

	// Setup projectors
	userProjector := projectors.NewUserProjector()
	userRoleBindingProjector := projectors.NewUserRoleBindingProjector()
	tenantProjector := projectors.NewTenantProjector()
	clusterProjector := projectors.NewClusterProjector()
	tenantClusterBindingProjector := projectors.NewTenantClusterBindingProjector()

	// Setup handler
	userProjectingHandler := eventhandler.NewProjectingEventHandler[*projections.User](userProjector, d.UserRepository)
	tenantProjectingHandler := eventhandler.NewProjectingEventHandler[*projections.Tenant](tenantProjector, d.TenantRepository)
	userRoleBindingProjectingHandler := eventhandler.NewProjectingEventHandler[*projections.UserRoleBinding](userRoleBindingProjector, d.UserRoleBindingRepository)
	clusterProjectingHandler := eventhandler.NewProjectingEventHandler[*projections.Cluster](clusterProjector, d.ClusterRepository)
	tenantClusterBindingProjectingHandler := eventhandler.NewProjectingEventHandler[*projections.TenantClusterBinding](tenantClusterBindingProjector, d.TenantClusterBindingRepository)

	// Setup middleware
	refreshDuration := time.Second * 30
	userHandlerChain := eventsourcing.UseEventHandlerMiddleware(userProjectingHandler, eventhandler.NewEventStoreReplayMiddleware(esClient), eventhandler.NewEventStoreRefreshMiddleware(esClient, refreshDuration))
	userRoleBindingHandlerChain := eventsourcing.UseEventHandlerMiddleware(userRoleBindingProjectingHandler, eventhandler.NewEventStoreReplayMiddleware(esClient), eventhandler.NewEventStoreRefreshMiddleware(esClient, refreshDuration))
	tenantHandlerChain := eventsourcing.UseEventHandlerMiddleware(tenantProjectingHandler, eventhandler.NewEventStoreReplayMiddleware(esClient), eventhandler.NewEventStoreRefreshMiddleware(esClient, refreshDuration))
	clusterHandlerChain := eventsourcing.UseEventHandlerMiddleware(clusterProjectingHandler, eventhandler.NewEventStoreReplayMiddleware(esClient), eventhandler.NewEventStoreRefreshMiddleware(esClient, refreshDuration))
	tenantClusterBindingHandlerChain := eventsourcing.UseEventHandlerMiddleware(tenantClusterBindingProjectingHandler, eventhandler.NewEventStoreReplayMiddleware(esClient), eventhandler.NewEventStoreRefreshMiddleware(esClient, refreshDuration))

	// Setup matcher for event bus
	userMatcher := eventBus.Matcher().MatchAggregateType(aggregates.User)
	userRoleBindingMatcher := eventBus.Matcher().MatchAggregateType(aggregates.UserRoleBinding)
	tenantMatcher := eventBus.Matcher().MatchAggregateType(aggregates.Tenant)
	clusterMatcher := eventBus.Matcher().MatchAggregateType(aggregates.Cluster)
	tenantClusterBindingMatcher := eventBus.Matcher().MatchAggregateType(aggregates.TenantClusterBinding)

	// Register event handler with event bus
	if err := eventBus.AddHandler(ctx, userHandlerChain, userMatcher); err != nil {
		return nil, err
	}
	if err := eventBus.AddHandler(ctx, userRoleBindingHandlerChain, userRoleBindingMatcher); err != nil {
		return nil, err
	}
	if err := eventBus.AddHandler(ctx, tenantHandlerChain, tenantMatcher); err != nil {
		return nil, err
	}
	if err := eventBus.AddHandler(ctx, clusterHandlerChain, clusterMatcher); err != nil {
		return nil, err
	}
	if err := eventBus.AddHandler(ctx, tenantClusterBindingHandlerChain, tenantClusterBindingMatcher); err != nil {
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
	if err := handler.WarmUp(ctx, esClient, aggregates.TenantClusterBinding, tenantClusterBindingHandlerChain); err != nil {
		return nil, err
	}

	return d, nil
}
