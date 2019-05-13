// Copyright Â© 2018. All rights reserved.
// Author: Alice Qio.
// Contacts: <qioalice@gmail.com>.
// License: https://opensource.org/licenses/MIT
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"),
// to deal in the Software without restriction, including without limitation
// the rights to use, copy, modify, merge, publish, distribute, sublicense,
// and/or sell copies of the Software, and to permit persons to whom
// the Software is furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS
// IN THE SOFTWARE.

package tgbot

import (
	"math"
)

// tViewIDEncoded is the internal type that represents an encoded View ID.
// It is used to more compact way storing an ID of Views.
//
// Also it is a logical part of tIKBActionEncoded:
// Each Telegram inline keyboard button created by this SDK, leads to some
// entity, called a View.
// Technically the inline keyboard button encoded action have this type bytes
// as part of yourself to perform that leading.
//
// More info: tViewID, tView, tIKBActionEncoded, tCtx, tSender.
type tViewIDEncoded uint32

// Predefined constants.
const (

	// Represents a nil encoded View ID and an indicator of some error.
	cViewIDEncodedNull tViewIDEncoded = 0

	// All encoded identifiers as numbers will be more than or equal to that value.
	// All less values are reserved for internal needs.
	cViewIDEncodedStartValue tViewIDEncoded = 100

	// All encoded identifiers as numbers will be less than that value.
	// All more than or equal to that values are reserver for internal needs.
	cViewIDEncodedMaxValue tViewIDEncoded = math.MaxUint32 - 1
)

// isValid returns true only if vide is valie tViewIDEncoded value.
//
// Valid encoded View ID must not be equal to the cViewIDEncodedNull
// and be in range [cViewIDEncodedStartValue, cViewIDEncodedMaxValue).
func (vide tViewIDEncoded) isValid() bool {

	return vide != cViewIDEncodedNull &&
		vide >= cViewIDEncodedStartValue && vide < cViewIDEncodedMaxValue
}
