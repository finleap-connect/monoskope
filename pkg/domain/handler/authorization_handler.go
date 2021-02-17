package handler

import (
	"context"

	projectionsApi "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/domain/projections"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/errors"
	metadata "gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/metadata"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/repositories"
	es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing"
)

type authorizationHandler struct {
	userRepo           repositories.ReadOnlyUserRepository
	nextHandlerInChain es.CommandHandler
}

// NewAuthorizationHandler creates a new CommandHandler which handles authorization.
func NewAuthorizationHandler(userRepo repositories.ReadOnlyUserRepository) *authorizationHandler {
	return &authorizationHandler{
		userRepo: userRepo,
	}
}

func (m *authorizationHandler) Middleware(h es.CommandHandler) es.CommandHandler {
	m.nextHandlerInChain = h
	return m
}

// HandleCommand implements the CommandHandler interface
func (h *authorizationHandler) HandleCommand(ctx context.Context, cmd es.Command) error {
	metadataMngr, err := metadata.NewDomainMetadataManager(ctx)
	if err != nil {
		return err
	}

	userInfo, err := metadataMngr.GetUserInformation()
	if err != nil {
		return err
	}

	user, err := h.userRepo.ByEmail(ctx, userInfo.Email)
	if err != nil && err != errors.ErrUserNotFound {
		return err
	}

	var userRoles []*projectionsApi.UserRoleBinding
	if user != nil {
		userRoles = user.GetRoles()
	}

	for _, policy := range cmd.Policies(ctx) {
		if policyAccepts(userInfo.Email, userRoles, policy) {
			if h.nextHandlerInChain != nil {
				return h.nextHandlerInChain.HandleCommand(ctx, cmd)
			} else {
				return nil
			}
		}
	}
	return errors.ErrUnauthorized
}

// policyAccepts validates the policy against a user
func policyAccepts(userEmail string, userRoleBindings []*projectionsApi.UserRoleBinding, policy es.Policy) bool {
	if userRoleBindings != nil {
		for _, roleBinding := range userRoleBindings {
			if policy.AcceptsRole(es.Role(roleBinding.Role)) &&
				policy.AcceptsScope(es.Scope(roleBinding.Scope)) &&
				policy.AcceptsResource(roleBinding.Resource) &&
				policy.AcceptsSubject(userEmail) {
				return true
			}
		}
	} else if policy.AcceptsSubject(userEmail) {
		return true
	}

	return false
}
