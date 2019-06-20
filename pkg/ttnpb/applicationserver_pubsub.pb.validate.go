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

// ValidateFields checks the field values on ApplicationPubSubIdentifiers with
// the rules defined in the proto definition for this message. If any rules
// are violated, an error is returned.
func (m *ApplicationPubSubIdentifiers) ValidateFields(paths ...string) error {
	if m == nil {
		return nil
	}

	if len(paths) == 0 {
		paths = ApplicationPubSubIdentifiersFieldPathsNested
	}

	for name, subs := range _processPaths(append(paths[:0:0], paths...)) {
		_ = subs
		switch name {
		case "application_ids":

			if v, ok := interface{}(&m.ApplicationIdentifiers).(interface{ ValidateFields(...string) error }); ok {
				if err := v.ValidateFields(subs...); err != nil {
					return ApplicationPubSubIdentifiersValidationError{
						field:  "application_ids",
						reason: "embedded message failed validation",
						cause:  err,
					}
				}
			}

		case "pubsub_id":

			if utf8.RuneCountInString(m.GetPubSubID()) > 36 {
				return ApplicationPubSubIdentifiersValidationError{
					field:  "pubsub_id",
					reason: "value length must be at most 36 runes",
				}
			}

			if !_ApplicationPubSubIdentifiers_PubSubID_Pattern.MatchString(m.GetPubSubID()) {
				return ApplicationPubSubIdentifiersValidationError{
					field:  "pubsub_id",
					reason: "value does not match regex pattern \"^[a-z0-9](?:[-]?[a-z0-9]){2,}$\"",
				}
			}

		default:
			return ApplicationPubSubIdentifiersValidationError{
				field:  name,
				reason: "invalid field path",
			}
		}
	}
	return nil
}

// ApplicationPubSubIdentifiersValidationError is the validation error returned
// by ApplicationPubSubIdentifiers.ValidateFields if the designated
// constraints aren't met.
type ApplicationPubSubIdentifiersValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e ApplicationPubSubIdentifiersValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e ApplicationPubSubIdentifiersValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e ApplicationPubSubIdentifiersValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e ApplicationPubSubIdentifiersValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e ApplicationPubSubIdentifiersValidationError) ErrorName() string {
	return "ApplicationPubSubIdentifiersValidationError"
}

