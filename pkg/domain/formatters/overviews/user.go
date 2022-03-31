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
	"github.com/finleap-connect/monoskope/pkg/domain/projections"
	"github.com/finleap-connect/monoskope/pkg/domain/repositories"
	"github.com/google/uuid"
	"strings"
	"time"
)

// userOverviewFormatter OverviewFormatter implementation for the user-aggregate
type userOverviewFormatter struct {
	userRepo    repositories.ReadOnlyUserRepository
	tenantRepo  repositories.ReadOnlyTenantRepository
	clusterRepo repositories.ReadOnlyClusterRepository
}

// NewUserOverviewFormatter creates a new overview formatter for the user-aggregate
func NewUserOverviewFormatter(userRepo repositories.ReadOnlyUserRepository, tenantRepo repositories.ReadOnlyTenantRepository, clusterRepo repositories.ReadOnlyClusterRepository) *userOverviewFormatter {
	return &userOverviewFormatter{
		userRepo:    userRepo,
		tenantRepo:  tenantRepo,
		clusterRepo: clusterRepo,
	}
}

// GetFormattedDetails formats the user-projection details in a human-readable format
func (f *userOverviewFormatter) GetFormattedDetails(ctx context.Context, user *projections.User) (string, error) {
	var details string

	id, err := uuid.Parse(user.CreatedById)
	if err != nil {
		return "", err
	}
	creator, err := f.userRepo.ByUserId(ctx, id)
	if err != nil {
		return "", err
	}
	details += fmt.Sprintf("“%s“ was created by “%s“ at “%s“", user.Email, creator.Email, user.Created.AsTime().Format(time.RFC822))

	if len(user.DeletedById) == 0 {
		return details, nil
	}

	id, err = uuid.Parse(user.DeletedById)
	if err != nil {
		return "", err
	}
	deleter, err := f.userRepo.ByUserId(ctx, id)
	if err != nil {
		return "", err
	}
	details += fmt.Sprintf(" and was deleted by “%s“ at “%s“", deleter.Email, user.Deleted.AsTime().Format(time.RFC822))

	return details, nil
}

// GetRolesDetails returns the roles details in general and the tenants and clusters details to which the roles apply
func (f *userOverviewFormatter) GetRolesDetails(ctx context.Context, user *projections.User) (string, string, string, error) {
	var rolesDetails, tenantsDetails, clustersDetails string
	for _, role := range user.Roles {
		roleDetails := fmt.Sprintf("- %s %s\n", role.Scope, role.Role)
		if !strings.Contains(rolesDetails, roleDetails) {
			rolesDetails += roleDetails
		}

		if len(role.Resource) == 0 {
			continue
		}

		tenant, err := f.tenantRepo.ByTenantId(ctx, role.Resource)
		if err == nil {
			tenantsDetails += fmt.Sprintf("- %s (%s)\n", tenant.Name, role.Role)
			continue // it's either a tenant or cluster
		}
		cluster, err := f.clusterRepo.ByClusterId(ctx, role.Resource)
		if err == nil {
			clustersDetails += fmt.Sprintf("- %s (%s)\n", cluster.DisplayName, role.Role)
		}
	}

	return strings.TrimSpace(rolesDetails), strings.TrimSpace(tenantsDetails), strings.TrimSpace(clustersDetails), nil
}
