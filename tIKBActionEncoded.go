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
	"unsafe"
)

// todo: Add negative index GetArg support.
// todo: Add named arguments.

// tIKBActionEncoded is the internal type that represents encoded
// a Telegram Inline Keyboard Button (IKB) action.
//
// Action must have an ID of action, a SSID this button links with and
// may contains arguments.
//
// At this moment (April, 2019) Telegram API (v4.1) allow represent
// action of Inline Keyboard Button as some string that must have length
// not more than 64 byte.
//
// The encode/decode algorithm described below.
//
// Encoded view of tIKBActionEncoded:
//
// < Action ID : sizeof(tIKBActionEncoded) (now 4 byte) >
// < Session ID : sizeof(tSessionID) (now 4 byte) >
// < Args count : 1 byte >
// < Index over last encoded argument : 1 byte >
// < Arg 1 Type : 1 byte >
// < Arg 1 Value : N bytes (depends by Arg 1 Type) > ...
//
// ATTENTION!
// DO NOT FORGET CALL init METHOD OF EACH NEW tIKBActionEncoded INSTANCE!
type tIKBActionEncoded [64]byte

// Predefined position constants that helps to perform encode/decode operations.
const (

	// Position in tIKBActionEncoded where View ID starts from.
	cIKBActionEncodedPosViewID byte = 0

	// Position in tIKBActionEncoded where Session ID starts from.
	cIKBActionEncodedPosSessionID byte = 4

	// Position in tIKBActionEncoded where encoded arguments' part starts from.
	cIKBActionEncodedPosArgs byte = 8

	// Position in tIKBActionEncoded where encoded arguments' counter is.
	cIKBActionEncodedPosArgsCount = cIKBActionEncodedPosArgs + 0

	// Position in tIKBActionEncoded where saved the position
	// starts from a next encoded argument can be placed.
	cIKBActionEncodedPosArgsFree = cIKBActionEncodedPosArgs + 1

	// Position in tIKBActionEncoded where encoded arguments' data starts from.
	cIKBActionEncodedPosArgsContent = cIKBActionEncodedPosArgs + 2

	// Max allowable position in tIKBActionEncoded.
	cIKBActionEncodedPosMax byte = 63

	// Error position value.
	// Returned from some methods.
	cIKBActionEncodedPosERROR byte = ^0
)

// Predefined index constants that means a special cases of encode/decode operations.
const (

	// Bad argument's index.
	// Each encoded argument in tIKBActionEncoded has its own index.
	// Some methods should return that index.
	// If any error occurred, this value indicates it.
	cIKBActionEncodedBadIndex int = -1
)

// Predefined argument's type constants that helps represents (encode/decode) arguments.
//
// ATTENTION!
// DO NOT FORGET ADD ALL NEW CONSTANTS TO THE argNextFromPos METHOD'S SWITCH!
// DO NOT FORGET ADD ALL NEW CONSTANTS TO THE argType2S METHOD'S SWITCH!
// DO NOT FORGET ADD BEHAVIOUR FOR NEW TYPES TO THE dump METHOD!
//
// ATTENTION!
// DO NOT OVERFLOW INT8 (1<<7) -1 (127). BECAUSE!
const (

	// Header of int8 argument
	cIKBActionEncodedArgTypeInt8 byte = 10

	// Header of int16 argument
	cIKBActionEncodedArgTypeInt16 byte = 11

	// Header of int32 argument
	cIKBActionEncodedArgTypeInt32 byte = 12

	// Header of int64 argument
	cIKBActionEncodedArgTypeInt64 byte = 13

	// Header of uint8 argument
	cIKBActionEncodedArgTypeUint8 byte = 14

	// Header of uint16 argument
	cIKBActionEncodedArgTypeUint16 byte = 15

	// Header of uint32 argument
	cIKBActionEncodedArgTypeUint32 byte = 16

	// Header of uint64 argument
	cIKBActionEncodedArgTypeUint64 byte = 17

	// Header of float32 argument
	cIKBActionEncodedArgTypeFloat32 byte = 18

	// Header of float64 argument
	cIKBActionEncodedArgTypeFloat64 byte = 19

	// Header of string argument
	cIKBActionEncodedArgTypeString byte = 20
)

