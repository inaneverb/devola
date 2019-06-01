// Copyright © 2019. All rights reserved.
// Author: Alice Qio.
// Contacts: <qioalice@gmail.com>.
// License: https://opensource.org/licenses/MIT

package event

// Type represents type of event that is occurred.
// This is a part of Event object, which is a part of backend Ctx.
//
// Thus you can always to figure out inside handler which kind of event
// your handler handling. Use predefined constants, presented below for this.
type Type uint8
