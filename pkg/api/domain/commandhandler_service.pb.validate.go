// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: api/domain/commandhandler_service.proto

package domain

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

// Validate checks the field values on PermissionModel with the rules defined
// in the proto definition for this message. If any rules are violated, the
// first error encountered is returned, or nil if there are no violations.
func (m *PermissionModel) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on PermissionModel with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// PermissionModelMultiError, or nil if none found.
func (m *PermissionModel) ValidateAll() error {
	return m.validate(true)
}

func (m *PermissionModel) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	for idx, item := range m.GetRoles() {
		_, _ = idx, item

		if len(item) > 60 {
			err := PermissionModelValidationError{
				field:  fmt.Sprintf("Roles[%v]", idx),
				reason: "value length must be at most 60 bytes",
			}
			if !all {
				return err
			}
			errors = append(errors, err)
		}

		if !_PermissionModel_Roles_Pattern.MatchString(item) {
			err := PermissionModelValidationError{
				field:  fmt.Sprintf("Roles[%v]", idx),
				reason: "value does not match regex pattern \"^[a-zA-Z][A-Za-z0-9_-]+$\"",
			}
			if !all {
				return err
			}
			errors = append(errors, err)
		}

	}

	for idx, item := range m.GetScopes() {
		_, _ = idx, item

		if len(item) > 60 {
			err := PermissionModelValidationError{
				field:  fmt.Sprintf("Scopes[%v]", idx),
				reason: "value length must be at most 60 bytes",
			}
			if !all {
				return err
			}
			errors = append(errors, err)
		}

		if !_PermissionModel_Scopes_Pattern.MatchString(item) {
			err := PermissionModelValidationError{
				field:  fmt.Sprintf("Scopes[%v]", idx),
				reason: "value does not match regex pattern \"^[a-zA-Z][A-Za-z0-9_-]+$\"",
			}
			if !all {
				return err
			}
			errors = append(errors, err)
		}

	}

	if len(errors) > 0 {
		return PermissionModelMultiError(errors)
	}
	return nil
}

// PermissionModelMultiError is an error wrapping multiple validation errors
// returned by PermissionModel.ValidateAll() if the designated constraints
// aren't met.
type PermissionModelMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m PermissionModelMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m PermissionModelMultiError) AllErrors() []error { return m }

// PermissionModelValidationError is the validation error returned by
// PermissionModel.Validate if the designated constraints aren't met.
type PermissionModelValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e PermissionModelValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e PermissionModelValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e PermissionModelValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e PermissionModelValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e PermissionModelValidationError) ErrorName() string { return "PermissionModelValidationError" }

// Error satisfies the builtin error interface
func (e PermissionModelValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sPermissionModel.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = PermissionModelValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = PermissionModelValidationError{}

var _PermissionModel_Roles_Pattern = regexp.MustCompile("^[a-zA-Z][A-Za-z0-9_-]+$")

var _PermissionModel_Scopes_Pattern = regexp.MustCompile("^[a-zA-Z][A-Za-z0-9_-]+$")

// Validate checks the field values on PolicyOverview with the rules defined in
// the proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *PolicyOverview) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on PolicyOverview with the rules defined
// in the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in PolicyOverviewMultiError,
// or nil if none found.
func (m *PolicyOverview) ValidateAll() error {
	return m.validate(true)
}

func (m *PolicyOverview) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	for idx, item := range m.GetPolicies() {
		_, _ = idx, item

		if all {
			switch v := interface{}(item).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, PolicyOverviewValidationError{
						field:  fmt.Sprintf("Policies[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, PolicyOverviewValidationError{
						field:  fmt.Sprintf("Policies[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return PolicyOverviewValidationError{
					field:  fmt.Sprintf("Policies[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if len(errors) > 0 {
		return PolicyOverviewMultiError(errors)
	}
	return nil
}

// PolicyOverviewMultiError is an error wrapping multiple validation errors
// returned by PolicyOverview.ValidateAll() if the designated constraints
// aren't met.
type PolicyOverviewMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m PolicyOverviewMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m PolicyOverviewMultiError) AllErrors() []error { return m }

// PolicyOverviewValidationError is the validation error returned by
// PolicyOverview.Validate if the designated constraints aren't met.
type PolicyOverviewValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e PolicyOverviewValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e PolicyOverviewValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e PolicyOverviewValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e PolicyOverviewValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e PolicyOverviewValidationError) ErrorName() string { return "PolicyOverviewValidationError" }

// Error satisfies the builtin error interface
func (e PolicyOverviewValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sPolicyOverview.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = PolicyOverviewValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = PolicyOverviewValidationError{}

// Validate checks the field values on Policy with the rules defined in the
// proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *Policy) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on Policy with the rules defined in the
// proto definition for this message. If any rules are violated, the result is
// a list of violation errors wrapped in PolicyMultiError, or nil if none found.
func (m *Policy) ValidateAll() error {
	return m.validate(true)
}

func (m *Policy) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if len(m.GetCommand()) > 60 {
		err := PolicyValidationError{
			field:  "Command",
			reason: "value length must be at most 60 bytes",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if !_Policy_Command_Pattern.MatchString(m.GetCommand()) {
		err := PolicyValidationError{
			field:  "Command",
			reason: "value does not match regex pattern \"^[a-zA-Z][A-Za-z0-9_-]+$\"",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if len(m.GetRole()) > 60 {
		err := PolicyValidationError{
			field:  "Role",
			reason: "value length must be at most 60 bytes",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if !_Policy_Role_Pattern.MatchString(m.GetRole()) {
		err := PolicyValidationError{
			field:  "Role",
			reason: "value does not match regex pattern \"^[a-zA-Z][A-Za-z0-9_-]+$\"",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if len(m.GetScope()) > 60 {
		err := PolicyValidationError{
			field:  "Scope",
			reason: "value length must be at most 60 bytes",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if !_Policy_Scope_Pattern.MatchString(m.GetScope()) {
		err := PolicyValidationError{
			field:  "Scope",
			reason: "value does not match regex pattern \"^[a-zA-Z][A-Za-z0-9_-]+$\"",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		return PolicyMultiError(errors)
	}
	return nil
}

// PolicyMultiError is an error wrapping multiple validation errors returned by
// Policy.ValidateAll() if the designated constraints aren't met.
type PolicyMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m PolicyMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m PolicyMultiError) AllErrors() []error { return m }

// PolicyValidationError is the validation error returned by Policy.Validate if
// the designated constraints aren't met.
type PolicyValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e PolicyValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e PolicyValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e PolicyValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e PolicyValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e PolicyValidationError) ErrorName() string { return "PolicyValidationError" }

// Error satisfies the builtin error interface
func (e PolicyValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sPolicy.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = PolicyValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = PolicyValidationError{}

var _Policy_Command_Pattern = regexp.MustCompile("^[a-zA-Z][A-Za-z0-9_-]+$")

var _Policy_Role_Pattern = regexp.MustCompile("^[a-zA-Z][A-Za-z0-9_-]+$")

var _Policy_Scope_Pattern = regexp.MustCompile("^[a-zA-Z][A-Za-z0-9_-]+$")
