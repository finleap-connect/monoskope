package scopes

import es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing"

// Scopes
const (
	// System scope
	System es.Scope = "system"

	// Tenant scope
	Tenant es.Scope = "tenant"
)
