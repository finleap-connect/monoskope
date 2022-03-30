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

package overviews

import (
	"context"
	"fmt"
	"github.com/finleap-connect/monoskope/pkg/api/domain/audit"
	"github.com/finleap-connect/monoskope/pkg/domain/projections"
	"github.com/finleap-connect/monoskope/pkg/domain/repositories"
	"github.com/google/uuid"
	"time"
)

// userOverviewFormatter TODO EventFormatter implementation for the user-aggregate
type userOverviewFormatter struct {
	userRepo       repositories.ReadOnlyUserRepository
	tenantRepo     repositories.ReadOnlyTenantRepository
	clusterRepo    repositories.ReadOnlyClusterRepository
}

// NewUserOverviewFormatter TODO creates a new event formatter for the user-aggregate
func NewUserOverviewFormatter(userRepo repositories.ReadOnlyUserRepository, tenantRepo repositories.ReadOnlyTenantRepository, clusterRepo repositories.ReadOnlyClusterRepository) *userOverviewFormatter {
	return &userOverviewFormatter{
		userRepo: userRepo,
		tenantRepo: tenantRepo,
		clusterRepo: clusterRepo,
	}
}

// GetFormattedDetails TODO formats the user-aggregate-events in a human-readable format
func (f *userOverviewFormatter) GetFormattedDetails(ctx context.Context, user *projections.User) (string, error) {
	creator, err := f.userRepo.ByUserId(ctx, uuid.MustParse(user.CreatedById))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s was created by %s at %s", user.Email, creator.Email, user.Created.AsTime().Format(time.RFC822)), nil
}

func (f *userOverviewFormatter) AddRolesDetails(ctx context.Context, user *projections.User, overview *audit.UserOverview) error {
	for _, role := range user.Roles {
		overview.Roles += fmt.Sprintf("- %s %s\n", role.Scope, role.Role)
		if len(role.Resource) == 0 {
			continue
		}

		tenant, err := f.tenantRepo.ByTenantId(ctx, role.Resource)
		if err == nil {
			overview.Tenants += fmt.Sprintf("- %s (%s)\n", tenant.Name, role.Role)
			continue // it's either a tenant or cluster
		}
		cluster, err := f.clusterRepo.ByClusterId(ctx, role.Resource)
		if err == nil {
			overview.Clusters += fmt.Sprintf("- %s (%s)\n", cluster.DisplayName, role.Role)
		}
	}

	return nil
}