// Copyright © 2018 The Things Network Foundation, The Things Industries B.V.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package errors

import (
	"encoding/json"
	"fmt"
	"strings"

	"google.golang.org/grpc/status"
)

type info struct {
	Message    string     `json:"error_message"`
	ID         string     `json:"error_id,omitempty"`
	Code       Code       `json:"error_code,omitempty"`
	Type       Type       `json:"error_type,omitempty"`
	Namespace  string     `json:"error_namespace,omitempty"`
	Attributes Attributes `json:"error_attributes,omitempty"`
}

// Impl implements Error
type Impl struct {
	descriptor *ErrDescriptor

	// info contains all the public information about the error, nested to
	// avoid name clashes.
	info info
}

// MarshalJSON implements json.Marshaler
func (i *Impl) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.info)
}

// UnmarshalJSON implements json.Unmarshaler
func (i *Impl) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &i.info)
}

// Error returns the formatted error message, prefixed with the error namespace
func (i *Impl) Error() string {
	prefix := ""
	if i.info.Namespace != "" {
		prefix = i.info.Namespace
	}

	if i.info.Code != NoCode {
		prefix = prefix + fmt.Sprintf("[%v]", i.info.Code)
	}

	message := i.info.Message

	if prefix != "" {
		message = strings.Trim(prefix, " ") + ": " + i.info.Message
	}

	if i.info.Code == NoCode && i.info.Attributes != nil && i.info.Attributes[causeKey] != nil {
		if cause, ok := i.info.Attributes[causeKey].(error); ok {
			message += fmt.Sprintf(" (%s)", cause.Error())
		}
	}

	return message
}

// Message returns the formatted error message
func (i *Impl) Message() string {
	return i.info.Message
}

// Code returns the error code
func (i *Impl) Code() Code {
	return i.info.Code
}

// Type returns the error type
func (i *Impl) Type() Type {
	return i.info.Type
}

// Attributes returns the error attributes
func (i *Impl) Attributes() map[string]interface{} {
	return i.info.Attributes
}

// Namespace returns the namespace of the error, which is usuallt the package it originates from.
func (i *Impl) Namespace() string {
	return i.info.Namespace
}

// Name returns the name of the error.
func (i *Impl) Name() string {
	return fmt.Sprintf("#%d", i.Code())
}

// ID returns the unique identifier of this error.
func (i *Impl) ID() string {
	return i.info.ID
}

// GRPCStatus returns a *status.Status from the error.
func (i *Impl) GRPCStatus() *status.Status {
	return CompatStatus(i)
}

// ToImpl creates an equivalent Impl for any Error
func ToImpl(err Error) *Impl {
	if i, ok := err.(*Impl); ok {
		return i
	}

	return normalize(&Impl{
		info: info{
			ID:         err.ID(),
			Message:    err.Message(),
			Code:       err.Code(),
			Type:       err.Type(),
			Attributes: err.Attributes(),
			Namespace:  err.Namespace(),
		},
	})
}

// normalize normalizes the error
func normalize(i *Impl) *Impl {
	if i.info.ID == "" {
		i.info.ID = NewID()
	}

	if i.descriptor == nil {
		i.descriptor = Get(i.info.Namespace, i.info.Code)
	}

	return i
}

// Fields implements log.Fielder.
func (i *Impl) Fields() map[string]interface{} {
	fields := make(map[string]interface{})

	for k, v := range i.Attributes() {
		fields[k] = v
	}

	fields["error_id"] = i.ID()
	fields["error"] = i.Message()
	fields["code"] = i.Code()
	fields["namespace"] = i.Namespace()
	fields["type"] = i.Type()

	return fields
}

// SafeImpl is the same as Impl, but only returns the safe attributes.
type SafeImpl struct {
	*Impl
}

// Attributes returns the safe attributes.
func (i *SafeImpl) Attributes() map[string]interface{} {
	if i.descriptor == nil {
		return i.Impl.Attributes()
	}

	attrs := i.Impl.Attributes()

	res := make(Attributes, len(i.descriptor.SafeAttributes))

	for _, key := range i.descriptor.SafeAttributes {
		if value, ok := attrs[key]; ok {
			res[key] = value
		}
	}

	return res
}

// Safe returns an error that only returns its safe attributes.
func Safe(err Error) Error {
	if i, ok := err.(*Impl); ok {
		return &SafeImpl{i}
	}

	return err
}