// ext1byte extracts 1 byte from encoded IKB action d starts from startPos
// and returns it.
func (d *tIKBActionEncoded) ext1byte(startPos byte) (v int8) {

	return int8(d[startPos])
}

// ext2bytes extracts 2 bytes from encoded IKB action d starts from startPos
// and returns it.
func (d *tIKBActionEncoded) ext2bytes(startPos byte) (v int16) {

	v |= int16(d[startPos+0] << 0)
	v |= int16(d[startPos+1] << 8)
	return v
}

// ext4bytes extracts 4 bytes from encoded IKB action d starts from startPos
// and returns it.
func (d *tIKBActionEncoded) ext4bytes(startPos byte) (v int32) {

	v |= int32(d[startPos+0] << 0)
	v |= int32(d[startPos+1] << 8)
	v |= int32(d[startPos+2] << 16)
	v |= int32(d[startPos+3] << 24)
	return v
}

// ext8bytes extracts 8 bytes from encoded IKB action d starts from startPos
// and returns it.
func (d *tIKBActionEncoded) ext8bytes(startPos byte) (v int64) {

	v |= int64(d[startPos+0] << 0)
	v |= int64(d[startPos+1] << 8)
	v |= int64(d[startPos+2] << 16)
	v |= int64(d[startPos+3] << 24)
	v |= int64(d[startPos+4] << 32)
	v |= int64(d[startPos+5] << 40)
	v |= int64(d[startPos+6] << 48)
	v |= int64(d[startPos+7] << 56)
	return v
}

// extNbytes extracts N bytes from encoded IKB action d starts from startPos
// and returns it.
func (d *tIKBActionEncoded) extNbytes(startPos, bytes byte) (v []byte) {

	v = make([]byte, bytes)
	for i := byte(0); i < bytes; i++ {
		v[i] = d[startPos+i]
	}
	return v
}

// put1byte puts 1 byte v to the encoded IKB action d starts from startPos.
func (d *tIKBActionEncoded) put1byte(startPos byte, v int8) {

	d[startPos] = byte(v)
}

// put2bytes puts 2 bytes v to the encoded IKB action d starts from startPos.
func (d *tIKBActionEncoded) put2bytes(startPos byte, v int16) {

	d[startPos+0] = byte(v >> 0)
	d[startPos+1] = byte(v >> 8)
}

// put4bytes puts 4 bytes v to the encoded IKB action d starts from startPos.
func (d *tIKBActionEncoded) put4bytes(startPos byte, v int32) {

	d[startPos+0] = byte(v >> 0)
	d[startPos+1] = byte(v >> 8)
	d[startPos+2] = byte(v >> 16)
	d[startPos+3] = byte(v >> 24)
}

// put8bytes puts 8 bytes v to the encoded IKB action d starts from startPos.
func (d *tIKBActionEncoded) put8bytes(startPos byte, v int64) {

	d[startPos+0] = byte(v >> 0)
	d[startPos+1] = byte(v >> 8)
	d[startPos+2] = byte(v >> 16)
	d[startPos+3] = byte(v >> 24)
	d[startPos+4] = byte(v >> 32)
	d[startPos+5] = byte(v >> 40)
	d[startPos+6] = byte(v >> 48)
	d[startPos+7] = byte(v >> 56)
}

// putNbytes puts N bytes v to the encoded IKB action d starts from startPos.
func (d *tIKBActionEncoded) putNbytes(startPos byte, v []byte) {

	for i, n := byte(0), byte(len(v)); i < n; i++ {
		d[startPos+i] = v[i]
	}
}

// PutViewID puts the encoded IKB action ID to the current IKB action d.
func (d *tIKBActionEncoded) PutViewID(id tViewIDEncoded) {

	d.put4bytes(cIKBActionEncodedPosViewID, int32(id))
}

// GetViewID extracts the encoded IKB action ID from the current IKB action d
// and returns it.
func (d *tIKBActionEncoded) GetViewID() (id tViewIDEncoded) {

	id = tViewIDEncoded(d.ext4bytes(cIKBActionEncodedPosViewID))
	return
}

