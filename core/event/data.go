// Copyright © 2019. All rights reserved.
// Author: Alice Qio.
// Contacts: <qioalice@gmail.com>.
// License: https://opensource.org/licenses/MIT

package event

// Data represents body of event that is occurred:
// This is a part of Event object, which is a part of ctx.Ctx.
//
// Thus you can always to get this inside handler of event your handler
// handling.
type Data string

// Constants of Data.
// Use these constants to represent special event cases.
const (

	// A marker of invalid data.
	CDataNil Data = ""
)
