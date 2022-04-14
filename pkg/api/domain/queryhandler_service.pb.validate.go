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

// Validate checks the field values on GetCountRequest with the rules defined
// in the proto definition for this message. If any rules are violated, the
// first error encountered is returned, or nil if there are no violations.
func (m *GetCountRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on GetCountRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// GetCountRequestMultiError, or nil if none found.
func (m *GetCountRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *GetCountRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for IncludeDeleted

	if len(errors) > 0 {
		return GetCountRequestMultiError(errors)
	}
	return nil
}

// GetCountRequestMultiError is an error wrapping multiple validation errors
// returned by GetCountRequest.ValidateAll() if the designated constraints
// aren't met.
type GetCountRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m GetCountRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m GetCountRequestMultiError) AllErrors() []error { return m }

// GetCountRequestValidationError is the validation error returned by
// GetCountRequest.Validate if the designated constraints aren't met.
type GetCountRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e GetCountRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e GetCountRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e GetCountRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e GetCountRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e GetCountRequestValidationError) ErrorName() string { return "GetCountRequestValidationError" }

// Error satisfies the builtin error interface
func (e GetCountRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sGetCountRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = GetCountRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = GetCountRequestValidationError{}

// Validate checks the field values on GetCountResult with the rules defined in
// the proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *GetCountResult) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on GetCountResult with the rules defined
// in the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in GetCountResultMultiError,
// or nil if none found.
func (m *GetCountResult) ValidateAll() error {
	return m.validate(true)
}

func (m *GetCountResult) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Count

	if len(errors) > 0 {
		return GetCountResultMultiError(errors)
	}
	return nil
}

// GetCountResultMultiError is an error wrapping multiple validation errors
// returned by GetCountResult.ValidateAll() if the designated constraints
// aren't met.
type GetCountResultMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m GetCountResultMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m GetCountResultMultiError) AllErrors() []error { return m }

// GetCountResultValidationError is the validation error returned by
// GetCountResult.Validate if the designated constraints aren't met.
type GetCountResultValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e GetCountResultValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e GetCountResultValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e GetCountResultValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e GetCountResultValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e GetCountResultValidationError) ErrorName() string { return "GetCountResultValidationError" }

// Error satisfies the builtin error interface
func (e GetCountResultValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sGetCountResult.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = GetCountResultValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = GetCountResultValidationError{}

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

// Validate checks the field values on GetByUserRequest with the rules defined
// in the proto definition for this message. If any rules are violated, the
// first error encountered is returned, or nil if there are no violations.
func (m *GetByUserRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on GetByUserRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// GetByUserRequestMultiError, or nil if none found.
func (m *GetByUserRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *GetByUserRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if all {
		switch v := interface{}(m.GetEmail()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, GetByUserRequestValidationError{
					field:  "Email",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, GetByUserRequestValidationError{
					field:  "Email",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetEmail()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return GetByUserRequestValidationError{
				field:  "Email",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if all {
		switch v := interface{}(m.GetDateRange()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, GetByUserRequestValidationError{
					field:  "DateRange",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, GetByUserRequestValidationError{
					field:  "DateRange",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetDateRange()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return GetByUserRequestValidationError{
				field:  "DateRange",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if len(errors) > 0 {
		return GetByUserRequestMultiError(errors)
	}
	return nil
}

// GetByUserRequestMultiError is an error wrapping multiple validation errors
// returned by GetByUserRequest.ValidateAll() if the designated constraints
// aren't met.
type GetByUserRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m GetByUserRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m GetByUserRequestMultiError) AllErrors() []error { return m }

// GetByUserRequestValidationError is the validation error returned by
// GetByUserRequest.Validate if the designated constraints aren't met.
type GetByUserRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e GetByUserRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e GetByUserRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e GetByUserRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e GetByUserRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e GetByUserRequestValidationError) ErrorName() string { return "GetByUserRequestValidationError" }

// Error satisfies the builtin error interface
func (e GetByUserRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sGetByUserRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = GetByUserRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = GetByUserRequestValidationError{}

// Validate checks the field values on GetUserActionsRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *GetUserActionsRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on GetUserActionsRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// GetUserActionsRequestMultiError, or nil if none found.
func (m *GetUserActionsRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *GetUserActionsRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if all {
		switch v := interface{}(m.GetEmail()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, GetUserActionsRequestValidationError{
					field:  "Email",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, GetUserActionsRequestValidationError{
					field:  "Email",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetEmail()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return GetUserActionsRequestValidationError{
				field:  "Email",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if all {
		switch v := interface{}(m.GetDateRange()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, GetUserActionsRequestValidationError{
					field:  "DateRange",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, GetUserActionsRequestValidationError{
					field:  "DateRange",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetDateRange()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return GetUserActionsRequestValidationError{
				field:  "DateRange",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if len(errors) > 0 {
		return GetUserActionsRequestMultiError(errors)
	}
	return nil
}

// GetUserActionsRequestMultiError is an error wrapping multiple validation
// errors returned by GetUserActionsRequest.ValidateAll() if the designated
// constraints aren't met.
type GetUserActionsRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m GetUserActionsRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m GetUserActionsRequestMultiError) AllErrors() []error { return m }

// GetUserActionsRequestValidationError is the validation error returned by
// GetUserActionsRequest.Validate if the designated constraints aren't met.
type GetUserActionsRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e GetUserActionsRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e GetUserActionsRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e GetUserActionsRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e GetUserActionsRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e GetUserActionsRequestValidationError) ErrorName() string {
	return "GetUserActionsRequestValidationError"
}

// Error satisfies the builtin error interface
func (e GetUserActionsRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sGetUserActionsRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = GetUserActionsRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = GetUserActionsRequestValidationError{}

// Validate checks the field values on GetUsersOverviewRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *GetUsersOverviewRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on GetUsersOverviewRequest with the
// rules defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// GetUsersOverviewRequestMultiError, or nil if none found.
func (m *GetUsersOverviewRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *GetUsersOverviewRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if all {
		switch v := interface{}(m.GetTimestamp()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, GetUsersOverviewRequestValidationError{
					field:  "Timestamp",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, GetUsersOverviewRequestValidationError{
					field:  "Timestamp",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetTimestamp()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return GetUsersOverviewRequestValidationError{
				field:  "Timestamp",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if len(errors) > 0 {
		return GetUsersOverviewRequestMultiError(errors)
	}
	return nil
}

// GetUsersOverviewRequestMultiError is an error wrapping multiple validation
// errors returned by GetUsersOverviewRequest.ValidateAll() if the designated
// constraints aren't met.
type GetUsersOverviewRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m GetUsersOverviewRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m GetUsersOverviewRequestMultiError) AllErrors() []error { return m }

// GetUsersOverviewRequestValidationError is the validation error returned by
// GetUsersOverviewRequest.Validate if the designated constraints aren't met.
type GetUsersOverviewRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e GetUsersOverviewRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e GetUsersOverviewRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e GetUsersOverviewRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e GetUsersOverviewRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e GetUsersOverviewRequestValidationError) ErrorName() string {
	return "GetUsersOverviewRequestValidationError"
}

// Error satisfies the builtin error interface
func (e GetUsersOverviewRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sGetUsersOverviewRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = GetUsersOverviewRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = GetUsersOverviewRequestValidationError{}
