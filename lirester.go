// Copyright © 2018. All rights reserved.
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

import "time"

// -- LiRester --
// Telegram [Li]mit [Rest]rictions Avoid[er]
// NOT THREAD SAFETY!
//
// There are three following rules of Lirester concept:
// 1. Check, whether message can be sent (Try!)
// 2. Confirm, that message is really sent and you want be protected by
// Lirester from Telegram limits overlow (Approve!)
// 3. Just wait when Lirester will allow you to send message again (Cleanup!)
//
// How it works.
// First of all, lirester objects starts and contains the some timer
// that contains channel. Each timer iteration, timer's channel receives
// time object. Thus, caller must subscribe to the 'MainLoop.C' channel
// and each time when this channel will have a new object will means
// that you can send message to Telegram servers.
// (Avoiding first restrictions: No more than X req per Y).
// Secondly, lirester object has map from chat id to the lirester chat object
// and lirester chat object contains info how much messages was sent already
// and when it was last time.
// (Third important things is decreasing counter of "how much messages was sent",
// but we'll speak about it later).
// So, here starts the Lirester Public API: 'Try' method, that takes
// the current time as unixnano timestamp (for decreasing time.Now() calls)
// and chat id we must check which.
// Method 'Try' just checks if counter of "how much messages was sent" is full.
// If it's so, it's returns false (signalize that we can't send message).
// If counter allows to send at least one more message, it returns true.
// Ok. We can send message. Do it. Go to next step.
// Now we must increase the counter of sent messages.
// Use 'Approve' method for this operation.
// 'Approve' just increases the counter and creates the cleanup rule.
// 'Approve' takes the current time as unixnano timestamp
// (also for decreasing time.Now() calls), chat id for that counter must be
// increased and cleanup rule must be created and bool flag that represents
// one small question: is it chat with a user or a channel/group/supergroup?
// Q: For what 'isUser' flag is exists?
// A: Because Telegram has different restrictions for chats 'user-bot' and
// 'bot-in-the-group'.
// And cleanup. The most important part.
// 'Approve' also created the cleanup rule. What is it? Just small object
// of two ints: chat id and unixnano timestamp.
// These cleanup rules will be applied in the 'Cleanup' method.
// So, you must call 'Cleanup' method as often as it possible for best results.
// 'Cleanup' method just extracts all accumulated cleanup rules as
// [chat id, unixnano timestamp] and if time when cleanup rule must be applied
// has come, the counter of specified chat id decreasing by 1.
type tLirester struct {
	// Consts section
	consts struct {
		mainLoopDelay          int64
		chatLifeTime           int64
		userChatMsgNumPerIter  uint8
		userChatCleanupDelay   int64
		groupChatMsgNumPerIter uint8
		groupChatCleanupDelay  int64
	}
	// Core section
	MainLoop *time.Timer
	core     map[int64]*tLiresterChat
	// Cleanup section
	chCleanup           chan tLiresterCleanupConfig
	isBackgroundCleanup bool
}

// 'tLiresterParam' is the alias to function that applying to the some
// lirester object.
// It uses as parameters for lirester constructor to initialize consts.
type tLiresterParam func(l *tLirester)

// 'tLiresterChat' represents the type of LiRester chat that contains
// three important things:
// - How much messages has been sent at this moment to that chat?
// - Is this chat with user or a channel/group/supergroup?
// - When this chat has been updated last time?
type tLiresterChat struct {
	data        uint8
	lastUpdated int64
}

// 'tLiresterCleanupConfig' represents the cleanup config that contains
// two important things:
// - Chat id of chat cleanup operation will perform about
// - Unixnano timestamp when cleanup operation will perform
type tLiresterCleanupConfig struct {
	who  int64
	when int64
}

// 'howMuch' returns how much messages already sent to the chat at this iter.
// Save all bits except high bit, 'cause high bit serves for indicate is this
// chat with a user or a (super) group.
// (See 'isUser' and 'isGroup' methods for details).
func (c *tLiresterChat) howMuch() uint8 { return uint8(c.data) & 0x7F }

// 'setHowMuch' sets the counter to the 'v' value in the current chat and
// returns changed value.
// Warning! 'v' physically must be in the range [-127..127].
// Warning! 'v' logically must be in the range [0..127].
// Otherwise, 'v mod 128' will be used as 'v'.
// Deprecated: Unneccessary
func (c *tLiresterChat) setHowMuch(v uint8) *tLiresterChat {
	uint8(c.data) &= 0x80
	uint8(c.data) |= v & 0x7F
	return c
}

