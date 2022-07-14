// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: api/domain/commanddata/clustermapping.proto

package commanddata

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"net/mail"
	"net/url"
	"regexp"
	"sort"
	"strings"
	"time"
	"unicode/utf8"

	"google.golang.org/protobuf/types/known/anypb"
)

// ensure the imports are used
var (
	_ = bytes.MinRead
	_ = errors.New("")
	_ = fmt.Print
	_ = utf8.UTFMax
	_ = (*regexp.Regexp)(nil)
	_ = (*strings.Reader)(nil)
	_ = net.IPv4len
	_ = time.Duration(0)
	_ = (*url.URL)(nil)
	_ = (*mail.Address)(nil)
	_ = anypb.Any{}
	_ = sort.Sort
)

// define the regex for a UUID once up-front
var _clustermapping_uuidPattern = regexp.MustCompile("^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$")

// Validate checks the field values on CreateTenantClusterBindingCommandData
// with the rules defined in the proto definition for this message. If any
// rules are violated, the first error encountered is returned, or nil if
// there are no violations.
func (m *CreateTenantClusterBindingCommandData) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on CreateTenantClusterBindingCommandData
// with the rules defined in the proto definition for this message. If any
// rules are violated, the result is a list of violation errors wrapped in
// CreateTenantClusterBindingCommandDataMultiError, or nil if none found.
func (m *CreateTenantClusterBindingCommandData) ValidateAll() error {
	return m.validate(true)
}

func (m *CreateTenantClusterBindingCommandData) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if err := m._validateUuid(m.GetTenantId()); err != nil {
		err = CreateTenantClusterBindingCommandDataValidationError{
			field:  "TenantId",
			reason: "value must be a valid UUID",
			cause:  err,
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if err := m._validateUuid(m.GetClusterId()); err != nil {
		err = CreateTenantClusterBindingCommandDataValidationError{
			field:  "ClusterId",
			reason: "value must be a valid UUID",
			cause:  err,
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		return CreateTenantClusterBindingCommandDataMultiError(errors)
	}

	return nil
}

func (m *CreateTenantClusterBindingCommandData) _validateUuid(uuid string) error {
	if matched := _clustermapping_uuidPattern.MatchString(uuid); !matched {
		return errors.New("invalid uuid format")
	}

	return nil
}

// CreateTenantClusterBindingCommandDataMultiError is an error wrapping
// multiple validation errors returned by
// CreateTenantClusterBindingCommandData.ValidateAll() if the designated
// constraints aren't met.
type CreateTenantClusterBindingCommandDataMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m CreateTenantClusterBindingCommandDataMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m CreateTenantClusterBindingCommandDataMultiError) AllErrors() []error { return m }

// CreateTenantClusterBindingCommandDataValidationError is the validation error
// returned by CreateTenantClusterBindingCommandData.Validate if the
// designated constraints aren't met.
type CreateTenantClusterBindingCommandDataValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e CreateTenantClusterBindingCommandDataValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e CreateTenantClusterBindingCommandDataValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e CreateTenantClusterBindingCommandDataValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e CreateTenantClusterBindingCommandDataValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e CreateTenantClusterBindingCommandDataValidationError) ErrorName() string {
	return "CreateTenantClusterBindingCommandDataValidationError"
}

// Error satisfies the builtin error interface
func (e CreateTenantClusterBindingCommandDataValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sCreateTenantClusterBindingCommandData.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = CreateTenantClusterBindingCommandDataValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = CreateTenantClusterBindingCommandDataValidationError{}
