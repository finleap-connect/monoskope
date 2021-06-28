package users

import (
	"fmt"

	"github.com/google/uuid"
)

const (
	BASE_DOMAIN = "monoskope.local"
)

var (
	CommandHandlerUser *SystemUser
	ReactorUser        *SystemUser
)

type SystemUser struct {
	ID    uuid.UUID
	Name  string
	Email string
}

func init() {
	CommandHandlerUser = NewSystemUser("commandhandler")
	ReactorUser = NewSystemUser("reactor")
}

func NewSystemUser(name string) *SystemUser {
	user := new(SystemUser)
	user.Name = name
	user.Email = GenerateSystemEmailAddress(name)
	user.ID = GenerateSystemUserUUID(name)
	return user
}

func GenerateSystemEmailAddress(name string) string {
	return fmt.Sprintf("%s@%s", name, BASE_DOMAIN)
}

// GenerateSystemUserUUID creates a reproducible UUID based on the name
func GenerateSystemUserUUID(name string) uuid.UUID {
	userMailAddress := GenerateSystemEmailAddress(name)
	return uuid.NewSHA1(uuid.NameSpaceURL, []byte(userMailAddress))
}
