package handler

import (
	"context"

	domainErrors "gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/errors"
	metadata "gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/metadata"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/repositories"
	es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/logger"
)

type userInformationHandler struct {
	log                logger.Logger
	userRepo           repositories.ReadOnlyUserRepository
	nextHandlerInChain es.CommandHandler
}

// NewUserInformationHandler creates a new CommandHandler which handles authorization.
func NewUserInformationHandler(userRepo repositories.ReadOnlyUserRepository) *userInformationHandler {
	return &userInformationHandler{
		log:      logger.WithName("user-information-middleware"),
		userRepo: userRepo,
	}
}

func (m *userInformationHandler) Middleware(h es.CommandHandler) es.CommandHandler {
	m.nextHandlerInChain = h
	return m
}

// HandleCommand implements the CommandHandler interface
func (h *userInformationHandler) HandleCommand(ctx context.Context, cmd es.Command) (*es.CommandReply, error) {
	// Gather context
	metadataManager, err := metadata.NewDomainMetadataManager(ctx)
	if err != nil {
		return nil, err
	}
	userInfo := metadataManager.GetUserInformation()

	if metadataManager.IsAuthorizationBypassed() {
		h.log.V(logger.WarnLevel).Info("Authorization bypass enabled.", "CommandType", cmd.CommandType(), "AggregateType", cmd.AggregateType(), "User", userInfo.Email)
		return h.nextHandlerInChain.HandleCommand(metadataManager.GetContext(), cmd)
	}

	// Check that user exists
	user, err := h.userRepo.ByEmail(ctx, userInfo.Email)
	if err != nil {
		h.log.Info("User does not exist -> unauthorized.", "CommandType", cmd.CommandType(), "AggregateType", cmd.AggregateType(), "User", userInfo.Email)
		return nil, domainErrors.ErrUnauthorized
	}

	// Enrich context with rolebindings and user information
	userRoleBindings := user.GetRoles()
	userInfo.Id = user.ID()
	metadataManager.SetUserInformation(userInfo)
	metadataManager.SetRoleBindings(userRoleBindings)

	// Run next handler in chain
	if h.nextHandlerInChain != nil {
		return h.nextHandlerInChain.HandleCommand(metadataManager.GetContext(), cmd)
	} else {
		return nil, nil
	}
}
