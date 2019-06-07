// Copyright © 2018. All rights reserved.
// Author: Alice Qio.
// Contacts: <qioalice@gmail.com>.
// License: https://opensource.org/licenses/MIT

package dechan

import (
	"unsafe"
)

// Dechan is an internal double-ended thread-safety queue for pointers.
// It's like golang chan but double-ended (pop front, push front)
// and it's like normal queue but with golang chan features.
type Dechan struct {

	// Base type and algos: https://github.com/gammazero/deque .
	// Copyright (c) 2018 Andrew J. Gillis

	// TODO: MAKE THREAD-SAFETY

	buf    []unsafe.Pointer
	head   int16
	tail   int16
	count  int16
	minCap int16
}

// Predefined constants.
const (

	// minCapacity is the smallest capacity that Dechan may have.
	// Must be power of 2 for bitwise modulus: x % n == x & (n - 1).
	minCapacity = 16
)

// IsEmpty returns true if Dechan is empty or nil.
func (ch *Dechan) IsEmpty() bool {
	return ch == nil || ch.count == 0
}

// Len returns the number of elements currently stored in the Dechan.
// Returns 0 if Dechan is nil.
func (ch *Dechan) Len() int16 {
	if ch == nil {
		return 0
	}
	return ch.count
}

// PushBack appends an element to the back of the queue. Implements FIFO when
// elements are removed with PopFront(), and LIFO when elements are removed
// with PopBack().
func (ch *Dechan) PushBack(ptr unsafe.Pointer) {
	ch.growIfFull()

	ch.buf[ch.tail] = ptr
	// Calculate new tail position.
	ch.tail = ch.next(ch.tail)
	ch.count++
}

// PushFront prepends an element to the front of the queue.
func (ch *Dechan) PushFront(ptr unsafe.Pointer) {
	ch.growIfFull()

	// Calculate new head position.
	ch.head = ch.prev(ch.head)
	ch.buf[ch.head] = ptr
	ch.count++
}

// PopFront removes and returns the element from the front of the queue.
// Implements FIFO when used with PushBack(). If the queue is empty, returns nil.
func (ch *Dechan) PopFront() unsafe.Pointer {
	if ch.count <= 0 {
		return nil
	}
	ret := ch.buf[ch.head]
	ch.buf[ch.head] = nil
	// Calculate new head position.
	ch.head = ch.next(ch.head)
	ch.count--

	ch.shrinkIfExcess()
	return ret
}

// PopBack removes and returns the element from the back of the queue.
// Implements LIFO when used with PushBack(). If the queue is empty, returns nil.
func (ch *Dechan) PopBack() unsafe.Pointer {
	if ch.count <= 0 {
		return nil
	}

	// Calculate new tail position
	ch.tail = ch.prev(ch.tail)

	// Remove value at tail.
	ret := ch.buf[ch.tail]
	ch.buf[ch.tail] = nil
	ch.count--

	ch.shrinkIfExcess()
	return ret
}

// Front returns the element at the front of the queue. This is the element
// that would be returned by PopFront(). This call returns nil if the queue is
// empty.
func (ch *Dechan) Front() unsafe.Pointer {
	if ch.count <= 0 {
		return nil
	}
	return ch.buf[ch.head]
}

// Back returns the element at the back of the queue. This is the element
// that would be returned by PopBack(). This call returns nil if the queue is
// empty.
func (ch *Dechan) Back() unsafe.Pointer {
	if ch.count <= 0 {
		return nil
	}
	return ch.buf[ch.prev(ch.tail)]
}

// At returns the element at index i in the queue without removing the element
// from the queue. At(0) refers to the first element and is the same as Front().
// At(Len()-1) refers to the last element and is the same as Back().
// If the index is invalid, returns nil.
//
// The purpose of At is to allow Dechan to serve as a more general purpose
// circular buffer, where items are only added to and removed from the ends of
// the Dechan, but may be read from any place within the Dechan. Consider the
// case of a fixed-size circular log buffer: A new entry is pushed onto one end
// and when full the oldest is popped from the other end. All the log entries
// in the buffer must be readable without altering the buffer contents.
func (ch *Dechan) At(i int16) unsafe.Pointer {
	if i < 0 || i >= ch.count {
		return nil
	}
	// bitwise modulus
	return ch.buf[(ch.head+i)&int16(len(ch.buf)-1)]
}

