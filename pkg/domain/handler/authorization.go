package handler

import (
	"context"

	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/aggregates"
	es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing"
)

type authorizationCommandHandler struct {
}

func NewAuthorizationHandler() es.CommandHandler {
	return &authorizationCommandHandler{}
}

func (h *authorizationCommandHandler) HandleCommand(ctx context.Context, cmd es.Command) error {
	// TODO:
	// Get current users rolebindings from ctx
	roleBindings := []aggregates.UserRoleBindingAggregate{}
	for _, policy := range cmd.Policies(ctx) {
		for _, roleBinding := range roleBindings {
			if policy.Role == roleBinding.Role() {
				return nil
			}
		}
	}
	return nil
}
