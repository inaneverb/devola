// Copyright © 2018. All rights reserved.
// Author: Alice Qio.
// Contacts: <qioalice@gmail.com>.
// License: https://opensource.org/licenses/MIT

package registrator

import (
	"strings"

	"github.com/qioalice/devola/core/errors"
)

// EBadCallback represents an SDK error that occurred while some incorrect callback
// has been tried to register as event handler or middleware.
//
// Returned by Handler, MainHandler, Middleware, MainMiddleware methods of
// Registrator type.
// Used for both kinds of errors: for handler's and for middleware's.
//
// You can figure out what kind of error is occurred using Code method:
//
// "e.Code() == ECBadHandler" for handler errors
// "e.Code() == ECBadMiddleware for middleware errors"
//
// or using IsIt method:
//
// "e.IsIt((*EBadCallback)(nil))" covers both handler and middleware errors,
// "e.IsIt(&EBadCallback{ IsMiddleware: false })" covers only handler errors,
// "e.IsIt(&EBadCallback{ IsMiddleware: true })" - only middleware's.
type EBadCallback struct {

	// IsMiddleware true if this EBadCallback reports
	// about middleware registration error, false - about handler's.
	IsMiddleware bool

	// IsNil true if this EBadCallback reports about nil handler or middleware
	// has been tried to be registered.
	IsNil bool

	// WantType is a string representation about what type of callback
	// SHOULD BE USED for requested registering operation.
	WantType string

	// HaveType is a string representation about what type of callback
	// WAS USED for requested registering operation.
	HaveType string

	// AppliedTo is a set of rules to which requested registering operation
	// has been performed and caused error.
	AppliedTo []rule
}

// Code returns ECBadHandler if e is error about bad handler type,
// ECBadMiddleware - bad middleware, and errors.ECOK if there's no error.
func (e *EBadCallback) Code() errors.Code {
	switch {
	case e == nil:
		return errors.ECOK
	case !e.IsMiddleware:
		return ECBadHandler
	case e.IsMiddleware:
		return ECBadMiddleware
	}
	return -1 // <-- never, but go requires
}

// What returns a different predefined descriptions about bad handler type or
// bad middleware type, or empty string if there's no error.
func (e *EBadCallback) What() string {
	switch {
	case e == nil:
		return ""
	case !e.IsMiddleware:
		return "The type of passed handler is incompatible with required type by context."
	case e.IsMiddleware:
		return "The type of passed middleware is incompatible with required type by context."
	}
	return "" // <-- never, but go requires
}

// IsIt returns true when e2 is EBadCallback but nil
// or not nil but its IsMiddleware field is the same as e's
// or pointer to e2 and pointer to e are equals.
func (e *EBadCallback) IsIt(e2 error) bool {
	e2t, ok := e2.(*EBadCallback)
	return ok && (e2t == e || e2t == nil ||
		e2t != nil && e != nil && e.IsMiddleware == e2t.IsMiddleware)
}

// Error returns a string representation of error in the following format:
// "<e.What()> Want type: <T1>, Have type: <T2>. Covered rules: <R>",
// where T1 - desired type, T2 - current type, R - rules to which that
// error is related.
//
// Returns an empty string if e is nil.
func (e *EBadCallback) Error() string {

	// TODO: Optimize algo, refuse to use strings algorithms, use once allocated []byte

	if e == nil {
		return ""
	}

	s := e.What()
	s += " Want type: " + e.WantType + ", Have type: " + e.HaveType + ". "

	if len(e.AppliedTo) != 0 {
		rs := make([]string, 0, len(e.AppliedTo))
		for _, rule := range e.AppliedTo {
			rs = append(rs, rule.String())
		}
		s += "Covered rules: " + strings.Join(rs, ", ")
	}

	return s
}

// String returns a string representation of error in the following format:
// "<e.What()> Want type: <T1>, Have type: <T2>. Covered rules: <R>",
// where T1 - desired type, T2 - current type, R - rules to which that
// error is related.
//
// Returns an empty string if e is nil.
func (e *EBadCallback) String() string {
	return e.Error()
}