// 'incHowMuch' increases the counter by the 'delta' value in the current chat
// and returns changed value.
// Note! If you want to decrease a counter, just use 'decHowMuch' method,
// or pass the negative delta to the current method.
// Warning! 'c.howMuch() + delta' physically must be in the range [-127..127].
// Warning! 'c.howMuch() + delta' logically must be in the range [0..127].
// Otherwise, there is no-op.
func (c *tLiresterChat) incHowMuch(delta uint8) *tLiresterChat {
	nv := (uint8(c.data) & 0x7F) + delta
	if nv&0x80 == 0 {
		uint8(c.data) = (uint8(c.data) & 0x80) | nv
	}
	return c
}

// 'decHowMuch' decreasses the counter by the 'delta' value in the current chat
// and returns changed value.
// See 'incHowMuch' method to understand the limitations, notes and warnings.
func (c *tLiresterChat) decHowMuch(delta uint8) *tLiresterChat {
	return c.incHowMuch(-delta)
}

// 'isUser' returns true only if the current chat is chat with a user.
// Cleared high bit means that the current chan is just a chat with user.
func (c *tLiresterChat) isUser() bool { return uint8(c.data)&0x80 == 0 }

// 'isGroup' returns true only if the current chat is a group or
// a supergroup, not chat with a user!
func (c *tLiresterChat) isGroup() bool { return uint8(c.data)&0x80 != 0 }

// 'setType' sets the flag of lirester chat type (with user or a group)
// by 'isUser' bool flag. Returns the current chat with changed flag.
func (c *tLiresterChat) setType(isUser bool) *tLiresterChat {
	uint8(c.data) &= 0x7F // cleanup prev value of flag
	if !isUser {
		uint8(c.data) |= 0x80
	} // if it's not user, set high bit
	return c
}

// 'setLastUpdated' sets the 'unixns' as new last updated time.
// Note! 'now' must be a unixnano timestamp.
func (c *tLiresterChat) setLastUpdated(now int64) {
	c.lastUpdated = now
}

// 'IsLifetimeOver' (believing that 'now' is the current unixnano timestamp)
// returns true if the lifetime of the lirester chat is over.
func (l *tLirester) isLifetimeOver(now int64, chat *tLiresterChat) bool {
	return chat.lastUpdated+l.consts.chatLifeTime < now
}

// 'Try' (believing that 'now' is the current unixnano timestamp)
// checks whether some message can be send to the chat with id 'chatId'
// right now.
// If it's possible, not nil lirester chat object is returned,
// and you can use it for manually increase the message counter, update
// last updated flag, etc.
// If message can't be sent to the chat with specified chat id right now,
// nil is returned.
func (l *tLirester) Try(now, chatID int64) (isAllow bool) {
	lirch := l.getChat(now, chatID)
	howmuch := lirch.howMuch()
	allowForUser := lirch.isUser() && howmuch < l.consts.userChatMsgNumPerIter
	allowForGroup := lirch.isGroup() && howmuch < l.consts.groupChatMsgNumPerIter
	// If howmuch == 0 it means that this chat has been created early or
	// if it's not, it doesn't matter, 0 means that we can send message
	// under any conditions.
	// Otherwise, if howmuch != 0, the flag of user must be set
	// Resolve condition by flag and howmuch by expressions above and
	// return the answer of question: "Can we send message to that chat?"
	return howmuch == 0 || allowForUser || allowForGroup
}

// 'getChat' (believing that 'now' is the current unixnano timestamp)
// returns the lirester chat object with the specified chat id.
// This method always returns not nil object.
func (l *tLirester) getChat(now, chatID int64) *tLiresterChat {
	lirch := l.core[chatID]
	if lirch == nil {
		lirch = l.makeChat(now, chatID)
	}
	return lirch
}

// 'Approve' is the alias of three important operations which applied together
// means that message to the specified chat is considered successfully sent.
// This method must be called only when message really sent and lirester
// must protect you from Telegram limits overflow.
func (l *tLirester) Approve(now, chatID int64, isUser bool) {
	l.getChat(now, chatID).setType(isUser).incHowMuch(1).setLastUpdated(now)
	l.addCleanup(now, chatID, isUser)
}

// 'LastUpdated' returns the unixnano timestamp of last update of
// lirester chat with specified chat id 'chatId'.
// If chat with the specified chat id not found in lirester object,
// -1 is returned.
func (l *tLirester) LastUpdated(chatID int64) int64 {
	lirch := l.core[chatID]
	if lirch == nil {
		return -1
	}
	return lirch.lastUpdated
}

