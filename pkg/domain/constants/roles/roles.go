package roles

import es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing"

// Roles
const (
	// User role
	User es.Role = "user"

	// Admin role
	Admin es.Role = "admin"

	// Agent role
	Agent es.Role = "agent"
)

// A list of all existing roles.
var AvailableRoles = []es.Role{
	User,
	Admin,
	Agent,
}
