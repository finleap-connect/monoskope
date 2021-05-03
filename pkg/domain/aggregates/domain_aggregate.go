package aggregates

import (
	"context"

	"github.com/google/uuid"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/constants/scopes"
	domainErrors "gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/errors"
	metadata "gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/metadata"
	es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing"
)

type DomainAggregateBase struct {
	*es.BaseAggregate
}

func (a *DomainAggregateBase) Authorize(ctx context.Context, cmd es.Command, expectedResource uuid.UUID) error {
	metadataManager, err := metadata.NewDomainMetadataManager(ctx)
	if err != nil {
		return err
	}

	if metadataManager.IsAuthorizationBypassed() {
		return nil
	}

	userRoleBindings := metadataManager.GetRoleBindings()
	for _, policy := range cmd.Policies(ctx) {
		for _, roleBinding := range userRoleBindings {
			if policy.AcceptsRole(es.Role(roleBinding.GetRole())) &&
				policy.AcceptsScope(es.Scope(roleBinding.GetScope())) {
				if roleBinding.GetScope() == scopes.System.String() {
					return nil
				} else if roleBinding.GetResource() == expectedResource.String() {
					return nil
				}
			}
		}
	}
	return domainErrors.ErrUnauthorized
}

func (a *DomainAggregateBase) Validate(ctx context.Context, cmd es.Command) error {
	if !a.Exists() {
		return domainErrors.ErrNotFound
	}
	if a.Deleted() {
		return domainErrors.ErrDeleted
	}
	return nil
}
