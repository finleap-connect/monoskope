package handler

import (
	"context"

	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/errors"
	metadata "gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/metadata"
	projections "gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/projections"
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
	if err != nil {
		return err
	}

	for _, policy := range cmd.Policies(ctx) {
		if policyAccepts(user, policy) {
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
func policyAccepts(user *projections.User, policy es.Policy) bool {
	for _, roleBinding := range user.GetRoles() {
		if policy.AcceptsRole(es.Role(roleBinding.Role)) &&
			policy.AcceptsScope(es.Scope(roleBinding.Scope)) &&
			policy.AcceptsResource(roleBinding.Resource) &&
			policy.AcceptsSubject(user.Email) {
			return true
		}
	}

	return false
}
