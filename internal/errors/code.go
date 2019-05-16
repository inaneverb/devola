// Copyright © 2018. All rights reserved.
// Author: Alice Qio.
// Contacts: <qioalice@gmail.com>.
// License: https://opensource.org/licenses/MIT

package errors

// Code represents error code type.
// The values of this type are returned by method Code of each type that
// implements Error interface, and also sometimes are returned directly
// from methods.
type Code int

// Predefined constants of error codes.
// First two letters of each consts, EC, means Error Code.
const (

	// ECOK means Error Code OK: Success, No errors.
	ECOK Code = 0
)