// PutSessionID puts the session ID IKB linked with to the current IKB action d.
func (d *tIKBActionEncoded) PutSessionID(ssid tSessionID) {

	d.put4bytes(cIKBActionEncodedPosSessionID, int32(ssid))
	return
}

// GetSessionID extract the encoded IKB session ID from the current IKB action d
// and returns it.
func (d *tIKBActionEncoded) GetSessionID() (ssid tSessionID) {

	ssid = tSessionID(d.ext4bytes(cIKBActionEncodedPosSessionID))
	return
}

// needForType returns the number of bytes that required to store
// an argument's value with type argType.
//
// WARNING! Only for any integer or float type.
// Calls with other type's constant will return a very big value.
func (*tIKBActionEncoded) argNeedForType(argType byte) (numBytes byte) {

	switch argType {

	case cIKBActionEncodedArgTypeInt8,
		cIKBActionEncodedArgTypeUint8:
		return 2

	case cIKBActionEncodedArgTypeInt16,
		cIKBActionEncodedArgTypeUint16:
		return 3

	case cIKBActionEncodedArgTypeInt32,
		cIKBActionEncodedArgTypeUint32,
		cIKBActionEncodedArgTypeFloat32:
		return 5

	case cIKBActionEncodedArgTypeInt64,
		cIKBActionEncodedArgTypeUint64,
		cIKBActionEncodedArgTypeFloat64:
		return 9

	default:
		return cIKBActionEncodedPosMax
	}
}

// argHaveFreeBytes returns true only if numBytes bytes of some argument
// can be saved into current encoded action. Otherwise false is returned.
func (d *tIKBActionEncoded) argHaveFreeBytes(numBytes byte) bool {

	return d[cIKBActionEncodedPosArgsFree]+numBytes <=
		cIKBActionEncodedPosMax
}

// argReserveForType reserves the number of bytes for argument with type argType
// (increases an internal free index position counter), and returns
// a position you can write bytes starting from which.
//
// WARNING! If argument with required type can't be stored (no more space),
// cIKBActionEncodedPosERROR is returned!
//
// WARNING! Only for any integer or float type.
// Otherwise cIKBActionEncodedPosERROR is returned!
func (d *tIKBActionEncoded) argReserveForType(argType byte) (startPos byte) {

	// Check whether d has as many free bytes as argType is required.
	requiredBytes := d.argNeedForType(argType)
	if requiredBytes >= cIKBActionEncodedPosMax {
		return cIKBActionEncodedPosERROR
	}

	// Extract current start pos and next if it will be correct
	startPos = d[cIKBActionEncodedPosArgsFree]
	nextStartPos := startPos + requiredBytes

	// Check whether nextStartPos <= max allowable position
	if nextStartPos > cIKBActionEncodedPosMax {
		return cIKBActionEncodedPosERROR
	}

	// Save arg type, inc start pos counter
	d[startPos] = argType
	d[cIKBActionEncodedPosArgsFree] = nextStartPos
	return startPos + 1
}

// argGet returns a position where argument's content with type argType
// starts from. The search begins from idx argument index.
//
// If index is too long, argument not exists or something wrong else,
// cIKBActionEncodedPosERROR is returned.
//
// Example:
//
// d has encoded arguments in that order:
// 0:int32, 1:int16, 2:string, 3:int8, 4:int8
//
// Calls:
//
// argGet(0, int32) == pos of content of 0 arg.
// argGet(0, string) == pos of content of 2 arg.
// argGet(4, int8) == pos of content of 4 arg.
//
// All types presented above are constants, of course.
func (d *tIKBActionEncoded) argGet(argIdx int, argType byte) (startPos byte) {

	// Check index is valid
	if d.ArgCount() <= argIdx {
		return cIKBActionEncodedPosERROR
	}

	// Skip unnecessary arguments (argIdx -1)
	startPos = cIKBActionEncodedPosArgsContent
	for argIdx--; argIdx > 0; argIdx-- {
		startPos = d.argNextFromPos(startPos)
	}

	// Try to find required argument
	nextFreeIndex := d[cIKBActionEncodedPosArgsFree]
	for startPos != cIKBActionEncodedPosERROR && startPos < nextFreeIndex {
		if d[startPos] == argType {
			// Found, return argument's content position
			return startPos + 1
		}
		// Go to next arg
		startPos = d.argNextFromPos(startPos)
	}

	// Not found
	return cIKBActionEncodedPosERROR
}

