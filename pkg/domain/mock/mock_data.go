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
	"github.com/finleap-connect/monoskope/pkg/domain/constants/roles"
	"github.com/finleap-connect/monoskope/pkg/domain/constants/scopes"
	"github.com/finleap-connect/monoskope/pkg/domain/projections"
	"github.com/google/uuid"
)

var (
	TestAdminUser        = projections.NewUserProjection(uuid.MustParse("00000000-0000-0000-0000-000000000000"))
	TestTenantAdminUser  = projections.NewUserProjection(uuid.MustParse("00000000-0000-0000-0000-000000000001"))
	TestExistingUser     = projections.NewUserProjection(uuid.MustParse("00000000-0000-0000-0000-000000000002"))
	TestNoneExistingUser = projections.NewUserProjection(uuid.MustParse("00000000-0000-0000-0000-000000000003"))

	TestMockUsers = []*projections.User{
		TestAdminUser,
		TestTenantAdminUser,
		TestExistingUser,
	}

	TestAdminUserRoleBinding       = projections.NewUserRoleBinding(uuid.MustParse("00000000-0000-0000-0001-000000000000"))
	TestTenantAdminUserRoleBinding = projections.NewUserRoleBinding(uuid.MustParse("00000000-0000-0000-0002-000000000000"))

	TestTenant  = projections.NewTenantProjection(uuid.MustParse("00000000-0000-0001-0000-000000000000"))
	TestCluster = projections.NewClusterProjection(uuid.MustParse("00000000-0000-0001-0000-000000000000"))

	TestTenantClusterBinding = projections.NewTenantClusterBindingProjection(uuid.MustParse("00000000-0001-0000-0000-000000000000"))
)

func init() {
	TestAdminUserRoleBinding.UserId = TestAdminUser.Id
	TestAdminUserRoleBinding.Role = string(roles.Admin)
	TestAdminUserRoleBinding.Scope = string(scopes.System)
	TestAdminUser.Roles = append(TestAdminUser.Roles, TestAdminUserRoleBinding.Proto())

	TestTenantAdminUserRoleBinding.UserId = TestTenantAdminUser.Id
	TestTenantAdminUserRoleBinding.Role = string(roles.Admin)
	TestTenantAdminUserRoleBinding.Scope = string(scopes.Tenant)
	TestTenantAdminUserRoleBinding.Resource = TestTenant.Id
	TestTenantAdminUser.Roles = append(TestTenantAdminUser.Roles, TestTenantAdminUserRoleBinding.Proto())

	TestCluster.Name = "test-cluster"
	TestCluster.ApiServerAddress = "https://somecluster.io"
	TestCluster.CaCertBundle = []byte("some-bundle")

	TestTenantClusterBinding.ClusterId = TestCluster.Id
	TestTenantClusterBinding.TenantId = TestTenant.Id

	TestAdminUser.Name = "admin"
	TestAdminUser.Email = "admin@monoskope.io"

	TestExistingUser.Name = "someone"
	TestExistingUser.Email = "someone@monoskope.io"

	TestNoneExistingUser.Name = "nobody"
	TestNoneExistingUser.Email = "nobody@monoskope.io"

	TestTenantAdminUser.Name = "tenant-admin"
	TestTenantAdminUser.Email = "tenant-admin@monoskope.io"

	TestTenant.Name = "test-tenant"
}