// Error satisfies the builtin error interface
func (e ApplicationPubSubIdentifiersValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sApplicationPubSubIdentifiers.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = ApplicationPubSubIdentifiersValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = ApplicationPubSubIdentifiersValidationError{}

var _ApplicationPubSubIdentifiers_PubSubID_Pattern = regexp.MustCompile("^[a-z0-9](?:[-]?[a-z0-9]){2,}$")

// ValidateFields checks the field values on ApplicationPubSub with the rules
// defined in the proto definition for this message. If any rules are
// violated, an error is returned.
func (m *ApplicationPubSub) ValidateFields(paths ...string) error {
	if m == nil {
		return nil
	}

	if len(paths) == 0 {
		paths = ApplicationPubSubFieldPathsNested
	}

	for name, subs := range _processPaths(append(paths[:0:0], paths...)) {
		_ = subs
		switch name {
		case "ids":

			if v, ok := interface{}(&m.ApplicationPubSubIdentifiers).(interface{ ValidateFields(...string) error }); ok {
				if err := v.ValidateFields(subs...); err != nil {
					return ApplicationPubSubValidationError{
						field:  "ids",
						reason: "embedded message failed validation",
						cause:  err,
					}
				}
			}

		case "created_at":

			if v, ok := interface{}(&m.CreatedAt).(interface{ ValidateFields(...string) error }); ok {
				if err := v.ValidateFields(subs...); err != nil {
					return ApplicationPubSubValidationError{
						field:  "created_at",
						reason: "embedded message failed validation",
						cause:  err,
					}
				}
			}

		case "updated_at":

			if v, ok := interface{}(&m.UpdatedAt).(interface{ ValidateFields(...string) error }); ok {
				if err := v.ValidateFields(subs...); err != nil {
					return ApplicationPubSubValidationError{
						field:  "updated_at",
						reason: "embedded message failed validation",
						cause:  err,
					}
				}
			}

		case "attributes":
			// no validation rules for Attributes
		case "format":
			// no validation rules for Format
		case "service":

			if _, ok := ApplicationPubSub_Service_name[int32(m.GetService())]; !ok {
				return ApplicationPubSubValidationError{
					field:  "service",
					reason: "value must be one of the defined enum values",
				}
			}

		case "downlink_push_topic":
			// no validation rules for DownlinkPushTopic
		case "downlink_replace_topic":
			// no validation rules for DownlinkReplaceTopic
		case "uplink_message":

			if v, ok := interface{}(m.GetUplinkMessage()).(interface{ ValidateFields(...string) error }); ok {
				if err := v.ValidateFields(subs...); err != nil {
					return ApplicationPubSubValidationError{
						field:  "uplink_message",
						reason: "embedded message failed validation",
						cause:  err,
					}
				}
			}

		case "join_accept":

			if v, ok := interface{}(m.GetJoinAccept()).(interface{ ValidateFields(...string) error }); ok {
				if err := v.ValidateFields(subs...); err != nil {
					return ApplicationPubSubValidationError{
						field:  "join_accept",
						reason: "embedded message failed validation",
						cause:  err,
					}
				}
			}

		case "downlink_ack":

			if v, ok := interface{}(m.GetDownlinkAck()).(interface{ ValidateFields(...string) error }); ok {
				if err := v.ValidateFields(subs...); err != nil {
					return ApplicationPubSubValidationError{
						field:  "downlink_ack",
						reason: "embedded message failed validation",
						cause:  err,
					}
				}
			}

		case "downlink_nack":

			if v, ok := interface{}(m.GetDownlinkNack()).(interface{ ValidateFields(...string) error }); ok {
				if err := v.ValidateFields(subs...); err != nil {
					return ApplicationPubSubValidationError{
						field:  "downlink_nack",
						reason: "embedded message failed validation",
						cause:  err,
					}
				}
			}

		case "downlink_sent":

			if v, ok := interface{}(m.GetDownlinkSent()).(interface{ ValidateFields(...string) error }); ok {
				if err := v.ValidateFields(subs...); err != nil {
					return ApplicationPubSubValidationError{
						field:  "downlink_sent",
						reason: "embedded message failed validation",
						cause:  err,
					}
				}
			}

		case "downlink_failed":

			if v, ok := interface{}(m.GetDownlinkFailed()).(interface{ ValidateFields(...string) error }); ok {
				if err := v.ValidateFields(subs...); err != nil {
					return ApplicationPubSubValidationError{
						field:  "downlink_failed",
						reason: "embedded message failed validation",
						cause:  err,
					}
				}
			}

		case "downlink_queued":

			if v, ok := interface{}(m.GetDownlinkQueued()).(interface{ ValidateFields(...string) error }); ok {
				if err := v.ValidateFields(subs...); err != nil {
					return ApplicationPubSubValidationError{
						field:  "downlink_queued",
						reason: "embedded message failed validation",
						cause:  err,
					}
				}
			}

		case "location_solved":

			if v, ok := interface{}(m.GetLocationSolved()).(interface{ ValidateFields(...string) error }); ok {
				if err := v.ValidateFields(subs...); err != nil {
					return ApplicationPubSubValidationError{
						field:  "location_solved",
						reason: "embedded message failed validation",
						cause:  err,
					}
				}
			}

		default:
			return ApplicationPubSubValidationError{
				field:  name,
				reason: "invalid field path",
			}
		}
	}
	return nil
}

// ApplicationPubSubValidationError is the validation error returned by
// ApplicationPubSub.ValidateFields if the designated constraints aren't met.
type ApplicationPubSubValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e ApplicationPubSubValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e ApplicationPubSubValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e ApplicationPubSubValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e ApplicationPubSubValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e ApplicationPubSubValidationError) ErrorName() string {
	return "ApplicationPubSubValidationError"
}

// Error satisfies the builtin error interface
func (e ApplicationPubSubValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sApplicationPubSub.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = ApplicationPubSubValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = ApplicationPubSubValidationError{}

// ValidateFields checks the field values on ApplicationPubSubs with the rules
// defined in the proto definition for this message. If any rules are
// violated, an error is returned.
func (m *ApplicationPubSubs) ValidateFields(paths ...string) error {
	if m == nil {
		return nil
	}

	if len(paths) == 0 {
		paths = ApplicationPubSubsFieldPathsNested
	}

	for name, subs := range _processPaths(append(paths[:0:0], paths...)) {
		_ = subs
		switch name {
		case "pubsubs":

			for idx, item := range m.GetPubsubs() {
				_, _ = idx, item

				if v, ok := interface{}(item).(interface{ ValidateFields(...string) error }); ok {
					if err := v.ValidateFields(subs...); err != nil {
						return ApplicationPubSubsValidationError{
							field:  fmt.Sprintf("pubsubs[%v]", idx),
							reason: "embedded message failed validation",
							cause:  err,
						}
					}
				}

			}

		default:
			return ApplicationPubSubsValidationError{
				field:  name,
				reason: "invalid field path",
			}
		}
	}
	return nil
}

// ApplicationPubSubsValidationError is the validation error returned by
// ApplicationPubSubs.ValidateFields if the designated constraints aren't met.
type ApplicationPubSubsValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e ApplicationPubSubsValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e ApplicationPubSubsValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e ApplicationPubSubsValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e ApplicationPubSubsValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e ApplicationPubSubsValidationError) ErrorName() string {
	return "ApplicationPubSubsValidationError"
}

