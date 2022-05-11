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
	"net/http"

	"github.com/elimity-com/scim"
	scim_errors "github.com/elimity-com/scim/errors"
	"github.com/finleap-connect/monoskope/internal/gateway/auth"
	"github.com/finleap-connect/monoskope/pkg/logger"
	"github.com/grpc-ecosystem/go-grpc-middleware/util/metautils"
	"google.golang.org/grpc/metadata"
)

type authHandler struct {
	nextHandler scim.ResourceHandler
	log         logger.Logger
}

// NewAuthHandler creates a new scim.ResourceHandler for handling auth headers as wrapper for other ResourceHandlers
func NewAuthHandler(nextHandler scim.ResourceHandler) scim.ResourceHandler {
	return &authHandler{
		nextHandler, logger.WithName("scim-auth-handler"),
	}
}

func (h *authHandler) withAuthContext(r *http.Request) (*http.Request, error) {
	token := r.Header.Get(auth.HeaderAuthorization)
	if token == "" {
		return nil, scim_errors.ScimError{
			Status: http.StatusUnauthorized,
			Detail: "Request unauthenticated with " + auth.AuthScheme,
		}
	}
	md := metadata.Pairs(auth.HeaderAuthorization, fmt.Sprintf("%s %v", auth.AuthScheme, token))
	nCtx := metautils.NiceMD(md).ToIncoming(r.Context())
	return r.WithContext(nCtx), nil
}

// Create stores given attributes. Returns a resource with the attributes that are stored and a (new) unique identifier.
func (h *authHandler) Create(r *http.Request, attributes scim.ResourceAttributes) (scim.Resource, error) {
	logDebug(h.log, r)
	r, err := h.withAuthContext(r)
	if err != nil {
		return scim.Resource{}, err
	}
	return h.nextHandler.Create(r, attributes)
}

// Get returns the resource corresponding with the given identifier.
func (h *authHandler) Get(r *http.Request, id string) (scim.Resource, error) {
	logDebug(h.log, r)
	r, err := h.withAuthContext(r)
	if err != nil {
		return scim.Resource{}, err
	}
	return h.nextHandler.Get(r, id)
}

// GetAll returns a paginated list of resources.
// An empty list of resources will be represented as `null` in the JSON response if `nil` is assigned to the
// Page.Resources. Otherwise, is an empty slice is assigned, an empty list will be represented as `[]`.
func (h *authHandler) GetAll(r *http.Request, params scim.ListRequestParams) (scim.Page, error) {
	logDebug(h.log, r)
	r, err := h.withAuthContext(r)
	if err != nil {
		return scim.Page{}, err
	}
	return h.nextHandler.GetAll(r, params)
}

// Replace replaces ALL existing attributes of the resource with given identifier. Given attributes that are empty
// are to be deleted. Returns a resource with the attributes that are stored.
func (h *authHandler) Replace(r *http.Request, id string, attributes scim.ResourceAttributes) (scim.Resource, error) {
	logDebug(h.log, r)
	r, err := h.withAuthContext(r)
	if err != nil {
		return scim.Resource{}, err
	}
	return h.nextHandler.Replace(r, id, attributes)
}

// Delete removes the resource with corresponding ID.
func (h *authHandler) Delete(r *http.Request, id string) error {
	logDebug(h.log, r)
	r, err := h.withAuthContext(r)
	if err != nil {
		return err
	}
	return h.nextHandler.Delete(r, id)

}

// Patch update one or more attributes of a SCIM resource using a sequence of
// operations to "add", "remove", or "replace" values.
// If you return no Resource.Attributes, a 204 No Content status code will be returned.
// This case is only valid in the following scenarios:
// 1. the Add/Replace operation should return No Content only when the value already exists AND is the same.
// 2. the Remove operation should return No Content when the value to be remove is already absent.
// More information in Section 3.5.2 of RFC 7644: https://tools.ietf.org/html/rfc7644#section-3.5.2
func (h *authHandler) Patch(r *http.Request, id string, operations []scim.PatchOperation) (scim.Resource, error) {
	logDebug(h.log, r)
	r, err := h.withAuthContext(r)
	if err != nil {
		return scim.Resource{}, err
	}
	return h.nextHandler.Patch(r, id, operations)
}
