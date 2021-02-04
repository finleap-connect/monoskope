package constants

import es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/event_sourcing"

// User
const (
	// Type for the UserAggregate
	UserType es.AggregateType = "User"

	// Command to create a new User
	CreateUserType es.CommandType = "CreateUser"

	// Event emitted when a User has been created
	UserCreatedType es.EventType = "UserCreated"
)
