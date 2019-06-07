// Copyright © 2018. All rights reserved.
// Author: Alice Qio.
// Contacts: <qioalice@gmail.com>.
// License: https://opensource.org/licenses/MIT

package lirester

import (
	"time"

	"github.com/qioalice/devola/core/chat"
)

// TODO: Add params to enable profiler

// Param is an alias to function that takes a Lirester object
// and changes its internal constants and values.
//
// Used for Lirester's constructor (MakeLirester) or Lirester.RequestRestartWith.
type Param func(l *Lirester)

// Params is the type of storage of Lirester params.
// Is a part of tParams type.
type Params struct {

	// Changes N, T of 1st Lirester restriction rule (see Lirester's docs).
	// There is no-op if n <= 0 or t <= 10ms and t/n must be in the range [100ms..1min].
	NPerTGlobal func(n int, t time.Duration) Param

	// Changes N'(i), T'(i) of 2nd Lirester restriction rule (see Lirester's docs).
	// There is no-op if typ >= core/chat.MaxTypeValue or n <= 0 or t <= 10 ms
	// and t/n must be in the range [100ms..1min].
	NPerTForType func(n int, t time.Duration, typ chat.Type) Param
}

// A storage of all Lirester params.
var vParams Params

// Initializes storage of all Lirester params.
func init() {

	vParams.NPerTGlobal = func(n int, t time.Duration) Param {
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

	vParams.NPerTForType = func(n int, t time.Duration, typ chat.Type) Param {
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
