// Package scopes sets the scope of permissions granted to a role: System, Tenant or Cluster.
// For the scopes Tenant and Cluster, a role binding will define to which specific tenant or cluster
// the role should be aplied for a given user.
package scopes

import (
	"fmt"

	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/errors"
	es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing"
)

// Scopes
const (
	// System scope
	System es.Scope = "system"

	// Tenant scope
	Tenant es.Scope = "tenant"

	// Cluster scope
	Cluster es.Scope = "cluster"
)

// A list of all existing scopes.
var AvailableScopes = []es.Scope{
	System,
	Tenant,
}

func ValidateScope(scope string) error {
	for _, v := range AvailableScopes {
		if v.String() == scope {
			return nil
		}
	}
	return errors.ErrInvalidArgument(fmt.Sprintf("Scope '%s' is invalid.", scope))
}
