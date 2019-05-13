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

// ID represents a RAW identifier of View.
//
// This type used only for readable format View ID representation.
// In internal SDK parts View ID represents by its encoded format using
// tViewIDEncoded type and tViewIDConverter to encode/decode operations.
//
// More info: tViewIDEncoded, tView, tViewIDConverter, tIKBActionEncoded,
// tCtx, tSender.
type ID string

// Predefined constants.
const (

	// Represents a nil View ID and an indicator of some error.
	cViewIDNull ID = ""
)

// IsValid returns true only if vid is valid ID value.
//
// Valid readable View ID must contain more than 2 chars and don't starts
// from double underscore (reserved for internal parts).
func (id ID) IsValid() bool {

	return vid != cViewIDNull && len(vid) > 2 && vid[:2] != "__"
}
