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

import (
	"time"
)

// -- LiRester --
// Telegram [Li]mit [Rest]rictions Avoid[er]
// NOT THREAD (GOROUTINE) SAFETY! (AND IT SHOULD NOT BE SO USED!)
//
//
// Telegram restrictions (as of 23 Apr 2019).
// https://core.telegram.org/bots/faq#my-bot-is-hitting-limits-how-do-i-avoid-this
//
// 1. Only 30 TOTAL requests from bot per second are allowed.
// 2. Only 20 messages from bot to the group chat per minute are allowed.
// 3. Only 1 message per second to the chat are allowed.
//
//
// Lirester concept.
//
// 1. Check, whether message can be sent (Try!)
// 2. Confirm, that message is really sent and you want to be protected by
// Lirester from Telegram limits overlow (Approve!)
// 3. Just wait when Lirester will allow you to send message again (Cleanup!)
//
//
// How it works.
//
// 1. tLirester objects starts and contains a Golang ticker.
// This ticker will tick every time (T) when it allows server to perform a few (N)
// Telegram API requests of sending messages.
// It means that only N per T operations are allowed (avoids 1st restriction).
//
// 2. tLirester object has a map from Telegram chat id (tChatID)
// to the lirester chat object (tLiresterChat), which contains already sent
// messages counter.
// If this counter decreases every some time (T'), and when this counter reaches
// a certain MAX value (X), the message won't be sent to the linked Telegram chat,
// it means, that only X per T' messages are allowed for linked chat.
// (avoids 2nd and 3rd restriction).
//
// 2.1. You should decrease the counters were mentioned above.
// Read more about it below.
//
//
// How to use.
//
// 1. Start from tLirester.Try.
// Pass a chat id to that method and get an answer (bool) to the question
// "Can I send at least one message to that chat right now?"
//
// 2. Use tLirester.Approve, if the message was sent successfully.
// Send message. If it failed, do nothing (method Try doesn't inc the counters).
// And if it was successful, let Lirester to know it.
// Pass a chat id to that method and a flag this chat is with user or is group.
//
// 3. Decrease counters using tLirester.Cleanup.
// After each Lirester main ticker iteration (see "How it works" section, p. 1)
// call that method.
// Approve creates a special "decreasing config" to internal counter that was
// increased and pass it to the special queue.
// Cleanup parses that queue and handles these configs.
//
//
// FAQ.
//
// Q: What is the isUser flag for?
// A: Because Telegram has a different restrictions for chats 'user-bot' and
// 'bot-in-the-chat-group'.
//
//
// NOTE!
// Almost all functions and methods of Lirester takes a Unixnano timestamp
// as first argument.
// It made to reduce time.Now() calls.
//
// More info: tLiresterChat, tLiresterCleanupConfig, tSender.
type tLirester struct {

	// TBot object, this tLirester object associated with.
	parent *TBot

	// Consts section.
	//
	// After starts tLirester object created and started,
	// it can be overwritten, but required to complete stop Lirester,
	// stop tSender main loop (depends by tLirester main ticker) and only
	// then these consts can be overwritten.
	// Otherwise the behaviour is undefined.
	consts struct {

		// todo: rename userChatMsgNumPerIter, userChatCleanupDelay as its consts
		// todo: rename groupChatMsgNumPerIter, groupChatCleanupDelay as its consts

		// Main loop ticker delay.
		// T / N of 1st Telegram restriction.
		// More info: tLirester "How it works" section, p.1.
		mainLoopDelay time.Duration

		// The sec value which should pass after last chat update before
		// that chat will be deleted from Lirester.
		chatLifeTime int64

		// X of 2nd Telegram restiction for chats with user.
		userChatMsgNumPerIter uint8
		// T' of 2nd Telegram restriction for chats with user.
		userChatCleanupDelay int64

		// X of 2nd Telegram restriction for group chats.
		groupChatMsgNumPerIter uint8
		// T' of 2nd Telegram restriction for group chats.
		groupChatCleanupDelay int64

		// Profiler enable flag for profiling tLirester.cleanupDecCounters.
		enP_CleanupDecCounters bool

		// Profiler enable flag for profiling tLirester.cleanupDestroyChats.
		enP_CleanupDestroyChats bool
	}

	// Main Lirester ticker.
	// You should create tLirester object and then try to perform sending
	// message request at that time which this ticker allows it.
	MainLoop *time.Ticker

	// Here stored all Lirester chats by its chat ids.
	core map[tChatID]*tLiresterChat

	// Cleanup rules queue.
	chCleanup chan tLiresterCleanupConfig

	restartRequested []tLiresterParam
}

