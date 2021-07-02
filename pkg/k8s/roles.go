package k8s

import (
	"fmt"

	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/errors"
)

// K8sRole is the name of a user's K8s role.
type K8sRole string

// K8s Roles
const (
	// User role
	DefaultRole K8sRole = "default"
	// Admin role
	AdminRole K8sRole = "admin"
	// OnCaller role
	OnCallRole K8sRole = "oncall"
)

// A list of all existing cluster roles.
var AvailableRoles = []K8sRole{
	DefaultRole,
	AdminRole,
	OnCallRole,
}

func ValidateRole(role string) error {
	for _, v := range AvailableRoles {
		if string(v) == role {
			return nil
		}
	}
	return errors.ErrInvalidArgument(fmt.Sprintf("Role '%s' is invalid.", role))
}
