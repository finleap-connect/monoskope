package scimserver

import (
	"github.com/elimity-com/scim"
	"github.com/elimity-com/scim/optional"
)

func NewProvierConfig() scim.ServiceProviderConfig {
	config := scim.ServiceProviderConfig{
		DocumentationURI: optional.NewString("www.example.com/scim"),
	}
	return config
}
