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
)
