// Copyright © 2019. All rights reserved.
// Author: Alice Qio.
// Contacts: <qioalice@gmail.com>.
// License: https://opensource.org/licenses/MIT

package sys

import (
	"unsafe"
)

type typedInterface struct {
	typ  uintptr
	word unsafe.Pointer
}

// FnPtr returns a real address of function f or nil if f is nil.
//
// If f is not a function, you still get the pointer, but it is unknown
// what that pointer points to.
//
// YOU CAN NOT USE RETURNED POINTER TO CALLING PASSED FUNCTION!
// For that purpose use FnPtrCallable.
func FnPtr(f interface{}) unsafe.Pointer {
	if f == nil {
		return nil
	}
	return (*GoInterface)(unsafe.Pointer(&f)).Word
}

// FnPtrCallable returns an address of function f using which you can call
// a function that was passed or nil if f is nil.
// To calling just convert returned untyped pointer to function-typed,
// dereference it and call.
//
// You can AVOID TYPE CHECKS using that way (wrong argument types, wrong
// return types) but in that way the BEHAVIOUR is UNDEFINED and do it only
// if you know what you're doing.
//
// If f is not a function, you still get the pointer, but it is unknown
// what that pointer points to.
func FnPtrCallable(f interface{}) unsafe.Pointer {
	// FnPtr and fnPtrReal2Callable has nil checks.
	return fnPtrReal2Callable(FnPtr(f))
}

// ptrReal2Callable converts a real function pointer to a pointer using which
// becomes possible to call a function ptr points to.
func fnPtrReal2Callable(ptr unsafe.Pointer) unsafe.Pointer {

	type fptr struct {
		ptr unsafe.Pointer
	}

	if ptr == nil {
		return nil
	}

	o := new(fptr)
	o.ptr = ptr
	return unsafe.Pointer(&o.ptr)
}

// // Ptr returns an address of function f or nil if f is nil.
// // If f is not a function, you still get the pointer, but it is unknown
// // what that pointer points to.
// func Ptr(f interface{}) unsafe.Pointer {

// 	// In Golang you can not just take an address of some func
// 	// only if it is not a func literal which is assigned to any var.
// 	//
// 	// But you can assign any kind of func to type-compatible variable
// 	// or interface{}.
// 	// And you can call then that assigned function through used variable.
// 	// But anyway that way still requires a type-compatible variable
// 	// or type-casting interface{} and you can not break that type-compatible
// 	// rules:
// 	// You can not assign func to some variable with wrong func signature,
// 	// You can not cast interface{} that stores func with one signature
// 	// to func with another.
// 	//
// 	// Unless you have an untyped pointer like void*.
// 	//
// 	// reflect.ValueOf(f).Pointer() returns a REAL function address,
// 	// but you CAN NOT use it to call.
// 	// I do not know how it is implemented in an internal Golang parts, but it is.
// 	//
// 	//
// 	// CAN:
// 	// f := func(){}
// 	// &f
// 	//
// 	// CAN NOT:
// 	// func f(){}
// 	// &f

// 	// But you can assign any kind of func to type-compatible variable
// 	// or interface{}.
