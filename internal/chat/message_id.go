// Copyright © 2018. All rights reserved.
// Author: Alice Qio.
// Contacts: <qioalice@gmail.com>.
// License: https://opensource.org/licenses/MIT

package chat

// MessageID represents some Telegram Bot Message's ID.
// Aliased type depends by type used in Telegram Bot API package.
type MessageID int

const (

	// CMessageIDNil is an identifer of bad message ID: incorrect, broken, nil,
	// not existed, etc.
	CMessageIDNil MessageID = -1
)
