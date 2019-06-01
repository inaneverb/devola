// Copyright © 2019. All rights reserved.
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
// More info: ID, IDConv.
type IDEnc uint32

// Predefined constants.
const (

	// Represents a nil encoded View ID and an indicator of some error.
	CIDEncNil IDEnc = 0

	// All encoded identifiers as numbers will be more than or equal to that value.
	// All less values are reserved for internal needs.
	CIDEncStartValue IDEnc = 100

	// All encoded identifiers as numbers will be less than that value.
	// All more than or equal to that values are reserved for internal needs.
	CIDEncMaxValue IDEnc = math.MaxUint32 - 1
)

// IsValid returns true only if id is valid IDEnc value.
//
// Valid encoded View ID must not be equal to the CIDEncNil
// and be in range [CIDEncodedStartValue, CIDEncodedMaxValue).
func (id IDEnc) IsValid() bool {
	return id != CIDEncNil &&
		id >= CIDEncStartValue && id < CIDEncMaxValue
}
