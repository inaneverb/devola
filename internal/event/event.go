// Copyright © 2018. All rights reserved.
// Author: Alice Qio.
// Contacts: <qioalice@gmail.com>.
// License: https://opensource.org/licenses/MIT

package event

import (
	"reflect"
	"unsafe"

	"./ikba"
)

// Event represents Telegram Bot Update type.
//
// When user somehow interact with Telegram Bot, it sends an api.Update
// object to the Telegram Bot.
// The Event object creates using api.Update object and answers to the
// following questions:
//
// - What kind of event is occurred?
// - What data has been received with occurred event?
//
// This is a part of ctx.Ctx.
//
// Thus you can always to get info about what event you're handling
// inside your handler.
//
// More info: Type, Data, ikba.Encoded, ctx.Ctx.
type Event struct {

	// TODO: Add arguments supporting (for example for commands)

	// The type of occurred event. Is one of predefined constant.
	Type Type `json:"type"`

	// The occurred event's data. Text, keyboard button text, inline keyboard
	// View ID, etc.
	Data Data `json:"data,omitempty"`

	// Encoded IKB action.
	// It's a pointer to avoid reallocate memory for ikba.Encoded object
	// and it is a pointer to the Data field with casted type.
	//
	// Not nil only if Type == CTypeInlineKeyboardButton.
	ikbae *ikba.Encoded `json:"-"`
}

// MakeEvent creates a new Event object with passed event type and event data,
// but also initializes IKB encoded action pointer if it is IKB event.
func MakeEvent(typ Type, data Data) *Event {

	event := &Event{
		Type: typ,
		Data: data,
	}

	// ikbae field will point to the Data field but with right type
	if typ == CTypeInlineKeyboardButton {
		dataHeader := (*reflect.StringHeader)(unsafe.Pointer(&event.Data))
		event.ikbae = (*ikba.Encoded)(unsafe.Pointer(dataHeader.Data))
	}

	return event
}
