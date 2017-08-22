// Copyright © 2017 The Things Network Foundation, distributed under the MIT license (see LICENSE file)

package errors

import "fmt"

// ErrDescriptor is a helper struct to easily build new Errors from and to be
// the authoritive information about error codes.
//
// The descriptor can be used to find out information about the error after it
// has been handed over between components
type ErrDescriptor struct {
	// MessageFormat is the format of the error message. Attributes will be filled
	// in when an error is created using New(). For example:
	//
	//   "This is an error about user {username}"
	//
	// when passed an atrtributes map with "username" set to "john" would interpolate to
	//
	//   "This is an error about user john"
	//
	// The idea about this message format is that is is localizable
	MessageFormat string

	// Code is the code of errors that are created by this descriptor
	Code Code

	// Type is the type of errors created by this descriptor
	Type Type

	// registered denotes wether or not the error has been registered
	// (by a call to Register)
	registered bool
}

// New creates a new error based on the error descriptor
func (err *ErrDescriptor) New(attributes Attributes) Error {
	if err.Code != NoCode && !err.registered {
		panic(fmt.Errorf("Error descriptor with code %v was not registered", err.Code))
	}

	return &Impl{
		message:    Format(err.MessageFormat, attributes),
		code:       err.Code,
		typ:        err.Type,
		attributes: attributes,
	}
}

// New creates a new error based on the error descriptor and adds a cause
func (err *ErrDescriptor) NewWithCause(attributes Attributes, cause error) Error {
	attr := make(map[string]interface{}, len(attributes)+1)
	for k, v := range attributes {
		attr[k] = v
	}

	attr[causeKey] = cause

	return err.New(attr)
}

// Register registers the descriptor
func (err *ErrDescriptor) Register() {
	Register(err)
}
