// Copyright © 2018. All rights reserved.
// Author: Alice Qio.
// Contacts: <qioalice@gmail.com>.
// License: https://opensource.org/licenses/MIT

package lirester

import (
	"time"

	"github.com/qioalice/devola/core/chat"
	"github.com/qioalice/devola/core/profiler"
)

// Lirester is a part of Sender's Devola SDK module and is an abbreviation of
// [Li]mit [Rest]rictions Avoid[er]
// NOT THREAD (GOROUTINE) SAFETY! (AND IT SHOULD NOT BE SO USED!).
//
// There only two main concepts:
//
// 1. Lirester restricts to send more than N messages per T at all.
// 2. Lirester restricts to send more than N'(i) messages per T'(i) to the same chat,
//    depending by chat type (no more types (i) than core/chat.MaxTypeValue.
//
// Full iteration of Lirester's usage:
//
// 1. Check, whether message can be sent.
//    Pass a chat id and chat type to the Try method and get an answer (bool)
//    to the question "Can I send at least one message to that chat right now?".
//
// 2. Use Approve, if the message was sent successfully.
//    Send message.
//    If it failed, do nothing (method Try just do checks and nothing more).
//    And if it was successful, let Lirester to know it - call Approve passing
//    chat id and chat type.
//
// 3. Decrease all internal counters using Cleanup method (*).
//    Because of this you can get true from Try to the same chat again.
//    Cleanup also performs restarting lirester object (**).
//
// NOTE!
// Methods Approve and Cleanup takes current unix NANO timestamp as 1st argument.
// It made to decrease time.Time() calls.
// BE SURE YOU PASSING UNIX NANO TIMESTAMP, NOT SIMPLE'S.
//
// ----------
//
// *  Call Cleanup method as often as your min(T'(i)), but it is not necessary
//    to call Cleanup after each Approve and even more it's redundant.
//
// ** Restarting is a operation of changing internal constants with which
//    Lirester object was created.
//    You can _REQUEST_ restart using RequestRestartWith method, but really
//    restart will be performed at the next Cleanup call.
//    While restart is requested but not completed, all another operations
//    are unavailable and their methods (Try, Approve) returns a special values.
type Lirester struct {

	// WARNING!
	// ALL TIMESTAMPS IN NANO SECONDS!

	// HOW IT WORKS.
	//
	// As you know from doc above, there are 2 restriction rules: global and chat's.
	//
	// Global rule is fulfilled by MainLoop field - a Golang ticker -
	// a chan that receives some value (*) each t time.
	//
	// And Sender (Devola SDK module) just should perform his main iteration
	// each time when MainLoop has a new value.
	//
	// The 2nd rule - chat's rule is a fulfilled by special small objects that
	// stores counter of sent messages for each chat ID and its types
	// and by small objects "cleanup rules" that are applied by Cleanup method
	// and decreases the counters of first small objects.
	//
	// So:
	// There is a map from chat ID to the chat objects, that stores:
	// 1. n - a value of how many messages have already been sent to that chat
	//    at the current time.
	// 2. typ - a type of chat (up to core/chat.MaxTypeValue).
	// 3. lastUpdated - an unix timestamp when counter has been updated last time.
	//
	// When you call Try it's just takes an associated chat object,
	// and checks if n <= allowable value for type of that chat.
	//
	// When you call Approve it's increases n by one in associated chat object,
	// and creates a new cleanup rule, that will decrease n after now + t,
	// where t is a T'(i) - a value depended on that chat's type.
	//
	// When you call Cleanup it's applies all accomplished cleanup rules
	// and just decreases all counters that must be decreased,
	// but applies only these rules, that must be applied (their time has come).

	// Consts section.
	// Can be overwritten by only RequestRestartWith method.
	consts struct {

		// Delay with which the ticker (MainLoop) will work.
		// Implements 1st Lirester rule, calculated: T/N.
		mainLoopDelay time.Duration

		// The sec value which should pass after last updating of chat object
		// before that chat will be deleted from Lirester.
		chatLifeTime int64

		Ns [chat.MaxTypeValue + 1]uint8 // N'(i) for chat types.
		Ts [chat.MaxTypeValue + 1]int64 // T'(i) for chat types.
	}

	// Main Lirester ticker.
	// Implements 1st Lirester restriction rule: Not more N messages per T at all.
	MainLoop *time.Ticker

	core      map[chat.ID]lirchat // all Lirester chat by its ids
	chCleanup chan cleanupConfig  // all cleanup configs that should be applied

	// Arguments of restart method.
	// Nil means restart isn't requested,
	// Empty slice means restart requested but w/o changing params
	// Otherwise represents a params with which lirester will be restarted.
	restartRequestedWith []interface{}

	prof                    *profiler.Profiler // profiler core's pointer
	profCleanupDecCounters  *profiler.Watcher  // profiler for cleanupDecCounters
	profCleanupDestroyChats *profiler.Watcher  // profiler for cleanupDestroyChats
}

