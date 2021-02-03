package event_sourcing

// Role is the name of a user's role.
type Role string

// Scope is the scope of a role.
type Scope string

// Policy describes which Role/Scope/Resource combination is allowed to execute a certain Command.
type Policy struct {
	// Role is the Role a user must have due to the policy.
	Role Role
	// Scope is the Scope of the Role a user must have due to the policy.
	Scope Scope
	// Resource is the Resource which has be part of the Role/Scope a user must have due to the policy.
	Resource string

	// Subject is the user one has to be due to the policy.
	Subject string
}
