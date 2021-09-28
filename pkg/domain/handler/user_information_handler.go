// Copyright 2021 Monoskope Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package handler

import (
	"context"

	domainErrors "github.com/finleap-connect/monoskope/pkg/domain/errors"
	metadata "github.com/finleap-connect/monoskope/pkg/domain/metadata"
	"github.com/finleap-connect/monoskope/pkg/domain/repositories"
	es "github.com/finleap-connect/monoskope/pkg/eventsourcing"
	"github.com/finleap-connect/monoskope/pkg/logger"
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
