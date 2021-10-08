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

package aggregates

import (
	"context"

	"github.com/finleap-connect/monoskope/pkg/api/domain/projections"
	"github.com/finleap-connect/monoskope/pkg/domain/constants/scopes"
	domainErrors "github.com/finleap-connect/monoskope/pkg/domain/errors"
	metadata "github.com/finleap-connect/monoskope/pkg/domain/metadata"
	es "github.com/finleap-connect/monoskope/pkg/eventsourcing"
	"github.com/google/uuid"
)

type DomainAggregateBase struct {
	*es.BaseAggregate
}

// Authorization authorizes the command against the issuing users rolebindings
func (a *DomainAggregateBase) Authorize(ctx context.Context, cmd es.Command, expectedResource uuid.UUID) error {
	// Extract domain context
	metadataManager, err := metadata.NewDomainMetadataManager(ctx)
	if err != nil {
		return err
	}

	// Check if authorization has been bypassed
	if metadataManager.IsAuthorizationBypassed() {
		return nil
	}

	// Validate rolebindings against command policies
	userRoleBindings := metadataManager.GetRoleBindings()
	for _, policy := range cmd.Policies(ctx) {
		for _, roleBinding := range userRoleBindings {
			if validatePolicy(roleBinding, policy, expectedResource) {
				return nil
			}
		}
	}

	// If no policy matches return unauthorized
	return domainErrors.ErrUnauthorized
}

// validatePolicy validates a certain rolebinding against a certain policy
func validatePolicy(roleBinding *projections.UserRoleBinding, policy es.Policy, expectedResource uuid.UUID) bool {
	if !policy.AcceptsRole(es.Role(roleBinding.GetRole())) {
		return false
	}
	if !policy.AcceptsScope(es.Scope(roleBinding.GetScope())) {
		return false
	}
	if roleBinding.GetScope() != scopes.System.String() && roleBinding.GetResource() != expectedResource.String() {
		return false
	}
	return true
}

// Validate validates if the aggregate exists and is not deleted
func (a *DomainAggregateBase) Validate(ctx context.Context, cmd es.Command) error {
	if !a.Exists() {
		return domainErrors.ErrNotFound
	}
	if a.Deleted() {
		return domainErrors.ErrDeleted
	}
	return nil
}
