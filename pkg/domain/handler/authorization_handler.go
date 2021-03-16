package handler

import (
	"context"
	"errors"

	projectionsApi "gitlab.figo.systems/platform/monoskope/monoskope/pkg/api/domain/projections"
	domainErrors "gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/errors"
	metadata "gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/metadata"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/repositories"
	es "gitlab.figo.systems/platform/monoskope/monoskope/pkg/eventsourcing"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/logger"
)

type authorizationHandler struct {
	log                logger.Logger
	userRepo           repositories.ReadOnlyUserRepository
	nextHandlerInChain es.CommandHandler
}

// NewAuthorizationHandler creates a new CommandHandler which handles authorization.
func NewAuthorizationHandler(userRepo repositories.ReadOnlyUserRepository) *authorizationHandler {
	return &authorizationHandler{
		log:      logger.WithName("authorization-middleware"),
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

	userInfo := metadataMngr.GetUserInformation()
	user, err := h.userRepo.ByEmail(ctx, userInfo.Email)
	if err != nil && !errors.Is(err, domainErrors.ErrUserNotFound) {
		return domainErrors.ErrUnauthorized
	}

	if user != nil {
		metadataMngr.SetUserId(user.ID().String())
	}

	var userRoles []*projectionsApi.UserRoleBinding
	if user != nil {
		userRoles = user.GetRoles()
	}

	h.log.Info("Checking command authorization...", "CommandType", cmd.CommandType(), "AggregateType", cmd.AggregateType(), "User", userInfo.Email)
	for _, policy := range cmd.Policies(ctx) {
		if policyAccepts(userInfo.Email, userRoles, policy) {
			h.log.Info("User is authorized.", "CommandType", cmd.CommandType(), "AggregateType", cmd.AggregateType(), "User", userInfo.Email)
			if h.nextHandlerInChain != nil {
				return h.nextHandlerInChain.HandleCommand(metadataMngr.GetOutgoingGrpcContext(), cmd)
			} else {
				return nil
			}
		}
	}
	h.log.Info("User is unauthorized.", "CommandType", cmd.CommandType(), "AggregateType", cmd.AggregateType(), "User", userInfo.Email)
	return domainErrors.ErrUnauthorized
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
	} else if policy.MustBeSubject(userEmail) {
		return true
	}

	return false
}
