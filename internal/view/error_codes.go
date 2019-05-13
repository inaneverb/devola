// Copyright © 2018. All rights reserved.
// Author: Alice Qio.
// Contacts: <qioalice@gmail.com>.
// License: https://opensource.org/licenses/MIT

package view

import (
	"../errors"
)

// Predefined error codes of all convert operations.
// These codes may be returned from IDConverter methods.
const (

	// No errors.
	EOK errors.Code = 0

	// Invalid View ID error.
	// Returned:
	// - From Encode method if invalid View ID is passed as argument
	// - From Register method if invalid View ID is passed as argument
	// - From Decode method if after decoding encoded View ID, the decoded
	//   View ID is invalid.
	EInvalidID errors.Code = 1

	// Invalid encoded View ID error.
	// Returned:
	// - From Decode method if invalid encoded View ID is passed as argument
	// - From Encode method if after encoding View ID, the encoded View ID
	//   is invalid.
	EInvalidIDEnc errors.Code = 2

	// Unregistered View ID error.
	// Returned:
	// - From Encode method if an unregistered by Register method View ID
	//   is passed as argument.
	// - From Decode method if passed an encoded View ID without decoded View ID
	//   pair (as a conclusion, there is no registered View ID with received
	//   encoded View ID).
	ENotRegistered errors.Code = 3

	// Already registered View ID error.
	// Returned:
	// - From Register method if passed View ID is already registered by
	//   one of previous Register calls.
	EAlreadyRegistered errors.Code = 4

	// Registered View IDs limit is reached error.
	// Returned:
	// - From Register method if there is limit of registered View IDs
	//   is reached and no one View ID can be registered anymore.
	ELimitReached errors.Code = 5
)
