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
	"net/http"

	"github.com/elimity-com/scim"
	scim_errors "github.com/elimity-com/scim/errors"
	"github.com/finleap-connect/monoskope/pkg/api/domain"
	"github.com/finleap-connect/monoskope/pkg/api/eventsourcing"
	"github.com/finleap-connect/monoskope/pkg/domain/constants/roles"
	es "github.com/finleap-connect/monoskope/pkg/eventsourcing"
	"github.com/finleap-connect/monoskope/pkg/logger"
	m8scim "github.com/finleap-connect/monoskope/pkg/scim"
	"github.com/google/uuid"
)

type groupHandler struct {
	cmdHandlerClient eventsourcing.CommandHandlerClient
	userClient       domain.UserClient
	log              logger.Logger
}

func NewGroupHandler(cmdHandlerClient eventsourcing.CommandHandlerClient, userClient domain.UserClient) scim.ResourceHandler {
	return &groupHandler{
		cmdHandlerClient, userClient, logger.WithName("scim-group-handler"),
	}
}

// Create stores given attributes. Returns a resource with the attributes that are stored and a (new) unique identifier.
func (h *groupHandler) Create(r *http.Request, attributes scim.ResourceAttributes) (scim.Resource, error) {
	logDebug(h.log, r, attributes, "", scim.ListRequestParams{})

	return scim.Resource{}, scim_errors.ScimError{
		Status: http.StatusNotImplemented,
	}
}

// Get returns the resource corresponding with the given identifier.
func (h *groupHandler) Get(r *http.Request, id string) (scim.Resource, error) {
	logDebug(h.log, r, scim.ResourceAttributes{}, id, scim.ListRequestParams{})

	return scim.Resource{}, scim_errors.ScimError{
		Status: http.StatusNotImplemented,
	}
}

// GetAll returns a paginated list of resources.
// An empty list of resources will be represented as `null` in the JSON response if `nil` is assigned to the
// Page.Resources. Otherwise, is an empty slice is assigned, an empty list will be represented as `[]`.
func (h *groupHandler) GetAll(r *http.Request, params scim.ListRequestParams) (scim.Page, error) {
	logDebug(h.log, r, scim.ResourceAttributes{}, "", params)

	roleCount := len(roles.AvailableRoles)

	// If count is less than one just return total count
	if params.Count < 1 {
		return scim.Page{
			TotalResults: roleCount,
		}, nil
	}

	resources := make([]scim.Resource, 0)
	// Seek through the stream
	i := 1
	for {
		if i > (params.StartIndex + params.Count - 1) {
			break // We're done
		}

		// Read next
		role := roles.AvailableRoles[i]

		// Skip users which are not in the current page
		if i >= params.StartIndex {
			resources = append(resources, toScimGroup(role))
		}
		i++
	}

	return scim.Page{
		TotalResults: int(roleCount),
		Resources:    resources,
	}, nil
}

// Replace replaces ALL existing attributes of the resource with given identifier. Given attributes that are empty
// are to be deleted. Returns a resource with the attributes that are stored.
func (h *groupHandler) Replace(r *http.Request, id string, attributes scim.ResourceAttributes) (scim.Resource, error) {
	logDebug(h.log, r, attributes, id, scim.ListRequestParams{})

	return scim.Resource{}, scim_errors.ScimError{
		Status: http.StatusNotImplemented,
	}
}

// Delete removes the resource with corresponding ID.
func (h *groupHandler) Delete(r *http.Request, id string) error {
	logDebug(h.log, r, scim.ResourceAttributes{}, id, scim.ListRequestParams{})

	return scim_errors.ScimError{
		Status: http.StatusNotImplemented,
	}
}

// Patch update one or more attributes of a SCIM resource using a sequence of
// operations to "add", "remove", or "replace" values.
// If you return no Resource.Attributes, a 204 No Content status code will be returned.
// This case is only valid in the following scenarios:
// 1. the Add/Replace operation should return No Content only when the value already exists AND is the same.
// 2. the Remove operation should return No Content when the value to be remove is already absent.
// More information in Section 3.5.2 of RFC 7644: https://tools.ietf.org/html/rfc7644#section-3.5.2
func (h *groupHandler) Patch(r *http.Request, id string, operations []scim.PatchOperation) (scim.Resource, error) {
	return scim.Resource{}, scim_errors.ScimError{
		Status: http.StatusNotImplemented,
	}
}

// toScimGroup converts a projections.UserRoleBinding to it's scim.Resource representation
func toScimGroup(role es.Role) scim.Resource {
	return scim.Resource{
		ID: uuid.NewSHA1(uuid.NameSpaceURL, []byte(role)).String(),
		Attributes: scim.ResourceAttributes{
			m8scim.DisplayNameAttribute: role,
		},
	}
}
