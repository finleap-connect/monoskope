// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: api/gateway/service.proto

package gateway

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

// Validate checks the field values on AuthState with the rules defined in the
// proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *AuthState) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on AuthState with the rules defined in
// the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in AuthStateMultiError, or nil
// if none found.
func (m *AuthState) ValidateAll() error {
	return m.validate(true)
}

func (m *AuthState) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for CallbackUrl

	if len(errors) > 0 {
		return AuthStateMultiError(errors)
	}
	return nil
}

// AuthStateMultiError is an error wrapping multiple validation errors returned
// by AuthState.ValidateAll() if the designated constraints aren't met.
type AuthStateMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m AuthStateMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m AuthStateMultiError) AllErrors() []error { return m }

// AuthStateValidationError is the validation error returned by
// AuthState.Validate if the designated constraints aren't met.
type AuthStateValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e AuthStateValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e AuthStateValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e AuthStateValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e AuthStateValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e AuthStateValidationError) ErrorName() string { return "AuthStateValidationError" }

// Error satisfies the builtin error interface
func (e AuthStateValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sAuthState.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = AuthStateValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = AuthStateValidationError{}

// Validate checks the field values on AuthInformation with the rules defined
// in the proto definition for this message. If any rules are violated, the
// first error encountered is returned, or nil if there are no violations.
func (m *AuthInformation) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on AuthInformation with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// AuthInformationMultiError, or nil if none found.
func (m *AuthInformation) ValidateAll() error {
	return m.validate(true)
}

func (m *AuthInformation) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for AuthCodeUrl

	// no validation rules for State

	if len(errors) > 0 {
		return AuthInformationMultiError(errors)
	}
	return nil
}

// AuthInformationMultiError is an error wrapping multiple validation errors
// returned by AuthInformation.ValidateAll() if the designated constraints
// aren't met.
type AuthInformationMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m AuthInformationMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m AuthInformationMultiError) AllErrors() []error { return m }

// AuthInformationValidationError is the validation error returned by
// AuthInformation.Validate if the designated constraints aren't met.
type AuthInformationValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e AuthInformationValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e AuthInformationValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e AuthInformationValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e AuthInformationValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e AuthInformationValidationError) ErrorName() string { return "AuthInformationValidationError" }

// Error satisfies the builtin error interface
func (e AuthInformationValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sAuthInformation.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = AuthInformationValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = AuthInformationValidationError{}

// Validate checks the field values on AuthCode with the rules defined in the
// proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *AuthCode) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on AuthCode with the rules defined in
// the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in AuthCodeMultiError, or nil
// if none found.
func (m *AuthCode) ValidateAll() error {
	return m.validate(true)
}

func (m *AuthCode) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Code

	// no validation rules for State

	// no validation rules for CallbackUrl

	if len(errors) > 0 {
		return AuthCodeMultiError(errors)
	}
	return nil
}

// AuthCodeMultiError is an error wrapping multiple validation errors returned
// by AuthCode.ValidateAll() if the designated constraints aren't met.
type AuthCodeMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m AuthCodeMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m AuthCodeMultiError) AllErrors() []error { return m }

// AuthCodeValidationError is the validation error returned by
// AuthCode.Validate if the designated constraints aren't met.
type AuthCodeValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e AuthCodeValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e AuthCodeValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e AuthCodeValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e AuthCodeValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e AuthCodeValidationError) ErrorName() string { return "AuthCodeValidationError" }

