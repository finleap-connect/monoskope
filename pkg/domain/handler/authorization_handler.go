package handler

import (
	"context"

	projections "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/domain/projections"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/constants/roles"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/constants/scopes"
	domainErrors "gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/errors"
	metadata "gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/metadata"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/repositories"
	es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/logger"
)

type authorizationHandler struct {
	log                 logger.Logger
	userRepo            repositories.ReadOnlyUserRepository
	nextHandlerInChain  es.CommandHandler
	bypassAuthorization bool
}

// NewAuthorizationHandler creates a new CommandHandler which handles authorization.
func NewAuthorizationHandler(userRepo repositories.ReadOnlyUserRepository) *authorizationHandler {
	return &authorizationHandler{
		log:      logger.WithName("authorization-middleware"),
		userRepo: userRepo,
	}
}

// BypassAuthorization disables authorization checks and returns a function to enable it again
func (h *authorizationHandler) BypassAuthorization() func() {
	h.bypassAuthorization = true
	h.log.Info("WARNING authorization bypass has been enabled.")

	return func() {
		h.bypassAuthorization = false
		h.log.Info("Authorization bypass has been disabled.")
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
	userInfo := metadataMngr.GetUserInformation()

	if h.bypassAuthorization {
		h.log.Info("WARNING authorization bypass enabled.", "CommandType", cmd.CommandType(), "AggregateType", cmd.AggregateType(), "User", userInfo.Email)
		if err := metadataMngr.SetRoleBindings([]*projections.UserRoleBinding{
			{
				Role:  roles.Admin.String(),
				Scope: scopes.System.String(),
			},
		}); err != nil {
			h.log.Error(err, "Error when setting rolebindings.", "CommandType", cmd.CommandType(), "AggregateType", cmd.AggregateType(), "User", userInfo.Email)
			return domainErrors.ErrUnauthorized
		}
		return h.nextHandlerInChain.HandleCommand(metadataMngr.GetOutgoingGrpcContext(), cmd)
	}

	user, err := h.userRepo.ByEmail(ctx, userInfo.Email)
	if err != nil {
		h.log.Info("User does not exist -> unauthorized.", "CommandType", cmd.CommandType(), "AggregateType", cmd.AggregateType(), "User", userInfo.Email)
		return domainErrors.ErrUnauthorized
	}

	userRoleBindings := user.GetRoles()
	userInfo.Id = user.ID()
	metadataMngr.SetUserInformation(userInfo)
	if err := metadataMngr.SetRoleBindings(userRoleBindings); err != nil {
		h.log.Error(err, "Error when setting rolebindings.", "CommandType", cmd.CommandType(), "AggregateType", cmd.AggregateType(), "User", userInfo.Email)
		return domainErrors.ErrUnauthorized
	}

	if h.nextHandlerInChain != nil {
		return h.nextHandlerInChain.HandleCommand(metadataMngr.GetOutgoingGrpcContext(), cmd)
	} else {
		return nil
	}
}
