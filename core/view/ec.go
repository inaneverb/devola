// Copyright © 2019. All rights reserved.
// Author: Alice Qio.
// Contacts: <qioalice@gmail.com>.
// License: https://opensource.org/licenses/MIT

package view

import (
	"github.com/qioalice/devola/core/errors"
)

// Predefined error codes of all convert operations.
// These codes may be returned from IDConv methods.
const (

	// Invalid View ID error.
	// Returned:
	// - From Encode method if invalid View ID is passed as argument
	// - From Decode method if after decoding encoded View ID, the decoded
	//   View ID is invalid.
	ECInvalidID errors.Code = 1

	// Invalid encoded View ID error.
	// Returned:
	// - From Decode method if invalid encoded View ID is passed as argument
	// - From Encode method if after encoding View ID, the encoded View ID
	//   is invalid.
	ECInvalidIDEnc errors.Code = 2

	// Unregistered View ID error.
	// Returned:
	// - From Encode method if an unregistered by reg method View ID
	//   is passed as argument.
	// - From Decode method if passed an encoded View ID without decoded View ID
	//   pair (as a conclusion, there is no registered View ID with received
	//   encoded View ID).
	ECNotRegistered errors.Code = 3
)
