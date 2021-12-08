// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: api/domain/eventdata/user.proto

package eventdata

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

	common "github.com/finleap-connect/monoskope/pkg/api/domain/common"
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

	_ = common.UserSource(0)
)

// Validate checks the field values on UserCreated with the rules defined in
// the proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *UserCreated) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on UserCreated with the rules defined in
// the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in UserCreatedMultiError, or
// nil if none found.
func (m *UserCreated) ValidateAll() error {
	return m.validate(true)
}

func (m *UserCreated) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Email

	// no validation rules for Name

	// no validation rules for Source

	if len(errors) > 0 {
		return UserCreatedMultiError(errors)
	}
	return nil
}

// UserCreatedMultiError is an error wrapping multiple validation errors
// returned by UserCreated.ValidateAll() if the designated constraints aren't met.
type UserCreatedMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m UserCreatedMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m UserCreatedMultiError) AllErrors() []error { return m }

// UserCreatedValidationError is the validation error returned by
// UserCreated.Validate if the designated constraints aren't met.
type UserCreatedValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e UserCreatedValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e UserCreatedValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e UserCreatedValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e UserCreatedValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e UserCreatedValidationError) ErrorName() string { return "UserCreatedValidationError" }

// Error satisfies the builtin error interface
func (e UserCreatedValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sUserCreated.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = UserCreatedValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = UserCreatedValidationError{}

// Validate checks the field values on UserRoleAdded with the rules defined in
// the proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *UserRoleAdded) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on UserRoleAdded with the rules defined
// in the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in UserRoleAddedMultiError, or
// nil if none found.
func (m *UserRoleAdded) ValidateAll() error {
	return m.validate(true)
}

func (m *UserRoleAdded) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for UserId

	// no validation rules for Role

	// no validation rules for Scope

	// no validation rules for Resource

	if len(errors) > 0 {
		return UserRoleAddedMultiError(errors)
	}
	return nil
}

// UserRoleAddedMultiError is an error wrapping multiple validation errors
// returned by UserRoleAdded.ValidateAll() if the designated constraints
// aren't met.
type UserRoleAddedMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m UserRoleAddedMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m UserRoleAddedMultiError) AllErrors() []error { return m }

// UserRoleAddedValidationError is the validation error returned by
// UserRoleAdded.Validate if the designated constraints aren't met.
type UserRoleAddedValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e UserRoleAddedValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e UserRoleAddedValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e UserRoleAddedValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e UserRoleAddedValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e UserRoleAddedValidationError) ErrorName() string { return "UserRoleAddedValidationError" }

// Error satisfies the builtin error interface
func (e UserRoleAddedValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sUserRoleAdded.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = UserRoleAddedValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = UserRoleAddedValidationError{}

// Validate checks the field values on UserUpdated with the rules defined in
// the proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *UserUpdated) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on UserUpdated with the rules defined in
// the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in UserUpdatedMultiError, or
// nil if none found.
func (m *UserUpdated) ValidateAll() error {
	return m.validate(true)
}

func (m *UserUpdated) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Name

	if len(errors) > 0 {
		return UserUpdatedMultiError(errors)
	}
	return nil
}

// UserUpdatedMultiError is an error wrapping multiple validation errors
// returned by UserUpdated.ValidateAll() if the designated constraints aren't met.
type UserUpdatedMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m UserUpdatedMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m UserUpdatedMultiError) AllErrors() []error { return m }

// UserUpdatedValidationError is the validation error returned by
// UserUpdated.Validate if the designated constraints aren't met.
type UserUpdatedValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e UserUpdatedValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e UserUpdatedValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e UserUpdatedValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e UserUpdatedValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e UserUpdatedValidationError) ErrorName() string { return "UserUpdatedValidationError" }

// Error satisfies the builtin error interface
func (e UserUpdatedValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sUserUpdated.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = UserUpdatedValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = UserUpdatedValidationError{}