// Clear removes all elements from the queue, but retains the current capacity.
// This is useful when repeatedly reusing the queue at high frequency to avoid
// GC during reuse. The queue will not be resized smaller as long as items are
// only added. Only when items are removed is the queue subject to getting
// resized smaller.
func (ch *Dechan) Clear() {
	// bitwise modulus
	modBits := int16(len(ch.buf) - 1)
	for h := ch.head; h != ch.tail; h = (h + 1) & modBits {
		ch.buf[h] = nil
	}
	ch.head = 0
	ch.tail = 0
	ch.count = 0
}

// Rotate rotates the Dechan n steps front-to-back. If n is negative, rotates
// back-to-front. Having Dechan provide Rotate() avoids resizing that could
// happen if implementing rotation using only Pop and Push methods.
func (ch *Dechan) Rotate(n int16) {
	if ch.count <= 1 {
		return
	}
	// Rotating a multiple of ch.count is same as no rotation.
	n %= ch.count
	if n == 0 {
		return
	}

	modBits := int16(len(ch.buf) - 1)
	// If no empty space in buffer, only move head and tail indexes.
	if ch.head == ch.tail {
		// Calculate new head and tail using bitwise modulus.
		ch.head = (ch.head + n) & modBits
		ch.tail = (ch.tail + n) & modBits
		return
	}

	if n < 0 {
		// Rotate back to front.
		for ; n < 0; n++ {
			// Calculate new head and tail using bitwise modulus.
			ch.head = (ch.head - 1) & modBits
			ch.tail = (ch.tail - 1) & modBits
			// Put tail value at head and remove value at tail.
			ch.buf[ch.head] = ch.buf[ch.tail]
			ch.buf[ch.tail] = nil
		}
		return
	}

	// Rotate front to back.
	for ; n > 0; n-- {
		// Put head value at tail and remove value at head.
		ch.buf[ch.tail] = ch.buf[ch.head]
		ch.buf[ch.head] = nil
		// Calculate new head and tail using bitwise modulus.
		ch.head = (ch.head + 1) & modBits
		ch.tail = (ch.tail + 1) & modBits
	}
}

// SetMinCapacity sets a minimum capacity of 2^minCap.
// If the value of the minimum capacity is less than or equal
// to the minimum allowed, then capacity is set to the minimum allowed.
// This may be called at anytime to set a new minimum capacity.
//
// Setting a larger minimum capacity may be used to prevent resizing
// when the number of stored items changes frequently across a wide range.
func (ch *Dechan) SetMinCapacity(minCap uint16) {
	if 1<<minCap > minCapacity {
		ch.minCap = 1 << minCap
	} else {
		ch.minCap = minCapacity
	}
}

// prev returns the previous buffer position wrapping around buffer.
func (ch *Dechan) prev(i int16) int16 {
	return (i - 1) & int16(len(ch.buf)-1) // bitwise modulus
}

// next returns the next buffer position wrapping around buffer.
func (ch *Dechan) next(i int16) int16 {
	return (i + 1) & int16(len(ch.buf)-1) // bitwise modulus
}

// growIfFull resizes up if the buffer is full.
func (ch *Dechan) growIfFull() {
	if ch.count == int16(len(ch.buf)) {
		ch.resize()
	}
}

// shrinkIfExcess resize down if the buffer 1/4 full.
func (ch *Dechan) shrinkIfExcess() {
	if l := int16(len(ch.buf)); l > ch.minCap && (ch.count<<2) == l {
		ch.resize()
	}
}

// resize resizes the Dechan to fit exactly twice its current contents.
// This is used to grow the queue when it is full,
// and also to shrink it when it is only a quarter full.
func (ch *Dechan) resize() {
	newBuf := make([]unsafe.Pointer, ch.count<<1)
	if ch.tail > ch.head {
		copy(newBuf, ch.buf[ch.head:ch.tail])
	} else {
		n := copy(newBuf, ch.buf[ch.head:])
		copy(newBuf[n:], ch.buf[:ch.tail])
	}

	ch.head = 0
	ch.tail = ch.count
	ch.buf = newBuf
}

// New creates a new Dechan object with passed capacity value.
// If passed capacity < minCapacity (16 by default), it will be overwritten
// by that value.
func New(cap int16) *Dechan {
	ch := new(Dechan)

	// don't need this check because of:
	// only Sender.consts.chchCap used as cap arg for this constructor
	// => Sender.consts.chchCap default value is minCapacity (sender.go:384),
	// => Sender.consts.chchCap changes only from Params.SetEXREOF and that param
	// has this check (params.go:93:95).
	// if cap < minCapacity {
	// 	cap = minCapacity
	// }

	ch.minCap = cap
	ch.buf = make([]unsafe.Pointer, ch.minCap)

	return ch
}
