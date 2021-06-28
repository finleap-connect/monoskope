package eventsourcing

import (
	"fmt"
)

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

const (
	AnyRole  Role  = "*"
	AnyScope Scope = "*"
)

type Policy interface {
	// Role returns the role this policy accepts.
	Role() Role
	// Scope returns the scope this policy accepts.
	Scope() Scope

	// WithRole sets the role this policy accepts to the given value.
	WithRole(Role) Policy
	// WithScope sets the scope this policy accepts to the given value.
	WithScope(Scope) Policy

	// AcceptsRole checks if the policy accepts the given role.
	AcceptsRole(Role) bool
	// AcceptsScope checks if the policy accepts the given scope.
	AcceptsScope(Scope) bool

	// String returns a string representation of the policy.
	String() string
}

// Policy describes which Role/Scope combination is allowed to execute a certain Command.
type policy struct {
	// Role is the Role a user must have due to the policy.
	role Role
	// Scope is the Scope of the Role a user must have due to the policy.
	scope Scope
}

// NewPolicy creates a new policy which accepts anything.
func NewPolicy() Policy {
	return &policy{
		role:  AnyRole,
		scope: AnyScope,
	}
}

// WithRole sets the role of this policy accepts to the given value.
func (p *policy) WithRole(role Role) Policy {
	p.role = role
	return p
}

// WithScope sets the scope of this policy accepts to the given value.
func (p *policy) WithScope(scope Scope) Policy {
	p.scope = scope
	return p
}

// Role returns the role this policy accepts.
func (p *policy) Role() Role {
	return p.role
}

// Scope returns the scope this policy accepts.
func (p *policy) Scope() Scope {
	return p.scope
}

// AcceptsRole checks if the policy accepts the given role.
func (p *policy) AcceptsRole(role Role) bool {
	return p.role == AnyRole || p.role == role
}

// AcceptsScope checks if the policy accepts the given scope.
func (p *policy) AcceptsScope(scope Scope) bool {
	return p.scope == AnyScope || p.scope == scope
}

// String returns a string representation of the policy.
func (p *policy) String() string {
	return fmt.Sprintf("RO:%s|SC:%s", p.role, p.scope)
}
