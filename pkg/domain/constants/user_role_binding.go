package constants

import es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/event_sourcing"

// UserRoleBinding
const (
	// Type for the UserRoleBindingAggregate
	UserRoleBindingType es.AggregateType = "UserRoleBinding"

	// Command to create a new UserRoleBinding
	CreateUserRoleBindingType es.CommandType = "CreateUserRoleBinding"

	// Event emitted when a new UserRoleBinding has been created
	UserRoleBindingCreatedType es.EventType = "UserRoleBindingCreated"
)
