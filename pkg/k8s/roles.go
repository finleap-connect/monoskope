// Copyright 2021 Monoskope Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package k8s

import (
	"fmt"

	"github.com/finleap-connect/monoskope/pkg/domain/errors"
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
