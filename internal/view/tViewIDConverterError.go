// Copyright Â© 2018. All rights reserved.
// Author: Alice Qio.
// Contacts: <qioalice@gmail.com>.
// License: https://opensource.org/licenses/MIT
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"),
// to deal in the Software without restriction, including without limitation
// the rights to use, copy, modify, merge, publish, distribute, sublicense,
// and/or sell copies of the Software, and to permit persons to whom
// the Software is furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS
// IN THE SOFTWARE.

package tgbot

import (
	"strconv"
)

// tViewIDConverterError is an SDK error type for representing returnable
// errors from tViewIDConverter methods.
type tViewIDConverterError struct {

	// What represents a body of error.
	What string `json:"error"`

	// Arg only for errors in Register method.
	// Meaning of arg depends by error.
	Arg tViewIDEncoded `json:"arg,omitempty"`
}

// arg makes the copy of e, sets the Arg field to arg in a made copy
// and then return that made copy.
func (e *tViewIDConverterError) arg(arg tViewIDEncoded) (copy *tViewIDConverterError) {

	copied := *e
	copy = &copied
	copy.Arg = arg
	return
}

// IsIt returns true if err have the same type and the same value as e.
// Implements iError.
func (e *tViewIDConverterError) IsIt(err error) bool {

	if e == nil || err == nil {
		return false
	}

	e2, ok := err.(*tViewIDConverterError)
	if !ok {
		return false
	}

	return e == e2 || e.What == e2.What && e.Arg == e2.Arg
}

// Error returns a string representing of e.
// Implements iError, error.
func (e *tViewIDConverterError) Error() string {

	if e == nil {
		return ""
	}

	s := e.What
	if e.Arg != cViewIDEncodedNull {
		s += " " + strconv.FormatUint(uint64(e.Arg), 10)
	}

	return s
}
