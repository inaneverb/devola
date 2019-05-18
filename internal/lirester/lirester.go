// Copyright © 2018. All rights reserved.
// Author: Alice Qio.
// Contacts: <qioalice@gmail.com>.
// License: https://opensource.org/licenses/MIT

package lirester

import (
	"time"

	"../chat"
	"../profiler"
)

// Lirester is an abbreviation of
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
// 1. Lirester objects starts and contains a Golang ticker.
// This ticker will tick every time (T) when it allows server to perform a few (N)
// Telegram API requests of sending messages.
// It means that only N per T operations are allowed (avoids 1st restriction).
//
// 2. Lirester object has a map from Telegram chat id (tChatID)
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
// 1. Start from Lirester.Try.
// Pass a chat id to that method and get an answer (bool) to the question
// "Can I send at least one message to that chat right now?"
//
// 2. Use Lirester.Approve, if the message was sent successfully.
// Send message. If it failed, do nothing (method Try doesn't inc the counters).
// And if it was successful, let Lirester to know it.
// Pass a chat id to that method and a flag this chat is with user or is group.
//
// 3. Decrease counters using Lirester.Cleanup.
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
// More info: chat, cleanupConfig, Sender.
type Lirester struct {

	// Profiler object with which profiling operations are performed.
	prof *profiler.Profiler

	// Consts section.
	//
	// After Lirester object is created and started,
	// it can be overwritten, but it is required to complete stop Lirester,
	// stop Sender main loop (depends by Lirester main ticker) and only
	// then these consts can be overwritten.
	// Otherwise the behaviour is undefined.
	consts struct {

		// Main loop ticker delay.
		// T / N of 1st Telegram restriction.
		// More info: Lirester "How it works" section, p.1.
		mainLoopDelay time.Duration

		// The sec value which should pass after last chat update before
		// that chat will be deleted from Lirester.
		chatLifeTime int64

		// X of 2nd Telegram restiction for chats with user.
		userChatN uint8
		// T' of 2nd Telegram restriction for chats with user.
		userChatT int64

		// X of 2nd Telegram restriction for group chats.
		groupChatN uint8
		// T' of 2nd Telegram restriction for group chats.
		groupChatT int64

		// Profiler enable flag for profiling Lirester.cleanupDecCounters.
		enProfCleanupDecCounters bool

		// Profiler enable flag for profiling Lirester.cleanupDestroyChats.
		enProfCleanupDestroyChats bool
	}

	// Main Lirester ticker.
	// You should create Lirester object and then try to perform sending
	// message request at that time which this ticker allows it.
	MainLoop *time.Ticker

	// Here stored all Lirester chats by its chat ids.
	core map[chat.ID]*lirchat

	// Cleanup rules queue.
	chCleanup chan cleanupConfig

	// Arguments of restart method.
	// Nil means restart isn't requested,
	// Empty slice means restart requested but w/o changing params
	// Otherwise represents a params with which lirester will be restarted.
	restartRequestedWith []interface{}
}

// Predefined default values of some important Lirester constants.
const (

	// Default main loop Lirester ticker delay.
	// T / N of 1st Telegram restriction.
	// More info: Lirester "How it works" section, p.1.
	cMainLoopDelay = 1 * time.Second / 30

	// Default Lirester chat lifetime.
	// After this time the process of chat destroying will be started.
	cChatLifetime = int64(10 * time.Minute)

	// Default X of 2nd Telegram restriction for chats with user.
	// More info: Lirester "How it works" section, p.2.
	cUserChatN = uint8(1)

	// Default T' of 2nd Telegram restriction for chats with user.
	// More info: Lirester "How it works" section, p.2.
	cUserChatT = int64(1 * time.Second)

	// Default X of 2nd Telegram restriction for group chats.
	// More info: Lirester "How it works" section, p.2.
	cGroupChatN = uint8(20)

	// Default T' of 2nd Telegram restriction for group chats.
	// More info: Lirester "How it works" section, p.2.
	cGroupChatT = int64(1 * time.Minute)

	// Length of Lirester cleanup rule's Golang chan.
	cCleanupChanLen = 16384 // 2^14
)

