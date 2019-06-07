// Copyright © 2019. All rights reserved.
// Author: Alice Qio.
// Contacts: <qioalice@gmail.com>.
// License: https://opensource.org/licenses/MIT

package ctxutils

import (
	"unsafe"
)

// ctxStringer is an alias to function that takes some backend context by its pointer
// and returns its string representation.
type ctxStringer func(ctx unsafe.Pointer) string

//
type ctxTransactionFinisher func(ctx unsafe.Pointer) error
