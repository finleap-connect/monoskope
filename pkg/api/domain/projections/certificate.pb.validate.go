// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: api/domain/projections/certificate.proto

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
var _certificate_uuidPattern = regexp.MustCompile("^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$")

// Validate checks the field values on Certificate with the rules defined in
// the proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *Certificate) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on Certificate with the rules defined in
// the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in CertificateMultiError, or
// nil if none found.
func (m *Certificate) ValidateAll() error {
	return m.validate(true)
}

func (m *Certificate) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if err := m._validateUuid(m.GetId()); err != nil {
		err = CertificateValidationError{
			field:  "Id",
			reason: "value must be a valid UUID",
			cause:  err,
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if err := m._validateUuid(m.GetReferencedAggregateId()); err != nil {
		err = CertificateValidationError{
			field:  "ReferencedAggregateId",
			reason: "value must be a valid UUID",
			cause:  err,
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if len(m.GetAggregateType()) > 60 {
		err := CertificateValidationError{
			field:  "AggregateType",
			reason: "value length must be at most 60 bytes",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if !_Certificate_AggregateType_Pattern.MatchString(m.GetAggregateType()) {
		err := CertificateValidationError{
			field:  "AggregateType",
			reason: "value does not match regex pattern \"^[a-zA-Z][A-Za-z0-9_-]+$\"",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	// no validation rules for Certificate

	// no validation rules for CaCertBundle

	if all {
		switch v := interface{}(m.GetMetadata()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, CertificateValidationError{
					field:  "Metadata",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, CertificateValidationError{
					field:  "Metadata",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetMetadata()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return CertificateValidationError{
				field:  "Metadata",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if len(errors) > 0 {
		return CertificateMultiError(errors)
	}
	return nil
}

func (m *Certificate) _validateUuid(uuid string) error {
	if matched := _certificate_uuidPattern.MatchString(uuid); !matched {
		return errors.New("invalid uuid format")
	}

	return nil
}

// CertificateMultiError is an error wrapping multiple validation errors
// returned by Certificate.ValidateAll() if the designated constraints aren't met.
type CertificateMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m CertificateMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m CertificateMultiError) AllErrors() []error { return m }

// CertificateValidationError is the validation error returned by
// Certificate.Validate if the designated constraints aren't met.
type CertificateValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e CertificateValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e CertificateValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e CertificateValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e CertificateValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e CertificateValidationError) ErrorName() string { return "CertificateValidationError" }

// Error satisfies the builtin error interface
func (e CertificateValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sCertificate.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = CertificateValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = CertificateValidationError{}

var _Certificate_AggregateType_Pattern = regexp.MustCompile("^[a-zA-Z][A-Za-z0-9_-]+$")
