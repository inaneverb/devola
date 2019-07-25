// Copyright © 2018. All rights reserved.
// Author: Alice Qio.
// Contacts: <qioalice@gmail.com>.
// License: https://opensource.org/licenses/MIT

package sender

import (
	"time"

	"github.com/qioalice/devola/core/chat"
	"github.com/qioalice/devola/core/math"
)

// param is an alias to function that takes a Sender object and changes
// only those internal constants and values, that affect only Sender.
type param func(s *Sender)

// paraml is an alias to function that takes a Sender object and changes
// only those internal constants and values, that affect Lirester - a Sender's part.
type paraml func(s *Sender)

// Params is the type of Sender params' set.
type Params struct {

	// DO NOT INSTANTIATE THIS OBJECT DIRECTLY!
	// IT DOES core/params PACKAGE!

	// Changes the number of sending message retry attempts.
	// N will be bounded above by a value 100 and below by -1.
	// -1 means infinity number of attempts.
	//
	// WARNING!
	// An infinity number of attempts very highly load the system.
	// Do it only if you really sure that you need it.
	RetryAttempts func(n int) param

	// Disables Lirester at all.
	//
	LiresterDisable func() param

	// Changes N, T of 1st Lirester restriction rule (see Lirester's docs).
	// There is no-op if n <= 0 or t <= 10ms and t/n must be in the range [100ms..1min].
	LiresterMainLoopDelay func(n int, t time.Duration) paraml

	// Changes N'(i), T'(i) of 2nd Lirester restriction rule (see Lirester's docs).
	// There is no-op if typ >= core/chat.MaxTypeValue or n <= 0 or t <= 10 ms
	// and t/n must be in the range [100ms..1min].
	LiresterAddChatRuleFor func(n int, t time.Duration, typ chat.Type) paraml

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
	SetEXREOF func(N, M int16) param
}

// A storage of all Sender params.
var vParams Params

// Initializes storage of all Sender params.
func init() {

	vParams.RetryAttempts =
		func(n int) param {
			return param(func(s *Sender) {
				retryAttempts = int8(math.ClampI(n, -1, 100))
			})
		}

	vParams.SetEXREOF =
		func(N, M int16) param {
			if M < tusent.cChChMinCapacity {
				M = tusent.cChChMinCapacity
			}
			return param(func(s *Sender) {
				chchUnusedCap = N
				chchCap = M
			})
		}

	vParams.LiresterMainLoopDelay = func(n int, t time.Duration) paraml {
		if n <= 0 || t <= 10*time.Microsecond {
			return nil
		}
		delay := time.Duration(int64(t) / int64(n))
		if delay < 100*time.Microsecond || delay > 1*time.Minute {
			return nil
		}
		return func(l *Lirester) {
			l.consts.mainLoopDelay = delay
		}
	}

	vParams.LiresterAddChatRuleFor = func(n int, t time.Duration, typ chat.Type) paraml {
		if typ >= chat.MaxTypeValue || n <= 0 || t <= 10*time.Microsecond {
			return nil
		}
		delay := time.Duration(int64(t) / int64(n))
		if delay < 100*time.Microsecond || delay > 1*time.Minute {
			return nil
		}
		return func(l *Lirester) {
			l.consts.Ns[typ] = uint8(n)
			l.consts.Ts[typ] = int64(t)
		}
	}
}
