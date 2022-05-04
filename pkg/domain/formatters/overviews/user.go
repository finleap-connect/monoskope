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
	"fmt"
	esApi "github.com/finleap-connect/monoskope/pkg/api/eventsourcing"
	"github.com/finleap-connect/monoskope/pkg/audit/formatters"
	"github.com/finleap-connect/monoskope/pkg/domain/projections"
	"github.com/finleap-connect/monoskope/pkg/domain/projectors"
	esErrors "github.com/finleap-connect/monoskope/pkg/eventsourcing/errors"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"strings"
	"time"
)

// userOverviewFormatter OverviewFormatter implementation for the user-aggregate
type userOverviewFormatter struct {
	*formatters.FormatterBase
}

// NewUserOverviewFormatter creates a new overview formatter for the user-aggregate
func NewUserOverviewFormatter(esClient esApi.EventStoreClient) *userOverviewFormatter {
	return &userOverviewFormatter{FormatterBase: &formatters.FormatterBase{EsClient: esClient}}
}

// GetFormattedDetails returns the user overview details in a human-readable format.
// The given timestamp is used to create the snapshots of the needed aggregates
func (f *userOverviewFormatter) GetFormattedDetails(ctx context.Context, user *projections.User, timestamp time.Time) (string, error) {
	var details string
	eventFilter := &esApi.EventFilter{MaxTimestamp: timestamppb.New(timestamp)}

	eventFilter.AggregateId = wrapperspb.String(user.CreatedById)
	creatorSnapshot, err := f.CreateSnapshot(ctx, projectors.NewUserProjector(), eventFilter)
	if err != nil {
		return "", err
	}
	creator, ok := creatorSnapshot.(*projections.User)
	if !ok {
		return "", esErrors.ErrInvalidProjectionType
	}

	details += fmt.Sprintf("“%s“ was created by “%s“ at “%s“", user.Email, creator.Email, user.Created.AsTime().Format(time.RFC822))

	if len(user.DeletedById) == 0 {
		return details, nil
	}

	eventFilter.AggregateId = wrapperspb.String(user.DeletedById)
	deleterSnapshot, err := f.CreateSnapshot(ctx, projectors.NewUserProjector(), eventFilter)
	if err != nil {
		return "", err
	}
	deleter, ok := deleterSnapshot.(*projections.User)
	if !ok {
		return "", esErrors.ErrInvalidProjectionType
	}

	details += fmt.Sprintf(" and was deleted by “%s“ at “%s“", deleter.Email, user.Deleted.AsTime().Format(time.RFC822))

	return details, nil
}

// GetRolesDetails returns the roles details in general and the tenants and clusters details to which the roles apply.
// The given timestamp is used to create the snapshots of the needed aggregates
func (f *userOverviewFormatter) GetRolesDetails(ctx context.Context, user *projections.User, timestamp time.Time) (string, string, string, error) {
	var rolesDetails, tenantsDetails, clustersDetails string
	eventFilter := &esApi.EventFilter{MaxTimestamp: timestamppb.New(timestamp)}

	for _, role := range user.Roles {
		roleDetails := fmt.Sprintf("- %s %s\n", role.Scope, role.Role)
		if !strings.Contains(rolesDetails, roleDetails) { // avoid having the same role multiple times
			rolesDetails += roleDetails
		}

		if len(role.Resource) == 0 {
			continue
		}

		eventFilter.AggregateId = wrapperspb.String(role.Resource)
		tenantSnapshot, err := f.CreateSnapshot(ctx, projectors.NewTenantProjector(), eventFilter)
		if err == nil {
			tenant, ok := tenantSnapshot.(*projections.Tenant)
			if ok {
				tenantsDetails += fmt.Sprintf("- %s (%s)\n", tenant.Name, role.Role)
				continue // it's either a tenant or cluster
			}
		}
		clusterSnapshot, err := f.CreateSnapshot(ctx, projectors.NewClusterProjector(), eventFilter)
		if err == nil {
			cluster, ok := clusterSnapshot.(*projections.Cluster)
			if ok {
				clustersDetails += fmt.Sprintf("- %s (%s)\n", cluster.DisplayName, role.Role)
			}
		}
	}

	return strings.TrimSpace(rolesDetails), strings.TrimSpace(tenantsDetails), strings.TrimSpace(clustersDetails), nil
}
