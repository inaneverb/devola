// Copyright © 2018. All rights reserved.
// Author: Alice Qio.
// Contacts: <qioalice@gmail.com>.
// License: https://opensource.org/licenses/MIT

package sender

import (
	"../math"
)

// Param is an alias to function that takes a Sender object
// and changes its internal constants and values.
//
// Used for Sender's constructor (MakeSender).
type Param func(s *Sender)

// Params is the type of storage of Sender params.
// Is a part of tParams type.
type Params struct {

	// Forces to always use MD (Markdown) to represent a text in
	// Telegram Bot responses.
	// Cancels previous AlwaysUseHTML if it was set.
	AlwaysUseMD func() Param

	// Forces to always use HTML to represent a text in Telegram Bot responses.
	// Cancels previous AlwaysUseMD if it was set.
	AlwaysUseHTML func() Param

	// Changes the number of sending message retry attempts.
	// N will be bounded above by a value 100 and below by -1.
	// -1 means infinity number of attempts.
	//
	// WARNING!
	// An infinity number of attempts very highly load the system.
	// Do it only if you really sure that you need it.
	RetryAttempts func(n int) Param

	// Changes the values of
	// EXperimental REallocation Optimization Feautre or disables it.
	// Pass any negative value or 0 as n if you want to disable it.
	//
	// Enabling that feature will cause decreasing memory allocate and GC
	// operations because of reusing allocated memory in some cases.
	//
	// M is the size of the buffer of
	// simultaneously unsent messages to the same chat.
	// By default this value is 16 and M can not be less than that value!
	//
	// WARNING!
	// Changing these values may cause UNEXPECTED MEMORY CONSUMPTION.
	// Decreasing allocate/free operations means that memory will be allocated
	// ALWAYS OR MOST of the TIME!
	// that feature means that sometimes the amount of allocated
	// RESERVED memory can REACH:
	// 8 + 4*N + N*(16 + 4*M) bytes for 32-bit, and
	// 16 + 8*N + N*(24 + 8*M) bytes for 64-bit!
	SetEXREOF func(N, M int16) Param
}

// A storage of all Sender params.
var vParams Params

// Initializes storage of all Sender params.
func init() {

	vParams.AlwaysUseMD =
		func() Param {
			return Param(func(s *Sender) {
				s.consts.isAlwaysUseMD = true
				s.consts.isAlwaysUseHTML = false
			})
		}

	vParams.AlwaysUseHTML =
		func() Param {
			return Param(func(s *Sender) {
				s.consts.isAlwaysUseHTML = true
				s.consts.isAlwaysUseMD = false
			})
		}

	vParams.RetryAttempts =
		func(n int) Param {
			return Param(func(s *Sender) {
				s.consts.retryAttempts = int8(math.ClampI(n, -1, 100))
			})
		}

	vParams.SetEXREOF =
		func(N, M int16) Param {
			if M < cChChMinCapacity {
				M = cChChMinCapacity
			}
			return Param(func(s *Sender) {
				s.consts.chchUnusedCap = N
				s.consts.chchCap = M
			})
		}
}
