package commands

import es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/event_sourcing"

const (
	// Event emitted when a User has been created
	CreateUser es.CommandType = "CreateUser"

	// Command to create a new UserRoleBinding
	CreateUserRoleBinding es.CommandType = "CreateUserRoleBinding"
)
