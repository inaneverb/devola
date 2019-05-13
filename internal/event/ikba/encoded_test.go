// Copyright © 2018. All rights reserved.
// Author: Alice Qio.
// Contacts: <qioalice@gmail.com>.
// License: https://opensource.org/licenses/MIT

package ikba

import (
	"testing"
	"unsafe"

	"../../chat"
	"../../view"

	"github.com/stretchr/testify/assert"
)

//
func TestConsts(t *testing.T) {

	// Checking whether real size of view.IDEnc type is compatible with
	// encode/decode algorithm in Encoded type.

	idenc := view.IDEnc(view.CIDEncNil)
	mustSize := cPosSessionID - cPosViewID

	assert.True(t,
		unsafe.Sizeof(idenc) == uintptr(mustSize),
		"Incompatible size of view.IDEnc type and this package's constants.",
	)

	// Checking whether real size of chat.SessionID type is compatible with
	// encode/decode algorithm in Encoded type.

	ssid := chat.SessionID(chat.CSessionIDNil)
	mustSize = cPosArgs - cPosSessionID

	assert.True(t,
		unsafe.Sizeof(ssid) == uintptr(mustSize),
		"Incompatible size of chat.SessionID type and this package's constants.",
	)
}