// Predefined default values of some important Lirester constants.
// These values used while constructing Lirester object, can be overwritten
// using tLiresterParam parameters or tLirester.RestartWith method.
const (

	// Default main loop Lirester ticker delay.
	// T / N of 1st Telegram restriction.
	// More info: tLirester "How it works" section, p.1.
	cLiresterMainLoopDelay = 1 * time.Second / 30

	// Default Lirester chat lifetime.
	// After this time the process of chat destroying will be started.
	cLiresterChatLifetime = int64(10 * time.Minute)

	// Default X of 2nd Telegram restriction for chats with user.
	// More info: tLirester "How it works" section, p.2.
	cLiresterUserChatN = uint8(1)

	// Default T' of 2nd Telegram restriction for chats with user.
	// More info: tLirester "How it works" section, p.2.
	cLiresterUserChatT = int64(1 * time.Second)

	// Default X of 2nd Telegram restriction for group chats.
	// More info: tLirester "How it works" section, p.2.
	cLiresterGroupChatN = uint8(20)

	// Default T' of 2nd Telegram restriction for group chats.
	// More info: tLirester "How it works" section, p.2.
	cLiresterGroupChatT = int64(1 * time.Minute)

	// Length of Lirester cleanup rule's Golang chan.
	cLiresterCleanupChanLen = 16384 // 2^14
)

// Profiler action constants.
const (

	// Profiler action name for tLirester.CleanupDecCounters.
	cPA_LiresterCleanupDecCounters tProfilerAction = "tLirester.Cleanup.DecreasingCounters"

	// Profiler action name for tLirester.CleanupDestroyChats.
	cPA_LiresterCleanupDestroyChats tProfilerAction = "tLirester.Cleanup.DestroyChats"
)

// Try checks whether some message can be send to the chat with the passed chat id
// and returns true if it is possible.
func (l *tLirester) Try(now int64, chatID tChatID) (isAllow bool) {

	if l.restartRequested {
		return false
	}

	var chat = l.getChat(now, chatID)
	var alreadySent = chat.howMuch()

	allowForUser := chat.isUser() && alreadySent < l.consts.userChatMsgNumPerIter
	allowForGroup := chat.isGroup() && alreadySent < l.consts.groupChatMsgNumPerIter

	return alreadySent == 0 || allowForUser || allowForGroup
}

// Approve lets Lirester to know that ONE message to the chat with the passed
// chat id is successfully sent.
//
// WARNING!
// Be sure that before using Approve for some chat id, Try returns true for the
// same chat id. Otherwise the behaviour is undefined.
func (l *tLirester) Approve(now int64, chatID tChatID, isUser bool) {

	if l.restartRequested {
		return
	}

	l.getChat(now, chatID).setType(isUser).incHowMuch(1).setLastUpdated(now)
	l.addCleanup(now, chatID, isUser)
}

// Cleanup parses all accumulated cleanup rules and tries to apply it.
// They will be applied only if their time has come.
// Not applied rules will be returned to the cleanup rules queue.
func (l *tLirester) Cleanup(now int64) {

	l.cleanupDecCounters(now)
	l.cleanupDestroyChats(now)

	// Perform restart if it was requested.
	//
	// Check nil, not a zero len, because zero len params means that
	// restart should be but without changing a parameters.
	if l.restartRequested != nil {
		l.restart()
	}
}

