package aggregates

import (
	"context"

	"github.com/google/uuid"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/domain/projections"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/constants/scopes"
	domainErrors "gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/errors"
	metadata "gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/metadata"
	es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing"
)

type DomainAggregateBase struct {
	*es.BaseAggregate
}

// Authorization authorizes the command against the issueing users rolebindings
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

// resetId must be called for each command that will create a new aggregate, so that
// subsequent events will not reference the user supplied (possibly empty) ID.
func (a *DomainAggregateBase) resetId() {
	atype := a.BaseAggregate.Type()
	a.BaseAggregate = es.NewBaseAggregate(atype, uuid.New())
}
