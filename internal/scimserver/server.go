package scimserver

import (
	"github.com/elimity-com/scim"
	"github.com/elimity-com/scim/optional"
	"github.com/elimity-com/scim/schema"
)

func NewServer(config scim.ServiceProviderConfig, userHandler, groupHandler scim.ResourceHandler) scim.Server {
	return scim.Server{
		Config: config,
		ResourceTypes: []scim.ResourceType{
			{
				ID:               optional.NewString("User"),
				Name:             "User",
				Endpoint:         "/Users",
				Description:      optional.NewString("User Account"),
				Schema:           schema.CoreUserSchema(),
				SchemaExtensions: []scim.SchemaExtension{{Schema: schema.ExtensionEnterpriseUser(), Required: true}},
				Handler:          userHandler,
			},
			{
				ID:          optional.NewString("Group"),
				Name:        "Group",
				Endpoint:    "/Groups",
				Description: optional.NewString("Group"),
				Schema:      schema.CoreGroupSchema(),
				Handler:     groupHandler,
			},
		},
	}
}