// Profiler action constants.
const (

	// Profiler action name for Lirester.CleanupDecCounters.
	cPACleanupDecCounters profiler.Action = "Lirester.Cleanup.DecreasingCounters"

	// Profiler action name for Lirester.CleanupDestroyChats.
	cPACleanupDestroyChats profiler.Action = "Lirester.Cleanup.DestroyChats"
)

// Try checks whether some message can be send to the chat with the passed chat id
// and returns true if it is possible.
func (l *Lirester) Try(now int64, chatID chat.ID) (isAllow bool) {

	if len(l.restartRequestedWith) != 0 {
		return false
	}

	var chat = l.getChat(now, chatID)
	var alreadySent = chat.howMuch()

	allowForUser := chat.isUser() && alreadySent < l.consts.userChatN
	allowForGroup := chat.isGroup() && alreadySent < l.consts.groupChatN

	return alreadySent == 0 || allowForUser || allowForGroup
}

// Approve lets Lirester to know that ONE message to the chat with the passed
// chat id is successfully sent.
//
// WARNING!
// Be sure that before using Approve for some chat id, Try returns true for the
// same chat id. Otherwise the behaviour is undefined.
func (l *Lirester) Approve(now int64, chatID chat.ID, isUser bool) {

	if len(l.restartRequestedWith) != 0 {
		return
	}

	l.getChat(now, chatID).setType(isUser).incHowMuch(1).setLastUpdated(now)
	l.addCleanup(now, chatID, isUser)
}

// Cleanup parses all accumulated cleanup rules and tries to apply it.
// They will be applied only if their time has come.
// Not applied rules will be returned to the cleanup rules queue.
func (l *Lirester) Cleanup(now int64) {

	l.cleanupDecCounters(now)
	l.cleanupDestroyChats(now)

	// Perform restart if it was requested.
	//
	// Check nil (not a zero len), because zero len params means that
	// restart should be but without changing a parameters.
	if len(l.restartRequestedWith) != 0 {
		l.restart()
	}
}

// LastUpdated returns the unixnano timestamp when the chat with the passed chat id
// was updated last time.
// If an info about that chat isn't specified in Lirester, -1 is returned.
func (l *Lirester) LastUpdated(chatID chat.ID) int64 {

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
func (l *Lirester) RequestRestartWith(params []interface{}) {

	if len(params) == 0 {
		params = make([]interface{}, 0, 0)
	}

	l.restartRequestedWith = params
}

// isLifetimeOver returns true if the lifetime of the Lirester chat is over.
func (l *Lirester) isLifetimeOver(now int64, chat *lirchat) bool {
	return chat.lastUpdated+l.consts.chatLifeTime < now
}

// findChat returns a chat object associated with the passed chat id.
// If a required chat is not exists in Lirester, nil is returned.
func (l *Lirester) findChat(chatID chat.ID) *lirchat {
	return l.core[chatID]
}

// getChat returns a chat object associated with the passed chat id.
// If a required chat is not exists in Lirester, it will be created
// and then returned.
// It guarantees, that getChat always returns not nil object.
func (l *Lirester) getChat(now int64, chatID chat.ID) *lirchat {

	var chat *lirchat

	if chat = l.findChat(chatID); chat == nil {
		chat = makeChat().setLastUpdated(now)
		l.core[chatID] = chat
	}

	return chat
}

// addCleanup creates a new cleanup rule for a chat with specified chat id
// that will be applied after now + delay (D).
// D depends from chat type (with user or group chat) and internal constants.
func (l *Lirester) addCleanup(now int64, chatID chat.ID, isUser bool) {

	var delay = l.consts.userChatT
	if !isUser {
		delay = l.consts.groupChatT
	}

	l.chCleanup <- makeCleanupConfig(chatID, delay+now)
}

// cleanupDecCounters performs first cleanup type operation:
// Decreasing the sent messages' counters.
//
// It checks an each cleanup config (rule) from the special queue
// (Golang channel) whether time of that config applying has come.
// If so, then that config will be applied and its counter will be decreased.
// Otherwise config will be returned to queue.
//
// @profiled: cPA_LiresterCleanupDecCounters
func (l *Lirester) cleanupDecCounters(now int64) {

	// Start profiler
	var watcher = l.prof.WatchIf(l.consts.enProfCleanupDecCounters).Start()

	var clcfg cleanupConfig
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
		chat := l.findChat(clcfg.chatID)

		// If counter isn't equal to zero, apply cleanup (decrease counter)
		// todo: Remove "if", because if cleanup config is, counter must be != 0
		if chat.howMuch() > 0 {
			chat.decHowMuch(1).setLastUpdated(now)
		}
	}

	// Stop profiler, flush results
	watcher.Stop().For(cPACleanupDecCounters)
}

