package handler

import (
	"context"

	projections "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/queryhandler"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/errors"
	metadata "gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/metadata"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/repositories"
	es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing"
)

type authorizationCommandHandler struct {
	userRepo repositories.ReadOnlyUserRepository
}

// NewAuthorizationHandler creates a new CommandHandler which handles authorization.
func NewAuthorizationHandler(userRepo repositories.ReadOnlyUserRepository) es.CommandHandler {
	return &authorizationCommandHandler{
		userRepo: userRepo,
	}
}

// HandleCommand implements the CommandHandler interface
func (h *authorizationCommandHandler) HandleCommand(ctx context.Context, cmd es.Command) error {
	metadataMngr := metadata.NewDomainMetadataManager(ctx)
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
			return nil
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