// argNextFromPos returns the next argument's position in d if pos is
// position of some argument.
//
// WARNING! Do not have any check!
func (d *tIKBActionEncoded) argNextFromPos(pos byte) (nextArgPos byte) {

	switch d[pos] {

	case cIKBActionEncodedArgTypeInt8,
		cIKBActionEncodedArgTypeUint8:
		return pos + 2

	case cIKBActionEncodedArgTypeInt16,
		cIKBActionEncodedArgTypeUint16:
		return pos + 3

	case cIKBActionEncodedArgTypeInt32,
		cIKBActionEncodedArgTypeUint32,
		cIKBActionEncodedArgTypeFloat32:
		return pos + 5

	case cIKBActionEncodedArgTypeInt64,
		cIKBActionEncodedArgTypeUint64,
		cIKBActionEncodedArgTypeFloat64:
		return pos + 9

	case cIKBActionEncodedArgTypeString:
		// d[pos] - arg type string, d[pos+1] - len of string
		return pos + 2 + d[pos+1]

	default:
		// THIS IS ERROR SWITCH BRANCH!
		// DO NOT "PUT" ANY CASES TO THIS BRANCH!
		//
		// it should never happen, but if it will, let it be safe
		// not just pos, because it may cause infinity loop in caller's
		// not pos + too big C, because it may cause seg fault
		return cIKBActionEncodedPosERROR
	}
}

// argType2S returns a string name of type argType.
func (d *tIKBActionEncoded) argType2S(argType byte) string {

	switch argType {

	case cIKBActionEncodedArgTypeInt8:
		return "int8"

	case cIKBActionEncodedArgTypeInt16:
		return "int16"

	case cIKBActionEncodedArgTypeInt32:
		return "int32"

	case cIKBActionEncodedArgTypeInt64:
		return "int64"

	case cIKBActionEncodedArgTypeUint8:
		return "uint8"

	case cIKBActionEncodedArgTypeUint16:
		return "uint16"

	case cIKBActionEncodedArgTypeUint32:
		return "uint32"

	case cIKBActionEncodedArgTypeUint64:
		return "uint64"

	case cIKBActionEncodedArgTypeFloat32:
		return "float32"

	case cIKBActionEncodedArgTypeFloat64:
		return "float64"

	case cIKBActionEncodedArgTypeString:
		return "string"

	default:
		return "UNKNOWN"
	}
}

// ArgCount returns the number of stored arguments in encoded IKB action d.
func (d *tIKBActionEncoded) ArgCount() (num int) {

	return int(d[cIKBActionEncodedPosArgsCount])
}

// argCountIncPostfix increases the number of stored arguments in encoded
// IKB action d and returns the value before increasing.
//
// (postfix increase operator that Golang don't have).
func (d *tIKBActionEncoded) argCountIncPostfix() (oldValue int) {

	oldValue = int(d[cIKBActionEncodedPosArgsCount])
	d[cIKBActionEncodedPosArgsCount]++
	return oldValue
}

// PutArgInt puts int argument v to the encoded IKB action d.
// If it was successfully, returns the index of that argument.
// Otherwise -1 is returned (argument has not been added).
func (d *tIKBActionEncoded) PutArgInt(v int) (argIdx int) {

	return d.PutArgInt32(int32(v))
}

// PutArgInt8 puts int8 argument v to the encoded IKB action d.
//
// If it was successfully, returns the index of that argument.
// Otherwise -1 is returned (argument has not been added).
func (d *tIKBActionEncoded) PutArgInt8(v int8) (argIdx int) {

	startPos := d.argReserveForType(cIKBActionEncodedArgTypeInt8)
	if startPos == cIKBActionEncodedPosERROR {
		return cIKBActionEncodedBadIndex
	}

	d.put1byte(startPos, v)
	return d.argCountIncPostfix()
}

