// Copyright © 2019. All rights reserved.
// Author: Alice Qio.
// Contacts: <qioalice@gmail.com>.
// License: https://opensource.org/licenses/MIT

package chat

// MessageID represents some chat message's ID.
type MessageID int

const (

	// CMessageIDNil is an identifier of bad message ID: incorrect, broken, nil,
	// not existed, etc.
	CMessageIDNil MessageID = -1
)

// IsValid returns true only if id is valid message ID value.
func (id MessageID) IsValid() bool {
	return id != CMessageIDNil
}
