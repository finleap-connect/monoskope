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
	"github.com/finleap-connect/monoskope/pkg/domain/projectors"
	"github.com/finleap-connect/monoskope/pkg/domain/repositories"
	"github.com/finleap-connect/monoskope/pkg/eventsourcing"
	"github.com/finleap-connect/monoskope/pkg/eventsourcing/eventhandler"
	esr "github.com/finleap-connect/monoskope/pkg/eventsourcing/repositories"
)

type QueryHandlerDomain struct {
	UserRoleBindingRepository      repositories.UserRoleBindingRepository
	UserRepository                 repositories.UserRepository
	TenantRepository               repositories.TenantRepository
	TenantUserRepository           repositories.ReadOnlyTenantUserRepository
	ClusterRepository              repositories.ClusterRepository
	CertificateRepository          repositories.CertificateRepository
	TenantClusterBindingRepository repositories.TenantClusterBindingRepository
	ClusterAccessRepo              repositories.ReadOnlyClusterAccessRepository
}

func NewQueryHandlerDomain(ctx context.Context, eventBus eventsourcing.EventBusConsumer, esClient eventsourcingApi.EventStoreClient) (*QueryHandlerDomain, error) {
	d := new(QueryHandlerDomain)

	// Setup repositories
	d.UserRoleBindingRepository = repositories.NewUserRoleBindingRepository(esr.NewInMemoryRepository())
	d.UserRepository = repositories.NewUserRepository(esr.NewInMemoryRepository(), d.UserRoleBindingRepository)
	d.TenantRepository = repositories.NewTenantRepository(esr.NewInMemoryRepository())
	d.TenantUserRepository = repositories.NewTenantUserRepository(d.UserRepository, d.UserRoleBindingRepository)
	d.ClusterRepository = repositories.NewClusterRepository(esr.NewInMemoryRepository())
	d.CertificateRepository = repositories.NewCertificateRepository(esr.NewInMemoryRepository())
	d.TenantClusterBindingRepository = repositories.NewTenantClusterBindingRepository(esr.NewInMemoryRepository())
	d.ClusterAccessRepo = repositories.NewClusterAccessRepository(d.TenantClusterBindingRepository, d.ClusterRepository, d.UserRoleBindingRepository)

	// Setup projectors
	userProjector := projectors.NewUserProjector()
	userRoleBindingProjector := projectors.NewUserRoleBindingProjector()
	tenantProjector := projectors.NewTenantProjector()
	clusterProjector := projectors.NewClusterProjector()
	certificateProjector := projectors.NewCertificateProjector()
	tenantClusterBindingProjector := projectors.NewTenantClusterBindingProjector()

	// Setup handler
	userProjectingHandler := eventhandler.NewProjectingEventHandler(userProjector, d.UserRepository)
	tenantProjectingHandler := eventhandler.NewProjectingEventHandler(tenantProjector, d.TenantRepository)
	userRoleBindingProjectingHandler := eventhandler.NewProjectingEventHandler(userRoleBindingProjector, d.UserRoleBindingRepository)
	clusterProjectingHandler := eventhandler.NewProjectingEventHandler(clusterProjector, d.ClusterRepository)
	certificateProjectingHandler := eventhandler.NewProjectingEventHandler(certificateProjector, d.CertificateRepository)
	tenantClusterBindingProjectingHandler := eventhandler.NewProjectingEventHandler(tenantClusterBindingProjector, d.TenantClusterBindingRepository)

	// Setup middleware
	refreshDuration := time.Second * 30
	userHandlerChain := eventsourcing.UseEventHandlerMiddleware(userProjectingHandler, eventhandler.NewEventStoreReplayMiddleware(esClient), eventhandler.NewEventStoreRefreshMiddleware(esClient, refreshDuration))
	userRoleBindingHandlerChain := eventsourcing.UseEventHandlerMiddleware(userRoleBindingProjectingHandler, eventhandler.NewEventStoreReplayMiddleware(esClient), eventhandler.NewEventStoreRefreshMiddleware(esClient, refreshDuration))
	tenantHandlerChain := eventsourcing.UseEventHandlerMiddleware(tenantProjectingHandler, eventhandler.NewEventStoreReplayMiddleware(esClient), eventhandler.NewEventStoreRefreshMiddleware(esClient, refreshDuration))
	clusterHandlerChain := eventsourcing.UseEventHandlerMiddleware(clusterProjectingHandler, eventhandler.NewEventStoreReplayMiddleware(esClient), eventhandler.NewEventStoreRefreshMiddleware(esClient, refreshDuration))
	certificateHandlerChain := eventsourcing.UseEventHandlerMiddleware(certificateProjectingHandler, eventhandler.NewEventStoreReplayMiddleware(esClient), eventhandler.NewEventStoreRefreshMiddleware(esClient, refreshDuration))
	tenantClusterBindingHandlerChain := eventsourcing.UseEventHandlerMiddleware(tenantClusterBindingProjectingHandler, eventhandler.NewEventStoreReplayMiddleware(esClient), eventhandler.NewEventStoreRefreshMiddleware(esClient, refreshDuration))

	// Setup matcher for event bus
	userMatcher := eventBus.Matcher().MatchAggregateType(aggregates.User)
	tenantMatcher := eventBus.Matcher().MatchAggregateType(aggregates.Tenant)
	userRoleBindingMatcher := eventBus.Matcher().MatchAggregateType(aggregates.UserRoleBinding)
	clusterMatcher := eventBus.Matcher().MatchAggregateType(aggregates.Cluster)
	certificateMatcher := eventBus.Matcher().MatchAggregateType(aggregates.Certificate)
	tenantClusterBindingMatcher := eventBus.Matcher().MatchAggregateType(aggregates.TenantClusterBinding)

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
	if err := handler.WarmUp(ctx, esClient, aggregates.Certificate, certificateHandlerChain); err != nil {
		return nil, err
	}
	if err := handler.WarmUp(ctx, esClient, aggregates.TenantClusterBinding, tenantClusterBindingHandlerChain); err != nil {
		return nil, err
	}

	return d, nil
}
