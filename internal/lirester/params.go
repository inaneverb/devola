// Copyright © 2018. All rights reserved.
// Author: Alice Qio.
// Contacts: <qioalice@gmail.com>.
// License: https://opensource.org/licenses/MIT

package lirester

import (
	"time"
)

// TODO: Add params to enable profiler

// Param is an alias to function that takes a Lirester object
// and changes its internal constants and values.
//
// Used for Lirester's consturctor (MakeLirester) or Lirester.RequestRestartWith.
type Param func(l *Lirester)

// Params is the type of storage of Lirester params.
// Is a part of tParams type.
type Params struct {

	// Changes the main loop ticker delay to the passed value.
	// Delay must be in the [100ms..1min] range.
	// More info: Lirester.consts.mainLoopDelay field.
	MainLoopDelay func(delay time.Duration) Param

	// Changes the Lirester internal chat's lifetime.
	// Lifetime must be in the [1min..1day] range.
	// More info: Lirester.consts.chatLifeTime field.
	LChatLifetime func(lifetime time.Duration) Param

	// Params to change some values about working Lirester around chats
	// with users.
	UserChat struct {

		// Changes the value how many messages can be send to user chat
		// per each iteration.
		// Value must be in the [1..99] range.
		// More info: Lirester.consts.cLiresterUserChatN field.
		MessagesPerIter func(num uint8) Param

		// Changes the duration of user chat's iteration.
		// Value mist be in the [100ms..1h] range.
		// More info: Lirester.consts.cLiresterUserChatT field.
		IterPeriod func(period time.Duration) Param
	}

	// Params to change some values about working Lirester around group chats,
	// channels, etc.
	GroupChat struct {

		// Changes the value how many messages can be send to user chat
		// per each iteration.
		// Value must be in the [1..99] range.
		// More info: Lirester.consts.cLiresterUserChatN field.
		MessagesPerIter func(num uint8) Param

		// Changes the duration of user chat's iteration.
		// Value mist be in the [100ms..1h] range.
		// More info: Lirester.consts.cLiresterUserChatT field.
		IterPeriod func(period time.Duration) Param
	}
}

// A storage of all Lirester params.
var vParams Params

// Initializes storage of all Lirester params.
func init() {

	vParams.MainLoopDelay =
		func(delay time.Duration) Param {
			if delay < 100*time.Microsecond || delay > 1*time.Minute {
				return Param(nil)
			}
			return Param(func(l *Lirester) {
				l.consts.mainLoopDelay = delay
			})
		}

	vParams.LChatLifetime =
		func(lifetime time.Duration) Param {
			if lifetime < 1*time.Minute || lifetime > 24*time.Hour {
				return Param(nil)
			}
			return Param(func(l *Lirester) {
				l.consts.chatLifeTime = lifetime.Nanoseconds()
			})
		}

	vParams.UserChat.MessagesPerIter =
		func(num uint8) Param {
			if num == 0 || num > 100 {
				return Param(nil)
			}
			return Param(func(l *Lirester) {
				l.consts.userChatN = num
			})
		}

	vParams.UserChat.IterPeriod =
		func(period time.Duration) Param {
			if period < 100*time.Microsecond || period > 1*time.Hour {
				return Param(nil)
			}
			return Param(func(l *Lirester) {
				l.consts.userChatT = period.Nanoseconds()
			})
		}

	vParams.GroupChat.MessagesPerIter =
		func(num uint8) Param {
			if num == 0 || num > 100 {
				return Param(nil)
			}
			return Param(func(l *Lirester) {
				l.consts.groupChatN = num
			})
		}

	vParams.GroupChat.IterPeriod =
		func(period time.Duration) Param {
			if period < 100*time.Microsecond || period > 1*time.Hour {
				return Param(nil)
			}
			return Param(func(l *Lirester) {
				l.consts.groupChatT = period.Nanoseconds()
			})
		}
}
