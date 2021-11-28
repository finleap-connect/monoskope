// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: api/domain/eventdata/cluster.proto

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

// Validate checks the field values on ClusterCreated with the rules defined in
// the proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *ClusterCreated) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on ClusterCreated with the rules defined
// in the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in ClusterCreatedMultiError,
// or nil if none found.
func (m *ClusterCreated) ValidateAll() error {
	return m.validate(true)
}

func (m *ClusterCreated) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if len(m.GetName()) > 150 {
		err := ClusterCreatedValidationError{
			field:  "Name",
			reason: "value length must be at most 150 bytes",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if len(m.GetLabel()) > 60 {
		err := ClusterCreatedValidationError{
			field:  "Label",
			reason: "value length must be at most 60 bytes",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if !_ClusterCreated_Label_Pattern.MatchString(m.GetLabel()) {
		err := ClusterCreatedValidationError{
			field:  "Label",
			reason: "value does not match regex pattern \"^[a-zA-Z][A-Za-z0-9_-]+$\"",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if err := m._validateHostname(m.GetApiServerAddress()); err != nil {
		err = ClusterCreatedValidationError{
			field:  "ApiServerAddress",
			reason: "value must be a valid hostname",
			cause:  err,
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	// no validation rules for CaCertificateBundle

	if len(errors) > 0 {
		return ClusterCreatedMultiError(errors)
	}
	return nil
}

func (m *ClusterCreated) _validateHostname(host string) error {
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

// ClusterCreatedMultiError is an error wrapping multiple validation errors
// returned by ClusterCreated.ValidateAll() if the designated constraints
// aren't met.
type ClusterCreatedMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m ClusterCreatedMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m ClusterCreatedMultiError) AllErrors() []error { return m }

// ClusterCreatedValidationError is the validation error returned by
// ClusterCreated.Validate if the designated constraints aren't met.
type ClusterCreatedValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e ClusterCreatedValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e ClusterCreatedValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e ClusterCreatedValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e ClusterCreatedValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e ClusterCreatedValidationError) ErrorName() string { return "ClusterCreatedValidationError" }

// Error satisfies the builtin error interface
func (e ClusterCreatedValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sClusterCreated.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = ClusterCreatedValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = ClusterCreatedValidationError{}

var _ClusterCreated_Label_Pattern = regexp.MustCompile("^[a-zA-Z][A-Za-z0-9_-]+$")

// Validate checks the field values on ClusterCreatedV2 with the rules defined
// in the proto definition for this message. If any rules are violated, the
// first error encountered is returned, or nil if there are no violations.
func (m *ClusterCreatedV2) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on ClusterCreatedV2 with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// ClusterCreatedV2MultiError, or nil if none found.
func (m *ClusterCreatedV2) ValidateAll() error {
	return m.validate(true)
}

func (m *ClusterCreatedV2) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if len(m.GetName()) > 60 {
		err := ClusterCreatedV2ValidationError{
			field:  "Name",
			reason: "value length must be at most 60 bytes",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if !_ClusterCreatedV2_Name_Pattern.MatchString(m.GetName()) {
		err := ClusterCreatedV2ValidationError{
			field:  "Name",
			reason: "value does not match regex pattern \"^[a-zA-Z][A-Za-z0-9_-]+$\"",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if len(m.GetDisplayName()) > 150 {
		err := ClusterCreatedV2ValidationError{
			field:  "DisplayName",
			reason: "value length must be at most 150 bytes",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if err := m._validateHostname(m.GetApiServerAddress()); err != nil {
		err = ClusterCreatedV2ValidationError{
			field:  "ApiServerAddress",
			reason: "value must be a valid hostname",
			cause:  err,
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	// no validation rules for CaCertificateBundle

	if len(errors) > 0 {
		return ClusterCreatedV2MultiError(errors)
	}
	return nil
}

func (m *ClusterCreatedV2) _validateHostname(host string) error {
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

// ClusterCreatedV2MultiError is an error wrapping multiple validation errors
// returned by ClusterCreatedV2.ValidateAll() if the designated constraints
// aren't met.
type ClusterCreatedV2MultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m ClusterCreatedV2MultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m ClusterCreatedV2MultiError) AllErrors() []error { return m }

// ClusterCreatedV2ValidationError is the validation error returned by
// ClusterCreatedV2.Validate if the designated constraints aren't met.
type ClusterCreatedV2ValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e ClusterCreatedV2ValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e ClusterCreatedV2ValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e ClusterCreatedV2ValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e ClusterCreatedV2ValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e ClusterCreatedV2ValidationError) ErrorName() string { return "ClusterCreatedV2ValidationError" }

// Error satisfies the builtin error interface
func (e ClusterCreatedV2ValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sClusterCreatedV2.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = ClusterCreatedV2ValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = ClusterCreatedV2ValidationError{}

var _ClusterCreatedV2_Name_Pattern = regexp.MustCompile("^[a-zA-Z][A-Za-z0-9_-]+$")

// Validate checks the field values on ClusterUpdated with the rules defined in
// the proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *ClusterUpdated) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on ClusterUpdated with the rules defined
// in the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in ClusterUpdatedMultiError,
// or nil if none found.
func (m *ClusterUpdated) ValidateAll() error {
	return m.validate(true)
}

func (m *ClusterUpdated) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if len(m.GetDisplayName()) > 150 {
		err := ClusterUpdatedValidationError{
			field:  "DisplayName",
			reason: "value length must be at most 150 bytes",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if err := m._validateHostname(m.GetApiServerAddress()); err != nil {
		err = ClusterUpdatedValidationError{
			field:  "ApiServerAddress",
			reason: "value must be a valid hostname",
			cause:  err,
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	// no validation rules for CaCertificateBundle

	if len(errors) > 0 {
		return ClusterUpdatedMultiError(errors)
	}
	return nil
}

func (m *ClusterUpdated) _validateHostname(host string) error {
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

// ClusterUpdatedMultiError is an error wrapping multiple validation errors
// returned by ClusterUpdated.ValidateAll() if the designated constraints
// aren't met.
type ClusterUpdatedMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m ClusterUpdatedMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m ClusterUpdatedMultiError) AllErrors() []error { return m }

// ClusterUpdatedValidationError is the validation error returned by
// ClusterUpdated.Validate if the designated constraints aren't met.
type ClusterUpdatedValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e ClusterUpdatedValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e ClusterUpdatedValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e ClusterUpdatedValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e ClusterUpdatedValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e ClusterUpdatedValidationError) ErrorName() string { return "ClusterUpdatedValidationError" }

// Error satisfies the builtin error interface
func (e ClusterUpdatedValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sClusterUpdated.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = ClusterUpdatedValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = ClusterUpdatedValidationError{}

// Validate checks the field values on ClusterBootstrapTokenCreated with the
// rules defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *ClusterBootstrapTokenCreated) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on ClusterBootstrapTokenCreated with the
// rules defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// ClusterBootstrapTokenCreatedMultiError, or nil if none found.
func (m *ClusterBootstrapTokenCreated) ValidateAll() error {
	return m.validate(true)
}

func (m *ClusterBootstrapTokenCreated) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Jwt

	if len(errors) > 0 {
		return ClusterBootstrapTokenCreatedMultiError(errors)
	}
	return nil
}

// ClusterBootstrapTokenCreatedMultiError is an error wrapping multiple
// validation errors returned by ClusterBootstrapTokenCreated.ValidateAll() if
// the designated constraints aren't met.
type ClusterBootstrapTokenCreatedMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m ClusterBootstrapTokenCreatedMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m ClusterBootstrapTokenCreatedMultiError) AllErrors() []error { return m }

// ClusterBootstrapTokenCreatedValidationError is the validation error returned
// by ClusterBootstrapTokenCreated.Validate if the designated constraints
// aren't met.
type ClusterBootstrapTokenCreatedValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e ClusterBootstrapTokenCreatedValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e ClusterBootstrapTokenCreatedValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e ClusterBootstrapTokenCreatedValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e ClusterBootstrapTokenCreatedValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e ClusterBootstrapTokenCreatedValidationError) ErrorName() string {
	return "ClusterBootstrapTokenCreatedValidationError"
}

// Error satisfies the builtin error interface
func (e ClusterBootstrapTokenCreatedValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sClusterBootstrapTokenCreated.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = ClusterBootstrapTokenCreatedValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = ClusterBootstrapTokenCreatedValidationError{}
