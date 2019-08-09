// Code generated by protoc-gen-fieldmask. DO NOT EDIT.

package ttnpb

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"net/mail"
	"net/url"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/gogo/protobuf/types"
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
	_ = types.DynamicAny{}
)

// ValidateFields checks the field values on ClaimEndDeviceRequest with the
// rules defined in the proto definition for this message. If any rules are
// violated, an error is returned.
func (m *ClaimEndDeviceRequest) ValidateFields(paths ...string) error {
	if m == nil {
		return nil
	}

	if len(paths) == 0 {
		paths = ClaimEndDeviceRequestFieldPathsNested
	}

	for name, subs := range _processPaths(append(paths[:0:0], paths...)) {
		_ = subs
		switch name {
		case "target_application_ids":

			if v, ok := interface{}(&m.TargetApplicationIDs).(interface{ ValidateFields(...string) error }); ok {
				if err := v.ValidateFields(subs...); err != nil {
					return ClaimEndDeviceRequestValidationError{
						field:  "target_application_ids",
						reason: "embedded message failed validation",
						cause:  err,
					}
				}
			}

		case "target_device_id":

			if utf8.RuneCountInString(m.GetTargetDeviceID()) > 36 {
				return ClaimEndDeviceRequestValidationError{
					field:  "target_device_id",
					reason: "value length must be at most 36 runes",
				}
			}

			if !_ClaimEndDeviceRequest_TargetDeviceID_Pattern.MatchString(m.GetTargetDeviceID()) {
				return ClaimEndDeviceRequestValidationError{
					field:  "target_device_id",
					reason: "value does not match regex pattern \"^[a-z0-9](?:[-]?[a-z0-9]){2,}$|^$\"",
				}
			}

		case "invalidate_authentication_code":
			// no validation rules for InvalidateAuthenticationCode
		case "source_device":
			if len(subs) == 0 {
				subs = []string{
					"authenticated_identifiers", "qr_code",
				}
			}
			for name, subs := range _processPaths(subs) {
				_ = subs
				switch name {
				case "authenticated_identifiers":

					if v, ok := interface{}(m.GetAuthenticatedIdentifiers()).(interface{ ValidateFields(...string) error }); ok {
						if err := v.ValidateFields(subs...); err != nil {
							return ClaimEndDeviceRequestValidationError{
								field:  "authenticated_identifiers",
								reason: "embedded message failed validation",
								cause:  err,
							}
						}
					}

				case "qr_code":

					if l := len(m.GetQRCode()); l < 1 || l > 1024 {
						return ClaimEndDeviceRequestValidationError{
							field:  "qr_code",
							reason: "value length must be between 1 and 1024 bytes, inclusive",
						}
					}

				}
			}
		default:
			return ClaimEndDeviceRequestValidationError{
				field:  name,
				reason: "invalid field path",
			}
		}
	}
	return nil
}

// ClaimEndDeviceRequestValidationError is the validation error returned by
// ClaimEndDeviceRequest.ValidateFields if the designated constraints aren't met.
type ClaimEndDeviceRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e ClaimEndDeviceRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e ClaimEndDeviceRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e ClaimEndDeviceRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e ClaimEndDeviceRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e ClaimEndDeviceRequestValidationError) ErrorName() string {
	return "ClaimEndDeviceRequestValidationError"
}

// Error satisfies the builtin error interface
func (e ClaimEndDeviceRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sClaimEndDeviceRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = ClaimEndDeviceRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = ClaimEndDeviceRequestValidationError{}

var _ClaimEndDeviceRequest_TargetDeviceID_Pattern = regexp.MustCompile("^[a-z0-9](?:[-]?[a-z0-9]){2,}$|^$")

// ValidateFields checks the field values on
// ClaimEndDeviceRequest_AuthenticatedIdentifiers with the rules defined in
// the proto definition for this message. If any rules are violated, an error
// is returned.
func (m *ClaimEndDeviceRequest_AuthenticatedIdentifiers) ValidateFields(paths ...string) error {
	if m == nil {
		return nil
	}

	if len(paths) == 0 {
		paths = ClaimEndDeviceRequest_AuthenticatedIdentifiersFieldPathsNested
	}

	for name, subs := range _processPaths(append(paths[:0:0], paths...)) {
		_ = subs
		switch name {
		case "join_eui":
			// no validation rules for JoinEUI
		case "dev_eui":
			// no validation rules for DevEUI
		case "authentication_code":

			if l := len(m.GetAuthenticationCode()); l < 1 || l > 8 {
				return ClaimEndDeviceRequest_AuthenticatedIdentifiersValidationError{
					field:  "authentication_code",
					reason: "value length must be between 1 and 8 bytes, inclusive",
				}
			}

		default:
			return ClaimEndDeviceRequest_AuthenticatedIdentifiersValidationError{
				field:  name,
				reason: "invalid field path",
			}
		}
	}
	return nil
}

// ClaimEndDeviceRequest_AuthenticatedIdentifiersValidationError is the
// validation error returned by
// ClaimEndDeviceRequest_AuthenticatedIdentifiers.ValidateFields if the
// designated constraints aren't met.
type ClaimEndDeviceRequest_AuthenticatedIdentifiersValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e ClaimEndDeviceRequest_AuthenticatedIdentifiersValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e ClaimEndDeviceRequest_AuthenticatedIdentifiersValidationError) Reason() string {
	return e.reason
}

// Cause function returns cause value.
func (e ClaimEndDeviceRequest_AuthenticatedIdentifiersValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e ClaimEndDeviceRequest_AuthenticatedIdentifiersValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e ClaimEndDeviceRequest_AuthenticatedIdentifiersValidationError) ErrorName() string {
	return "ClaimEndDeviceRequest_AuthenticatedIdentifiersValidationError"
}

// Error satisfies the builtin error interface
func (e ClaimEndDeviceRequest_AuthenticatedIdentifiersValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sClaimEndDeviceRequest_AuthenticatedIdentifiers.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = ClaimEndDeviceRequest_AuthenticatedIdentifiersValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = ClaimEndDeviceRequest_AuthenticatedIdentifiersValidationError{}
