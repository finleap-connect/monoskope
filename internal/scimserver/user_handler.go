// Copyright 2022 Monoskope Authors
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
	"fmt"
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
	"github.com/finleap-connect/monoskope/pkg/domain/constants/roles"
	"github.com/finleap-connect/monoskope/pkg/domain/constants/users"
	"github.com/finleap-connect/monoskope/pkg/domain/errors"
	es_errors "github.com/finleap-connect/monoskope/pkg/eventsourcing/errors"
	"github.com/finleap-connect/monoskope/pkg/logger"
	m8scim "github.com/finleap-connect/monoskope/pkg/scim"
	"github.com/google/uuid"
	"github.com/scim2/filter-parser/v2"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type userHandler struct {
	cmdHandlerClient eventsourcing.CommandHandlerClient
	userClient       domain.UserClient
	log              logger.Logger
}

// NewUserHandler creates a new scim.ResourceHandler for handling User resources
func NewUserHandler(cmdHandlerClient eventsourcing.CommandHandlerClient, userClient domain.UserClient) scim.ResourceHandler {
	return &userHandler{
		cmdHandlerClient, userClient, logger.WithName("scim-user-handler"),
	}
}

func (h *userHandler) getBy(f func() (*projections.User, error)) (scim.Resource, error) {
	user, err := f()
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
	if user.Metadata.Deleted != nil {
		return scim.Resource{}, scim_errors.ScimError{
			Status: http.StatusNotFound,
		}
	}
	return toScimUser(user), nil
}

// Create stores given attributes. Returns a resource with the attributes that are stored and a (new) unique identifier.
func (h *userHandler) Create(r *http.Request, attributes scim.ResourceAttributes) (scim.Resource, error) {
	logDebug(h.log, r)

	ctx, err := users.CreateUserContextGrpc(r.Context(), users.SCIMServerUser)
	if err != nil {
		return scim.Resource{}, scim_errors.ScimError{
			Status: http.StatusInternalServerError,
			Detail: err.Error(),
		}
	}

	userAttributes, err := m8scim.NewUserAttribute(attributes)
	if err != nil {
		return scim.Resource{}, scim_errors.ScimError{
			Status: http.StatusInternalServerError,
			Detail: err.Error(),
		}
	}

	command := cmd.NewCommandWithData(uuid.Nil, commandTypes.CreateUser, &cmdData.CreateUserCommandData{
		Email: userAttributes.UserName,
		Name:  userAttributes.DisplayName,
	})
	if err != nil {
		return scim.Resource{}, scim_errors.ScimError{
			Status: http.StatusInternalServerError,
			Detail: err.Error(),
		}
	}

	reply, err := h.cmdHandlerClient.Execute(ctx, command)
	if err != nil {
		err = errors.TranslateFromGrpcError(err)
		if err == errors.ErrUserAlreadyExists {
			return h.getBy(func() (*projections.User, error) {
				return h.userClient.GetByEmail(r.Context(), wrapperspb.String(userAttributes.UserName))
			})
		}
		return scim.Resource{}, scim_errors.ScimError{
			Status: http.StatusInternalServerError,
			Detail: err.Error(),
		}
	}

	return scim.Resource{
		ID:         reply.AggregateId,
		Attributes: attributes,
	}, nil
}

// Get returns the resource corresponding with the given identifier.
func (h *userHandler) Get(r *http.Request, id string) (scim.Resource, error) {
	logDebug(h.log, r)
	return h.getBy(func() (*projections.User, error) {
		return h.userClient.GetById(r.Context(), wrapperspb.String(id))
	})
}

// GetAll returns a paginated list of resources.
// An empty list of resources will be represented as `null` in the JSON response if `nil` is assigned to the
// Page.Resources. Otherwise, is an empty slice is assigned, an empty list will be represented as `[]`.
func (h *userHandler) GetAll(r *http.Request, params scim.ListRequestParams) (scim.Page, error) {
	logDebug(h.log, r)

	// Get total user count initially
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

	var filterByName string
	if params.Filter != nil {
		switch e := params.Filter.(type) {
		case *filter.AttributeExpression:
			if e.AttributePath.AttributeName == m8scim.UserNameAttribute && e.Operator == filter.EQ {
				filterByName = e.CompareValue.(string)
			}
		default:
			err := fmt.Errorf("unknown expression type: %s", e)
			h.log.Error(err, "unknown expression type", "type", e)
			return scim.Page{}, scim_errors.ScimError{
				Status: http.StatusInternalServerError,
				Detail: err.Error(),
			}
		}
	}

	resources := make([]scim.Resource, 0)
	if len(filterByName) > 0 {
		user, err := h.userClient.GetByEmail(r.Context(), wrapperspb.String(filterByName))
		if err != nil {
			err = errors.TranslateFromGrpcError(err)
			if err == errors.ErrUserNotFound {
				return scim.Page{
					TotalResults: int(userCount.Count),
					Resources:    resources,
				}, nil
			}
			return scim.Page{}, scim_errors.ScimError{
				Status: http.StatusInternalServerError,
				Detail: err.Error(),
			}
		}
		resources = append(resources, toScimUser(user))
	} else {
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
	}

	return scim.Page{
		TotalResults: int(userCount.Count),
		Resources:    resources,
	}, nil
}