// Error satisfies the builtin error interface
func (e AuthCodeValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sAuthCode.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = AuthCodeValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = AuthCodeValidationError{}

// Validate checks the field values on AuthResponse with the rules defined in
// the proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *AuthResponse) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on AuthResponse with the rules defined
// in the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in AuthResponseMultiError, or
// nil if none found.
func (m *AuthResponse) ValidateAll() error {
	return m.validate(true)
}

func (m *AuthResponse) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for AccessToken

	if all {
		switch v := interface{}(m.GetExpiry()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, AuthResponseValidationError{
					field:  "Expiry",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, AuthResponseValidationError{
					field:  "Expiry",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetExpiry()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return AuthResponseValidationError{
				field:  "Expiry",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	// no validation rules for Username

	if len(errors) > 0 {
		return AuthResponseMultiError(errors)
	}
	return nil
}

// AuthResponseMultiError is an error wrapping multiple validation errors
// returned by AuthResponse.ValidateAll() if the designated constraints aren't met.
type AuthResponseMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m AuthResponseMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m AuthResponseMultiError) AllErrors() []error { return m }

// AuthResponseValidationError is the validation error returned by
// AuthResponse.Validate if the designated constraints aren't met.
type AuthResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e AuthResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e AuthResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e AuthResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e AuthResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e AuthResponseValidationError) ErrorName() string { return "AuthResponseValidationError" }

// Error satisfies the builtin error interface
func (e AuthResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sAuthResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = AuthResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = AuthResponseValidationError{}

// Validate checks the field values on ClusterAuthTokenRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *ClusterAuthTokenRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on ClusterAuthTokenRequest with the
// rules defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// ClusterAuthTokenRequestMultiError, or nil if none found.
func (m *ClusterAuthTokenRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *ClusterAuthTokenRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for ClusterId

	// no validation rules for Role

	if len(errors) > 0 {
		return ClusterAuthTokenRequestMultiError(errors)
	}
	return nil
}

// ClusterAuthTokenRequestMultiError is an error wrapping multiple validation
// errors returned by ClusterAuthTokenRequest.ValidateAll() if the designated
// constraints aren't met.
type ClusterAuthTokenRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m ClusterAuthTokenRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m ClusterAuthTokenRequestMultiError) AllErrors() []error { return m }

// ClusterAuthTokenRequestValidationError is the validation error returned by
// ClusterAuthTokenRequest.Validate if the designated constraints aren't met.
type ClusterAuthTokenRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e ClusterAuthTokenRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e ClusterAuthTokenRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e ClusterAuthTokenRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e ClusterAuthTokenRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e ClusterAuthTokenRequestValidationError) ErrorName() string {
	return "ClusterAuthTokenRequestValidationError"
}

// Error satisfies the builtin error interface
func (e ClusterAuthTokenRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sClusterAuthTokenRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = ClusterAuthTokenRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = ClusterAuthTokenRequestValidationError{}

// Validate checks the field values on ClusterAuthTokenResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *ClusterAuthTokenResponse) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on ClusterAuthTokenResponse with the
// rules defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// ClusterAuthTokenResponseMultiError, or nil if none found.
func (m *ClusterAuthTokenResponse) ValidateAll() error {
	return m.validate(true)
}

func (m *ClusterAuthTokenResponse) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for AccessToken

	if all {
		switch v := interface{}(m.GetExpiry()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, ClusterAuthTokenResponseValidationError{
					field:  "Expiry",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, ClusterAuthTokenResponseValidationError{
					field:  "Expiry",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetExpiry()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return ClusterAuthTokenResponseValidationError{
				field:  "Expiry",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if len(errors) > 0 {
		return ClusterAuthTokenResponseMultiError(errors)
	}
	return nil
}

// ClusterAuthTokenResponseMultiError is an error wrapping multiple validation
// errors returned by ClusterAuthTokenResponse.ValidateAll() if the designated
// constraints aren't met.
type ClusterAuthTokenResponseMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m ClusterAuthTokenResponseMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m ClusterAuthTokenResponseMultiError) AllErrors() []error { return m }

// ClusterAuthTokenResponseValidationError is the validation error returned by
// ClusterAuthTokenResponse.Validate if the designated constraints aren't met.
type ClusterAuthTokenResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e ClusterAuthTokenResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e ClusterAuthTokenResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e ClusterAuthTokenResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e ClusterAuthTokenResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e ClusterAuthTokenResponseValidationError) ErrorName() string {
	return "ClusterAuthTokenResponseValidationError"
}

// Error satisfies the builtin error interface
func (e ClusterAuthTokenResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sClusterAuthTokenResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = ClusterAuthTokenResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = ClusterAuthTokenResponseValidationError{}