// PutArgInt16 puts int16 argument v to the encoded IKB action d.
//
// If it was successfully, returns the index of that argument.
// Otherwise -1 is returned (argument has not been added).
func (d *tIKBActionEncoded) PutArgInt16(v int16) (argIdx int) {

	startPos := d.argReserveForType(cIKBActionEncodedArgTypeInt16)
	if startPos == cIKBActionEncodedPosERROR {
		return cIKBActionEncodedBadIndex
	}

	d.put2bytes(startPos, v)
	return d.argCountIncPostfix()
}

// PutArgInt32 puts int32 argument v to the encoded IKB action d.
//
// If it was successfully, returns the index of that argument.
// Otherwise -1 is returned (argument has not been added).
func (d *tIKBActionEncoded) PutArgInt32(v int32) (argIdx int) {

	startPos := d.argReserveForType(cIKBActionEncodedArgTypeInt32)
	if startPos == cIKBActionEncodedPosERROR {
		return cIKBActionEncodedBadIndex
	}

	d.put4bytes(startPos, v)
	return d.argCountIncPostfix()
}

// PutArgInt64 puts int64 argument v to the encoded IKB action d.
//
// If it was successfully, returns the index of that argument.
// Otherwise -1 is returned (argument has not been added).
func (d *tIKBActionEncoded) PutArgInt64(v int64) (argIdx int) {

	startPos := d.argReserveForType(cIKBActionEncodedArgTypeInt64)
	if startPos == cIKBActionEncodedPosERROR {
		return cIKBActionEncodedBadIndex
	}

	d.put8bytes(startPos, v)
	return d.argCountIncPostfix()
}

// PutArgUint puts uint argument v to the encoded IKB action d.
//
// If it was successfully, returns the index of that argument.
// Otherwise -1 is returned (argument has not been added).
func (d *tIKBActionEncoded) PutArgUint(v uint) (argIdx int) {

	return d.PutArgUint32(uint32(v))
}

// PutArgUint8 puts uint8 argument v to the encoded IKB action d.
//
// If it was successfully, returns the index of that argument.
// Otherwise -1 is returned (argument has not been added).
func (d *tIKBActionEncoded) PutArgUint8(v uint8) (argIdx int) {

	startPos := d.argReserveForType(cIKBActionEncodedArgTypeUint8)
	if startPos == cIKBActionEncodedPosERROR {
		return cIKBActionEncodedBadIndex
	}

	d.put1byte(startPos, int8(v))
	return d.argCountIncPostfix()
}

// PutArgUint16 puts uint16 argument v to the encoded IKB action d.
//
// If it was successfully, returns the index of that argument.
// Otherwise -1 is returned (argument has not been added).
func (d *tIKBActionEncoded) PutArgUint16(v uint16) (argIdx int) {

	startPos := d.argReserveForType(cIKBActionEncodedArgTypeUint16)
	if startPos == cIKBActionEncodedPosERROR {
		return cIKBActionEncodedBadIndex
	}

	d.put2bytes(startPos, int16(v))
	return d.argCountIncPostfix()
}

// PutArgUint32 puts uint32 argument v to the encoded IKB action d.
//
// If it was successfully, returns the index of that argument.
// Otherwise -1 is returned (argument has not been added).
func (d *tIKBActionEncoded) PutArgUint32(v uint32) (argIdx int) {

	startPos := d.argReserveForType(cIKBActionEncodedArgTypeUint32)
	if startPos == cIKBActionEncodedPosERROR {
		return cIKBActionEncodedBadIndex
	}

	d.put4bytes(startPos, int32(v))
	return d.argCountIncPostfix()
}

// PutArgUint64 puts uint64 argument v to the encoded IKB action d.
//
// If it was successfully, returns the index of that argument.
// Otherwise -1 is returned (argument has not been added).
func (d *tIKBActionEncoded) PutArgUint64(v uint64) (argIdx int) {

	startPos := d.argReserveForType(cIKBActionEncodedArgTypeUint64)
	if startPos == cIKBActionEncodedPosERROR {
		return cIKBActionEncodedBadIndex
	}
	d.put8bytes(startPos, int64(v))

	return d.argCountIncPostfix()
}

