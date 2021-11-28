// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: api/domain/projections/user.proto

package projections

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
var _user_uuidPattern = regexp.MustCompile("^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$")

// Validate checks the field values on User with the rules defined in the proto
// definition for this message. If any rules are violated, the first error
// encountered is returned, or nil if there are no violations.
func (m *User) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on User with the rules defined in the
// proto definition for this message. If any rules are violated, the result is
// a list of violation errors wrapped in UserMultiError, or nil if none found.
func (m *User) ValidateAll() error {
	return m.validate(true)
}

func (m *User) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if err := m._validateUuid(m.GetId()); err != nil {
		err = UserValidationError{
			field:  "Id",
			reason: "value must be a valid UUID",
			cause:  err,
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if len(m.GetName()) > 150 {
		err := UserValidationError{
			field:  "Name",
			reason: "value length must be at most 150 bytes",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if err := m._validateEmail(m.GetEmail()); err != nil {
		err = UserValidationError{
			field:  "Email",
			reason: "value must be a valid email address",
			cause:  err,
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	for idx, item := range m.GetRoles() {
		_, _ = idx, item

		if all {
			switch v := interface{}(item).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, UserValidationError{
						field:  fmt.Sprintf("Roles[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, UserValidationError{
						field:  fmt.Sprintf("Roles[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return UserValidationError{
					field:  fmt.Sprintf("Roles[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if all {
		switch v := interface{}(m.GetMetadata()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, UserValidationError{
					field:  "Metadata",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, UserValidationError{
					field:  "Metadata",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetMetadata()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return UserValidationError{
				field:  "Metadata",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if len(errors) > 0 {
		return UserMultiError(errors)
	}
	return nil
}

func (m *User) _validateHostname(host string) error {
	s := strings.ToLower(strings.TrimSuffix(host, "."))

	if len(host) > 253 {
		return errors.New("hostname cannot exceed 253 characters")
	}

	for _, part := range strings.Split(s, ".") {
		if l := len(part); l == 0 || l > 63 {
			return errors.New("hostname part must be non-empty and cannot exceed 63 characters")
		}

		if part[0] == '-' {
			return errors.New("hostname parts cannot begin with hyphens")
		}

		if part[len(part)-1] == '-' {
			return errors.New("hostname parts cannot end with hyphens")
		}

		for _, r := range part {
			if (r < 'a' || r > 'z') && (r < '0' || r > '9') && r != '-' {
				return fmt.Errorf("hostname parts can only contain alphanumeric characters or hyphens, got %q", string(r))
			}
		}
	}

	return nil
}

func (m *User) _validateEmail(addr string) error {
	a, err := mail.ParseAddress(addr)
	if err != nil {
		return err
	}
	addr = a.Address

	if len(addr) > 254 {
		return errors.New("email addresses cannot exceed 254 characters")
	}

	parts := strings.SplitN(addr, "@", 2)

	if len(parts[0]) > 64 {
		return errors.New("email address local phrase cannot exceed 64 characters")
	}

	return m._validateHostname(parts[1])
}

func (m *User) _validateUuid(uuid string) error {
	if matched := _user_uuidPattern.MatchString(uuid); !matched {
		return errors.New("invalid uuid format")
	}

	return nil
}

// UserMultiError is an error wrapping multiple validation errors returned by
// User.ValidateAll() if the designated constraints aren't met.
type UserMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m UserMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m UserMultiError) AllErrors() []error { return m }

// UserValidationError is the validation error returned by User.Validate if the
// designated constraints aren't met.
type UserValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e UserValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e UserValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e UserValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e UserValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e UserValidationError) ErrorName() string { return "UserValidationError" }

// Error satisfies the builtin error interface
func (e UserValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sUser.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = UserValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = UserValidationError{}

// Validate checks the field values on UserRoleBinding with the rules defined
// in the proto definition for this message. If any rules are violated, the
// first error encountered is returned, or nil if there are no violations.
func (m *UserRoleBinding) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on UserRoleBinding with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// UserRoleBindingMultiError, or nil if none found.
func (m *UserRoleBinding) ValidateAll() error {
	return m.validate(true)
}

func (m *UserRoleBinding) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if err := m._validateUuid(m.GetId()); err != nil {
		err = UserRoleBindingValidationError{
			field:  "Id",
			reason: "value must be a valid UUID",
			cause:  err,
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if err := m._validateUuid(m.GetUserId()); err != nil {
		err = UserRoleBindingValidationError{
			field:  "UserId",
			reason: "value must be a valid UUID",
			cause:  err,
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if len(m.GetRole()) > 60 {
		err := UserRoleBindingValidationError{
			field:  "Role",
			reason: "value length must be at most 60 bytes",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if !_UserRoleBinding_Role_Pattern.MatchString(m.GetRole()) {
		err := UserRoleBindingValidationError{
			field:  "Role",
			reason: "value does not match regex pattern \"^[a-zA-Z_]+$\"",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if len(m.GetScope()) > 60 {
		err := UserRoleBindingValidationError{
			field:  "Scope",
			reason: "value length must be at most 60 bytes",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if !_UserRoleBinding_Scope_Pattern.MatchString(m.GetScope()) {
		err := UserRoleBindingValidationError{
			field:  "Scope",
			reason: "value does not match regex pattern \"^[a-zA-Z_]+$\"",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if len(m.GetResource()) > 60 {
		err := UserRoleBindingValidationError{
			field:  "Resource",
			reason: "value length must be at most 60 bytes",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if !_UserRoleBinding_Resource_Pattern.MatchString(m.GetResource()) {
		err := UserRoleBindingValidationError{
			field:  "Resource",
			reason: "value does not match regex pattern \"^[a-zA-Z_]+$\"",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if all {
		switch v := interface{}(m.GetMetadata()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, UserRoleBindingValidationError{
					field:  "Metadata",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, UserRoleBindingValidationError{
					field:  "Metadata",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetMetadata()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return UserRoleBindingValidationError{
				field:  "Metadata",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if len(errors) > 0 {
		return UserRoleBindingMultiError(errors)
	}
	return nil
}

func (m *UserRoleBinding) _validateUuid(uuid string) error {
	if matched := _user_uuidPattern.MatchString(uuid); !matched {
		return errors.New("invalid uuid format")
	}

	return nil
}

// UserRoleBindingMultiError is an error wrapping multiple validation errors
// returned by UserRoleBinding.ValidateAll() if the designated constraints
// aren't met.
type UserRoleBindingMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m UserRoleBindingMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m UserRoleBindingMultiError) AllErrors() []error { return m }

// UserRoleBindingValidationError is the validation error returned by
// UserRoleBinding.Validate if the designated constraints aren't met.
type UserRoleBindingValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e UserRoleBindingValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e UserRoleBindingValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e UserRoleBindingValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e UserRoleBindingValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e UserRoleBindingValidationError) ErrorName() string { return "UserRoleBindingValidationError" }

// Error satisfies the builtin error interface
func (e UserRoleBindingValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sUserRoleBinding.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = UserRoleBindingValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = UserRoleBindingValidationError{}

var _UserRoleBinding_Role_Pattern = regexp.MustCompile("^[a-zA-Z_]+$")

var _UserRoleBinding_Scope_Pattern = regexp.MustCompile("^[a-zA-Z_]+$")

var _UserRoleBinding_Resource_Pattern = regexp.MustCompile("^[a-zA-Z_]+$")
