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

// ValidateFields checks the field values on Rights with the rules defined in
// the proto definition for this message. If any rules are violated, an error
// is returned.
func (m *Rights) ValidateFields(paths ...string) error {
	if m == nil {
		return nil
	}

	if len(paths) == 0 {
		paths = RightsFieldPathsNested
	}

	for name, subs := range _processPaths(append(paths[:0:0], paths...)) {
		_ = subs
		switch name {
		case "rights":

		default:
			return RightsValidationError{
				field:  name,
				reason: "invalid field path",
			}
		}
	}
	return nil
}

// RightsValidationError is the validation error returned by
// Rights.ValidateFields if the designated constraints aren't met.
type RightsValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e RightsValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e RightsValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e RightsValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e RightsValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e RightsValidationError) ErrorName() string { return "RightsValidationError" }

// Error satisfies the builtin error interface
func (e RightsValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sRights.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = RightsValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = RightsValidationError{}

// ValidateFields checks the field values on APIKey with the rules defined in
// the proto definition for this message. If any rules are violated, an error
// is returned.
func (m *APIKey) ValidateFields(paths ...string) error {
	if m == nil {
		return nil
	}

	if len(paths) == 0 {
		paths = APIKeyFieldPathsNested
	}

	for name, subs := range _processPaths(append(paths[:0:0], paths...)) {
		_ = subs
		switch name {
		case "id":
			// no validation rules for ID
		case "key":
			// no validation rules for Key
		case "name":
			// no validation rules for Name
		case "rights":

		default:
			return APIKeyValidationError{
				field:  name,
				reason: "invalid field path",
			}
		}
	}
	return nil
}

// APIKeyValidationError is the validation error returned by
// APIKey.ValidateFields if the designated constraints aren't met.
type APIKeyValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e APIKeyValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e APIKeyValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e APIKeyValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e APIKeyValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e APIKeyValidationError) ErrorName() string { return "APIKeyValidationError" }

// Error satisfies the builtin error interface
func (e APIKeyValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sAPIKey.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = APIKeyValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = APIKeyValidationError{}

// ValidateFields checks the field values on APIKeys with the rules defined in
// the proto definition for this message. If any rules are violated, an error
// is returned.
func (m *APIKeys) ValidateFields(paths ...string) error {
	if m == nil {
		return nil
	}

	if len(paths) == 0 {
		paths = APIKeysFieldPathsNested
	}

	for name, subs := range _processPaths(append(paths[:0:0], paths...)) {
		_ = subs
		switch name {
		case "api_keys":

			for idx, item := range m.GetAPIKeys() {
				_, _ = idx, item

				if v, ok := interface{}(item).(interface{ ValidateFields(...string) error }); ok {
					if err := v.ValidateFields(subs...); err != nil {
						return APIKeysValidationError{
							field:  fmt.Sprintf("api_keys[%v]", idx),
							reason: "embedded message failed validation",
							cause:  err,
						}
					}
				}

			}

		default:
			return APIKeysValidationError{
				field:  name,
				reason: "invalid field path",
			}
		}
	}
	return nil
}

// APIKeysValidationError is the validation error returned by
// APIKeys.ValidateFields if the designated constraints aren't met.
type APIKeysValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e APIKeysValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e APIKeysValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e APIKeysValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e APIKeysValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e APIKeysValidationError) ErrorName() string { return "APIKeysValidationError" }

// Error satisfies the builtin error interface
func (e APIKeysValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sAPIKeys.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = APIKeysValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = APIKeysValidationError{}

// ValidateFields checks the field values on Collaborator with the rules
// defined in the proto definition for this message. If any rules are
// violated, an error is returned.
func (m *Collaborator) ValidateFields(paths ...string) error {
	if m == nil {
		return nil
	}

	if len(paths) == 0 {
		paths = CollaboratorFieldPathsNested
	}

	for name, subs := range _processPaths(append(paths[:0:0], paths...)) {
		_ = subs
		switch name {
		case "ids":

			if v, ok := interface{}(&m.OrganizationOrUserIdentifiers).(interface{ ValidateFields(...string) error }); ok {
				if err := v.ValidateFields(subs...); err != nil {
					return CollaboratorValidationError{
						field:  "ids",
						reason: "embedded message failed validation",
						cause:  err,
					}
				}
			}

		case "rights":

		default:
			return CollaboratorValidationError{
				field:  name,
				reason: "invalid field path",
			}
		}
	}
	return nil
}

// CollaboratorValidationError is the validation error returned by
// Collaborator.ValidateFields if the designated constraints aren't met.
type CollaboratorValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e CollaboratorValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e CollaboratorValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e CollaboratorValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e CollaboratorValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e CollaboratorValidationError) ErrorName() string { return "CollaboratorValidationError" }

// Error satisfies the builtin error interface
func (e CollaboratorValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sCollaborator.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = CollaboratorValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = CollaboratorValidationError{}

// ValidateFields checks the field values on Collaborators with the rules
// defined in the proto definition for this message. If any rules are
// violated, an error is returned.
func (m *Collaborators) ValidateFields(paths ...string) error {
	if m == nil {
		return nil
	}

	if len(paths) == 0 {
		paths = CollaboratorsFieldPathsNested
	}

	for name, subs := range _processPaths(append(paths[:0:0], paths...)) {
		_ = subs
		switch name {
		case "collaborators":

			for idx, item := range m.GetCollaborators() {
				_, _ = idx, item

				if v, ok := interface{}(item).(interface{ ValidateFields(...string) error }); ok {
					if err := v.ValidateFields(subs...); err != nil {
						return CollaboratorsValidationError{
							field:  fmt.Sprintf("collaborators[%v]", idx),
							reason: "embedded message failed validation",
							cause:  err,
						}
					}
				}

			}

		default:
			return CollaboratorsValidationError{
				field:  name,
				reason: "invalid field path",
			}
		}
	}
	return nil
}

// CollaboratorsValidationError is the validation error returned by
// Collaborators.ValidateFields if the designated constraints aren't met.
type CollaboratorsValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e CollaboratorsValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e CollaboratorsValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e CollaboratorsValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e CollaboratorsValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e CollaboratorsValidationError) ErrorName() string { return "CollaboratorsValidationError" }

// Error satisfies the builtin error interface
func (e CollaboratorsValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sCollaborators.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = CollaboratorsValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = CollaboratorsValidationError{}
