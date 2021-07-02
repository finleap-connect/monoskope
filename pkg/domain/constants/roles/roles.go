package roles

import (
	"fmt"

	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/errors"
	es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing"
)

// Roles
const (
	// Admin role
	Admin es.Role = "admin"

	// User role
	User es.Role = "user"

	// K8sOperator role
	K8sOperator es.Role = "k8soperator"
)

// A list of all existing roles.
var AvailableRoles = []es.Role{
	Admin,
	User,
	K8sOperator,
}

func ValidateRole(role string) error {
	for _, v := range AvailableRoles {
		if v.String() == role {
			return nil
		}
	}
	return errors.ErrInvalidArgument(fmt.Sprintf("Role '%s' is invalid.", role))
}
