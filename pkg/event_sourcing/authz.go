package event_sourcing

// Role is the name of a user's role.
type Role string

// Scope is the scope of a role.
type Scope string

// MetaRole describes which Role/Scope/Resource combination is allowed to execute a certain Command.
type MetaRole struct {
	Role     Role
	Scope    Scope
	Resource string

	// Explicit user allowed
	Subject string
}