// Error satisfies the builtin error interface
func (e ApplicationPubSubsValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sApplicationPubSubs.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = ApplicationPubSubsValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = ApplicationPubSubsValidationError{}

// ValidateFields checks the field values on ApplicationPubSubFormats with the
// rules defined in the proto definition for this message. If any rules are
// violated, an error is returned.
func (m *ApplicationPubSubFormats) ValidateFields(paths ...string) error {
	if m == nil {
		return nil
	}

	if len(paths) == 0 {
		paths = ApplicationPubSubFormatsFieldPathsNested
	}

	for name, subs := range _processPaths(append(paths[:0:0], paths...)) {
		_ = subs
		switch name {
		case "formats":
			// no validation rules for Formats
		default:
			return ApplicationPubSubFormatsValidationError{
				field:  name,
				reason: "invalid field path",
			}
		}
	}
	return nil
}

// ApplicationPubSubFormatsValidationError is the validation error returned by
// ApplicationPubSubFormats.ValidateFields if the designated constraints
// aren't met.
type ApplicationPubSubFormatsValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e ApplicationPubSubFormatsValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e ApplicationPubSubFormatsValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e ApplicationPubSubFormatsValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e ApplicationPubSubFormatsValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e ApplicationPubSubFormatsValidationError) ErrorName() string {
	return "ApplicationPubSubFormatsValidationError"
}

// Error satisfies the builtin error interface
func (e ApplicationPubSubFormatsValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sApplicationPubSubFormats.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = ApplicationPubSubFormatsValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = ApplicationPubSubFormatsValidationError{}

// ValidateFields checks the field values on GetApplicationPubSubRequest with
// the rules defined in the proto definition for this message. If any rules
// are violated, an error is returned.
func (m *GetApplicationPubSubRequest) ValidateFields(paths ...string) error {
	if m == nil {
		return nil
	}

	if len(paths) == 0 {
		paths = GetApplicationPubSubRequestFieldPathsNested
	}

	for name, subs := range _processPaths(append(paths[:0:0], paths...)) {
		_ = subs
		switch name {
		case "ids":

			if v, ok := interface{}(&m.ApplicationPubSubIdentifiers).(interface{ ValidateFields(...string) error }); ok {
				if err := v.ValidateFields(subs...); err != nil {
					return GetApplicationPubSubRequestValidationError{
						field:  "ids",
						reason: "embedded message failed validation",
						cause:  err,
					}
				}
			}

		case "field_mask":

			if v, ok := interface{}(&m.FieldMask).(interface{ ValidateFields(...string) error }); ok {
				if err := v.ValidateFields(subs...); err != nil {
					return GetApplicationPubSubRequestValidationError{
						field:  "field_mask",
						reason: "embedded message failed validation",
						cause:  err,
					}
				}
			}

		default:
			return GetApplicationPubSubRequestValidationError{
				field:  name,
				reason: "invalid field path",
			}
		}
	}
	return nil
}

// GetApplicationPubSubRequestValidationError is the validation error returned
// by GetApplicationPubSubRequest.ValidateFields if the designated constraints
// aren't met.
type GetApplicationPubSubRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e GetApplicationPubSubRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e GetApplicationPubSubRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e GetApplicationPubSubRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e GetApplicationPubSubRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e GetApplicationPubSubRequestValidationError) ErrorName() string {
	return "GetApplicationPubSubRequestValidationError"
}

