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

// tEventRegistered is the part of tReceiver, extends tEvent type.
//
// More info: tReceiver.
type tEventRegistered struct {

	// A base type.
	tEvent `json:",inline"`

	// A set of "current" ViewIDs when this registering event should be reacted.
	// If empty, registering event will be reacted anytime, but if it's not,
	// the registering event will be handled only when current session's
	// View ID is the same as any View ID from this field.
	When []tViewID `json:"when,omitempty"`
}

// cp creates a copy of current registering event and returns it.
func (e *tEventRegistered) cp() *tEventRegistered {

	// See https://github.com/go101/go101/wiki (about slice copying below)
	copy := makeEventRegistered(e.Type, e.Data, append(e.When[:0:0], e.When...))
	return &copy
}

// makeEventRegistered creates a new tEventRegistered object in which the
// base class object tEvent will be initialized with passed event type and
// event data and another field When will be initialized by the same name arg.
func makeEventRegistered(typ tEventType, data tEventData, when []tViewID) *tEventRegistered {

	return &tEventRegistered{
		tEvent: *makeEvent(typ, data),
		When:   when,
	}
}
