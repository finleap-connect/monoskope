// Copyright 2022 Monoskope Authors
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

package users

import (
	"context"
	"fmt"

	"github.com/finleap-connect/monoskope/pkg/domain/constants/roles"
	"github.com/finleap-connect/monoskope/pkg/domain/constants/scopes"
	"github.com/finleap-connect/monoskope/pkg/domain/metadata"
	"github.com/finleap-connect/monoskope/pkg/domain/projections"
	"github.com/google/uuid"
)

const (
	BASE_DOMAIN = "monoskope.local"
)

var (
	// CommandHandlerUser is the system user representing the CommandHandler
	CommandHandlerUser *projections.User
	// SCIMServerUser is the system user representing the SCIM server
	SCIMServerUser        *projections.User
	GitRepoReconcilerUser *projections.User
)

// A maps of all existing system users.
var AvailableSystemUsers map[uuid.UUID]*projections.User

func init() {
	CommandHandlerUser = NewSystemUser("commandhandler")
	SCIMServerUser = NewSystemUser("scimserver")
	GitRepoReconcilerUser = NewSystemUser("gitreporeconciler")

	AvailableSystemUsers = map[uuid.UUID]*projections.User{
		CommandHandlerUser.ID():    CommandHandlerUser,
		SCIMServerUser.ID():        SCIMServerUser,
		GitRepoReconcilerUser.ID(): GitRepoReconcilerUser,
	}
}

// NewSystemUser creates a new system user with a reproducible name based on the name and an admin rolebinding
func NewSystemUser(name string) *projections.User {
	userId := generateSystemUserUUID(name)

	// Create admin rolebinding
	adminRoleBinding := projections.NewUserRoleBinding(uuid.Nil)
	adminRoleBinding.UserId = userId.String()
	adminRoleBinding.Role = string(roles.Admin)
	adminRoleBinding.Scope = string(scopes.System)

	// Create system user
	user := projections.NewUserProjection(userId)
	user.Name = name
	user.Email = generateSystemEmailAddress(name)
	user.Roles = append(user.Roles, adminRoleBinding.UserRoleBinding)

	return user
}

// generateSystemEmailAddress generates an email address with the name and the base domain constant
func generateSystemEmailAddress(name string) string {
	return fmt.Sprintf("%s@%s", name, BASE_DOMAIN)
}

// generateSystemUserUUID creates a reproducible UUID based on the name
func generateSystemUserUUID(name string) uuid.UUID {
	userMailAddress := generateSystemEmailAddress(name)
	return uuid.NewSHA1(uuid.NameSpaceURL, []byte(userMailAddress))
}

func createUserContext(ctx context.Context, user *projections.User) (*metadata.DomainMetadataManager, error) {
	metadataManager, err := metadata.NewDomainMetadataManager(ctx)
	if err != nil {
		return nil, err
	}
	userInfo := metadataManager.GetUserInformation()
	userInfo.Id = user.ID()
	userInfo.Name = user.Name
	userInfo.Email = user.Email
	metadataManager.SetUserInformation(userInfo)
	return metadataManager, nil
}

func CreateUserContext(ctx context.Context, user *projections.User) (context.Context, error) {
	mdManager, err := createUserContext(ctx, user)
	if err != nil {
		return nil, err
	}
	return mdManager.GetContext(), err
}

func CreateUserContextGrpc(ctx context.Context, user *projections.User) (context.Context, error) {
	mdManager, err := createUserContext(ctx, user)
	if err != nil {
		return nil, err
	}
	return mdManager.GetOutgoingGrpcContext(), err
}
