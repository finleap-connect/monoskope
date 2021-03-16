package events

import es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing"

const (
	// Event emitted when a User has been created
	UserCreated es.EventType = "UserCreated"
	// Event emitted when a new UserRoleBinding has been created
	UserRoleBindingCreated es.EventType = "UserRoleBindingCreated"
	// Event emitted when a User has been created
	TenantCreated es.EventType = "TenantCreated"
)
