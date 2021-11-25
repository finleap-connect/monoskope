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

package scimserver

import (
	"github.com/elimity-com/scim"
	"github.com/elimity-com/scim/optional"
)

// NewProvierConfig create a new scim.ServiceProviderConfig for Monoskope
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
