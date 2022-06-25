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

package overviews

import (
	"context"
	"github.com/finleap-connect/monoskope/pkg/domain/constants/users"
	"github.com/google/uuid"
	"strings"
	"time"

	esApi "github.com/finleap-connect/monoskope/pkg/api/eventsourcing"
	"github.com/finleap-connect/monoskope/pkg/audit/formatters"
	fConsts "github.com/finleap-connect/monoskope/pkg/domain/constants/formatters"
	"github.com/finleap-connect/monoskope/pkg/domain/projections"
	"github.com/finleap-connect/monoskope/pkg/domain/projectors"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

// userOverviewFormatter OverviewFormatter implementation for the user-aggregate
type userOverviewFormatter struct {
	esClient esApi.EventStoreClient
}

// NewUserOverviewFormatter creates a new overview formatter for the user-aggregate
func NewUserOverviewFormatter(esClient esApi.EventStoreClient) *userOverviewFormatter {
	return &userOverviewFormatter{esClient}
}

// GetFormattedDetails returns the user overview details in a human-readable format.
// The given timestamp is used to create the snapshots of the needed aggregates
func (f *userOverviewFormatter) GetFormattedDetails(ctx context.Context, user *projections.User, timestamp time.Time) (string, error) {
	var details string
	eventFilter := &esApi.EventFilter{MaxTimestamp: timestamppb.New(timestamp)}
	userProjector := projectors.NewUserProjector()
	snapshotter := formatters.NewSnapshotter(f.esClient, userProjector)

	eventFilter.AggregateId = wrapperspb.String(user.GetCreatedById())
	creator, err := snapshotter.CreateSnapshot(ctx, eventFilter)
	if err != nil { // possible if user was created by a system user that was created manually (no events)
		creator = userProjector.NewProjection(uuid.MustParse(eventFilter.AggregateId.Value))
		creator.Email = "system@" + users.BASE_DOMAIN
	}
	details += fConsts.UserCreatedOverviewDetailsFormat.Sprint(user.Email, creator.Email, user.GetCreated().AsTime().Format(fConsts.TimeFormat))

	if len(user.GetDeletedById()) == 0 {
		return details, nil
	}

	eventFilter.AggregateId = wrapperspb.String(user.GetDeletedById())
	deleter, err := snapshotter.CreateSnapshot(ctx, eventFilter)
	if err != nil {
		deleter = userProjector.NewProjection(uuid.MustParse(eventFilter.AggregateId.Value))
		deleter.Email = "system@" + users.BASE_DOMAIN
	}
	details += fConsts.UserDeletedOverviewDetailsFormat.Sprint(deleter.Email, user.GetDeleted().AsTime().Format(fConsts.TimeFormat))

	return details, nil
}

// GetRolesDetails returns the roles details in general and the tenants and clusters details to which the roles apply.
// The given timestamp is used to create the snapshots of the needed aggregates
func (f *userOverviewFormatter) GetRolesDetails(ctx context.Context, user *projections.User, timestamp time.Time) (string, string, string, error) {
	var rolesDetails, tenantsDetails, clustersDetails string
	eventFilter := &esApi.EventFilter{MaxTimestamp: timestamppb.New(timestamp)}

	for _, role := range user.Roles {
		if role.Metadata.Deleted != nil {
			continue
		}

		roleDetails := fConsts.UserRoleBindingOverviewDetailsFormat.Sprint(role.Scope, role.Role)
		if !strings.Contains(rolesDetails, roleDetails) { // avoid having the same role multiple times
			rolesDetails += roleDetails
		}

		if len(role.Resource) == 0 {
			continue
		}

		eventFilter.AggregateId = wrapperspb.String(role.Resource)

		tenantSnapshotter := formatters.NewSnapshotter(f.esClient, projectors.NewTenantProjector())
		tenant, err := tenantSnapshotter.CreateSnapshot(ctx, eventFilter)
		if err == nil {
			tenantsDetails += fConsts.TenantUserRoleBindingOverviewDetailsFormat.Sprint(tenant.Name, role.Role)
			continue // it's either a tenant or cluster
		}

		clusterSnapshotter := formatters.NewSnapshotter(f.esClient, projectors.NewClusterProjector())
		cluster, err := clusterSnapshotter.CreateSnapshot(ctx, eventFilter)
		if err == nil {
			clustersDetails += fConsts.ClusterUserRoleBindingOverviewDetailsFormat.Sprint(cluster.DisplayName, role.Role)
		}
	}

	return strings.TrimSpace(rolesDetails), strings.TrimSpace(tenantsDetails), strings.TrimSpace(clustersDetails), nil
}
