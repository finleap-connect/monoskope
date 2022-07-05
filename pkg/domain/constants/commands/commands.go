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

package commands

import es "github.com/finleap-connect/monoskope/pkg/eventsourcing"

const (
	// Event emitted when a User has been created
	CreateUser es.CommandType = "CreateUser"
	// Command to delete a User
	DeleteUser es.CommandType = "DeleteUser"
	// Command to update a User
	UpdateUser es.CommandType = "UpdateUser"

	// Command to create a new UserRoleBinding
	CreateUserRoleBinding es.CommandType = "CreateUserRoleBinding"
	// Command to delete a UserRoleBinding
	DeleteUserRoleBinding es.CommandType = "DeleteUserRoleBinding"

	// Event emitted when a Tenant has been created
	CreateTenant es.CommandType = "CreateTenant"
	// Command to update a Tenant
	UpdateTenant es.CommandType = "UpdateTenant"
	// Command to delete a Tenant
	DeleteTenant es.CommandType = "DeleteTenant"

	// Command to create a Cluster
	CreateCluster es.CommandType = "CreateCluster"
	// Command to delete a Cluster
	DeleteCluster es.CommandType = "DeleteCluster"
	// Command to update a Cluster
	UpdateCluster es.CommandType = "UpdateCluster"

	// Command to allow a tenant to access a certain cluster
	CreateTenantClusterBinding es.CommandType = "CreateTenantClusterBinding"
	// Command to remove access of a tenant to a certain cluster
	DeleteTenantClusterBinding es.CommandType = "DeleteTenantClusterBinding"
)

var (
	UserCommands = []es.CommandType{
		CreateUser,
		DeleteUser,
		UpdateUser,
	}

	UserRoleBindingCommands = []es.CommandType{
		CreateUserRoleBinding,
		DeleteUserRoleBinding,
	}

	TenantCommands = []es.CommandType{
		CreateTenant,
		UpdateTenant,
		DeleteTenant,
	}

	ClusterCommands = []es.CommandType{
		CreateCluster,
		DeleteCluster,
		UpdateCluster,
	}

	TenantClusterBindingCommands = []es.CommandType{
		CreateTenantClusterBinding,
		DeleteTenantClusterBinding,
	}

	CommandTypes = map[string][]es.CommandType{
		"User":                 UserCommands,
		"UserRoleBinding":      UserRoleBindingCommands,
		"Tenant":               TenantCommands,
		"Cluster":              ClusterCommands,
		"TenantClusterBinding": TenantClusterBindingCommands,
	}
)
