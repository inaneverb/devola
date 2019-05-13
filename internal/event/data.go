// Copyright © 2018. All rights reserved.
// Author: Alice Qio.
// Contacts: <qioalice@gmail.com>.
// License: https://opensource.org/licenses/MIT

package event

// Data represents body of event that is occurred:
// - for Keyboard Button it is the button text
// - fot Inline Keyboard Button it is the View ID (not encoded)
// - for Text it is the text
// - for Command it is the body of command (w/o '/')
// etc.
//
// This is a part of Event object, which is a part of ctx.Ctx.
//
// Thus you can always to get this inside handler of event your handler
// handling.
//
// More info: ikba.Encoded, view.ID, view.IDConverter, Event, ctx.Ctx.
type Data string

// Constants of Data.
// Use these constants to represent special event cases.
const (

	// A marker of invalid data.
	cDataNil Data = ""
)
