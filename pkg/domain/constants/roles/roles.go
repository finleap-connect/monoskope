// Copyright 2022 Monoskope Authors
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

package roles

import (
	"fmt"

	"github.com/finleap-connect/monoskope/pkg/domain/errors"
	es "github.com/finleap-connect/monoskope/pkg/eventsourcing"
	"github.com/google/uuid"
)

// Roles
const (
	// Admin role
	Admin es.Role = "admin"

	// User role
	User es.Role = "user"

	// OnCall role
	OnCall es.Role = "oncall"

	// K8sOperator role
	K8sOperator es.Role = "k8soperator"
)

// A list of all existing roles.
var AvailableRoles = []es.Role{
	Admin,
	User,
	K8sOperator,
	OnCall,
}

// A list of all existing roles.
var AvailableRolesMap = make(map[uuid.UUID]es.Role)

func init() {
	for _, role := range AvailableRoles {
		AvailableRolesMap[uuid.NewSHA1(uuid.NameSpaceURL, []byte(role))] = role
	}
}

func ValidateRole(role string) error {
	for _, v := range AvailableRoles {
		if v.String() == role {
			return nil
		}
	}
	return errors.ErrInvalidArgument(fmt.Sprintf("Role '%s' is invalid.", role))
}

func IdFromRole(role es.Role) uuid.UUID {
	return uuid.NewSHA1(uuid.NameSpaceURL, []byte(role))
}

func ToRole(role string) (es.Role, error) {
	err := ValidateRole(role)
	if err != nil {
		return "", err
	}
	return es.Role(role), nil
}
