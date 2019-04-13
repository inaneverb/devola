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
	"reflect"
	"unsafe"
)

// todo: Add arguments supporting (for example for commands)

// tEvent represents Telegram Bot Update type.
//
// When user somehow interact with Telegram Bot, it sends an api.Update
// object to the Telegram Bot.
// The tEvent object creates using api.Update object and answers to the
// following questions:
//
// - What kind of event is occurred?
// - What data has been received with occurred event?
//
// This is a part of TCtx.
//
// Thus you can always to get info about what event you're handling
// inside your handler.
//
// More info: tEventType, tEventData, tIKBActionEncoded, TCtx.
type tEvent struct {

	// The type of occurred event. Is one of predefined constant.
	Type tEventType `json:"type"`

	// The occurred event's data. Text, keyboard button text, inline keyboard
	// View ID, etc.
	Data tEventData `json:"data"`

	// Encoded IKB action.
	// It's a pointer to avoid reallocate memory for tIKBActionEncoded object
	// and it is a pointert to the Data field with casted type.
	//
	// Not nil only if Type == CEventTypeInlineKeyboardButton.
	ikbae *tIKBActionEncoded
}

// makeEvent creates a new tEvent object with passed event type and event data,
// but also initializes IKB encoded action pointer if it is IKB event.
func makeEvent(typ tEventType, data tEventData) *tEvent {

	event := &tEvent{
		Type: typ,
		Data: data,
	}

	// ikbae field will point to the Data field but with right type
	if typ == CEventTypeInlineKeyboardButton {
		dataHeader := (*reflect.StringHeader)(unsafe.Pointer(&event.Data))
		event.ikbae = (*tIKBActionEncoded)(unsafe.Pointer(dataHeader.Data))
	}

	return event
}