// LastUpdated returns the unixnano timestamp when the chat with the passed chat id
// was updated last time.
// If an info about that chat isn't specified in Lirester, -1 is returned.
func (l *tLirester) LastUpdated(chatID int64) int64 {

	if chat := l.findChat(chatID); chat != nil {
		return chat.lastUpdated
	}
	return -1
}

// RequestRestartWith requests restart Lirester with a passed new Lirester params.
// It guarantees, that restart will be done, but NOT NOW!
//
// Generally, restart will be done after next cleanup operation.
// If no one param will be passed, restart will be, but without changing
// the parameters.
func (l *tLirester) RequestRestartWith(params ...tLiresterParam) {

	if len(params) == 0 {
		params = make([]tLiresterParam, 0, 0)
	}

	l.restartRequested = params
}

// isLifetimeOver returns true if the lifetime of the Lirester chat is over.
func (l *tLirester) isLifetimeOver(now int64, chat *tLiresterChat) bool {

	return chat.lastUpdated+l.consts.chatLifeTime < now
}

// findChat returns a chat object associated with the passed chat id.
// If a required chat is not exists in Lirester, nil is returned.
func (l *tLirester) findChat(chatID tChatID) *tLiresterChat {

	return l.core[chatID]
}

// getChat returns a chat object associated with the passed chat id.
// If a required chat is not exists in Lirester, it will be created
// and then returned.
// It guarantees, that getChat always returns not nil object.
func (l *tLirester) getChat(now int64, chatID tChatID) *tLiresterChat {

	var chat *tLiresterChat

	if chat = l.findChat(chatID); chat == nil {
		chat = makeLiresterChat().setLastUpdated(now)
		l.core[chatID] = chat
	}

	return chat
}

// addCleanup creates a new cleanup rule for a chat with specified chat id
// that will be applied after now + delay (D).
// D depends from chat type (with user or group chat) and internal constants.
func (l *tLirester) addCleanup(now int64, chatID tChatID, isUser bool) {

	var delay = l.consts.userChatCleanupDelay
	if !isUser {
		delay = l.consts.groupChatCleanupDelay
	}

	l.chCleanup <- makeLiresterCleanupConfig(chatID, delay+now)
}

// @profiled: cPA_LiresterCleanupDecCounters
//
// cleanupDecCounters performs first cleanup type operation:
// Decreasing the sent messages' counters.
//
// It checks an each cleanup config (rule) from the special queue
// (Golang channel) whether time of that config applying has come.
// If so, then that config will be applied and its counter will be decreased.
// Otherwise config will be returned to queue.
func (l *tLirester) cleanupDecCounters(now int64) {

	// Start profiler
	var watcher = l.parent.prof.WatchIf(l.consts.enP_CleanupDecCounters).Start()

	var clcfg tLiresterCleanupConfig
	var isQueueOk bool

	// Perform decreasing counters
	for i, n := 0, len(l.chCleanup); i < n; i++ {

		// Get one rule.
		// If queue is not ok (chan is closed), do nothing, go out.
		if clcfg, isQueueOk = <-l.chCleanup; !isQueueOk {
			break
		}

		// If cleanup config must be executed later than now,
		// just return it to the cleanup configs channel and go to next config
		if clcfg.when > now {
			l.chCleanup <- clcfg
			continue
		}

		// If cleanup config for some chat is existed, the chat in map
		// MUST BE.
		// If nil will be returned here, it is SDK error and it should be fixed.
		chat := l.findChat(clcfg.who)

		// If counter isn't equal to zero, apply cleanup (decrease counter)
		// todo: Remove "if", because if cleanup config is, counter must be != 0
		if chat.howMuch() > 0 {
			chat.decHowMuch(1).setLastUpdated(now)
		}
	}

	// Stop profiler, flush results
	watcher.Stop().For(cPA_LiresterCleanupDecCounters)
}

