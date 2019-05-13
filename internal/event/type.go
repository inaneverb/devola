// Copyright © 2018. All rights reserved.
// Author: Alice Qio.
// Contacts: <qioalice@gmail.com>.
// License: https://opensource.org/licenses/MIT

package event

// Type represents type of event that is occurred.
// This is a part of Event object, which is a part of ctx.Ctx.
//
// Thus you can always to figure out inside handler which kind of event
// your handler handling. Use predefined constants, presented below for this.
//
// More info: Event, ctx.Ctx.
type Type uint8

// Constants of Type.
// Use these constants to figure out what kind of event is occurred
// (by comparing Type).
const (

	// Marker of invalid type.
	CTypeInvalid Type = 0 + iota

	// Chat text command.
	// tEvent's Data field represents a lowercase command without arguments.
	CTypeCommand Type = 100

	// Pressed keyboard button.
	// Technically this is a text (in a chat), but a text sent by
	// pressing to the keyboard button.
	// tEvent's Data field represents this keyboard button data.
	CTypeKeyboardButton Type = 101

	// Typed text.
	// tEvent's Data field represents the whole text but with trimmed
	// leading and trailing spaces.
	CTypeText Type = 102

	// Pressed inline keyboard button.
	// tEvent's Data field stored TAction value, representing the your
	// action causes occured event.
	CTypeInlineKeyboardButton Type = 200
)

// IsValid returns true only if ht is public Type predefined constant.
// Otherwise false is returned.
func (t Type) IsValid() bool {
	return t == CTypeCommand || t == CTypeText ||
		t == CTypeKeyboardButton || t == CTypeInlineKeyboardButton
}
