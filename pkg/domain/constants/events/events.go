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

package events

import es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing"

const (
	// Event emitted when a User has been created
	UserCreated es.EventType = "UserCreated"
	// Event emitted when a User has been deleted
	UserDeleted es.EventType = "UserDeleted"
	// Event emitted when a new UserRoleBinding has been created
	UserRoleBindingCreated es.EventType = "UserRoleBindingCreated"
	// Event emitted when a UserRoleBinding has been deleted
	UserRoleBindingDeleted es.EventType = "UserRoleBindingDeleted"

	// Event emitted when a User has been created
	TenantCreated es.EventType = "TenantCreated"
	// Event emitted when a Tenant has been updated
	TenantUpdated es.EventType = "TenantUpdated"
	// Event emitted when a Tenant has been deleted
	TenantDeleted es.EventType = "TenantDeleted"

	// Event emitted when a Cluster has been created
	ClusterCreated   es.EventType = "ClusterCreated"
	ClusterCreatedV2 es.EventType = "ClusterCreatedV2"
	// Event emitted when a Cluster has been created
	ClusterUpdated es.EventType = "ClusterUpdated"
	// Event emitted when a Cluster has been deleted
	ClusterDeleted es.EventType = "ClusterDeleted"
	// Event emitted when a bootstrap token has been created
	ClusterBootstrapTokenCreated es.EventType = "ClusterBootstrapTokenCreated"

	// Event emitted when a certificate has been requested
	CertificateRequested es.EventType = "CertificateRequested"
	// Event emitted when a certificate request has been issued
	CertificateRequestIssued es.EventType = "CertificateRequestIssued"
	// Event emitted when a certificate has been issued
	CertificateIssued es.EventType = "CertificateIssued"
	// Event emitted when a certificate could not be issued
	CertificateIssueingFailed es.EventType = "CertificateIssueingFailed"
)
