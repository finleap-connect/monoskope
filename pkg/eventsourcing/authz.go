package eventsourcing

// Role is the name of a user's role.
type Role string

func (r Role) String() string {
	return string(r)
}

// Scope is the scope of a role.
type Scope string

func (s Scope) String() string {
	return string(s)
}

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
