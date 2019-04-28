// Copyright © 2018. All rights reserved.
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

// tErrors is the type of storage of all SDK errors grouped by its classes.
// More info: Errors.
type tErrors struct {

	// Errors that may be occurred while convertng from readable view id
	// to the internal view id type and vice-versa.
	// Read more: tViewIDConverter.
	ViewIDConverter tViewIDConverterErrors
}

// Errors is a storage of all SDK errors grouped by SDK classes.
//
// All these errors are instances of a special its types, and implements
// Golang error iface and iError.
//
// You can read more about what these errors represents from doc to their
// storage types (all errors of some type has its type's errors storage).
// Read these field's types docs.
var Errors tErrors

// By expr above memory already allocated for tErrors struct and all nested.
// Errors' fields (nested structs) initialized by its init functions.
