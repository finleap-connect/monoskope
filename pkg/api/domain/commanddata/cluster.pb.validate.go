// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: api/domain/commanddata/cluster.proto

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

// Validate checks the field values on CreateCluster with the rules defined in
// the proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *CreateCluster) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on CreateCluster with the rules defined
// in the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in CreateClusterMultiError, or
// nil if none found.
func (m *CreateCluster) ValidateAll() error {
	return m.validate(true)
}

func (m *CreateCluster) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if utf8.RuneCountInString(m.GetName()) > 60 {
		err := CreateClusterValidationError{
			field:  "Name",
			reason: "value length must be at most 60 runes",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if !_CreateCluster_Name_Pattern.MatchString(m.GetName()) {
		err := CreateClusterValidationError{
			field:  "Name",
			reason: "value does not match regex pattern \"^[a-zA-Z][A-Za-z0-9_-]+$\"",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if utf8.RuneCountInString(m.GetDisplayName()) > 150 {
		err := CreateClusterValidationError{
			field:  "DisplayName",
			reason: "value length must be at most 150 runes",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if !_CreateCluster_ApiServerAddress_Pattern.MatchString(m.GetApiServerAddress()) {
		err := CreateClusterValidationError{
			field:  "ApiServerAddress",
			reason: "value does not match regex pattern \"^(https?://)?[^\\\\s/$.?#/_].[^\\\\s_]*$\"",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if !bytes.HasPrefix(m.GetCaCertBundle(), []uint8{0x2D, 0x2D, 0x2D, 0x2D, 0x2D, 0x42, 0x45, 0x47, 0x49, 0x4E, 0x20, 0x43, 0x45, 0x52, 0x54, 0x49, 0x46, 0x49, 0x43, 0x41, 0x54, 0x45, 0x2D, 0x2D, 0x2D, 0x2D, 0x2D}) {
		err := CreateClusterValidationError{
			field:  "CaCertBundle",
			reason: "value does not have prefix \"\\x2D\\x2D\\x2D\\x2D\\x2D\\x42\\x45\\x47\\x49\\x4E\\x20\\x43\\x45\\x52\\x54\\x49\\x46\\x49\\x43\\x41\\x54\\x45\\x2D\\x2D\\x2D\\x2D\\x2D\"",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if !bytes.HasSuffix(m.GetCaCertBundle(), []uint8{0x2D, 0x2D, 0x2D, 0x2D, 0x2D, 0x45, 0x4E, 0x44, 0x20, 0x43, 0x45, 0x52, 0x54, 0x49, 0x46, 0x49, 0x43, 0x41, 0x54, 0x45, 0x2D, 0x2D, 0x2D, 0x2D, 0x2D}) {
		err := CreateClusterValidationError{
			field:  "CaCertBundle",
			reason: "value does not have suffix \"\\x2D\\x2D\\x2D\\x2D\\x2D\\x45\\x4E\\x44\\x20\\x43\\x45\\x52\\x54\\x49\\x46\\x49\\x43\\x41\\x54\\x45\\x2D\\x2D\\x2D\\x2D\\x2D\"",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		return CreateClusterMultiError(errors)
	}

	return nil
}

// CreateClusterMultiError is an error wrapping multiple validation errors
// returned by CreateCluster.ValidateAll() if the designated constraints
// aren't met.
type CreateClusterMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m CreateClusterMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m CreateClusterMultiError) AllErrors() []error { return m }

// CreateClusterValidationError is the validation error returned by
// CreateCluster.Validate if the designated constraints aren't met.
type CreateClusterValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e CreateClusterValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e CreateClusterValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e CreateClusterValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e CreateClusterValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e CreateClusterValidationError) ErrorName() string { return "CreateClusterValidationError" }

// Error satisfies the builtin error interface
func (e CreateClusterValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sCreateCluster.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = CreateClusterValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = CreateClusterValidationError{}

var _CreateCluster_Name_Pattern = regexp.MustCompile("^[a-zA-Z][A-Za-z0-9_-]+$")

var _CreateCluster_ApiServerAddress_Pattern = regexp.MustCompile("^(https?://)?[^\\s/$.?#/_].[^\\s_]*$")

// Validate checks the field values on UpdateCluster with the rules defined in
// the proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *UpdateCluster) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on UpdateCluster with the rules defined
// in the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in UpdateClusterMultiError, or
// nil if none found.
func (m *UpdateCluster) ValidateAll() error {
	return m.validate(true)
}

func (m *UpdateCluster) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if wrapper := m.GetDisplayName(); wrapper != nil {

		if utf8.RuneCountInString(wrapper.GetValue()) > 150 {
			err := UpdateClusterValidationError{
				field:  "DisplayName",
				reason: "value length must be at most 150 runes",
			}
			if !all {
				return err
			}
			errors = append(errors, err)
		}

	}

	if wrapper := m.GetApiServerAddress(); wrapper != nil {

		if !_UpdateCluster_ApiServerAddress_Pattern.MatchString(wrapper.GetValue()) {
			err := UpdateClusterValidationError{
				field:  "ApiServerAddress",
				reason: "value does not match regex pattern \"^(https?://)?[^\\\\s/$.?#/_].[^\\\\s_]*$\"",
			}
			if !all {
				return err
			}
			errors = append(errors, err)
		}

	}

	// no validation rules for CaCertBundle

	if len(errors) > 0 {
		return UpdateClusterMultiError(errors)
	}

	return nil
}

// UpdateClusterMultiError is an error wrapping multiple validation errors
// returned by UpdateCluster.ValidateAll() if the designated constraints
// aren't met.
type UpdateClusterMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m UpdateClusterMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m UpdateClusterMultiError) AllErrors() []error { return m }

// UpdateClusterValidationError is the validation error returned by
// UpdateCluster.Validate if the designated constraints aren't met.
type UpdateClusterValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e UpdateClusterValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e UpdateClusterValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e UpdateClusterValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e UpdateClusterValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e UpdateClusterValidationError) ErrorName() string { return "UpdateClusterValidationError" }

// Error satisfies the builtin error interface
func (e UpdateClusterValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sUpdateCluster.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = UpdateClusterValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = UpdateClusterValidationError{}

var _UpdateCluster_ApiServerAddress_Pattern = regexp.MustCompile("^(https?://)?[^\\s/$.?#/_].[^\\s_]*$")