// PutArgFloat32 puts float32 argument v to the encoded IKB action d.
//
// If it was successfully, returns the index of that argument.
// Otherwise -1 is returned (argument has not been added).
func (d *tIKBActionEncoded) PutArgFloat32(v float32) (argIdx int) {

	startPos := d.argReserveForType(cIKBActionEncodedArgTypeFloat32)
	if startPos == cIKBActionEncodedPosERROR {
		return cIKBActionEncodedBadIndex
	}

	d.put4bytes(startPos, *(*int32)(unsafe.Pointer(&v)))
	return d.argCountIncPostfix()
}

// PutArgFloat64 puts float64 argument v to the encoded IKB action d.
//
// If it was successfully, returns the index of that argument.
// Otherwise -1 is returned (argument has not been added).
func (d *tIKBActionEncoded) PutArgFloat64(v float64) (argIdx int) {

	startPos := d.argReserveForType(cIKBActionEncodedArgTypeFloat64)
	if startPos == cIKBActionEncodedPosERROR {
		return cIKBActionEncodedBadIndex
	}

	d.put8bytes(startPos, *(*int64)(unsafe.Pointer(&v)))
	return d.argCountIncPostfix()
}

// PutArgString puts string argument v to the encoded IKB action d.
//
// If it was successfully, returns the index of that argument.
// Otherwise -1 is returned (argument has not been added).
func (d *tIKBActionEncoded) PutArgString(v string) (argIdx int) {

	// String encoding: Arg Type byte, string len byte, string content
	strlen := byte(len(v))
	if !d.argHaveFreeBytes(2 + strlen) {
		return cIKBActionEncodedBadIndex
	}

	// Get start pos, update free index for next argument
	startPos := d[cIKBActionEncodedPosArgsFree]
	d[cIKBActionEncodedPosArgsFree] += strlen + 2

	// Save arg type, save string len
	d[startPos+0] = cIKBActionEncodedArgTypeString
	d[startPos+1] = strlen

	// Save string content
	d.putNbytes(startPos+2, []byte(v))
	return d.argCountIncPostfix()
}

// GetArgInt extracts int argument from encoded IKB action d,
// starting search from startIdx argument's index.
//
// Returns it and true as success if it is, or zero value and false if error.
func (d *tIKBActionEncoded) GetArgInt(startIdx int) (v int, success bool) {

	var v_ int32
	v_, success = d.GetArgInt32(startIdx)
	return int(v_), success
}

// GetArgInt8 extracts int8 argument from encoded IKB action d,
// starting search from startIdx argument's index.
//
// Returns it and true as success if it is, or zero value and false if error.
func (d *tIKBActionEncoded) GetArgInt8(startIdx int) (v int8, success bool) {

	startPos := d.argGet(startIdx, cIKBActionEncodedArgTypeInt8)
	if startPos == cIKBActionEncodedPosERROR {
		return 0, false
	}
	return d.ext1byte(startPos), true
}

// GetArgInt16 extracts int16 argument from encoded IKB action d,
// starting search from startIdx argument's index.
//
// Returns it and true as success if it is, or zero value and false if error.
func (d *tIKBActionEncoded) GetArgInt16(startIdx int) (v int16, success bool) {

	startPos := d.argGet(startIdx, cIKBActionEncodedArgTypeInt16)
	if startPos == cIKBActionEncodedPosERROR {
		return 0, false
	}
	return d.ext2bytes(startPos), true
}

// GetArgInt32 extracts int32 argument from encoded IKB action d,
// starting search from startIdx argument's index.
//
// Returns it and true as success if it is, or zero value and false if error.
func (d *tIKBActionEncoded) GetArgInt32(startIdx int) (v int32, success bool) {

	startPos := d.argGet(startIdx, cIKBActionEncodedArgTypeInt32)
	if startPos == cIKBActionEncodedPosERROR {
		return 0, false
	}
	return d.ext4bytes(startPos), true
}

// GetArgInt64 extracts int64 argument from encoded IKB action d,
// starting search from startIdx argument's index.
//
// Returns it and true as success if it is, or zero value and false if error.
func (d *tIKBActionEncoded) GetArgInt64(startIdx int) (v int64, success bool) {

	startPos := d.argGet(startIdx, cIKBActionEncodedArgTypeInt64)
	if startPos == cIKBActionEncodedPosERROR {
		return 0, false
	}
	return d.ext8bytes(startPos), true
}