// 'Cleanup' (believing that 'now' is the current unixnano timestamp)
// once takes all accumulated cleanup rules and tries to apply it
// only if the required time (by rule) has come.
// All not applied rules will be returned to the cleanup rules queue (channel).
func (l *tLirester) Cleanup(now int64) {
	for i, n := 0, len(l.chCleanup); i < n; i++ {
		clcfg := <-l.chCleanup
		// If cleanup config must be executed later than now,
		// just return it to the cleanup configs channel and go to next config
		if clcfg.when > now {
			l.chCleanup <- clcfg
			continue
		}
		lirch := l.core[clcfg.who]
		// If counter isn't equal to zero, apply cleanup (decrease counter)
		if lirch.howMuch() > 0 {
			lirch.decHowMuch(1).setLastUpdated(now)
		}
		// If lirester chat lifetime is over, just cleanup it completely
		if l.isLifetimeOver(now, lirch) {
			delete(l.core, clcfg.who)
		}
	}
}

// 'makeChat' creates the new lirester chat with the specified values.
// It also appends created chat to the current lirester object
// and returns created chat.
// It guaranteed, that 'makeChat' call always returns not nil object.
func (l *tLirester) makeChat(now, chatID int64) *tLiresterChat {
	ch := &tLiresterChat{lastUpdated: now}
	l.core[chatID] = ch
	return ch
}

// 'addCleanup' creates the new cleanup rule for chat with specified chat id
// that will be applied after delay will pass since the 'now' unixnano timestamp.
func (l *tLirester) addCleanup(now, chatID int64, isUser bool) {
	delay := l.consts.userChatCleanupDelay
	if !isUser {
		delay = l.consts.groupChatCleanupDelay
	}
	l.chCleanup <- tLiresterCleanupConfig{who: chatID, when: delay + now}
}

// 'newLirester' creates the new lirester object with the specified lirester
// params, starts the main lirester loop, starts the lirester background
// cleanup and return the created lirester object.
// 'params' might be only 'tLiresterParam' type. Values of any other type
// will be ignored.
func newLirester(params ...interface{}) *tLirester {
	l := &tLirester{}
	// The default values of lirester params
	const (
		defMainLoopDelay          = int64(1 * time.Second / 30)
		defChatLifetime           = int64(10 * time.Minute)
		defUserChatMsgNumPerIter  = uint8(1)
		defUserChatCleanupDelay   = int64(1 * time.Second)
		defGroupChatMsgNumPerIter = uint8(20)
		defGroupChatCleanupDelay  = int64(1 * time.Minute)
		defCleanupChanLen         = 16384 // 2^14
		defCleanupDelay           = 100 * time.Millisecond
	)
	// Apply default lirester params.
	// It'll be overwritten later by passed lirester params.
	l.consts.mainLoopDelay = defMainLoopDelay
	l.consts.chatLifeTime = defChatLifetime
	l.consts.userChatMsgNumPerIter = defUserChatMsgNumPerIter
	l.consts.userChatCleanupDelay = defUserChatCleanupDelay
	l.consts.groupChatMsgNumPerIter = defGroupChatMsgNumPerIter
	l.consts.groupChatCleanupDelay = defGroupChatCleanupDelay
	// Apply lirester params
	for _, param := range params {
		if param, ok := param.(tLiresterParam); ok && param != nil {
			param(l)
		}
	}
	// Allocate mem for main core lirester and cleanup channel
	l.core = make(map[int64]*tLiresterChat)
	l.chCleanup = make(chan tLiresterCleanupConfig, defCleanupChanLen)
	// Create and start main loop timer
	l.MainLoop = time.NewTimer(time.Duration(l.consts.mainLoopDelay))
	// Finish, return created and ready lirester object
	return l
}

// 'wipeLirester' stops the lirester background cleanup, stops the lirester
// main loop timer, closes cleanup chan, and nulls the variable in which
// pointer to the lirester is store.
// Note! Lirester object must be passed by double pointer for avoiding
// to use lirester object in the callable code after calling this function.
// Warning! Lirester object can't be used after calling this function!
// DO NOT TRY TO BREAK THIS RULE! OTHERWISE PANIC GUARANTEED!
func wipeLirester(l **tLirester) {
	if l == nil || *l == nil {
		return
	}
	_l := *l
	*l = nil
	_l.isBackgroundCleanup = false // stop background cleanup
	_l.MainLoop.Stop()             // stop main loop
	time.Sleep(1 * time.Second)    // wait background cleanup
	close(_l.chCleanup)            // close cleanup rule's chan
	_l.chCleanup, _l.MainLoop, _l.core = nil, nil, nil
}

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
