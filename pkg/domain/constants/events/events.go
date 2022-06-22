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

package events

import es "github.com/finleap-connect/monoskope/pkg/eventsourcing"

const (
	// UserCreated event emitted when a User has been created
	UserCreated es.EventType = "UserCreated"
	// UserUpdated event emitted when a User has been updated
	UserUpdated es.EventType = "UserUpdated"
	// UserDeleted event emitted when a User has been deleted
	UserDeleted es.EventType = "UserDeleted"
	// UserRoleBindingCreated event emitted when a new UserRoleBinding has been created
	UserRoleBindingCreated es.EventType = "UserRoleBindingCreated"
	// UserRoleBindingDeleted event emitted when a UserRoleBinding has been deleted
	UserRoleBindingDeleted es.EventType = "UserRoleBindingDeleted"

	// TenantCreated event emitted when a User has been created
	TenantCreated es.EventType = "TenantCreated"
	// TenantUpdated event emitted when a Tenant has been updated
	TenantUpdated es.EventType = "TenantUpdated"
	// TenantDeleted event emitted when a Tenant has been deleted
	TenantDeleted es.EventType = "TenantDeleted"

	// ClusterCreated event emitted when a Cluster has been created
	ClusterCreated   es.EventType = "ClusterCreated"
	ClusterCreatedV2 es.EventType = "ClusterCreatedV2"
	// ClusterUpdated event emitted when a Cluster has been created
	ClusterUpdated es.EventType = "ClusterUpdated"
	// ClusterDeleted event emitted when a Cluster has been deleted
	ClusterDeleted es.EventType = "ClusterDeleted"
	// IGNORED: ClusterBootstrapTokenCreated event emitted when a bootstrap token has been created
	ClusterBootstrapTokenCreated es.EventType = "ClusterBootstrapTokenCreated"

	// CertificateRequested event emitted when a certificate has been requested
	CertificateRequested es.EventType = "CertificateRequested"
	// CertificateRequestIssued event emitted when a certificate request has been issued
	CertificateRequestIssued es.EventType = "CertificateRequestIssued"
	// CertificateIssued event emitted when a certificate has been issued
	CertificateIssued es.EventType = "CertificateIssued"
	// CertificateIssueingFailed event emitted when a certificate could not be issued
	CertificateIssueingFailed es.EventType = "CertificateIssueingFailed"

	// TenantClusterBindingCreated event emitted when a tenant was given access to a certain cluster
	TenantClusterBindingCreated es.EventType = "TenantClusterBindingCreated"
	// TenantClusterBindingDeleted event emitted when a tenant's access to a cluster has been revoked
	TenantClusterBindingDeleted es.EventType = "TenantClusterBindingDeleted"
)

var (
	UserEvents = []es.EventType{
		UserCreated,
		UserUpdated,
		UserDeleted,
		UserRoleBindingCreated,
		UserRoleBindingDeleted,
	}

	TenantEvents = []es.EventType{
		TenantCreated,
		TenantUpdated,
		TenantDeleted,
		TenantClusterBindingCreated,
		TenantClusterBindingDeleted,
	}

	ClusterEvents = []es.EventType{
		ClusterCreated,
		ClusterCreatedV2,
		ClusterUpdated,
		ClusterDeleted,
	}

	CertificateEvents = []es.EventType{
		CertificateRequested,
		CertificateRequestIssued,
		CertificateIssued,
		CertificateIssueingFailed,
	}
)
