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

func NewAuthorizationHandler(userRepo repositories.ReadOnlyUserRepository) es.CommandHandler {
	return &authorizationCommandHandler{
		userRepo: userRepo,
	}
}

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
		for _, roleBinding := range user.GetRoles() {
			if isAuthrozied(roleBinding, policy) {
				return nil
			}
		}
	}
	return errors.ErrUnauthorized
}

func isAuthrozied(roleBinding *projections.UserRoleBinding, policy es.Policy) bool {
	return false
}
