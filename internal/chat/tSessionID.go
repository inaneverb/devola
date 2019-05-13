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

// tSessionID represents ID of some session. It's a part of tSession object.
//
// A session's (tSession's) uniqueness is achieved using a combination
// of two parameters: a chat ID (tChatID) and session ID (tSessionID).
// Because of this, uint32 is enough to represent session ID.
type tSessionID uint32

// Predefined session ID constants.
const (
	// An identifier of some bad session: nil session, broken session,
	// not existed session, etc.
	// No valid session can have this identifier.
	cSessionIDNil tSessionID = 0
)

// isValid returns true only if current session ID is valid session ID
// and isn't cSessionIDNil or some another bad const (in the future).
func (id tSessionID) isValid() bool {
	return id != cSessionIDNil
}