// Replace replaces ALL existing attributes of the resource with given identifier. Given attributes that are empty
// are to be deleted. Returns a resource with the attributes that are stored.
func (h *userHandler) Replace(r *http.Request, id string, attributes scim.ResourceAttributes) (scim.Resource, error) {
	logDebug(h.log, r)

	ctx, err := users.CreateUserContextGrpc(r.Context(), users.SCIMServerUser)
	if err != nil {
		return scim.Resource{}, scim_errors.ScimError{
			Status: http.StatusInternalServerError,
			Detail: err.Error(),
		}
	}

	user, err := h.userClient.GetById(ctx, wrapperspb.String(id))
	if err != nil {
		err = errors.TranslateFromGrpcError(err)
		if err == errors.ErrUserNotFound || err == es_errors.ErrProjectionNotFound {
			return scim.Resource{}, scim_errors.ScimError{
				Status: http.StatusNotFound,
			}
		}
		return scim.Resource{}, scim_errors.ScimError{
			Status: http.StatusInternalServerError,
			Detail: err.Error(),
		}
	}

	userAttributes, err := m8scim.NewUserAttribute(attributes)
	if err != nil {
		return scim.Resource{}, scim_errors.ScimError{
			Status: http.StatusInternalServerError,
			Detail: err.Error(),
		}
	}

	command := cmd.NewCommandWithData(uuid.MustParse(user.Id), commandTypes.UpdateUser, &cmdData.UpdateUserCommandData{
		Name: wrapperspb.String(userAttributes.DisplayName),
	})
	if err != nil {
		return scim.Resource{}, scim_errors.ScimError{
			Status: http.StatusInternalServerError,
			Detail: err.Error(),
		}
	}

	_, err = h.cmdHandlerClient.Execute(ctx, command)
	if err != nil {
		return scim.Resource{}, scim_errors.ScimError{
			Status: http.StatusInternalServerError,
			Detail: err.Error(),
		}
	}

	user.Name = userAttributes.DisplayName
	return toScimUser(user), nil
}

// Delete removes the resource with corresponding ID.
func (h *userHandler) Delete(r *http.Request, id string) error {
	logDebug(h.log, r)

	uid, err := uuid.Parse(id)
	if err != nil {
		h.log.Error(err, "Failed to parse id")
		return scim_errors.ScimError{
			Status: http.StatusBadRequest,
			Detail: err.Error(),
		}
	}

	ctx, err := users.CreateUserContextGrpc(r.Context(), users.SCIMServerUser)
	if err != nil {
		return scim_errors.ScimError{
			Status: http.StatusInternalServerError,
			Detail: err.Error(),
		}
	}

	_, err = h.cmdHandlerClient.Execute(ctx, cmd.NewCommand(uid, commandTypes.DeleteUser))
	if err != nil {
		err = errors.TranslateFromGrpcError(err)
		if err == errors.ErrDeleted {
			return nil
		}

		return scim_errors.ScimError{
			Status: http.StatusInternalServerError,
			Detail: err.Error(),
		}
	}

	return nil
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

// toScimUser converts a projections.User to it's scim.Resource representation
func toScimUser(user *projections.User) scim.Resource {
	created := user.Metadata.Created.AsTime()
	lastModified := user.Metadata.LastModified.AsTime()

	groups := make([]map[string]string, 0)
	for _, r := range user.Roles {
		role, err := roles.ToRole(r.Role)
		if err != nil {
			continue
		}

		groups = append(groups, map[string]string{
			m8scim.GroupMemberValueAttribute:   roles.IdFromRole(role).String(),
			m8scim.GroupMemberDisplayAttribute: r.Role,
		})
	}

	return scim.Resource{
		ID: user.Id,
		Meta: scim.Meta{
			Created:      &created,
			LastModified: &lastModified,
		},
		Attributes: scim.ResourceAttributes{
			m8scim.UserNameAttribute:    user.Email,
			m8scim.DisplayNameAttribute: user.Name,
			m8scim.GroupAttribute:       groups,
		},
	}
}
