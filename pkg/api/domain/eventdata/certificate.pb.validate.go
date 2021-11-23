// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: api/domain/eventdata/certificate.proto

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

// define the regex for a UUID once up-front
var _certificate_uuidPattern = regexp.MustCompile("^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$")

// Validate checks the field values on CertificateRequested with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *CertificateRequested) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on CertificateRequested with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// CertificateRequestedMultiError, or nil if none found.
func (m *CertificateRequested) ValidateAll() error {
	return m.validate(true)
}

func (m *CertificateRequested) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if err := m._validateUuid(m.GetReferencedAggregateId()); err != nil {
		err = CertificateRequestedValidationError{
			field:  "ReferencedAggregateId",
			reason: "value must be a valid UUID",
			cause:  err,
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if len(m.GetReferencedAggregateType()) > 60 {
		err := CertificateRequestedValidationError{
			field:  "ReferencedAggregateType",
			reason: "value length must be at most 60 bytes",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if !_CertificateRequested_ReferencedAggregateType_Pattern.MatchString(m.GetReferencedAggregateType()) {
		err := CertificateRequestedValidationError{
			field:  "ReferencedAggregateType",
			reason: "value does not match regex pattern \"^[a-zA-Z_]+$\"",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if !bytes.HasPrefix(m.GetSigningRequest(), []uint8{0x2D, 0x2D, 0x2D, 0x2D, 0x2D, 0x42, 0x45, 0x47, 0x49, 0x4E, 0x20, 0x43, 0x45, 0x52, 0x54, 0x49, 0x46, 0x49, 0x43, 0x41, 0x54, 0x45, 0x20, 0x52, 0x45, 0x51, 0x55, 0x45, 0x53, 0x54, 0x2D, 0x2D, 0x2D, 0x2D, 0x2D}) {
		err := CertificateRequestedValidationError{
			field:  "SigningRequest",
			reason: "value does not have prefix \"\\x2D\\x2D\\x2D\\x2D\\x2D\\x42\\x45\\x47\\x49\\x4E\\x20\\x43\\x45\\x52\\x54\\x49\\x46\\x49\\x43\\x41\\x54\\x45\\x20\\x52\\x45\\x51\\x55\\x45\\x53\\x54\\x2D\\x2D\\x2D\\x2D\\x2D\"",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if !bytes.HasSuffix(m.GetSigningRequest(), []uint8{0x2D, 0x2D, 0x2D, 0x2D, 0x2D, 0x45, 0x4E, 0x44, 0x20, 0x43, 0x45, 0x52, 0x54, 0x49, 0x46, 0x49, 0x43, 0x41, 0x54, 0x45, 0x20, 0x52, 0x45, 0x51, 0x55, 0x45, 0x53, 0x54, 0x2D, 0x2D, 0x2D, 0x2D, 0x2D}) {
		err := CertificateRequestedValidationError{
			field:  "SigningRequest",
			reason: "value does not have suffix \"\\x2D\\x2D\\x2D\\x2D\\x2D\\x45\\x4E\\x44\\x20\\x43\\x45\\x52\\x54\\x49\\x46\\x49\\x43\\x41\\x54\\x45\\x20\\x52\\x45\\x51\\x55\\x45\\x53\\x54\\x2D\\x2D\\x2D\\x2D\\x2D\"",
		}
		if !all {
			return err
		}
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		return CertificateRequestedMultiError(errors)
	}
	return nil
}

func (m *CertificateRequested) _validateUuid(uuid string) error {
	if matched := _certificate_uuidPattern.MatchString(uuid); !matched {
		return errors.New("invalid uuid format")
	}

	return nil
}

// CertificateRequestedMultiError is an error wrapping multiple validation
// errors returned by CertificateRequested.ValidateAll() if the designated
// constraints aren't met.
type CertificateRequestedMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m CertificateRequestedMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m CertificateRequestedMultiError) AllErrors() []error { return m }

// CertificateRequestedValidationError is the validation error returned by
// CertificateRequested.Validate if the designated constraints aren't met.
type CertificateRequestedValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e CertificateRequestedValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e CertificateRequestedValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e CertificateRequestedValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e CertificateRequestedValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e CertificateRequestedValidationError) ErrorName() string {
	return "CertificateRequestedValidationError"
}

// Error satisfies the builtin error interface
func (e CertificateRequestedValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sCertificateRequested.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = CertificateRequestedValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = CertificateRequestedValidationError{}

var _CertificateRequested_ReferencedAggregateType_Pattern = regexp.MustCompile("^[a-zA-Z_]+$")

// Validate checks the field values on CertificateIssued with the rules defined
// in the proto definition for this message. If any rules are violated, the
// first error encountered is returned, or nil if there are no violations.
func (m *CertificateIssued) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on CertificateIssued with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// CertificateIssuedMultiError, or nil if none found.
func (m *CertificateIssued) ValidateAll() error {
	return m.validate(true)
}

func (m *CertificateIssued) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if all {
		switch v := interface{}(m.GetCertificate()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, CertificateIssuedValidationError{
					field:  "Certificate",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, CertificateIssuedValidationError{
					field:  "Certificate",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetCertificate()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return CertificateIssuedValidationError{
				field:  "Certificate",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	if len(errors) > 0 {
		return CertificateIssuedMultiError(errors)
	}
	return nil
}

// CertificateIssuedMultiError is an error wrapping multiple validation errors
// returned by CertificateIssued.ValidateAll() if the designated constraints
// aren't met.
type CertificateIssuedMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m CertificateIssuedMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m CertificateIssuedMultiError) AllErrors() []error { return m }

// CertificateIssuedValidationError is the validation error returned by
// CertificateIssued.Validate if the designated constraints aren't met.
type CertificateIssuedValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e CertificateIssuedValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e CertificateIssuedValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e CertificateIssuedValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e CertificateIssuedValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e CertificateIssuedValidationError) ErrorName() string {
	return "CertificateIssuedValidationError"
}

// Error satisfies the builtin error interface
func (e CertificateIssuedValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sCertificateIssued.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = CertificateIssuedValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = CertificateIssuedValidationError{}
