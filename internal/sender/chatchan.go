// Copyright © 2018. All rights reserved.
// Author: Alice Qio.
// Contacts: <qioalice@gmail.com>.
// License: https://opensource.org/licenses/MIT

package sender

// chatchan is an internal double-ended thread-safety queue for ToSend pointers.
// It's like golang chan but double-ended (pop front, push front)
// and it's like normal chatchan but with golang chan features.
type chatchan struct {

	// Base type and algos: https://github.com/gammazero/deque .
	// Copyright (c) 2018 Andrew J. Gillis

	// TODO: MAKE THREAD-SAFETY

	buf    []*ToSend
	head   int16
	tail   int16
	count  int16
	minCap int16
}

// Predefined constants.
const (

	// cChChMinCapacity is the smallest capacity that deque may have.
	// Must be power of 2 for bitwise modulus: x % n == x & (n - 1).
	cChChMinCapacity = 16
)

// IsEmpty returns true of chatchat is empty or nil.
func (ch *chatchan) IsEmpty() bool {
	return ch == nil || ch.count == 0
}

// Len returns the number of elements currently stored in the chatchan.
// Returns 0 if chatchat is nil.
func (ch *chatchan) Len() int16 {
	if ch == nil {
		return 0
	}
	return ch.count
}

// PushBack appends an element to the back of the queue. Implements FIFO when
// elements are removed with PopFront(), and LIFO when elements are removed
// with PopBack().
func (ch *chatchan) PushBack(ptr *ToSend) {
	ch.growIfFull()

	ch.buf[ch.tail] = ptr
	// Calculate new tail position.
	ch.tail = ch.next(ch.tail)
	ch.count++
}

// PushFront prepends an element to the front of the queue.
func (ch *chatchan) PushFront(ptr *ToSend) {
	ch.growIfFull()

	// Calculate new head position.
	ch.head = ch.prev(ch.head)
	ch.buf[ch.head] = ptr
	ch.count++
}

// PopFront removes and returns the element from the front of the queue.
// Implements FIFO when used with PushBack(). If the queue is empty, returns nil.
func (ch *chatchan) PopFront() *ToSend {
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
func (ch *chatchan) PopBack() *ToSend {
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
func (ch *chatchan) Front() *ToSend {
	if ch.count <= 0 {
		return nil
	}
	return ch.buf[ch.head]
}

// Back returns the element at the back of the queue. This is the element
// that would be returned by PopBack(). This call returns nil if the queue is
// empty.
func (ch *chatchan) Back() *ToSend {
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
// The purpose of At is to allow chatchan to serve as a more general purpose
// circular buffer, where items are only added to and removed from the ends of
// the chatchan, but may be read from any place within the chatchan. Consider the
// case of a fixed-size circular log buffer: A new entry is pushed onto one end
// and when full the oldest is popped from the other end. All the log entries
// in the buffer must be readable without altering the buffer contents.
func (ch *chatchan) At(i int16) *ToSend {
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
func (ch *chatchan) Clear() {
	// bitwise modulus
	modBits := int16(len(ch.buf) - 1)
	for h := ch.head; h != ch.tail; h = (h + 1) & modBits {
		ch.buf[h] = nil
	}
	ch.head = 0
	ch.tail = 0
	ch.count = 0
}

// Rotate rotates the chatchan n steps front-to-back. If n is negative, rotates
// back-to-front. Having chatchan provide Rotate() avoids resizing that could
// happen if implementing rotation using only Pop and Push methods.
func (ch *chatchan) Rotate(n int16) {
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

// SetcChChMinCapacity sets a minimum capacity of 2^cChChMinCapacityExp.
// If the value of the minimum capacity is less than or equal
// to the minimum allowed, then capacity is set to the minimum allowed.
// This may be called at anytime to set a new minimum capacity.
//
// Setting a larger minimum capacity may be used to prevent resizing
// when the number of stored items changes frequently across a wide range.
func (ch *chatchan) SetcChChMinCapacity(cChChMinCapacityExp uint16) {
	if 1<<cChChMinCapacityExp > cChChMinCapacity {
		ch.minCap = 1 << cChChMinCapacityExp
	} else {
		ch.minCap = cChChMinCapacity
	}
}

// prev returns the previous buffer position wrapping around buffer.
func (ch *chatchan) prev(i int16) int16 {
	return (i - 1) & int16(len(ch.buf)-1) // bitwise modulus
}

// next returns the next buffer position wrapping around buffer.
func (ch *chatchan) next(i int16) int16 {
	return (i + 1) & int16(len(ch.buf)-1) // bitwise modulus
}

// growIfFull resizes up if the buffer is full.
func (ch *chatchan) growIfFull() {
	if len(ch.buf) == 0 {
		if ch.minCap == 0 {
			ch.minCap = cChChMinCapacity
		}
		ch.buf = make([]*ToSend, ch.minCap)
		return
	}
	if ch.count == int16(len(ch.buf)) {
		ch.resize()
	}
}

// shrinkIfExcess resize down if the buffer 1/4 full.
func (ch *chatchan) shrinkIfExcess() {
	if l := int16(len(ch.buf)); l > ch.minCap && (ch.count<<2) == l {
		ch.resize()
	}
}

// resize resizes the chatchan to fit exactly twice its current contents.
// This is used to grow the queue when it is full,
// and also to shrink it when it is only a quarter full.
func (ch *chatchan) resize() {
	newBuf := make([]*ToSend, ch.count<<1)
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
