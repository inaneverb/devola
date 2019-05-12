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

// tEventData represents data (body) of event that is occurred:
// - for Keyboard Button it is the button text
// - fot Inline Keyboard Button it is the View ID (not encoded)
// - for Text it is the text
// - for Command it is the body of command (w/o '/')
// etc.
//
// This is a part of tEvent object, which is a part of TCtx.
//
// Thus you can always to get this inside handler of event your handler
// handling.
//
// More info: tIKBActionEncoded, tViewID, tViewIDConverter, tEvent, TCtx.
type tEventData string

// Constants of tEventData.
// Use these constants to represent special event cases.
const (
	
	// A marker of invalid data.
	cEventDataNull tEventData = ""
)
