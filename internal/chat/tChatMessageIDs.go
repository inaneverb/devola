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
	"reflect"
	"unsafe"
)

// tChatMessageIDs is an alias to slice of tChatMessageID.
// Represents a dynamic array of Telegram Bot Message's IDs.
//
// Used as a way to store all sent messages from Telegram Bot to the some chat,
// as an example.
//
// More info: private/session.tSession.
//
// NOTE! BE CAREFUL!
// All methods changes the receiver object, not its copy!
// (Of course if not stated otherwise in methods' docs).
type tChatMessageIDs []tChatMessageID

// Len returns the length of s chat message ids' array.
// Nilcallable (0 is returned).
func (s *tChatMessageIDs) Len() int {
	if s == nil {
		return 0
	}
	return len(*s)
}

// SetLen changes the length of current chat message ids' array
// to the new value, depended by count.
// If count > 0 and count < current length, the count is new length.
// If count < 0 the new length is current length - |count|
// (the same as "s.SetLen(s.Len()-|count|)" ).
// Otherwise there is no-op.
// Nilcallable (nil is returned).
//
// WARNING!
// Keep in mind! Memory "leak" (holding) possible!
// Changing the length don't cause freeing allocated memory by all
// "unused" message id's from now on.
// If you want it, call FlushLen method after.
func (s *tChatMessageIDs) SetLen(count int) (this *tChatMessageIDs) {
	this = s
	if this == nil || count < 0 {
		return
	}
	header := (*reflect.SliceHeader)(unsafe.Pointer(this))
	header.Len = count
	*this = *(*tChatMessageIDs)(unsafe.Pointer(header))
	return
}

// FlushLen "flushes" the length of current chat message ids' array.
// So, it fixes the situations when capacity is more than real slice's length
// (internal Golang parts).
// All real used memory will be reallocated, all data will be copied
// and new memory will be saved as part of the current chat message ids' object.
// If capacity and real length are equal, there is no-op.
// Nilcallable (nil is returned).
//
// NOTICE! It's reallocating and copying and it may takes a time.
//
// WARNING! Reassign receiver!
// It's too important to reassign receiver in caller code by returned value.
//
func (s *tChatMessageIDs) FlushLen() (this *tChatMessageIDs) {
	this = s
	if this == nil {
		return
	}
	header := (*reflect.SliceHeader)(unsafe.Pointer(this))
	if header.Len == header.Cap {
		return
	}
	copied := makeChatSentMessageIDs(header.Len)
	// because copy using .Len field, not a .Cap
	copiedHeader := (*reflect.SliceHeader)(unsafe.Pointer(copied))
	copiedHeader.Len = copiedHeader.Cap
	copy(*copied, *this)

}

// IsEmpty returns true if s chat message ids' array is empty or nil.
// Otherwise false is returned.
// Nilcallable (true is returned).
func (s *tChatMessageIDs) IsEmpty() bool {
	return s.Len() == 0
}

// Clean is SetLen(0) and PeekAll calls at the same time.
// Returns all message ids that has been stored and then cleans
// the current s chat message ids' array.
func (s *tChatMessageIDs) Clean() (messageIDs []tChatMessageID) {
	this = s
	if s.IsEmpty() {
		return
	}
	*this = *this.SetLen(0)
	return
}

//
func (s *tChatMessageIDs) Push(messageID int) (this *tChatMessageIDs) {
	this = s
	if this == nil {
		this = makeChatSentMessageIDs()
	}
	*this = append(*this, messageID)
	return
}

//
func (s *tChatMessageIDs) Peek() (messageID tChatMessageID) {
	if s.IsEmpty() {
		return cChatSentMessageIDsNoMessages
	}
	messageID = (*s)[len(*s)-1]
	return
}

//
func (s *tChatMessageIDs) PeekAll() (messageIDs []tChatMessageID) {
	if s.IsEmpty() {
		messageIDs = nil
		return
	}
	messageIDs = (*s)[:]
	return
}

//
func (s *tChatMessageIDs) PeekN(count int) (messageIDs []tChatMessageID) {
	slen := s.Len()
	if slen == 0 || count == 0 {
		messageIDs = nil
	}
	count = mathClampI(count, -slen, slen) // -slen <= count <= slen
	if count > 0 {
		messageIDs = (*s)[slen-count:]
	}
	if count < 0 {
		messageIDs = (*s)[-count:]
	}
	return
}

//
func (s *tChatMessageIDs) Pop() (messageID tChatMessageID) {
	messageID = s.Peek()
	s.SetLen(s.Len()-1)
	return
}

//
func (s *tChatMessageIDs) PopAll() (messageIDs []tChatMessageID) {
	messageIDs = s.PeekAll()
	s.Clean()
	return
}

//
func (s *tChatMessageIDs) PopN(count int) (messageIDs []tChatMessageID) {
	slen := s.Len()
	if slen == 0 || count == 0 {
		messageIDs = nil
		return
	}
	count = mathClampI(count, -slen, slen) // -slen <= count <= slen
	header := (*reflect.SliceHeader)(unsafe.Pointer(s))
	if count > 0 {
		messageIDs = (*s)[slen-count:]
		header.Len = slen-count
	}
	if count < 0 {
		messageIDs = (*s)[-count:]
		header.Len = -count
	}
	*s = *(*tChatMessageIDs)(unsafe.Pointer(header))
	return
}

//
func makeChatSentMessageIDs(capacity int) *tChatMessageIDs {

}