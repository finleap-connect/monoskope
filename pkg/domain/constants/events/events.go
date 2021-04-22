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

	// Event emitted when a Cluster Registration has been requested
	ClusterRegistrationRequested es.EventType = "ClusterRegistrationRequested"
	// Event emitted when a Cluster Registration has been registered
	ClusterRegistrationApproved es.EventType = "ClusterRegistrationApproved"
	// Event emitted when a Cluster Registration has been denied
	ClusterRegistrationDenied es.EventType = "ClusterRegistrationDenied"

	// Event emitted when a Cluster has been created
	ClusterCreated es.EventType = "ClusterCreated"
	// Event emitted when a Cluster has been deleted
	ClusterDeleted es.EventType = "ClusterDeleted"
)
