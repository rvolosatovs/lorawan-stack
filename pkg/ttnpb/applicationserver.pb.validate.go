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

// define the regex for a UUID once up-front
var _applicationserver_uuidPattern = regexp.MustCompile("^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$")

// ValidateFields checks the field values on ApplicationLink with the rules
// defined in the proto definition for this message. If any rules are
// violated, an error is returned.
func (m *ApplicationLink) ValidateFields(paths ...string) error {
	if m == nil {
		return nil
	}

	if len(paths) == 0 {
		paths = ApplicationLinkFieldPathsNested
	}

	for name, subs := range _processPaths(append(paths[:0:0], paths...)) {
		_ = subs
		switch name {
		case "default_formatters":

			if v, ok := interface{}(m.GetDefaultFormatters()).(interface{ ValidateFields(...string) error }); ok {
				if err := v.ValidateFields(subs...); err != nil {
					return ApplicationLinkValidationError{
						field:  "default_formatters",
						reason: "embedded message failed validation",
						cause:  err,
					}
				}
			}

		case "tls":
			// no validation rules for TLS
		case "skip_payload_crypto":

			if v, ok := interface{}(m.GetSkipPayloadCrypto()).(interface{ ValidateFields(...string) error }); ok {
				if err := v.ValidateFields(subs...); err != nil {
					return ApplicationLinkValidationError{
						field:  "skip_payload_crypto",
						reason: "embedded message failed validation",
						cause:  err,
					}
				}
			}

		default:
			return ApplicationLinkValidationError{
				field:  name,
				reason: "invalid field path",
			}
		}
	}
	return nil
}

// ApplicationLinkValidationError is the validation error returned by
// ApplicationLink.ValidateFields if the designated constraints aren't met.
type ApplicationLinkValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e ApplicationLinkValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e ApplicationLinkValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e ApplicationLinkValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e ApplicationLinkValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e ApplicationLinkValidationError) ErrorName() string { return "ApplicationLinkValidationError" }

// Error satisfies the builtin error interface
func (e ApplicationLinkValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sApplicationLink.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = ApplicationLinkValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = ApplicationLinkValidationError{}

// ValidateFields checks the field values on GetApplicationLinkRequest with the
// rules defined in the proto definition for this message. If any rules are
// violated, an error is returned.
func (m *GetApplicationLinkRequest) ValidateFields(paths ...string) error {
	if m == nil {
		return nil
	}

	if len(paths) == 0 {
		paths = GetApplicationLinkRequestFieldPathsNested
	}

	for name, subs := range _processPaths(append(paths[:0:0], paths...)) {
		_ = subs
		switch name {
		case "application_ids":

			if v, ok := interface{}(&m.ApplicationIdentifiers).(interface{ ValidateFields(...string) error }); ok {
				if err := v.ValidateFields(subs...); err != nil {
					return GetApplicationLinkRequestValidationError{
						field:  "application_ids",
						reason: "embedded message failed validation",
						cause:  err,
					}
				}
			}

		case "field_mask":

			if v, ok := interface{}(&m.FieldMask).(interface{ ValidateFields(...string) error }); ok {
				if err := v.ValidateFields(subs...); err != nil {
					return GetApplicationLinkRequestValidationError{
						field:  "field_mask",
						reason: "embedded message failed validation",
						cause:  err,
					}
				}
			}

		default:
			return GetApplicationLinkRequestValidationError{
				field:  name,
				reason: "invalid field path",
			}
		}
	}
	return nil
}

// GetApplicationLinkRequestValidationError is the validation error returned by
// GetApplicationLinkRequest.ValidateFields if the designated constraints
// aren't met.
type GetApplicationLinkRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e GetApplicationLinkRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e GetApplicationLinkRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e GetApplicationLinkRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e GetApplicationLinkRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e GetApplicationLinkRequestValidationError) ErrorName() string {
	return "GetApplicationLinkRequestValidationError"
}