// @profiled: cPA_LiresterCleanupDestroyChats
//
// cleanupDestroyChats performs second cleanup type operation:
// Destroying the chats, lifetime of which is over.
//
// It checks an each Lirester chat whether its lifetime is over and if so,
// destroy it.
func (l *tLirester) cleanupDestroyChats(now int64) {

	// Start profiler
	var watcher = l.parent.prof.WatchIf(l.consts.enP_CleanupDestroyChats).Start()

	// Remove chats, lifetime of which is over
	for chatID, chat := range l.core {

		if l.isLifetimeOver(now, chat) {
			delete(l.core, chatID)
		}
	}

	// Stop profiler, flush results
	watcher.Stop().For(cPA_LiresterCleanupDestroyChats)
}

// restart performs restart Lirester.
//
// It overwrites used Lirester consts
// by params passed to the tLirester.RequestRestartWith, restarts main loop ticker.
func (l *tLirester) restart() {

	l.RequestRestartWith(params)
	// todo: Add recreate cleanup rules queue.

	// Stop main loop ticker, discard cleanup rules queue.
	// restart performs in the end of tLirester.Cleanup, what means that
	// some cleanup rules already applied.
	l.MainLoop.Stop()
	// close(l.chCleanup)

	// The rest cleanup rules (which were discard) will be emulated
	// by force zeroing rest of counters.
	// for _, chat := range l.core {
	// 	chat.setHowMuch(0)
	// }

	// Apply new values.
	l.applyParams(params)

	// Fetch probably new values that can't be changed using params.
	l.consts.enP_CleanupDecCounters = l.parent.prof.isEnabledFor(cPA_LiresterCleanupDecCounters)
	l.consts.enP_CleanupDestroyChats = l.parent.prof.isEnabledFor(cPA_LiresterCleanupDestroyChats)

	// Create cleanup rules queue, start main loop ticker.
	// l.chCleanup = make(chan tLiresterCleanupConfig, l.consts.??)

	// Create and start new main loop ticker.
	l.MainLoop = time.NewTicker(l.consts.mainLoopDelay)

	// Lirester restarted.
	l.restartRequested = nil
}

// makeLirester creates a new tLirester object, the passed params will be applied to.
// Also starts main tLirester loop (MainLoop field).
//
// Parameters:
//
// Only tLiresterParam type of arguments is allowed.
// Values of any other type will be ignored.
func makeLirester(params ...interface{}) *tLirester {

	var l tLirester

	// Apply default lirester params.
	// It'll be overwritten later by passed lirester params.
	l.consts.mainLoopDelay = cLiresterMainLoopDelay
	l.consts.chatLifeTime = cLiresterChatLifetime
	l.consts.userChatMsgNumPerIter = cLiresterUserChatN
	l.consts.userChatCleanupDelay = cLiresterUserChatT
	l.consts.groupChatMsgNumPerIter = cLiresterGroupChatN
	l.consts.groupChatCleanupDelay = cLiresterGroupChatT

	// Apply passed lirester params.
	for _, param := range params {
		if param, ok := param.(tLiresterParam); ok && param != nil {
			param(&l)
		}
	}

	// Allocate mem for main core lirester and cleanup queue (chan).
	l.core = make(map[tChatID]*tLiresterChat)
	l.chCleanup = make(chan tLiresterCleanupConfig, cLiresterCleanupChanLen)

	// Create and start main loop ticker.
	l.MainLoop = time.NewTicker(l.consts.mainLoopDelay)

	return &l
}

// destroyLirester stops a Lirester background cleanup, stops a Lirester's
// main loop ticker, closes cleanup queue (chan), and nulls the variable in which
// pointer to the that tLirester object is store.
//
// NOTE!
// tLirester object must be passed by double pointer for avoiding
// next tLirester object using in the callable code after calling this function.
//
// WARNING!
// tLirester object can't be used after calling this function!
// DO NOT TRY TO BREAK THIS RULE!
// OTHERWISE PANIC GUARANTEED!
func destroyLirester(l **tLirester) {

	// todo: What is this function for? Remove.

	if l == nil || *l == nil {
		return
	}

	_l := *l
	*l = nil

	_l.MainLoop.Stop()          // stop main loop
	time.Sleep(1 * time.Second) // wait background cleanup
	close(_l.chCleanup)         // close cleanup rule's chan
	_l.chCleanup, _l.MainLoop, _l.core = nil, nil, nil
}
