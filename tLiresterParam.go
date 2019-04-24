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

// 'tLiresterParam' is the alias to function that applying to the some
// lirester object.
// It uses as parameters for lirester constructor to initialize consts.
type tLiresterParam func(l *tLirester)

// 'LiresterMainLoopDelay' creates the new lirester creator parameter
// that can be used as argument for lirester constructor.
// This parameter means, how often lirester main loop timer will be tick.
// Allowable values: (100ms..1min).
func LiresterMainLoopDelay(delay time.Duration) tLiresterParam {
	if delay <= 100*time.Microsecond || delay >= 1*time.Minute {
		return nil
	}
	return func(l *tLirester) {
		l.consts.mainLoopDelay = delay.Nanoseconds()
	}
}

// 'LiresterMainLoopDelay' creates the new lirester creator parameter
// that can be used as argument for lirester constructor.
// This parameter means, how much time the lirester chat will exist in the
// internal lirester structure.
// Less means more CPU load and less RAM load,
// More means less CPU load and more RAM load.
// Allowable values: (1min..1day).
func LiresterChatLifetime(lifetime time.Duration) tLiresterParam {
	if lifetime <= 1*time.Minute || lifetime >= 24*time.Hour {
		return nil
	}
	return func(l *tLirester) {
		l.consts.chatLifeTime = lifetime.Nanoseconds()
	}
}

// 'LiresterMainLoopDelay' creates the new lirester creator parameter
// that can be used as argument for lirester constructor.
// This parameter means, how much messages can be allowed by lirester
// for user chats without cleanup until the limit is reached.
// It's not recommend to set this value more than 1, 'cause in another cases
// you can see the 429 Telegram error.
// Allowable values: (0..100].
func LiresterUserChatMessageNumberPerIter(num uint8) tLiresterParam {
	if num == 0 || num > 100 {
		return nil
	}
	return func(l *tLirester) {
		l.consts.userChatMsgNumPerIter = num
	}
}

// 'LiresterMainLoopDelay' creates the new lirester creator parameter
// that can be used as argument for lirester constructor.
// This parameter means how much time must be passed until cleanup
// isn't started for some sent message to user chats.
// Thus, increasing this value means increasing the one 'iteration' range.
// It's not recommend to set this value less than 1sec, 'cause in another cases
// you can see the 429 Telegram error.
// Allowable values: (100ms..1hr).
func LiresterUserChatCleanupDelay(delay time.Duration) tLiresterParam {
	if delay <= 100*time.Microsecond || delay >= 1*time.Hour {
		return nil
	}
	return func(l *tLirester) {
		l.consts.userChatCleanupDelay = delay.Nanoseconds()
	}
}

// 'LiresterMainLoopDelay' creates the new lirester creator parameter
// that can be used as argument for lirester constructor.
// This parameter means, how much messages can be allowed by lirester
// for group chats without cleanup until the limit is reached.
// It's not recommend to set this value more than 20, 'cause in another cases
// you can see the 429 Telegram error.
// Allowable values: (0..100].
func LiresterGroupChatMessageNumberPerIter(num uint8) tLiresterParam {
	if num > 100 {
		return nil
	}
	return func(l *tLirester) {
		l.consts.groupChatMsgNumPerIter = num
	}
}

// 'LiresterMainLoopDelay' creates the new lirester creator parameter
// that can be used as argument for lirester constructor.
// This parameter means how much time must be passed until cleanup
// isn't started for some sent message to group chats.
// Thus, increasing this value means increasing the one 'iteration' range.
// It's not recommend to set this value less than 1min, 'cause in another cases
// you can see the 429 Telegram error.
// Allowable values: (100ms..1hr).
func LiresterGroupChatCleanupDelay(delay time.Duration) tLiresterParam {
	if delay <= 100*time.Microsecond || delay >= 1*time.Hour {
		return nil
	}
	return func(l *tLirester) {
		l.consts.groupChatCleanupDelay = delay.Nanoseconds()
	}
}
