// Copyright © 2018. All rights reserved.
// Author: Alice Qio.
// Contacts: <qioalice@gmail.com>.
// License: https://opensource.org/licenses/MIT

package chat

// SessionID represents ID of some session.
// It's a part of Session object.
//
// A session's (Session) uniqueness is achieved using a combination
// of two parameters: a chat ID (ID) and session ID (SessionID).
// Because of this, uint32 is enough to represent session ID.
type SessionID uint32

// Predefined session ID constants.
const (

	// An identifier of some bad session: nil session, broken session,
	// not existed session, etc.
	// No valid session can have this identifier.
	CSessionIDNil SessionID = 0
)

// IsValid returns true only if current session ID is valid session ID
// and isn't cSessionIDNil or some another bad const (in the future).
func (id SessionID) IsValid() bool {
	return id != cSessionIDNil
}
