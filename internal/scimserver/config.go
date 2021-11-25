package scimserver

import (
	"github.com/elimity-com/scim"
	"github.com/elimity-com/scim/optional"
)

func NewProvierConfig() scim.ServiceProviderConfig {
	config := scim.ServiceProviderConfig{
		DocumentationURI: optional.NewString("https://github.com/finleap-connect/monoskope/tree/main/docs/operation/scim"),
		AuthenticationSchemes: []scim.AuthenticationScheme{
			{
				Type:        scim.AuthenticationTypeOauthBearerToken,
				Name:        "bearer",
				Description: "Monoskope Bearer Token Authentication",
				Primary:     true,
			},
		},
		SupportFiltering: false,
		SupportPatch:     true,
	}
	return config
}