// GetArgUint extracts uint argument from encoded IKB action d,
// starting search from startIdx argument's index.
//
// Returns it and true as success if it is, or zero value and false if error.
func (d *tIKBActionEncoded) GetArgUint(startIdx int) (v uint, success bool) {

	var v_ uint32
	v_, success = d.GetArgUint32(startIdx)
	return uint(v_), success
}

// GetArgUint8 extracts uint8 argument from encoded IKB action d,
// starting search from startIdx argument's index.
//
// Returns it and true as success if it is, or zero value and false if error.
func (d *tIKBActionEncoded) GetArgUint8(startIdx int) (v uint8, success bool) {

	startPos := d.argGet(startIdx, cIKBActionEncodedArgTypeUint8)
	if startPos == cIKBActionEncodedPosERROR {
		return 0, false
	}
	v = uint8(d.ext1byte(startPos))
	return v, true
}

// GetArgUint16 extracts uint16 argument from encoded IKB action d,
// starting search from startIdx argument's index.
//
// Returns it and true as success if it is, or zero value and false if error.
func (d *tIKBActionEncoded) GetArgUint16(startIdx int) (v uint16, success bool) {

	startPos := d.argGet(startIdx, cIKBActionEncodedArgTypeUint16)
	if startPos == cIKBActionEncodedPosERROR {
		return 0, false
	}
	v = uint16(d.ext2bytes(startPos))
	return v, true
}

// GetArgUint32 extracts uint32 argument from encoded IKB action d,
// starting search from startIdx argument's index.
//
// Returns it and true as success if it is, or zero value and false if error.
func (d *tIKBActionEncoded) GetArgUint32(startIdx int) (v uint32, success bool) {

	startPos := d.argGet(startIdx, cIKBActionEncodedArgTypeUint32)
	if startPos == cIKBActionEncodedPosERROR {
		return 0, false
	}
	v = uint32(d.ext4bytes(startPos))
	return v, true
}

// GetArgUint64 extracts uint64 argument from encoded IKB action d,
// starting search from startIdx argument's index.
//
// Returns it and true as success if it is, or zero value and false if error.
func (d *tIKBActionEncoded) GetArgUint64(startIdx int) (v uint64, success bool) {

	startPos := d.argGet(startIdx, cIKBActionEncodedArgTypeUint64)
	if startPos == cIKBActionEncodedPosERROR {
		return 0, false
	}
	v = uint64(d.ext8bytes(startPos))
	return v, true
}

// GetArgFloat32 extracts float32 argument from encoded IKB action d,
// starting search from startIdx argument's index.
//
// Returns it and true as success if it is, or zero value and false if error.
func (d *tIKBActionEncoded) GetArgFloat32(startIdx int) (v float32, success bool) {

	startPos := d.argGet(startIdx, cIKBActionEncodedArgTypeFloat32)
	if startPos == cIKBActionEncodedPosERROR {
		return 0, false
	}
	v_ := d.ext4bytes(startPos)
	return *(*float32)(unsafe.Pointer(&v_)), true
}

// GetArgFloat64 extracts float64 argument from encoded IKB action d,
// starting search from startIdx argument's index.
//
// Returns it and true as success if it is, or zero value and false if error.
func (d *tIKBActionEncoded) GetArgFloat64(startIdx int) (v float64, success bool) {

	startPos := d.argGet(startIdx, cIKBActionEncodedArgTypeFloat64)
	if startPos == cIKBActionEncodedPosERROR {
		return 0, false
	}
	v_ := d.ext8bytes(startPos)
	return *(*float64)(unsafe.Pointer(&v_)), true
}

// GetArgString extracts string argument from encoded IKB action d,
// starting search from startIdx argument's index.
//
// Returns it and true as success if it is, or zero value and false if error.
func (d *tIKBActionEncoded) GetArgString(startIdx int) (v string, success bool) {

	startPos := d.argGet(startIdx, cIKBActionEncodedArgTypeString)
	if startPos == cIKBActionEncodedPosERROR {
		v, success = "", false
		return
	}

	// startPos - arg type,
	// startPos + 1 - strlen, startPos + 2,... - string content
	return string(d.extNbytes(startPos+2, startPos+1)), true
}

