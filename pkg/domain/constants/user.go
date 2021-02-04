package constants

import es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/event_sourcing"

// User
const (
	// Type for the User aggregate
	UserType es.AggregateType = "User"

	// Command to create a new user
	CreateUserType es.CommandType = "CreateUser"

	// Event emitted when a user has been created
	UserCreatedType es.EventType = "UserCreated"
)
