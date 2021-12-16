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

package users

import (
	"fmt"

	"github.com/finleap-connect/monoskope/pkg/api/domain/common"
	"github.com/finleap-connect/monoskope/pkg/domain/projections"
	"github.com/google/uuid"
)

const (
	BASE_DOMAIN = "monoskope.local"
)

var (
	// CommandHandlerUser is the system user representing the CommandHandler
	CommandHandlerUser *projections.User
	// ReactorUser is the system user representing any Reactor
	ReactorUser *projections.User
)

// A maps of all existing system users.
var AvailableSystemUsers map[uuid.UUID]*projections.User

func init() {
	CommandHandlerUser = newSystemUser("commandhandler")
	ReactorUser = newSystemUser("reactor")

	AvailableSystemUsers = map[uuid.UUID]*projections.User{
		CommandHandlerUser.ID(): CommandHandlerUser,
		ReactorUser.ID():        ReactorUser,
	}
}

// newSystemUser creates a new system user with a reproducible name based on the name and an admin rolebinding
func newSystemUser(name string) *projections.User {
	userId := generateSystemUserUUID(name)

	// Create admin rolebinding
	adminRoleBinding := projections.NewUserRoleBinding(uuid.Nil)
	adminRoleBinding.UserId = userId.String()
	adminRoleBinding.Role = string(common.Role_admin.String())
	adminRoleBinding.Scope = string(common.Scope_system.String())

	// Create system user
	user := projections.NewUserProjection(userId).(*projections.User)
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