// cleanupDestroyChats performs second cleanup type operation:
// Destroying the chats, lifetime of which is over.
//
// It checks an each Lirester chat whether its lifetime is over and if so,
// destroy it.
//
// @profiled: cPA_LiresterCleanupDestroyChats
func (l *Lirester) cleanupDestroyChats(now int64) {

	// Start profiler
	var watcher = l.prof.WatchIf(l.consts.enProfCleanupDestroyChats).Start()

	// Remove chats, lifetime of which is over
	for chatID, chat := range l.core {

		if l.isLifetimeOver(now, chat) {
			delete(l.core, chatID)
		}
	}

	// Stop profiler, flush results
	watcher.Stop().For(cPACleanupDestroyChats)
}

// restart performs restart Lirester.
//
// It overwrites used Lirester consts
// by params passed to the Lirester.RequestRestartWith, restarts main loop ticker.
func (l *Lirester) restart() {

	// TODO: Recreate cleanup rules queue (but APPLY all already accumulated).
	// close(l.chCleanup)

	// Stop main loop ticker, discard cleanup rules queue.
	// restart performs in the end of Lirester.Cleanup, what means that
	// some cleanup rules already applied.
	l.MainLoop.Stop()

	// The rest cleanup rules (which were discard) will be emulated
	// by force zeroing rest of counters.
	// for _, chat := range l.core {
	// 	chat.setHowMuch(0)
	// }

	// Apply new values.
	l.applyParams(l.restartRequestedWith)

	// Fetch probably new values that can't be changed using params.
	l.consts.enProfCleanupDecCounters = l.prof.IsEnabledFor(cPACleanupDecCounters)
	l.consts.enProfCleanupDestroyChats = l.prof.IsEnabledFor(cPACleanupDestroyChats)

	// Create cleanup rules queue, start main loop ticker.
	// l.chCleanup = make(chan tLiresterCleanupConfig, l.consts.??)

	// Create and start new main loop ticker.
	l.MainLoop = time.NewTicker(l.consts.mainLoopDelay)

	// Lirester restarted.
	l.restartRequestedWith = nil
}

// applyParams applies each param from params slice to the current Lirester.
// It overwrites alreay applied parameters.
func (l *Lirester) applyParams(params []interface{}) {
	// TODO: Maybe something more must be here?

	for _, param := range params {

		if typedParam, ok := param.(Param); ok && typedParam != nil {
			typedParam(l)

		} else if paramGen, ok := param.(func() Param); ok && paramGen != nil {
			if typedParam := paramGen(); typedParam != nil {
				typedParam(l)
			}
		}
	}
}

// MakeLirester creates a new Lirester object, the passed params will be applied to.
// Also starts main Lirester loop (MainLoop field).
func MakeLirester(params []interface{}) *Lirester {

	var l Lirester

	// Apply default lirester params.
	// It'll be overwritten later by passed lirester params.
	l.consts.mainLoopDelay = cMainLoopDelay
	l.consts.chatLifeTime = cChatLifetime
	l.consts.userChatN = cUserChatN
	l.consts.userChatT = cUserChatT
	l.consts.groupChatN = cGroupChatN
	l.consts.groupChatT = cGroupChatT

	l.applyParams(params)

	// Allocate mem for main core lirester and cleanup queue (chan).
	l.core = make(map[chat.ID]*lirchat)
	l.chCleanup = make(chan cleanupConfig, cCleanupChanLen)

	// Create and start main loop ticker.
	l.MainLoop = time.NewTicker(l.consts.mainLoopDelay)

	return &l
}
