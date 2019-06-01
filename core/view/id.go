// Copyright © 2019. All rights reserved.
// Author: Alice Qio.
// Contacts: <qioalice@gmail.com>.
// License: https://opensource.org/licenses/MIT

package view

// ID represents a RAW identifier of View.
//
// This type used only for readable format View ID representation.
// In internal SDK parts View ID represents by its encoded format using
// IDEnc type and IDConv to encode/decode operations.
//
// More info: IDEnc, IDConv
type ID string

// Predefined constants.
const (

	// Represents a nil View ID and an indicator of some error.
	CIDNil ID = ""
)

// IsValid returns true only if vid is valid ID value.
//
// Valid readable View ID must contain more than 2 chars and don't starts
// from double underscore (reserved for internal parts).
func (id ID) IsValid() bool {
	return id != CIDNil && len(id) > 2 && id[:2] != "__"
}
