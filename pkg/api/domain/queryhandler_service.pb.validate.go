// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: api/domain/queryhandler_service.proto

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

// Validate checks the field values on GetAllRequest with the rules defined in
// the proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *GetAllRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on GetAllRequest with the rules defined
// in the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in GetAllRequestMultiError, or
// nil if none found.
func (m *GetAllRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *GetAllRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for IncludeDeleted

	if len(errors) > 0 {
		return GetAllRequestMultiError(errors)
	}
	return nil
}

// GetAllRequestMultiError is an error wrapping multiple validation errors
// returned by GetAllRequest.ValidateAll() if the designated constraints
// aren't met.
type GetAllRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m GetAllRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m GetAllRequestMultiError) AllErrors() []error { return m }

// GetAllRequestValidationError is the validation error returned by
// GetAllRequest.Validate if the designated constraints aren't met.
type GetAllRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e GetAllRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e GetAllRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e GetAllRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e GetAllRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e GetAllRequestValidationError) ErrorName() string { return "GetAllRequestValidationError" }

// Error satisfies the builtin error interface
func (e GetAllRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sGetAllRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = GetAllRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = GetAllRequestValidationError{}

// Validate checks the field values on GetCertificateRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *GetCertificateRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on GetCertificateRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// GetCertificateRequestMultiError, or nil if none found.
func (m *GetCertificateRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *GetCertificateRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for AggregateId

	// no validation rules for AggregateType

	if len(errors) > 0 {
		return GetCertificateRequestMultiError(errors)
	}
	return nil
}

// GetCertificateRequestMultiError is an error wrapping multiple validation
// errors returned by GetCertificateRequest.ValidateAll() if the designated
// constraints aren't met.
type GetCertificateRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m GetCertificateRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m GetCertificateRequestMultiError) AllErrors() []error { return m }

// GetCertificateRequestValidationError is the validation error returned by
// GetCertificateRequest.Validate if the designated constraints aren't met.
type GetCertificateRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e GetCertificateRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e GetCertificateRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e GetCertificateRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e GetCertificateRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e GetCertificateRequestValidationError) ErrorName() string {
	return "GetCertificateRequestValidationError"
}

// Error satisfies the builtin error interface
func (e GetCertificateRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sGetCertificateRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = GetCertificateRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = GetCertificateRequestValidationError{}

// Validate checks the field values on GetClusterMappingRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *GetClusterMappingRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on GetClusterMappingRequest with the
// rules defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// GetClusterMappingRequestMultiError, or nil if none found.
func (m *GetClusterMappingRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *GetClusterMappingRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for TenantId

	// no validation rules for ClusterId

	if len(errors) > 0 {
		return GetClusterMappingRequestMultiError(errors)
	}
	return nil
}

// GetClusterMappingRequestMultiError is an error wrapping multiple validation
// errors returned by GetClusterMappingRequest.ValidateAll() if the designated
// constraints aren't met.
type GetClusterMappingRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m GetClusterMappingRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m GetClusterMappingRequestMultiError) AllErrors() []error { return m }

// GetClusterMappingRequestValidationError is the validation error returned by
// GetClusterMappingRequest.Validate if the designated constraints aren't met.
type GetClusterMappingRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e GetClusterMappingRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e GetClusterMappingRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e GetClusterMappingRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e GetClusterMappingRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e GetClusterMappingRequestValidationError) ErrorName() string {
	return "GetClusterMappingRequestValidationError"
}

// Error satisfies the builtin error interface
func (e GetClusterMappingRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sGetClusterMappingRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = GetClusterMappingRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = GetClusterMappingRequestValidationError{}

// Validate checks the field values on GetAuditLogByDateRangeRequest with the
// rules defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *GetAuditLogByDateRangeRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on GetAuditLogByDateRangeRequest with
// the rules defined in the proto definition for this message. If any rules
// are violated, the result is a list of violation errors wrapped in
// GetAuditLogByDateRangeRequestMultiError, or nil if none found.
func (m *GetAuditLogByDateRangeRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *GetAuditLogByDateRangeRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if all {
		switch v := interface{}(m.GetMinTimestamp()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, GetAuditLogByDateRangeRequestValidationError{
					field:  "MinTimestamp",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, GetAuditLogByDateRangeRequestValidationError{
					field:  "MinTimestamp",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetMinTimestamp()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return GetAuditLogByDateRangeRequestValidationError{
				field:  "MinTimestamp",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if all {
		switch v := interface{}(m.GetMaxTimestamp()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, GetAuditLogByDateRangeRequestValidationError{
					field:  "MaxTimestamp",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, GetAuditLogByDateRangeRequestValidationError{
					field:  "MaxTimestamp",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetMaxTimestamp()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return GetAuditLogByDateRangeRequestValidationError{
				field:  "MaxTimestamp",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if len(errors) > 0 {
		return GetAuditLogByDateRangeRequestMultiError(errors)
	}
	return nil
}

// GetAuditLogByDateRangeRequestMultiError is an error wrapping multiple
// validation errors returned by GetAuditLogByDateRangeRequest.ValidateAll()
// if the designated constraints aren't met.
type GetAuditLogByDateRangeRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m GetAuditLogByDateRangeRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m GetAuditLogByDateRangeRequestMultiError) AllErrors() []error { return m }

// GetAuditLogByDateRangeRequestValidationError is the validation error
// returned by GetAuditLogByDateRangeRequest.Validate if the designated
// constraints aren't met.
type GetAuditLogByDateRangeRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e GetAuditLogByDateRangeRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e GetAuditLogByDateRangeRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e GetAuditLogByDateRangeRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e GetAuditLogByDateRangeRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e GetAuditLogByDateRangeRequestValidationError) ErrorName() string {
	return "GetAuditLogByDateRangeRequestValidationError"
}

// Error satisfies the builtin error interface
func (e GetAuditLogByDateRangeRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sGetAuditLogByDateRangeRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = GetAuditLogByDateRangeRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = GetAuditLogByDateRangeRequestValidationError{}
