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
	"encoding/json"
	"io"
	"net/http"

	"github.com/elimity-com/scim"
	scim_errors "github.com/elimity-com/scim/errors"
	"github.com/finleap-connect/monoskope/pkg/api/domain"
	cmdData "github.com/finleap-connect/monoskope/pkg/api/domain/commanddata"
	"github.com/finleap-connect/monoskope/pkg/api/domain/projections"
	"github.com/finleap-connect/monoskope/pkg/api/eventsourcing"
	cmd "github.com/finleap-connect/monoskope/pkg/domain/commands"
	commandTypes "github.com/finleap-connect/monoskope/pkg/domain/constants/commands"
	"github.com/finleap-connect/monoskope/pkg/domain/errors"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type userHandler struct {
	cmdHandlerClient eventsourcing.CommandHandlerClient
	userClient       domain.UserClient
}

// NewUserHandler creates a new scim.ResourceHandler for handling User resources
func NewUserHandler(cmdHandlerClient eventsourcing.CommandHandlerClient, userClient domain.UserClient) scim.ResourceHandler {
	return &userHandler{
		cmdHandlerClient, userClient,
	}
}

// Create stores given attributes. Returns a resource with the attributes that are stored and a (new) unique identifier.
func (h *userHandler) Create(r *http.Request, attributes scim.ResourceAttributes) (scim.Resource, error) {
	var err error
	var commandData *cmdData.CreateUserCommandData

	command := cmd.CreateCommand(uuid.Nil, commandTypes.CreateUser)

	if commandData, err = toCreateUserCommandDataFromAttributes(attributes); err != nil {
		return scim.Resource{}, scim_errors.ScimError{
			Status: http.StatusInternalServerError,
			Detail: err.Error(),
		}
	}

	_, err = cmd.AddCommandData(command, commandData)
	if err != nil {
		return scim.Resource{}, scim_errors.ScimError{
			Status: http.StatusInternalServerError,
			Detail: err.Error(),
		}
	}

	_, err = h.cmdHandlerClient.Execute(r.Context(), command)
	if err != nil {
		return scim.Resource{}, scim_errors.ScimError{
			Status: http.StatusInternalServerError,
			Detail: err.Error(),
		}
	}

	return scim.Resource{}, nil
}

// Get returns the resource corresponding with the given identifier.
func (h *userHandler) Get(r *http.Request, id string) (scim.Resource, error) {
	user, err := h.userClient.GetById(r.Context(), wrapperspb.String(id))
	if err != nil {
		err = errors.TranslateFromGrpcError(err)
		if err == errors.ErrUserNotFound {
			return scim.Resource{}, scim_errors.ScimError{
				Status: http.StatusNotFound,
			}
		}
		return scim.Resource{}, scim_errors.ScimError{
			Status: http.StatusInternalServerError,
			Detail: err.Error(),
		}
	}
	return toScimUser(user), nil
}

// GetAll returns a paginated list of resources.
// An empty list of resources will be represented as `null` in the JSON response if `nil` is assigned to the
// Page.Resources. Otherwise, is an empty slice is assigned, an empty list will be represented as `[]`.
func (h *userHandler) GetAll(r *http.Request, params scim.ListRequestParams) (scim.Page, error) {
	// Get total user count intially
	userCount, err := h.userClient.GetCount(r.Context(), &domain.GetCountRequest{IncludeDeleted: true})
	if err != nil {
		return scim.Page{}, scim_errors.ScimError{
			Status: http.StatusInternalServerError,
			Detail: err.Error(),
		}
	}

	// If count is less than one just return total count
	if params.Count < 1 {
		return scim.Page{
			TotalResults: int(userCount.Count),
		}, nil
	}

	// Get stream of users
	userStream, err := h.userClient.GetAll(r.Context(), &domain.GetAllRequest{IncludeDeleted: true})
	if err != nil {
		err = errors.TranslateFromGrpcError(err)
		return scim.Page{}, scim_errors.ScimError{
			Status: http.StatusInternalServerError,
			Detail: err.Error(),
		}
	}

	// Seek through the stream
	resources := make([]scim.Resource, 0)
	i := 1
	for {
		if i > (params.StartIndex + params.Count - 1) {
			break // We're done
		}

		// Read next
		user, err := userStream.Recv()

		// End of stream
		if err == io.EOF {
			break // No further users to query
		}

		if err != nil { // Some other error
			return scim.Page{}, scim_errors.ScimError{
				Status: http.StatusInternalServerError,
				Detail: err.Error(),
			}
		}

		// Skip users which are not in the current page
		if i >= params.StartIndex {
			resources = append(resources, toScimUser(user))
		}
		i++
	}

	return scim.Page{
		TotalResults: int(userCount.Count),
		Resources:    resources,
	}, nil
}

// Replace replaces ALL existing attributes of the resource with given identifier. Given attributes that are empty
// are to be deleted. Returns a resource with the attributes that are stored.
func (h *userHandler) Replace(r *http.Request, id string, attributes scim.ResourceAttributes) (scim.Resource, error) {
	return scim.Resource{}, scim_errors.ScimError{
		Status: http.StatusNotImplemented,
	}
}

// Delete removes the resource with corresponding ID.
func (h *userHandler) Delete(r *http.Request, id string) error {
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
func (h *userHandler) Patch(r *http.Request, id string, operations []scim.PatchOperation) (scim.Resource, error) {
	return scim.Resource{}, scim_errors.ScimError{
		Status: http.StatusNotImplemented,
	}
}

const (
	SCIM_USER_USERNAME = "userName"
)

// toScimUser converts a projections.User to it's scim.Resource representation
func toScimUser(user *projections.User) scim.Resource {
	created := user.Metadata.Created.AsTime()
	lastModified := user.Metadata.LastModified.AsTime()
	deleted := user.Metadata.Deleted.IsValid()
	return scim.Resource{
		ID: user.Id,
		Meta: scim.Meta{
			Created:      &created,
			LastModified: &lastModified,
		},
		Attributes: scim.ResourceAttributes{
			SCIM_USER_USERNAME: user.Name,
			"active":           !deleted,
			"emails": []interface{}{
				map[string]interface{}{
					"primary": true,
					"value":   user.Email,
				},
			},
		},
	}
}

type emailsResource struct {
	Primary bool   `json:"primary"`
	Value   string `json:"value"`
}

type userResource struct {
	UserName string           `json:"userName"`
	Active   bool             `json:"active"`
	Emails   []emailsResource `json:"emails"`
}

// newUserResourceFromAttributes converts the attributes to json and back to an instance of the userResource struct
func newUserResourceFromAttributes(attributes scim.ResourceAttributes) (*userResource, error) {
	attributesJson, err := json.Marshal(attributes)
	if err != nil {
		return nil, err
	}
	userResource := new(userResource)
	err = json.Unmarshal(attributesJson, userResource)
	if err != nil {
		return nil, err
	}
	return userResource, nil
}

// toScimUser converts a projections.User to it's scim.Resource representation
func toCreateUserCommandDataFromAttributes(attributes scim.ResourceAttributes) (*cmdData.CreateUserCommandData, error) {
	userResource, err := newUserResourceFromAttributes(attributes)
	if err != nil {
		return nil, err
	}

	emailAddress := ""
	for _, email := range userResource.Emails {
		emailAddress = email.Value
		if email.Primary {
			break
		}
	}

	return &cmdData.CreateUserCommandData{
		Name:  userResource.UserName,
		Email: emailAddress,
	}, err
}
