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

// tEventType represents type of event that is occurred.
//
// This is a part of tEvent object, which is a part of TCtx.
//
// Thus you can always to figure out inside handler which kind of event
// your handler handling. Use predefined constants, presented below for this.
//
// More info: tEvent, TCtx.
type tEventType uint8

// Constants of tEventType.
// Use these constants to figure out what kind of event is occurred
// (by comparing tEventType).
const (
	// Internal type.
	// A marked invalid type.
	cEventTypeInvalid tEventType = 0 + iota

	// Internal type.
	// A not fully determined type,
	// but it is either keyboard button type or text type.
	cEventTypeKeyboardButtonOrText tEventType = 2

	// A chat text command.
	// tEvent's Data field represents a lowercase command without arguments.
	CEventTypeCommand tEventType = 100

	// A pressed keyboard button.
	// Technically this is a text (in a chat), but a text sent by
	// pressing to the keyboard button.
	// tEvent's Data field represents this keyboard button data.
	CEventTypeKeyboardButton tEventType = 101

	// A typed text.
	// tEvent's Data field represents the whole text but with trimmed
	// leading and trailing spaces.
	CEventTypeText tEventType = 102

	// A pressed inline keyboard button.
	// tEvent's Data field stored TAction value, representing the your
	// action causes occured event.
	CEventTypeInlineKeyboardButton tEventType = 200
)

// isValid returns true only if ht is public tEventType predefined constant.
// Otherwise false is returned.
func (et tEventType) isValid() bool {
	return et == CEventTypeCommand ||
		et == CEventTypeKeyboardButton ||
		et == CEventTypeText ||
		et == CEventTypeInlineKeyboardButton
}
