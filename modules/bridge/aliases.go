// Copyright © 2019. All rights reserved.
// Author: Alice Qio.
// Contacts: <qioalice@gmail.com>.
// License: https://opensource.org/licenses/MIT

package bridge

import (
	"unsafe"
)

//
type Handler func(ctx unsafe.Pointer)

//
type Middleware func(ctx unsafe.Pointer) (isAllowed bool)

//
type OnSuccessFinisher func(ctx, o unsafe.Pointer)

//
type OnErrorFinisher func(ctx unsafe.Pointer, err error)
