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
	"io/ioutil"
	"net/http"

	"github.com/elimity-com/scim"
	"github.com/elimity-com/scim/optional"
	"github.com/elimity-com/scim/schema"
	"github.com/finleap-connect/monoskope/pkg/logger"
	m8scim "github.com/finleap-connect/monoskope/pkg/scim"
)

func NewServer(config scim.ServiceProviderConfig, userHandler scim.ResourceHandler, groupHandler scim.ResourceHandler) scim.Server {
	resourceTypes := []scim.ResourceType{
		{
			ID:          optional.NewString("User"),
			Name:        "User",
			Endpoint:    "/Users",
			Description: optional.NewString("User Account"),
			Schema:      m8scim.MonoskopeUserSchema(),
			Handler:     userHandler,
		},
		{
			ID:          optional.NewString("Group"),
			Name:        "Group",
			Endpoint:    "/Groups",
			Description: optional.NewString("User Groups"),
			Schema:      schema.CoreGroupSchema(),
			Handler:     groupHandler,
		},
	}
	return scim.Server{
		Config:        config,
		ResourceTypes: resourceTypes,
	}
}

func logDebug(log logger.Logger, r *http.Request, attributes scim.ResourceAttributes, id string, params scim.ListRequestParams) {
	var err error
	var body []byte
	if r.Body != nil {
		body, err = ioutil.ReadAll(r.Body)
		if err != nil {
			log.V(logger.DebugLevel).Error(err, "Error reading body", "RequestURI", r.RequestURI)
			return
		}
	}
	log.V(logger.DebugLevel).Info("Received request", "RequestURI", r.RequestURI, "RequestBody", body, "RemoteAddr", r.RemoteAddr, "Referer", r.Referer(), "Attributes", attributes, "ID", id, "Params", params)
}