// Error satisfies the builtin error interface
func (e GetApplicationLinkRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sGetApplicationLinkRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = GetApplicationLinkRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = GetApplicationLinkRequestValidationError{}

// ValidateFields checks the field values on SetApplicationLinkRequest with the
// rules defined in the proto definition for this message. If any rules are
// violated, an error is returned.
func (m *SetApplicationLinkRequest) ValidateFields(paths ...string) error {
	if m == nil {
		return nil
	}

	if len(paths) == 0 {
		paths = SetApplicationLinkRequestFieldPathsNested
	}

	for name, subs := range _processPaths(append(paths[:0:0], paths...)) {
		_ = subs
		switch name {
		case "application_ids":

			if v, ok := interface{}(&m.ApplicationIdentifiers).(interface{ ValidateFields(...string) error }); ok {
				if err := v.ValidateFields(subs...); err != nil {
					return SetApplicationLinkRequestValidationError{
						field:  "application_ids",
						reason: "embedded message failed validation",
						cause:  err,
					}
				}
			}

		case "link":

			if v, ok := interface{}(&m.ApplicationLink).(interface{ ValidateFields(...string) error }); ok {
				if err := v.ValidateFields(subs...); err != nil {
					return SetApplicationLinkRequestValidationError{
						field:  "link",
						reason: "embedded message failed validation",
						cause:  err,
					}
				}
			}

		case "field_mask":

			if v, ok := interface{}(&m.FieldMask).(interface{ ValidateFields(...string) error }); ok {
				if err := v.ValidateFields(subs...); err != nil {
					return SetApplicationLinkRequestValidationError{
						field:  "field_mask",
						reason: "embedded message failed validation",
						cause:  err,
					}
				}
			}

		default:
			return SetApplicationLinkRequestValidationError{
				field:  name,
				reason: "invalid field path",
			}
		}
	}
	return nil
}

// SetApplicationLinkRequestValidationError is the validation error returned by
// SetApplicationLinkRequest.ValidateFields if the designated constraints
// aren't met.
type SetApplicationLinkRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e SetApplicationLinkRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e SetApplicationLinkRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e SetApplicationLinkRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e SetApplicationLinkRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e SetApplicationLinkRequestValidationError) ErrorName() string {
	return "SetApplicationLinkRequestValidationError"
}

// Error satisfies the builtin error interface
func (e SetApplicationLinkRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sSetApplicationLinkRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = SetApplicationLinkRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = SetApplicationLinkRequestValidationError{}

// ValidateFields checks the field values on ApplicationLinkStats with the
// rules defined in the proto definition for this message. If any rules are
// violated, an error is returned.
func (m *ApplicationLinkStats) ValidateFields(paths ...string) error {
	if m == nil {
		return nil
	}

	if len(paths) == 0 {
		paths = ApplicationLinkStatsFieldPathsNested
	}

	for name, subs := range _processPaths(append(paths[:0:0], paths...)) {
		_ = subs
		switch name {
		case "linked_at":

			if v, ok := interface{}(m.GetLinkedAt()).(interface{ ValidateFields(...string) error }); ok {
				if err := v.ValidateFields(subs...); err != nil {
					return ApplicationLinkStatsValidationError{
						field:  "linked_at",
						reason: "embedded message failed validation",
						cause:  err,
					}
				}
			}

		case "network_server_address":

			if !_ApplicationLinkStats_NetworkServerAddress_Pattern.MatchString(m.GetNetworkServerAddress()) {
				return ApplicationLinkStatsValidationError{
					field:  "network_server_address",
					reason: "value does not match regex pattern \"^(?:(?:[a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9\\\\-]*[a-zA-Z0-9])\\\\.)*(?:[A-Za-z0-9]|[A-Za-z0-9][A-Za-z0-9\\\\-]*[A-Za-z0-9])(?::[0-9]{1,5})?$|^$\"",
				}
			}

		case "last_up_received_at":

			if v, ok := interface{}(m.GetLastUpReceivedAt()).(interface{ ValidateFields(...string) error }); ok {
				if err := v.ValidateFields(subs...); err != nil {
					return ApplicationLinkStatsValidationError{
						field:  "last_up_received_at",
						reason: "embedded message failed validation",
						cause:  err,
					}
				}
			}

		case "up_count":
			// no validation rules for UpCount
		case "last_downlink_forwarded_at":

			if v, ok := interface{}(m.GetLastDownlinkForwardedAt()).(interface{ ValidateFields(...string) error }); ok {
				if err := v.ValidateFields(subs...); err != nil {
					return ApplicationLinkStatsValidationError{
						field:  "last_downlink_forwarded_at",
						reason: "embedded message failed validation",
						cause:  err,
					}
				}
			}

		case "downlink_count":
			// no validation rules for DownlinkCount
		default:
			return ApplicationLinkStatsValidationError{
				field:  name,
				reason: "invalid field path",
			}
		}
	}
	return nil
}

// ApplicationLinkStatsValidationError is the validation error returned by
// ApplicationLinkStats.ValidateFields if the designated constraints aren't met.
type ApplicationLinkStatsValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e ApplicationLinkStatsValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e ApplicationLinkStatsValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e ApplicationLinkStatsValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e ApplicationLinkStatsValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e ApplicationLinkStatsValidationError) ErrorName() string {
	return "ApplicationLinkStatsValidationError"
}