// Try returns true if it is possible to send at least one message to the chat
// with passed its id and type (as combined idt value).
//
// Always returns true if Lirester is disabled.
func (l *Lirester) Try(idt chat.IDT) (isAllow bool) {

	// Lirester can be nil (disabled), but code may still have lirester calls.
	if l == nil {
		return true
	}

	// Don't allow to send if restart is requested
	// (until it is executed in next Cleanup call).
	if l.restartRequestedWith == nil {
		return false
	}

	if ch, ok := l.core[idt.ID()]; ok {
		return ch.n < l.consts.Ns[idt.Type()]
	}

	return true
}

// Approve lets Lirester to know that ONE message is successfully sent to the chat
// with passed its id and type (as combined IDT value).
func (l *Lirester) Approve(now int64, idt chat.IDT) {

	// Note.
	// 2nd arg is chat.IDT (not a just chat.ID), because addCleanup requires chat.Type.

	// WARNING!
	// ALL TIMESTAMPS IN NANO SECONDS!

	// Lirester can be nil (disabled), but code may still have lirester calls.
	// Don't allow to send if restart is requested
	// (until it is executed in next Cleanup call).
	if l == nil || l.restartRequestedWith == nil {
		return
	}

	// if core don't have chat with that chat IDT, ch is an object with zero values.
	ch := l.core[idt.ID()]

	// update values and save an updated object back to the core
	ch.n++
	ch.lastUpdated = now
	l.core[idt.ID()] = ch

	l.addCleanup(now, idt)
}

// Cleanup parses all accumulated cleanup rules and tries to apply it.
// They will be applied only if their time has come.
// Not applied rules will be returned to the cleanup rules queue.
func (l *Lirester) Cleanup(now int64) {

	// WARNING!
	// ALL TIMESTAMPS IN NANO SECONDS!

	// Lirester can be nil (disabled), but code may still have lirester calls.
	if l == nil {
		return
	}

	l.cleanupDecCounters(now)
	l.cleanupDestroyChats(now)

	// Perform restart if it was requested.
	//
	// Check nil (not a zero len), because zero len params means that
	// restart should be but without changing a parameters.
	if l.restartRequestedWith != nil {
		l.restart()
	}
}

// LastUpdated returns the unixnano timestamp when the chat
// with passed its id was updated last time.
//
// It returns -1 if that info can't be obtained.
func (l *Lirester) LastUpdated(id chat.ID) int64 {

	// Lirester can be nil (disabled), but code may still have lirester calls.
	if l == nil || l.restartRequestedWith != nil {
		return -1
	}

	if ch, ok := l.core[id]; ok {
		return ch.lastUpdated
	}
	return -1
}

// RequestRestartWith requests restart Lirester with a passed new Lirester params.
// Restart will be done after next cleanup operation.
// If no one param will be passed, restart will be, but without changing
// the parameters.
func (l *Lirester) RequestRestartWith(params []interface{}) {

	// Lirester can be nil (disabled), but code may still have lirester calls.
	if l == nil {
		return
	}

	if len(params) == 0 {
		params = make([]interface{}, 0, 0)
	}

	l.restartRequestedWith = params
}

