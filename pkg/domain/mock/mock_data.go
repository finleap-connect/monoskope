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

package mock

import (
	"context"

	"github.com/finleap-connect/monoskope/pkg/domain/constants/roles"
	"github.com/finleap-connect/monoskope/pkg/domain/constants/scopes"
	"github.com/finleap-connect/monoskope/pkg/domain/projections"
	"github.com/finleap-connect/monoskope/pkg/domain/repositories"
	"github.com/google/uuid"
)

var (
	TestAdminUser        = projections.NewUserProjection(uuid.MustParse("00000000-0000-0000-0000-000000000001"))
	TestTenantAdminUser  = projections.NewUserProjection(uuid.MustParse("00000000-0000-0000-0000-000000000002"))
	TestExistingUser     = projections.NewUserProjection(uuid.MustParse("00000000-0000-0000-0000-000000000003"))
	TestNoneExistingUser = projections.NewUserProjection(uuid.MustParse("00000000-0000-0000-0000-000000000004"))

	TestAdminUserRoleBinding       = projections.NewUserRoleBinding(uuid.MustParse("79e61bbb-9373-4905-885b-a24eadc04bb7"))
	TestTenantAdminUserRoleBinding = projections.NewUserRoleBinding(uuid.MustParse("9d13bed6-00b5-4010-a787-60044dbc709a"))

	TestTenant  = projections.NewTenantProjection(uuid.MustParse("00000000-0000-0000-0001-000000000001"))
	TestCluster = projections.NewClusterProjection(uuid.MustParse("00000000-0000-0000-0002-000000000001"))

	TestTenantClusterBinding = projections.NewTenantClusterBindingProjection(uuid.MustParse("ad3e7523-3f32-433f-a6ed-5e0fe4ea6848"))
)

func init() {
	TestAdminUser.Name = "admin"
	TestAdminUser.Email = "admin@monoskope.io"

	TestExistingUser.Name = "someone"
	TestExistingUser.Email = "someone@monoskope.io"

	TestNoneExistingUser.Name = "nobody"
	TestNoneExistingUser.Email = "nobody@monoskope.io"

	TestTenantAdminUser.Name = "tenant-admin"
	TestTenantAdminUser.Email = "tenant-admin@monoskope.io"

	TestAdminUserRoleBinding.UserId = TestAdminUser.Id
	TestAdminUserRoleBinding.Role = string(roles.Admin)
	TestAdminUserRoleBinding.Scope = string(scopes.System)

	TestTenantAdminUserRoleBinding.UserId = TestTenantAdminUser.Id
	TestTenantAdminUserRoleBinding.Role = string(roles.Admin)
	TestTenantAdminUserRoleBinding.Scope = string(scopes.Tenant)
	TestTenantAdminUserRoleBinding.Resource = TestTenant.Id

	TestCluster.Name = "test-cluster"
	TestCluster.DisplayName = "Test Cluster"
	TestCluster.ApiServerAddress = "https://somecluster.io"
	TestCluster.CaCertBundle = []byte("some-bundle")

	TestTenantClusterBinding.ClusterId = TestCluster.Id
	TestTenantClusterBinding.TenantId = TestTenant.Id
}

func AddMockUsers(ctx context.Context, repo repositories.UserRepository) error {
	if err := repo.Upsert(ctx, TestAdminUser); err != nil {
		return err
	}
	if err := repo.Upsert(ctx, TestTenantAdminUser); err != nil {
		return err
	}
	if err := repo.Upsert(ctx, TestExistingUser); err != nil {
		return err
	}
	return nil
}

func AddMockUserRoleBindings(ctx context.Context, repo repositories.UserRoleBindingRepository) error {
	if err := repo.Upsert(ctx, TestAdminUserRoleBinding); err != nil {
		return err
	}
	if err := repo.Upsert(ctx, TestTenantAdminUserRoleBinding); err != nil {
		return err
	}
	return nil
}

func AddMockClusters(ctx context.Context, repo repositories.ClusterRepository) error {
	if err := repo.Upsert(ctx, TestCluster); err != nil {
		return err
	}
	return nil
}

func AddMockTenantClusterBindings(ctx context.Context, repo repositories.TenantClusterBindingRepository) error {
	if err := repo.Upsert(ctx, TestTenantClusterBinding); err != nil {
		return err
	}
	return nil
}

func AddMockTenants(ctx context.Context, repo repositories.TenantRepository) error {
	if err := repo.Upsert(ctx, TestTenant); err != nil {
		return err
	}
	return nil
}
