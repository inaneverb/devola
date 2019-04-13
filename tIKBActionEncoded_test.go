// Copyright © 2018. All rights reserved.
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
	"testing"
	"unsafe"

	testify "github.com/stretchr/testify/assert"
)

// todo: Add another tests

// TestTEventIKBActionEncoded_ValidTEventIKBIDEncoded checks whether
// real size of tViewIDEncoded type is compatible with
// encode/decode algorithm in tIKBActionEncoded type.
func TestTEventIKBActionEncoded_ValidTEventIKBIDEncoded(t *testing.T) {

	ikbid := tViewIDEncoded(cViewIDEncodedNull)
	mustSize := cIKBActionEncodedPosSessionID - cIKBActionEncodedPosViewID

	testMsg := "Sizeof tViewIDEncoded is incompatible with predefined " +
		"cIKBActionEncodedPosViewID and cIKBActionEncodedPosSessionID " +
		"position constants."

	testify.True(t, unsafe.Sizeof(ikbid) == uintptr(mustSize), testMsg)
}

// TestTEventIKBActionEncoded_ValidTSessionID checks whether
// real size of tSessionID type is compatible with
// encode/decode algorithm in tIKBActionEncoded type.
func TestTEventIKBActionEncoded_ValidTSessionID(t *testing.T) {

	ssid := tSessionID(cSessionIDNil)
	mustSize := cIKBActionEncodedPosArgs - cIKBActionEncodedPosSessionID

	testMsg := "Sizeof tSessionID is incompatible with predefined " +
		"cIKBActionEncodedPosSessionID and cIKBActionEncodedPosArgs " +
		"position constants."

	testify.True(t, unsafe.Sizeof(ssid) == uintptr(mustSize), testMsg)
}
