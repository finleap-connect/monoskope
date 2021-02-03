package authz

import es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/event_sourcing"

const (
	Admin es.Role = "admin"
	User  es.Role = "user"

	System es.Scope = "system"
	Tenant es.Scope = "tenant"
)
