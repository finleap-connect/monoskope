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

// Package scopes sets the scope of permissions granted to a role: System, Tenant or Cluster.
// For the scopes Tenant and Cluster, a role binding will define to which specific tenant or cluster
// the role should be applied for a given user.
package scopes

import (
	"fmt"

	"github.com/finleap-connect/monoskope/pkg/domain/errors"
	es "github.com/finleap-connect/monoskope/pkg/eventsourcing"
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
