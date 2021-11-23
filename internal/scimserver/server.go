package scimserver

import (
	"github.com/elimity-com/scim"
	"github.com/elimity-com/scim/optional"
)

func NewServer() scim.Server {
	return scim.Server{
		ResourceTypes: []scim.ResourceType{
			{
				ID:          optional.NewString("User"),
				Name:        "User",
				Endpoint:    "/Users",
				Description: optional.NewString("User Account"),
				Schema:      UserSchema,
				Handler:     NewUserHandler(),
			},
			{
				ID:          optional.NewString("Group"),
				Name:        "Group",
				Endpoint:    "/Groups",
				Description: optional.NewString("Group"),
				Schema:      GroupSchema,
				Handler:     NewGroupHandler(),
			},
		},
	}
}
