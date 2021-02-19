package eventsourcing

import "fmt"

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
	AnyRole     Role   = "*"
	AnyScope    Scope  = "*"
	AnyResource string = "*"
	AnySubject  string = "*"
)

type Policy interface {
	// Role returns the role this policy accepts.
	Role() Role
	// Scope returns the scope this policy accepts.
	Scope() Scope
	// Resource returns the resource this policy accepts.
	Resource() string
	// Subject returns the subject this policy accepts.
	Subject() string

	// WithRole sets the role of this policy accepts to the given value.
	WithRole(Role) Policy
	// WithScope sets the scope of this policy accepts to the given value.
	WithScope(Scope) Policy
	// WithResource sets the resource of this policy accepts to the given value.
	WithResource(string) Policy
	// WithSubject sets the subject of this policy accepts to the given value.
	WithSubject(string) Policy

	// AcceptsRole checks if the policy accepts the given role.
	AcceptsRole(Role) bool
	// AcceptsScope checks if the policy accepts the given scope.
	AcceptsScope(Scope) bool
	// AcceptsResource checks if the policy accepts the given resource.
	AcceptsResource(string) bool
	// AcceptsSubject checks if the policy accepts the given subject.
	AcceptsSubject(string) bool
	// MustBeSubject checks if the policy enforces the given subject.
	MustBeSubject(string) bool
}

// Policy describes which Role/Scope/Resource combination is allowed to execute a certain Command.
type policy struct {
	// Role is the Role a user must have due to the policy.
	role Role
	// Scope is the Scope of the Role a user must have due to the policy.
	scope Scope
	// Resource is the Resource which has be part of the Role/Scope a user must have due to the policy.
	resource string

	// Subject is the user one has to be due to the policy.
	subject string
}

// NewPolicy creates a new policy which accepts anything.
func NewPolicy() Policy {
	return &policy{
		role:     AnyRole,
		scope:    AnyScope,
		resource: AnyResource,
		subject:  AnySubject,
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

// WithResource sets the resource of this policy accepts to the given value.
func (p *policy) WithResource(resource string) Policy {
	p.resource = resource
	return p
}

// WithSubject sets the subject of this policy accepts to the given value.
func (p *policy) WithSubject(subject string) Policy {
	p.subject = subject
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

// Resource returns the resource this policy accepts.
func (p *policy) Resource() string {
	return p.resource
}

// Subject returns the subject this policy accepts.
func (p *policy) Subject() string {
	return p.subject
}

// AcceptsRole checks if the policy accepts the given role.
func (p *policy) AcceptsRole(role Role) bool {
	return p.role == AnyRole || p.role == role
}

// AcceptsScope checks if the policy accepts the given scope.
func (p *policy) AcceptsScope(scope Scope) bool {
	return p.scope == AnyScope || p.scope == scope
}

// AcceptsResource checks if the policy accepts the given resource.
func (p *policy) AcceptsResource(resource string) bool {
	return p.resource == AnyResource || p.resource == resource
}

// AcceptsSubject checks if the policy accepts the given subject.
func (p *policy) AcceptsSubject(subject string) bool {
	return p.subject == AnySubject || p.subject == subject
}

// MustBeSubject checks if the policy enforces the given subject.
func (p *policy) MustBeSubject(subject string) bool {
	return p.subject == subject
}

// String returns a string representation of the policy.
func (p *policy) String() string {
	return fmt.Sprintf("RO:%s|SC:%s|RE:%s|SU:%s", p.role, p.scope, p.resource, p.subject)
}
