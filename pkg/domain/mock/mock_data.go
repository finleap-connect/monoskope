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
	TestAdminUser        = projections.NewUserProjection(uuid.Nil)
	TestTenantAdminUser  = projections.NewUserProjection(uuid.Nil)
	TestExistingUser     = projections.NewUserProjection(uuid.Nil)
	TestNoneExistingUser = projections.NewUserProjection(uuid.Nil)

	TestMockUsers = []*projections.User{
		TestAdminUser,
		TestTenantAdminUser,
		TestExistingUser,
	}

	TestAdminUserRoleBinding       = projections.NewUserRoleBinding(uuid.Nil)
	TestTenantAdminUserRoleBinding = projections.NewUserRoleBinding(uuid.Nil)

	TestTenant  = projections.NewTenantProjection(uuid.Nil)
	TestCluster = projections.NewClusterProjection(uuid.Nil)

	TestTenantClusterBinding = projections.NewTenantClusterBindingProjection(uuid.Nil)
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
	TestAdminUser.Roles = append(TestAdminUser.Roles, TestAdminUserRoleBinding.Proto())

	TestTenantAdminUserRoleBinding.UserId = TestTenantAdminUser.Id
	TestTenantAdminUserRoleBinding.Role = string(roles.Admin)
	TestTenantAdminUserRoleBinding.Scope = string(scopes.Tenant)
	TestTenantAdminUserRoleBinding.Resource = TestTenant.Id
	TestTenantAdminUser.Roles = append(TestTenantAdminUser.Roles, TestTenantAdminUserRoleBinding.Proto())

	TestCluster.Name = "test-cluster"
	TestCluster.DisplayName = "Test Cluster"
	TestCluster.ApiServerAddress = "https://somecluster.io"
	TestCluster.CaCertBundle = []byte("some-bundle")

	TestTenantClusterBinding.ClusterId = TestCluster.Id
	TestTenantClusterBinding.TenantId = TestTenant.Id
}