// copy returns a copy of the current encoded IKB action d.
func (d *tIKBActionEncoded) copy() (copy *tIKBActionEncoded) {

	copied := *d
	return &copied
}

// dump returns a complete debug information about encoded IKB action d.
// Each slice element represent one entity of encoded IKB action d.
func (d *tIKBActionEncoded) dump() []tIKBActionEncodedDumpNode {

	// make result slice with length == len of encoded args +
	// id + ssid + args counter + args free index
	argCount := d.ArgCount()
	dumpRes := make([]tIKBActionEncodedDumpNode, argCount+4)

	// Reflect ID
	dumpRes[0].Type = "Encoded View ID"
	dumpRes[0].Pos = cIKBActionEncodedPosViewID
	dumpRes[0].Value = d.GetViewID()

	// Reflect SSID
	dumpRes[1].Type = "Session ID (SSID)"
	dumpRes[1].Pos = cIKBActionEncodedPosSessionID
	dumpRes[1].Value = d.GetSessionID()

	// Reflect args counter
	dumpRes[2].Type = "Arguments counter"
	dumpRes[2].Pos = cIKBActionEncodedPosArgsCount
	dumpRes[2].Value = argCount

	// Reflect args free index
	dumpRes[3].Type = "Arguments next free position"
	dumpRes[3].Pos = cIKBActionEncodedPosArgsFree
	dumpRes[3].Value = d[cIKBActionEncodedPosArgsFree]

	// Save info about arguments
	pos := cIKBActionEncodedPosArgsContent
	for i := 0; i < argCount; i++ {

		dumpRes[4+i].Type = "Argument (" + d.argType2S(d[pos]) + ")"
		dumpRes[4+i].Pos = pos
		dumpRes[4+i].PosType = pos
		dumpRes[4+i].TypeHeader = d[pos]

		// Save content position
		// By default content starts with offset 1
		// Exceptions: strings
		switch d[pos] {

		case cIKBActionEncodedArgTypeString:
			// pos+0 - arg type, pos+1 - strlen, pos+2,... - content
			dumpRes[4+i].PosContent = pos + 2

		default:
			dumpRes[4+i].PosContent = pos + 1
		}

		// Save values
		// By default content starts with offset 1
		// Exceptions: strings
		switch d[pos] {

		case cIKBActionEncodedArgTypeInt8:
			dumpRes[4+i].Value = d.ext1byte(pos + 1)

		case cIKBActionEncodedArgTypeInt16:
			dumpRes[4+i].Value = d.ext2bytes(pos + 1)

		case cIKBActionEncodedArgTypeInt32:
			dumpRes[4+i].Value = d.ext4bytes(pos + 1)

		case cIKBActionEncodedArgTypeInt64:
			dumpRes[4+i].Value = d.ext8bytes(pos + 1)

		case cIKBActionEncodedArgTypeUint8:
			dumpRes[4+i].Value = uint8(d.ext1byte(pos + 1))

		case cIKBActionEncodedArgTypeUint16:
			dumpRes[4+i].Value = uint16(d.ext2bytes(pos + 1))

		case cIKBActionEncodedArgTypeUint32:
			dumpRes[4+i].Value = uint32(d.ext4bytes(pos + 1))

		case cIKBActionEncodedArgTypeUint64:
			dumpRes[4+i].Value = uint64(d.ext8bytes(pos + 1))

		case cIKBActionEncodedArgTypeFloat32:
			v := int32(d.ext4bytes(pos + 1))
			dumpRes[4+i].Value = *(*float32)(unsafe.Pointer(&v))

		case cIKBActionEncodedArgTypeFloat64:
			v := int64(d.ext8bytes(pos + 1))
			dumpRes[4+i].Value = *(*float64)(unsafe.Pointer(&v))

		case cIKBActionEncodedArgTypeString:
			// pos+1 - strlen, pos+2,... - string content
			dumpRes[4+i].Value = string(d.extNbytes(pos+2, pos+1))

		default:
			dumpRes[4+i].Value = nil
		}
	}

	// Dump completed
	return dumpRes
}

// init initializes the current encoded IKB action object d.
func (d *tIKBActionEncoded) init() {
	d[cIKBActionEncodedPosArgsFree] =
		cIKBActionEncodedPosArgsContent
}
