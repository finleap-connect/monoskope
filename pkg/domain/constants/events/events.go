package events

import es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing"

const (
	// Event emitted when a User has been created
	UserCreated es.EventType = "UserCreated"
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
	ClusterCreated es.EventType = "ClusterCreated"
	// Event emitted when a Cluster has been created
	ClusterUpdate es.EventType = "ClusterUpdated"
	// Event emitted when a Cluster has been deleted
	ClusterDeleted es.EventType = "ClusterDeleted"
	// Event emitted when a bootstrap token has been created
	ClusterBootstrapTokenCreated es.EventType = "ClusterBootstrapTokenCreated"
	// Event emitted when a certificate has been requested for a cluster
	ClusterCertificateRequested es.EventType = "ClusterCertificateRequested"
	// Event emitted when a certificate has been issued for a cluster
	ClusterOperatorCertificateIssued es.EventType = "ClusterOperatorCertificateIssued"
)
