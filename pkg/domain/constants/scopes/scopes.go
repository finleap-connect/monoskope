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
