package users

import (
	"fmt"

	"github.com/google/uuid"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/constants/roles"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/constants/scopes"
	"gitlab.figo.systems/platform/monoskope/monoskope/pkg/domain/projections"
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
	CommandHandlerUser = NewSystemUser("commandhandler")
	ReactorUser = NewSystemUser("reactor")

	AvailableSystemUsers = map[uuid.UUID]*projections.User{
		CommandHandlerUser.ID(): CommandHandlerUser,
		ReactorUser.ID():        ReactorUser,
	}
}

// NewSystemUser creates a new system user with a reproducible name based on the name and an admin rolebinding
func NewSystemUser(name string) *projections.User {
	userId := GenerateSystemUserUUID(name)

	// Create admin rolebinding
	adminRoleBinding := projections.NewUserRoleBinding(uuid.Nil)
	adminRoleBinding.UserId = userId.String()
	adminRoleBinding.Role = string(roles.Admin)
	adminRoleBinding.Scope = string(scopes.System)

	// Create system user
	user := projections.NewUserProjection(userId).(*projections.User)
	user.Name = name
	user.Email = GenerateSystemEmailAddress(name)
	user.Roles = append(user.Roles, adminRoleBinding.UserRoleBinding)

	return user
}

// Generates an email address with the name and the base domain constant
func GenerateSystemEmailAddress(name string) string {
	return fmt.Sprintf("%s@%s", name, BASE_DOMAIN)
}

// GenerateSystemUserUUID creates a reproducible UUID based on the name
func GenerateSystemUserUUID(name string) uuid.UUID {
	userMailAddress := GenerateSystemEmailAddress(name)
	return uuid.NewSHA1(uuid.NameSpaceURL, []byte(userMailAddress))
}
