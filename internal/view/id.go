// Copyright Â© 2018. All rights reserved.
// Author: Alice Qio.
// Contacts: <qioalice@gmail.com>.
// License: https://opensource.org/licenses/MIT

package tgbot

// ID represents a RAW identifier of View.
//
// This type used only for readable format View ID representation.
// In internal SDK parts View ID represents by its encoded format using
// tViewIDEncoded type and tViewIDConverter to encode/decode operations.
//
// More info: tViewIDEncoded, tView, tViewIDConverter, tIKBActionEncoded,
// tCtx, tSender.
type ID string

// Predefined constants.
const (

	// Represents a nil View ID and an indicator of some error.
	cViewIDNull ID = ""
)

// IsValid returns true only if vid is valid ID value.
//
// Valid readable View ID must contain more than 2 chars and don't starts
// from double underscore (reserved for internal parts).
func (id ID) IsValid() bool {

	return vid != cViewIDNull && len(vid) > 2 && vid[:2] != "__"
}