// Error satisfies the builtin error interface
func (e ApplicationLinkStatsValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sApplicationLinkStats.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = ApplicationLinkStatsValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = ApplicationLinkStatsValidationError{}

var _ApplicationLinkStats_NetworkServerAddress_Pattern = regexp.MustCompile("^(?:(?:[a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9\\-]*[a-zA-Z0-9])\\.)*(?:[A-Za-z0-9]|[A-Za-z0-9][A-Za-z0-9\\-]*[A-Za-z0-9])(?::[0-9]{1,5})?$|^$")

// ValidateFields checks the field values on NsAsHandleUplinkRequest with the
// rules defined in the proto definition for this message. If any rules are
// violated, an error is returned.
func (m *NsAsHandleUplinkRequest) ValidateFields(paths ...string) error {
	if m == nil {
		return nil
	}

	if len(paths) == 0 {
		paths = NsAsHandleUplinkRequestFieldPathsNested
	}

	for name, subs := range _processPaths(append(paths[:0:0], paths...)) {
		_ = subs
		switch name {
		case "application_ups":

			if len(m.GetApplicationUps()) < 1 {
				return NsAsHandleUplinkRequestValidationError{
					field:  "application_ups",
					reason: "value must contain at least 1 item(s)",
				}
			}

			for idx, item := range m.GetApplicationUps() {
				_, _ = idx, item

				if v, ok := interface{}(item).(interface{ ValidateFields(...string) error }); ok {
					if err := v.ValidateFields(subs...); err != nil {
						return NsAsHandleUplinkRequestValidationError{
							field:  fmt.Sprintf("application_ups[%v]", idx),
							reason: "embedded message failed validation",
							cause:  err,
						}
					}
				}

			}

		default:
			return NsAsHandleUplinkRequestValidationError{
				field:  name,
				reason: "invalid field path",
			}
		}
	}
	return nil
}

// NsAsHandleUplinkRequestValidationError is the validation error returned by
// NsAsHandleUplinkRequest.ValidateFields if the designated constraints aren't met.
type NsAsHandleUplinkRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e NsAsHandleUplinkRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e NsAsHandleUplinkRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e NsAsHandleUplinkRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e NsAsHandleUplinkRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e NsAsHandleUplinkRequestValidationError) ErrorName() string {
	return "NsAsHandleUplinkRequestValidationError"
}

// Error satisfies the builtin error interface
func (e NsAsHandleUplinkRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sNsAsHandleUplinkRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = NsAsHandleUplinkRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = NsAsHandleUplinkRequestValidationError{}
