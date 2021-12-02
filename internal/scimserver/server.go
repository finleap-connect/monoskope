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
	m8scim "github.com/finleap-connect/monoskope/pkg/scim"
)

func NewServer(config scim.ServiceProviderConfig, userHandler scim.ResourceHandler) scim.Server {
	resourceTypes := []scim.ResourceType{
		{
			ID:          optional.NewString("User"),
			Name:        "User",
			Endpoint:    "/Users",
			Description: optional.NewString("User Account"),
			Schema:      m8scim.MonoskopeUserSchema(),
			Handler:     userHandler,
		},
	}
	return scim.Server{
		Config:        config,
		ResourceTypes: resourceTypes,
	}
}
