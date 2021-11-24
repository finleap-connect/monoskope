package scimserver

import (
	"net/http"

	"github.com/elimity-com/scim"
	"github.com/elimity-com/scim/errors"
)

type userHandler struct {
}

func NewUserHandler() scim.ResourceHandler {
	userHandler := new(userHandler)
	return userHandler
}

// Create stores given attributes. Returns a resource with the attributes that are stored and a (new) unique identifier.
func (h *userHandler) Create(r *http.Request, attributes scim.ResourceAttributes) (scim.Resource, error) {
	return scim.Resource{}, errors.ScimError{
		Status: http.StatusNotImplemented,
	}
}

// Get returns the resource corresponding with the given identifier.
func (h *userHandler) Get(r *http.Request, id string) (scim.Resource, error) {
	return scim.Resource{}, errors.ScimError{
		Status: http.StatusNotImplemented,
	}
}

// GetAll returns a paginated list of resources.
// An empty list of resources will be represented as `null` in the JSON response if `nil` is assigned to the
// Page.Resources. Otherwise, is an empty slice is assigned, an empty list will be represented as `[]`.
func (h *userHandler) GetAll(r *http.Request, params scim.ListRequestParams) (scim.Page, error) {
	return scim.Page{}, errors.ScimError{
		Status: http.StatusNotImplemented,
	}
}

// Replace replaces ALL existing attributes of the resource with given identifier. Given attributes that are empty
// are to be deleted. Returns a resource with the attributes that are stored.
func (h *userHandler) Replace(r *http.Request, id string, attributes scim.ResourceAttributes) (scim.Resource, error) {
	return scim.Resource{}, errors.ScimError{
		Status: http.StatusNotImplemented,
	}
}

// Delete removes the resource with corresponding ID.
func (h *userHandler) Delete(r *http.Request, id string) error {
	return errors.ScimError{
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
	return scim.Resource{}, errors.ScimError{
		Status: http.StatusNotImplemented,
	}
}
