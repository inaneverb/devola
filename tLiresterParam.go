// Copyright Â© 2018. All rights reserved.
// Author: Alice Qio.
// Contacts: <qioalice@gmail.com>.
// License: https://opensource.org/licenses/MIT
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"),
// to deal in the Software without restriction, including without limitation
// the rights to use, copy, modify, merge, publish, distribute, sublicense,
// and/or sell copies of the Software, and to permit persons to whom
// the Software is furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS
// IN THE SOFTWARE.

package tgbot

import (
	"time"
)

// tLiresterParam is an alias to function that takes a Lirester object
// and changes its internal constants and values.
// 
// Used for Lirester's consturctor (makeLirester) or tLirester.RequestRestartWith.
type tLiresterParam func(l *tLirester)

// tLiresterParams is the type of storage of tLirester params.
// Is a part of tParams type.
type tLiresterParams struct {

	// Changes the main loop ticker delay to the passed value.
	// Delay must be in the [100ms..1min] range.
	// More info: tLirester.consts.mainLoopDelay field.
	MainLoopDelay func(delay time.Duration) tLiresterParam

	// Changes the Lirester internal chat's lifetime.
	// Lifetime must be in the [1min..1day] range.
	// More info: tLirester.consts.chatLifeTime field.
	LChatLifetime func(lifetime time.Duration) tLiresterParam

	// Params to change some values about working Lirester around chats
	// with users.
	UserChat struct {

		// Changes the value how many messages can be send to user chat 
		// per each iteration.
		// Value must be in the [1..99] range.
		// More info: tLirester.consts.cLiresterUserChatN field.
		MessagesPerIter func(num uint8) tLiresterParam

		// Changes the duration of user chat's iteration.
		// Value mist be in the [100ms..1h] range.
		// More info: tLirester.consts.cLiresterUserChatT field.
		IterPeriod func(period time.Duration) tLiresterParam
	}

	// Params to change some values about working Lirester around group chats,
	// channels, etc.
	GroupChat struct {

		// Changes the value how many messages can be send to user chat 
		// per each iteration.
		// Value must be in the [1..99] range.
		// More info: tLirester.consts.cLiresterUserChatN field.
		MessagesPerIter func(num uint8) tLiresterParam

		// Changes the duration of user chat's iteration.
		// Value mist be in the [100ms..1h] range.
		// More info: tLirester.consts.cLiresterUserChatT field.
		IterPeriod func(period time.Duration) tLiresterParam
	}
}

// paramsLirester is the storage of tLirster params.
// Is a part of Params object.
var paramsLirester = tLiresterParams{}

// Initializes paramsLirester object.
///
// There is no in-place initialization because tLiresterParams have
// nested structs and a separate initialization reduces the amount of code
// and increases readability.
func init() {

	paramsLirester.MainLoopDelay = 
	func(delay time.Duration) tLiresterParam {
		if delay < 100 * time.Microsecond || delay > 1 * time.Minute {
			return tLiresterParam(nil)
		}
		return tLiresterParam(func(l *tLirester) {
			l.consts.mainLoopDelay = delay
		})
	}

	params.Lirester.LChatLifetime = 
	func(lifetime time.Duration) tLiresterParam {
		if lifetime < 1 * time.Minute || lifetime > 24 * time.Hour {
			return tLiresterParam(nil)
		}
		return tLiresterParam(func(l *tLirester) {
			l.consts.chatLifeTime = lifetime.Nanoseconds()
		})
	}

	paramsLirester.UserChat.MessagesPerIter = 
	func(num uint8) tLiresterParam {
		if num == 0 || num > 100 {
			return tLiresterParam(nil)
		}
		return tLiresterParam(func(l *tLirester) {
			l.consts.userChatMsgNumPerIter = num
		})
	}

	paramsLirester.UserChat.IterPeriod = 
	func(period time.Duration) tLiresterParam {
		if delay < 100 * time.Microsecond || delay > 1 * time.Hour {
			return tLiresterParam(nil)
		}
		return tLiresterParam(func(l *tLirester) {
			l.consts.userChatCleanupDelay = delay.Nanoseconds()
		})
	}

	paramsLirester.GroupChat.MessagesPerIter = 
	func(num uint8) tLiresterParam {
		if num == 0 || num > 100 {
			return tLiresterParam(nil)
		}
		return tLiresterParam(func(l *tLirester) {
			l.consts.groupChatMsgNumPerIter = num
		})
	}

	paramsLirester.GroupChat.IterPeriod = 
	func(period time.Duration) tLiresterParam {
		if delay < 100 * time.Microsecond || delay > 1 * time.Hour {
			return tLiresterParam(nil)
		}
		return tLiresterParam(func(l *tLirester) {
			l.consts.groupChatCleanupDelay = delay.Nanoseconds()
		})
	}
}