// addCleanup creates a new cleanup rule that will be applied after now + delay (D)
// to a chat which id and type are passed (as combined idt value).
//
// D depends from chat's type and its delay value that's computed at the
// Lirester's initialization.
func (l *Lirester) addCleanup(now int64, idt chat.IDT) {

	// WARNING!
	// ALL TIMESTAMPS IN NANO SECONDS!

	l.chCleanup <- makeCleanupConfig(idt.ID(), l.consts.Ts[idt.Type()]+now)
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

	// WARNING!
	// ALL TIMESTAMPS IN NANO SECONDS!

	// Start profiler
	l.profCleanupDecCounters.Start()

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
		if ch, ok := l.core[clcfg.chatID]; ok && ch.n > 0 {
			ch.n, ch.lastUpdated = ch.n-1, now
			l.core[clcfg.chatID] = ch
		}
	}

	// Stop profiler, flush results
	l.profCleanupDecCounters.Stop()
}

// cleanupDestroyChats performs second cleanup type operation:
// Destroying the chats, lifetime of which is over.
//
// It checks an each Lirester chat whether its lifetime is over and if so,
// destroy it.
//
// @profiled: cPA_LiresterCleanupDestroyChats
func (l *Lirester) cleanupDestroyChats(now int64) {

	// WARNING!
	// ALL TIMESTAMPS IN NANO SECONDS!

	// Start profiler
	l.profCleanupDestroyChats.Start()

	// Remove chats, lifetime of which is over
	for chatID, ch := range l.core {

		// if chat's lifetime is over
		if ch.lastUpdated+l.consts.chatLifeTime < now {
			delete(l.core, chatID)
		}
	}

	// Stop profiler, flush results
	l.profCleanupDestroyChats.Stop()
}

// restart performs restart Lirester.
//
// It overwrites used Lirester consts
// by params passed to the Lirester.RequestRestartWith, restarts main loop ticker.
func (l *Lirester) restart() {

	// WARNING!
	// ALL TIMESTAMPS IN NANO SECONDS!

	// TODO: Recreate cleanup rules queue (but APPLY all already accumulated).
	// close(l.chCleanup)

	// Stop main loop ticker, discard cleanup rules queue.
	// restart performs in the end of Lirester.Cleanup, what means that
	// some cleanup rules already applied.
	if l.MainLoop != nil {
		l.MainLoop.Stop()
	}

	// The rest cleanup rules (which were discard) will be emulated
	// by force zeroing rest of counters.
	// for _, chat := range l.core {
	// 	chat.setHowMuch(0)
	// }

	// Apply new values.
	for _, param := range l.restartRequestedWith {

		if typedParam, ok := param.(Param); ok && typedParam != nil {
			typedParam(l)

		} else if paramGen, ok := param.(func() Param); ok && paramGen != nil {
			if typedParam := paramGen(); typedParam != nil {
				typedParam(l)
			}

		} else if prof, ok := param.(*profiler.Profiler); ok && prof != nil {
			l.prof = prof
		}
	}

	// Probably create profiler's watchers.
	l.profCleanupDecCounters = l.prof.For("Lirester.Cleanup.DecCounters")
	l.profCleanupDestroyChats = l.prof.For("Lirester.Cleanup.DestroyChats")

	// Create cleanup rules queue, start main loop ticker.
	// l.chCleanup = make(chan tLiresterCleanupConfig, l.consts.??)

	// Create and start new main loop ticker.
	l.MainLoop = time.NewTicker(l.consts.mainLoopDelay)

	// Lirester restarted.
	l.restartRequestedWith = nil
}

// MakeLirester creates a new Lirester object, the passed params will be applied to.
// Also starts main Lirester loop (MainLoop field).
func MakeLirester(params []interface{}) *Lirester {

	// WARNING!
	// ALL TIMESTAMPS IN NANO SECONDS!

	var l Lirester

	// Apply default lirester params.
	// It'll be overwritten later by passed lirester params.
	l.consts.mainLoopDelay = 1 * time.Second / 1000
	l.consts.chatLifeTime = int64(10 * time.Minute) // the same as .Nanoseconds() call

	// Apply default lirester N'(i) and T'(i) values
	for i, n := 0, len(l.consts.Ns); i < n; i++ {
		l.consts.Ns[i] = 100
		l.consts.Ts[i] = int64(1 * time.Minute)
	}

	l.restartRequestedWith = params
	l.restart()

	// Allocate mem for main core lirester and cleanup queue (chan).
	l.core = make(map[chat.ID]lirchat)
	l.chCleanup = make(chan cleanupConfig, 16384) // TODO: smaller?

	return &l
}
