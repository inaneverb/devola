// Copyright © 2019. All rights reserved.
// Author: Alice Qio.
// Contacts: <qioalice@gmail.com>.
// License: https://opensource.org/licenses/MIT

package chat

import (
	"reflect"
	"unsafe"

	"github.com/qioalice/devola/core/math"
)

// MessageIDs is an alias to the slice of MessageID.
// Represents a dynamic array of backend chat message's IDs.
//
// NOTE.
// All methods changes the receiver object, not its copy!
// (Of course if not stated otherwise in methods' docs).
type MessageIDs []MessageID

// Predefined constants.
const (

	// Default MessageIDs capacity.
	cMessageIDsDefCap int = 10
)

// Len returns the length of s chat message ids' array.
func (s *MessageIDs) Len() int {
	if s == nil {
		return 0
	}
	return len(*s)
}

// SetLen changes the length to the value, depended by n.
// If |n| > len , n will be len * sign(n) and then:
// 1. N >= 0 then N is a new length.
// 2. N < 0 then (length - N) is a new length.
//
// WARNING!
// MEMORY LEAK (HOLDING) POSSIBLE!
// Changing the length don't cause freeing allocated memory by all "unused"
// message id's from now on. If you want it, call FlushLen method after.
func (s *MessageIDs) SetLen(n int) *MessageIDs {

	if s == nil || n < 0 {
		return nil
	}

	header := (*reflect.SliceHeader)(unsafe.Pointer(s))
	switch n := math.ClampI(n, -header.Len, header.Len); {

	case n >= 0 && n < header.Len:
		header.Len = n

	case n < 0 && math.AbsI(n) <= header.Len:
		header.Len += n // +a + -b == a - b
	}

	return s
}

// FlushLen "flushes" the length of current chat message ids' array.
// So, it fixes the situations when capacity is more than real slice's length
// (internal Golang parts).
// All real used memory will be reallocated, all data will be copied
// and new memory will be saved as part of the current chat message ids' object.
// If capacity and real length are equal, there is no-op.
//
// NOTE.
// It's reallocating and copying. This operation may take a time.
//
// WARNING!
// IT IS TOO IMPORTANT TO REASSIGN RECEIVER IN CALLER CODE BY RETURNED VALUE.
func (s *MessageIDs) FlushLen() *MessageIDs {

	header := (*reflect.SliceHeader)(unsafe.Pointer(s))
	if header == nil || header.Len == header.Cap {
		return s
	}

	copied := makeMessageIDs(header.Len)

	// because copy using .Len field, not a .Cap
	copiedHeader := (*reflect.SliceHeader)(unsafe.Pointer(&copied))
	copiedHeader.Len = copiedHeader.Cap

	copy(copied, *s)
	return s
}

// IsEmpty returns true if s chat message ids' array is empty or nil.
func (s *MessageIDs) IsEmpty() bool {
	return s.Len() == 0
}

// Clean is SetLen(0) and PeekAll calls at the same time.
func (s *MessageIDs) Clean() []MessageID {
	if s.IsEmpty() {
		return nil
	}
	returned := (*s)[:]
	*s = *s.SetLen(0)
	return returned
}

// Push pushes id to the end of s and then returns s.
//
// You can use it when s is nil, in reassign context:
// var s = MessageIDs(nil)
// s = s.Push(100) // s now is [100].
func (s *MessageIDs) Push(id MessageID) *MessageIDs {

	if s == nil {
		buf := makeMessageIDs(cMessageIDsDefCap)
		//noinspection GoAssignmentToReceiver
		s = &buf
	}

	*s = append(*s, id)
	return s
}

// Peek returns a last id of s without removing it from s.
// Returns CMessageIDNil if s is empty.
func (s *MessageIDs) Peek() MessageID {
	if s.IsEmpty() {
		return CMessageIDNil
	}
	return (*s)[len(*s)-1]
}

// PeekAll returns a slice of all ids in s without removing them from s.
// Returns CMessageIDNil if s is empty.
func (s *MessageIDs) PeekAll() []MessageID {
	if s.IsEmpty() {
		return nil
	}
	return (*s)[:]
}

// PeekN returns a slice with a count of ids depends by n without removing them from s.
// If |n| > len , n will be len * sign(n) and then:
// 1. N > 0 then n is a how many last ids of s will be returned.
// 2. N < 0 then (len - |n|) ids of s will be returned.
// Returns CMessageIDNil if s is empty or n == 0.
func (s *MessageIDs) PeekN(n int) []MessageID {

	slen := s.Len()

	if slen == 0 || n == 0 {
		return nil
	}

	switch n = math.ClampI(n, -slen, slen); {

	case n > 0:
		return (*s)[slen-n:]

	case n < 0:
		return (*s)[-n:]
	}

	return nil // <-- never, but Go requires
}

// Pop removes a last id of s and returns it.
// Returns CMessageIDNil if s is empty.
func (s *MessageIDs) Pop() MessageID {
	id := s.Peek()
	s.SetLen(-1)
	return id
}

// PopAll removes all ids from s and returns it.
// Returns CMessageIDNil if s is empty.
func (s *MessageIDs) PopAll() []MessageID {
	ids := s.PeekAll()
	s.Clean()
	return ids
}

// PopN removes and returns M ids of s (depends by n).
// If |n| > len , n will be len * sign(n) and then:
// 1. N > 0 then M = n.
// 2. N < 0 then M = (len - |n|).
// Returns CMessageIDNil if s is empty or n == 0.
func (s *MessageIDs) PopN(n int) []MessageID {
	ids := s.PeekN(n)
	s.SetLen(-n)
	return ids
}

// makeMessageIDs creates a new MessageIDs object with a passed value of capacity.
func makeMessageIDs(cap int) MessageIDs {
	return make([]MessageID, 0, cap)
}
