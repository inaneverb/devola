// Copyright Â© 2018. All rights reserved.
// Author: Alice Qio.
// Contacts: <qioalice@gmail.com>.
// License: https://opensource.org/licenses/MIT

package view

import (
	"math"
)

// IDEnc is the internal type that represents an encoded View ID.
// It is used to more compact way storing an ID of Views.
//
// Also it is a logical part of tIKBActionEncoded:
// Each Telegram inline keyboard button created by this SDK, leads to some
// entity, called a View.
// Technically the inline keyboard button encoded action have this type bytes
// as part of yourself to perform that leading.
//
// More info: ID, tView, tIKBActionEncoded, tCtx, tSender.
type IDEnc uint32

// Predefined constants.
const (

	// Represents a nil encoded View ID and an indicator of some error.
	CIDEncNil IDEnc = 0

	// All encoded identifiers as numbers will be more than or equal to that value.
	// All less values are reserved for internal needs.
	CIDEncStartValue IDEnc = 100

	// All encoded identifiers as numbers will be less than that value.
	// All more than or equal to that values are reserver for internal needs.
	CIDEncMaxValue IDEnc = math.MaxUint32 - 1
)

// IsValid returns true only if vide is valie IDEnc value.
//
// Valid encoded View ID must not be equal to the CIDEncNil
// and be in range [CIDEncodedStartValue, CIDEncodedMaxValue).
func (idenc IDEnc) IsValid() bool {

	return idenc != CIDEncNil &&
		idenc >= CIDEncStartValue && idenc < CIDEncMaxValue
}
