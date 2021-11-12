// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: api/eventsourcing/commands/command.proto

package commands

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

// Validate checks the field values on Command with the rules defined in the
// proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *Command) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on Command with the rules defined in the
// proto definition for this message. If any rules are violated, the result is
// a list of violation errors wrapped in CommandMultiError, or nil if none found.
func (m *Command) ValidateAll() error {
	return m.validate(true)
}

func (m *Command) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Id

	// no validation rules for Type

	if all {
		switch v := interface{}(m.GetData()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, CommandValidationError{
					field:  "Data",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, CommandValidationError{
					field:  "Data",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetData()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return CommandValidationError{
				field:  "Data",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if len(errors) > 0 {
		return CommandMultiError(errors)
	}
	return nil
}

// CommandMultiError is an error wrapping multiple validation errors returned
// by Command.ValidateAll() if the designated constraints aren't met.
type CommandMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m CommandMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m CommandMultiError) AllErrors() []error { return m }

// CommandValidationError is the validation error returned by Command.Validate
// if the designated constraints aren't met.
type CommandValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e CommandValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e CommandValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e CommandValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e CommandValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e CommandValidationError) ErrorName() string { return "CommandValidationError" }

// Error satisfies the builtin error interface
func (e CommandValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sCommand.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = CommandValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = CommandValidationError{}

// Validate checks the field values on TestCommandData with the rules defined
// in the proto definition for this message. If any rules are violated, the
// first error encountered is returned, or nil if there are no violations.
func (m *TestCommandData) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on TestCommandData with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// TestCommandDataMultiError, or nil if none found.
func (m *TestCommandData) ValidateAll() error {
	return m.validate(true)
}

func (m *TestCommandData) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Test

	// no validation rules for TestCount

	if len(errors) > 0 {
		return TestCommandDataMultiError(errors)
	}
	return nil
}

// TestCommandDataMultiError is an error wrapping multiple validation errors
// returned by TestCommandData.ValidateAll() if the designated constraints
// aren't met.
type TestCommandDataMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m TestCommandDataMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m TestCommandDataMultiError) AllErrors() []error { return m }

// TestCommandDataValidationError is the validation error returned by
// TestCommandData.Validate if the designated constraints aren't met.
type TestCommandDataValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e TestCommandDataValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e TestCommandDataValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e TestCommandDataValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e TestCommandDataValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e TestCommandDataValidationError) ErrorName() string { return "TestCommandDataValidationError" }

// Error satisfies the builtin error interface
func (e TestCommandDataValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sTestCommandData.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = TestCommandDataValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = TestCommandDataValidationError{}
