// Copyright © 2018. All rights reserved.
// Author: Alice Qio.
// Contacts: <qioalice@gmail.com>.
// License: https://opensource.org/licenses/MIT

package errors

import (
	"strconv"
)

// Error represents an error's interface - abstract type of representing
// this SDK some error.
// Compatible with golang error type.
type Error interface {

	// Code returns an Error Code of occurred error.
	Code() Code

	// What returns a message of occurred error.
	What() string

	// IsIt returns true if e represents the same "kind of error"
	// that represents current object.
	IsIt(e error) bool

	// Error returns a string representation of error object.
	// Provides compatibility with embedded golang error's type error.
	Error() string

	// String returns a string representation of error object.
	// Provides compatibility with golang type fmt.Stringer.
	String() string
}

// BaseError is default implementation of Error interface.
// In some places used as base type for extra SDK error types.
type BaseError struct {

	// TODO: Replace strconv algoritms by more faster's

	code Code
	what string
}

// Code returns an Error Code of occurred error (ECOK if be is nil)
func (be *BaseError) Code() Code {
	if be == nil {
		return ECOK
	}
	return be.code
}

// What returns a message (without error code) of occurred error
// (empty message if be is nil or code == ECOK).
func (be *BaseError) What() string {
	if be == nil || be.code == ECOK {
		return ""
	}
	return be.what
}

// IsIt returns true if e is *BaseError and their error codes and messages
// are the same.
func (be *BaseError) IsIt(e error) bool {
	ebe, ok := e.(*BaseError)
	return ok && (be == ebe ||
		be != nil && ebe != nil && be.code == ebe.code && be.what == ebe.what)
}

// Error returns an error code and message of occurred error
// in the following format: "(<code>): <message>".
// Works the same as String method.
//
// Returns "(<code>)" string if message is empty.
// Returns an empty string if be is nil or code == ECOK.
func (be *BaseError) Error() string {

	// TODO: Optimize algo, refuse to use strconv algo, use once allocated []byte
	// with calculated sizes to store EC, format symbols (parentheses, colon),
	// and what string.

	if be == nil || be.code == ECOK {
		return ""
	}

	// be.code != ECOK
	if be.what == "" {
		return "(" + strconv.Itoa(int(be.code)) + ")"
	}

	// be.code != EOK && be.what != ""
	return "(" + strconv.Itoa(int(be.code)) + "): " + be.what
}

// Error returns an error code and message of occurred error
// in the following format: "(<code>): <message>".
// Works the same as Error method.
//
// Returns "(<code>)" string if message is empty.
// Returns an empty string if be is nil or code == ECOK.
func (be *BaseError) String() string {
	return be.Error()
}
