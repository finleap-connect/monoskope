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
