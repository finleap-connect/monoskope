package aggregates

import (
	"context"

	domainErrors "gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/errors"
	metadata "gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/metadata"
	es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing"
)

type DomainAggregateBase struct {
	*es.BaseAggregate
}

func (a *DomainAggregateBase) Authorize(ctx context.Context, cmd es.Command) error {
	metadataManager, err := metadata.NewDomainMetadataManager(ctx)
	if err != nil {
		return err
	}

	userRoleBindings, err := metadataManager.GetRoleBindings()
	if err != nil {
		return err
	}

	for _, policy := range cmd.Policies(ctx) {
		for _, roleBinding := range userRoleBindings {
			if policy.AcceptsRole(es.Role(roleBinding.Role)) &&
				policy.AcceptsScope(es.Scope(roleBinding.Scope)) {
				if !policy.ResourceMatch() || roleBinding.Resource == a.ID().String() {
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
