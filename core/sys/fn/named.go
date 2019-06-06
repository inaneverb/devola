// Copyright © 2019. All rights reserved.
// Author: Alice Qio.
// Contacts: <qioalice@gmail.com>.
// License: https://opensource.org/licenses/MIT

package fn

import (
	"unsafe"
)

// Named is an object that stores a pointer to a function and its name.
// It is assumed that pointer will be obtained using TakeCallableAddr func.
//
// You can also get callable pointer of function (implicitly) and do name it
// (construct object of this class) using MakeNamed function.
type Named struct {
	Name string
	Ptr  unsafe.Pointer
}

// MakeNamed creates a Named object using passed name and function object.
//
// MakeNamed includes TakeCallableAddr call and you do not need to do it:
// you can pass your function object directly to this constructor.
//
// However, you can also pass pointer obtained by TakeCallableAddr.
// In that case, as you expect, there is no TakeCallableAddr implicitly call.
//
// WARNING!
// DO NOT PASS A POINTER OBTAINED BY TakeRealAddr FUNC. OTHERWISE THERE IS UB!
func MakeNamed(name string, fn interface{}) Named {

	var ptr unsafe.Pointer

	// using a second argument disables panic if fn is not unsafe.Pointer
	if ptr, _ = fn.(unsafe.Pointer); ptr == nil {
		ptr = TakeCallableAddr(fn)
	}

	return Named{name, ptr}
}
