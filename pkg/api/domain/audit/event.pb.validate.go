// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: api/domain/audit/event.proto

package audit

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

// Validate checks the field values on HumanReadableEvent with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *HumanReadableEvent) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on HumanReadableEvent with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// HumanReadableEventMultiError, or nil if none found.
func (m *HumanReadableEvent) ValidateAll() error {
	return m.validate(true)
}

func (m *HumanReadableEvent) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if all {
		switch v := interface{}(m.GetTimestamp()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, HumanReadableEventValidationError{
					field:  "Timestamp",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, HumanReadableEventValidationError{
					field:  "Timestamp",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetTimestamp()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return HumanReadableEventValidationError{
				field:  "Timestamp",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	// no validation rules for Issuer

	// no validation rules for IssuerId

	// no validation rules for EventType

	// no validation rules for Details

	if len(errors) > 0 {
		return HumanReadableEventMultiError(errors)
	}

	return nil
}

// HumanReadableEventMultiError is an error wrapping multiple validation errors
// returned by HumanReadableEvent.ValidateAll() if the designated constraints
// aren't met.
type HumanReadableEventMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m HumanReadableEventMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m HumanReadableEventMultiError) AllErrors() []error { return m }

// HumanReadableEventValidationError is the validation error returned by
// HumanReadableEvent.Validate if the designated constraints aren't met.
type HumanReadableEventValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e HumanReadableEventValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e HumanReadableEventValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e HumanReadableEventValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e HumanReadableEventValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e HumanReadableEventValidationError) ErrorName() string {
	return "HumanReadableEventValidationError"
}

// Error satisfies the builtin error interface
func (e HumanReadableEventValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sHumanReadableEvent.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = HumanReadableEventValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = HumanReadableEventValidationError{}