// Error satisfies the builtin error interface
func (e GetApplicationPubSubRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sGetApplicationPubSubRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = GetApplicationPubSubRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = GetApplicationPubSubRequestValidationError{}

// ValidateFields checks the field values on ListApplicationPubSubsRequest with
// the rules defined in the proto definition for this message. If any rules
// are violated, an error is returned.
func (m *ListApplicationPubSubsRequest) ValidateFields(paths ...string) error {
	if m == nil {
		return nil
	}

	if len(paths) == 0 {
		paths = ListApplicationPubSubsRequestFieldPathsNested
	}

	for name, subs := range _processPaths(append(paths[:0:0], paths...)) {
		_ = subs
		switch name {
		case "application_ids":

			if v, ok := interface{}(&m.ApplicationIdentifiers).(interface{ ValidateFields(...string) error }); ok {
				if err := v.ValidateFields(subs...); err != nil {
					return ListApplicationPubSubsRequestValidationError{
						field:  "application_ids",
						reason: "embedded message failed validation",
						cause:  err,
					}
				}
			}

		case "field_mask":

			if v, ok := interface{}(&m.FieldMask).(interface{ ValidateFields(...string) error }); ok {
				if err := v.ValidateFields(subs...); err != nil {
					return ListApplicationPubSubsRequestValidationError{
						field:  "field_mask",
						reason: "embedded message failed validation",
						cause:  err,
					}
				}
			}

		default:
			return ListApplicationPubSubsRequestValidationError{
				field:  name,
				reason: "invalid field path",
			}
		}
	}
	return nil
}

// ListApplicationPubSubsRequestValidationError is the validation error
// returned by ListApplicationPubSubsRequest.ValidateFields if the designated
// constraints aren't met.
type ListApplicationPubSubsRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e ListApplicationPubSubsRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e ListApplicationPubSubsRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e ListApplicationPubSubsRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e ListApplicationPubSubsRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e ListApplicationPubSubsRequestValidationError) ErrorName() string {
	return "ListApplicationPubSubsRequestValidationError"
}

// Error satisfies the builtin error interface
func (e ListApplicationPubSubsRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sListApplicationPubSubsRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = ListApplicationPubSubsRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = ListApplicationPubSubsRequestValidationError{}

// ValidateFields checks the field values on SetApplicationPubSubRequest with
// the rules defined in the proto definition for this message. If any rules
// are violated, an error is returned.
func (m *SetApplicationPubSubRequest) ValidateFields(paths ...string) error {
	if m == nil {
		return nil
	}

	if len(paths) == 0 {
		paths = SetApplicationPubSubRequestFieldPathsNested
	}

	for name, subs := range _processPaths(append(paths[:0:0], paths...)) {
		_ = subs
		switch name {
		case "pubsub":

			if v, ok := interface{}(&m.ApplicationPubSub).(interface{ ValidateFields(...string) error }); ok {
				if err := v.ValidateFields(subs...); err != nil {
					return SetApplicationPubSubRequestValidationError{
						field:  "pubsub",
						reason: "embedded message failed validation",
						cause:  err,
					}
				}
			}

		case "field_mask":

			if v, ok := interface{}(&m.FieldMask).(interface{ ValidateFields(...string) error }); ok {
				if err := v.ValidateFields(subs...); err != nil {
					return SetApplicationPubSubRequestValidationError{
						field:  "field_mask",
						reason: "embedded message failed validation",
						cause:  err,
					}
				}
			}

		default:
			return SetApplicationPubSubRequestValidationError{
				field:  name,
				reason: "invalid field path",
			}
		}
	}
	return nil
}

// SetApplicationPubSubRequestValidationError is the validation error returned
// by SetApplicationPubSubRequest.ValidateFields if the designated constraints
// aren't met.
type SetApplicationPubSubRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e SetApplicationPubSubRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e SetApplicationPubSubRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e SetApplicationPubSubRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e SetApplicationPubSubRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e SetApplicationPubSubRequestValidationError) ErrorName() string {
	return "SetApplicationPubSubRequestValidationError"
}

// Error satisfies the builtin error interface
func (e SetApplicationPubSubRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sSetApplicationPubSubRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = SetApplicationPubSubRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = SetApplicationPubSubRequestValidationError{}

// ValidateFields checks the field values on ApplicationPubSub_Message with the
// rules defined in the proto definition for this message. If any rules are
// violated, an error is returned.
func (m *ApplicationPubSub_Message) ValidateFields(paths ...string) error {
	if m == nil {
		return nil
	}

	if len(paths) == 0 {
		paths = ApplicationPubSub_MessageFieldPathsNested
	}

	for name, subs := range _processPaths(append(paths[:0:0], paths...)) {
		_ = subs
		switch name {
		case "topic":
			// no validation rules for Topic
		default:
			return ApplicationPubSub_MessageValidationError{
				field:  name,
				reason: "invalid field path",
			}
		}
	}
	return nil
}

// ApplicationPubSub_MessageValidationError is the validation error returned by
// ApplicationPubSub_Message.ValidateFields if the designated constraints
// aren't met.
type ApplicationPubSub_MessageValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e ApplicationPubSub_MessageValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e ApplicationPubSub_MessageValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e ApplicationPubSub_MessageValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e ApplicationPubSub_MessageValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e ApplicationPubSub_MessageValidationError) ErrorName() string {
	return "ApplicationPubSub_MessageValidationError"
}

// Error satisfies the builtin error interface
func (e ApplicationPubSub_MessageValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sApplicationPubSub_Message.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = ApplicationPubSub_MessageValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = ApplicationPubSub_MessageValidationError{}
