// Copyright © 2018. All rights reserved.
// Author: Alice Qio.
// Contacts: <qioalice@gmail.com>.
// License: https://opensource.org/licenses/MIT

package sender

import (
	"math"
	"sync"
	"time"
	"unsafe"

	"go.uber.org/zap"

	"github.com/qioalice/devola/core/chat"
	"github.com/qioalice/devola/core/logger"
	"github.com/qioalice/devola/core/sys/deque"
	"github.com/qioalice/devola/modules/bridge"
)

// Sender is a module of Devola SDK that performs async sending responses
// to the occurred events and calls registered callbacks after.
type Sender struct {

	// HOW IT WORKS.
	//
	// Sender works with AT LEAST 3 goroutines or more.
	// But there is only 3 types of goroutines:
	// - (A) User goroutine in which SendAsync method is called
	//       (may be as more as you want),
	// - (B) Goroutine of real sending responses using backend's API (only one),
	// - (C) Goroutine(s) of calling onSuccess, onError sending callbacks
	//       (may be as more as required).
	//
	// The core of Sender is preparedResponses field.
	// It is map of core/chat.IDT to its chat queues (see core/deque.DequePtr docs).
	//
	// In A thread(s) you call any method of SendAsync's set,
	// passing Tusent that should to be sent.
	// It moves it to the undisposedResponses (op protected by mu).
	//
	// Main loop (B thread), moves objects from undisposedResponses
	// to preparedResponses with corresponding rules (op protected by mu).
	//
	// Then (next in the iteration of main loop) (still B thread),
	// the random non-empty chat queue from preparedResponses is selected
	// and if lirester enabled and allows to send message to the chat
	// of selected queue, Tusents are extracted from that queue
	// and operation of sending is performed.
	//
	// Sending operation may complete successfully or with error.
	// Depends on that, the method MakeSuccess or MakeError is called on
	// current core/tusent.Tusent object (still B thread).
	//
	// And then that object is sent to the handledResponses queue.
	// These objects will be processed by C thread(s),
	// which calls onSuccess and onError callbacks
	// and finishes transactions (if its need) using core/bridge.Bridge module.

	// WHAT IS A LIRESTER?
	//
	// Earlier Lirester was a separate class of Devola SDK.
	// Since 21 Jun 2019 code review it's a part of Sender module.
	//
	// Lirester is an abbreviation of [Li]mit [Rest]rictions Avoid[er]
	// and earlier it was intended to comply with restrictions of Telegram API
	// (when Devola was in development as only Telegram-oriented framework).
	//
	// Now it's also used so in Telegram backend, but also it's improved
	// and can be used for another purposes.
	//
	// First, you can restrict how many times (N) in a certain period (T)
	// main loop iteration (method iter) should be called
	// (so you adjust the call frequency of backend-depended sending method).
	// This delay is consts.mainLoopDelay and it's T/N.
	//
	// Second, you can restrict how many (N'(i)) backend-depended sending
	// operations in a certain period (T'(i)) to the same chat should be performed,
	// and you can adjust a different values (N, T) for different chat types
	// (number of types is up to core/chat.MaxTypeValue).
	//
	// It may be useful at backend development for a different purposes.
	// Do you need an API call restrictions like Telegram backend?
	// Do you want embed that feature to backend and make different chats
	// like chat with admins, chat with employees, etc?
	// Do you need lirester at all? You can disable it using Sender parameter.

	// HOW LIRESTER IS EMBEDDED TO SENDER?
	//
	// Keep in mind, Lirester works ONLY in B thread!
	// It does not affect neither A thread routines nor C's.
	//
	// 1. Lirester is disabled if consts.disableLirester is true.
	//
	// 2. Each new iteration begins, when consts.mainLoopDelay time (in ns)
	//    has been passed since prev iteration.
	//    Value consts.mainLoopDelay is 0 if lirester is disabled,
	//    and iteration's code is simplified (no time checks at all).
	//
	// 3. When you trying to  send a Tusent to specific chat with type t,
	//    it checks whether no more than consts.Ns[t] messages sent atm to that chat.
	//    If it's so, it sends Tusent (using backend API) with
	//    increasing internal counter of sent messages to that chat by 1,
	//    and registers cleanup rule (which will decrease that counter by 1)
	//    which will be applied after consts.Ts[t] ns.
	//
	// 4. If you want to send message to chat with filled counter of sent messages
	//    (is equal to consts.Ns[t], where t is type of that chat),
	//    that message will be deferring using chat queues (and applied as soon,
	//    as it will be allowed by counter).

	bridge *bridge.Bridge

	consts struct {
		disableLirester bool

		cqCap            uint8  // 2^cqCap is the capacity of each chat q in preparedResponses
		cqReuseBufferLen uint16 // len of buffer of chat qs that will be reused
		cqLifetime       int64  // delay (ns) since last updated chat q before it'll be del

		mainLoopDelay time.Duration // iteration's delay of thread B main loop

		Ns [chat.MaxTypeValue + 1]uint8 // N'(i) for chat types.
		Ts [chat.MaxTypeValue + 1]int64 // T'(i) for chat types.

		retryAttempts       int8 // num of attempts to send Tusent before raise an error
		retryAttemptsInfMax int8 // max of "infinity" num of sending attempts

		undisposedResponsesLenExp uint8 // 2^that is len of undisposedResponses
		handledResponsesLenExp    uint8 // 2^that is len of handledResponses

		threadCGN uint8 // goroutine numbers of thread C

		cleanupRulesLenExp uint8 // 2^that is len of cleanupRules todo: param
	}

	// mutex for protecting:
	// undisposedResponses, handledResponses, restartRequestedWith, isStopped.
	mu sync.Mutex

	wg sync.WaitGroup // waiter for threads B, C to complete

	undisposedResponses *deque.DequePtr // it's static queue (A thread(s))

	preparedResponses map[chat.IDT]*chatQueue // B thread's core (and Sender's at all)
	cqReuseBuffer     []*chatQueue            // buffer of chat qs that may be reused

	handledResponses *deque.DequePtr // accessed in B,C threads

	cleanupRules *deque.Deque128 // queue of lirester cleanup rules

	restartRequestedWith []interface{} // not nil (empty or not) -> restart is requested
	isStopped            bool          // true if completely stop or restart is requested
}

// chatQueue is a Tusent's deque and associated lirester values:
// n - how much messages has been sent at this time to this chat?
// la - when this chat has been updated last time?
//
// Fully controlled (constructing, changing) by Sender's cqGet, cqDel, iter methods.
type chatQueue struct {
	q  deque.DequePtr
	n  uint8
	la int64
}

// SendAsync moves passed Tusent config to the chat queue whose IDT
// is in passed config and does it asynchronously.
func (s *Sender) SendAsync(cfg *Tusent) {

	if cfg == nil {
		return
	}

	s.mu.Lock()

	if !s.isStopped {
		s.undisposedResponses.PushBack(cfg)
	}
	s.mu.Unlock()
}

// SendAsync2 moves both passed Tusent configs to the chat queues whose IDTs
// are in passed configs and does it asynchronously.
func (s *Sender) SendAsync2(cfg1, cfg2 *Tusent) {

	if cfg1 == nil && cfg2 == nil {
		return
	}

	s.mu.Lock()
	if !s.isStopped {
		if cfg1 != nil {
			s.undisposedResponses.PushBack(cfg1)
		}
		if cfg2 != nil {
			s.undisposedResponses.PushBack(cfg2)
		}
	}
	s.mu.Unlock()
}

// SendAsync3 moves all three passed Tusent configs to the chat queues whose IDTs
// are in passed configs and does it asynchronously.
func (s *Sender) SendAsync3(cfg1, cfg2, cfg3 *Tusent) {

	if cfg1 == nil && cfg2 == nil && cfg3 == nil {
		return
	}

	s.mu.Lock()
	if !s.isStopped {
		if cfg1 != nil {
			s.undisposedResponses.PushBack(cfg1)
		}
		if cfg2 != nil {
			s.undisposedResponses.PushBack(cfg2)
		}
		if cfg3 != nil {
			s.undisposedResponses.PushBack(cfg3)
		}
	}
	s.mu.Unlock()
}

// SendAsync3 moves all three passed Tusent configs to the chat queues whose IDTs
// are in passed configs and does it asynchronously.
func (s *Sender) SendAsync4(cfg1, cfg2, cfg3, cfg4 *Tusent) {

	if cfg1 == nil && cfg2 == nil && cfg3 == nil && cfg4 == nil {
		return
	}

	s.mu.Lock()
	if !s.isStopped {
		if cfg1 != nil {
			s.undisposedResponses.PushBack(cfg1)
		}
		if cfg2 != nil {
			s.undisposedResponses.PushBack(cfg2)
		}
		if cfg3 != nil {
			s.undisposedResponses.PushBack(cfg3)
		}
		if cfg4 != nil {
			s.undisposedResponses.PushBack(cfg4)
		}
	}
	s.mu.Unlock()
}

// SendAsyncN moves all passed Tusent configs to the chat queues whose IDTs
// are in passed configs and does it asynchronously.
func (s *Sender) SendAsyncN(cfgs []*Tusent) {

	if len(cfgs) == 0 {
		return
	}

	s.mu.Lock()
	if !s.isStopped {
		for _, cfg := range cfgs {
			if cfg != nil {
				s.undisposedResponses.PushBack(cfg)
			}
		}
	}
	s.mu.Unlock()
}

// RequestRestartWith requests restart Sender (and probably embedded Lirester)
// with a passed new Sender (and Lirester) params. Restart will be done ASAP.
//
// If no one param will be passed, restart will be, but without changing
// the parameters.
func (s *Sender) RequestRestartWith(params []interface{}) (accepted bool) {

	s.mu.Lock()

	if !s.isStopped && s.restartRequestedWith == nil {
		s.restartRequestedWith = params
		accepted = true
	}

	s.mu.Unlock()
	return
}

// cqGet (chat queue getter) returns a chat queue associated with passed chat ID.
// It will be created, saved and then returned if it doesn't exist.
func (s *Sender) cqGet(chatIDT chat.IDT) *chatQueue {

	var cq *chatQueue

	// already exists
	if cq = s.preparedResponses[chatIDT]; cq != nil {
		return cq
	}

	// have to create a new chatQueue object for requested chat ID
	// or reuse already created (reallocation optimization).

	// reallocation optimization
	if l := len(s.cqReuseBuffer); l != 0 {
		cq = s.cqReuseBuffer[l-1]

		// avoid leaving ptr to reused cq in reuse storage
		// even in inaccessible place (slice will be shrinked)
		s.cqReuseBuffer[l-1] = nil
		s.cqReuseBuffer = s.cqReuseBuffer[:l-1]

		cq.q.Clear() // prepare cq to reuse

	} else {
		// nothing to reuse, allocate a new instance
		cq = new(chatQueue)
		cq.q.InitUnsafe(s.consts.cqCap, s.consts.cqCap)
		cq.n = 0
	}

	// associate with requested chat ID in main DB and return it
	s.preparedResponses[chatIDT] = cq
	return cq
}

// cqDel (chat queue delete) deletes a chat queue associated with passed chat ID.
//
// If that chat was updated a long ago than its allowed by cqLifetime const.
// If that chat is relatively new, it won't be deleted.
//
// Forcing cleanup.
// To force cleanup (avoid last update time checks) pass any value <= 0 as now arg.
func (s *Sender) cqDel(now int64, chatIDT chat.IDT) {

	var cq *chatQueue

	// exit if cq is not exists
	if cq = s.preparedResponses[chatIDT]; cq == nil {
		return
	}

	// exit if force cleanup disabled and chat still has time to live
	if now >= 0 && cq.la+s.consts.cqLifetime+s.consts.Ts[chatIDT.Type()] > now {
		return
	}

	// reallocation optimization (if it's enabled)
	if s.cqReuseBuffer != nil && uint16(len(s.cqReuseBuffer)) < s.consts.cqReuseBufferLen {
		s.cqReuseBuffer = append(s.cqReuseBuffer, cq)
		// don't call cq.Clear, will be Cleared in cqGet method

	} else {
		cq.q.Clear() // GC Tusents in cq
	}

	s.preparedResponses[chatIDT] = nil
	delete(s.preparedResponses, chatIDT)
}

// iter performs one sending iteration that contains following actions:
//
// 1. Moves all (or almost all*) core/tusent.Tusent objects
//    from undisposedResponses to preparedResponses.
//
// 2. Selects random chat queue from preparedResponses and if it is allowed to send
//    to that chat by lirester, does it.
//
// 3. Sending a message using backend methods (through modules/bridge.bridge).
//
// 4a. If message has been sent successfully and there is registered onSuccess
//     callbacks, moves a current core/tusent.Tusent object to the handledResponses.
//
// 4b. If message has not been sent successfully and there are still attempts,
//     returns a current core/tusent.Tusent to its chat queue and goes to next iter.
//
// 4c. If message has not been sent successfully and there are no attempts,
//     but registered onError callbacks, does the same as in 3a.
//
// It always returns false, except when Sender should be restarted
// and there is no Tusent to applying in preparedResponses field.
func (s *Sender) iter(now int64) (forceBreakLoop bool) {

	// WARNING!
	// ALL TIMESTAMPS IN NANO SECONDS!

	// TODO: Go to another chatQueue if there is an error of sending.

	// Variables of current iteration.
	// Filled by 2 action for 3+ actions, in 1 action used as temp var.
	var (
		ct *Tusent    // Current Tusent
		cq *chatQueue // Current chat queue
	)

	// 1 action (method's docs).
	// Do it only if request has not been requested.
	s.mu.Lock()
	if s.restartRequestedWith == nil {
		// TODO: Make n upper bounded by some const (can be changed)
		for i, n := int16(0), s.undisposedResponses.Len(); i < n; i++ {
			ct = (*Tusent)(s.undisposedResponses.PopFront())
			s.cqGet(ct.ChatIDT).q.PushBack(ct)
		}
	}
	s.mu.Unlock()

	// 2 action (method's docs).
	// Golang maps guarantees that order of this loop will be randomized
	// each time.
	for chatID, chatQueue := range s.preparedResponses {

		// maybe we have to cleanup chat if it's empty?
		if chatQueue.q.IsEmpty() {
			s.cqDel(now, chatID)
			continue // go to next chat and its queue
		}

		// do lirester (if it's enabled) allow us to send message to that chat?
		if s.cleanupRules != nil && !(chatQueue.n < s.consts.Ns[chatID.Type()]) {
			continue
		}

		ct = (*Tusent)(chatQueue.q.PopFront())
		cq = chatQueue
		break

		// there is no need nil check before breaking loop in assigning above,
		// because the ONLY public method (SendAsync) has nil check, which can't
		// allow to pass nil to undisposedResponses, and undisposedResponses
		// is the only one way how object can be passed to preparedResponses map.
	}

	// Variables of 3 action (real sending).
	// Declared before goto according to the Go rules.
	var (
		res        unsafe.Pointer
		err        error
		isFinalErr bool
	)

	// If true after loop above, there is nothing to do,
	// because there is no one config which prepared to sending.
	// TODO: Break anyway is stop or restart requested and save all Tusents? (how?)
	if ct == nil {
		s.mu.Lock()
		if s.isStopped || s.restartRequestedWith != nil {
			s.mu.Unlock()
			return true
		}
		s.mu.Unlock()
		goto exit
	}

	// Set retry attempts, if it wasn't done while generating Tusent
	// and there is "default" value of retry attempts
	if ct.RetryAttempts == 0 {
		ct.RetryAttempts = s.consts.retryAttempts
	}

	// 3 action (method's docs).
	res, err, isFinalErr = s.bridge.DoSend(ct.Config)
	ct.RetryAttempts--
	//noinspection GoNilness
	cq.la = now
	switch {
	// DO NOT FORGET TO SYNC LOGICALLY CASE CONDITIONS IF YOU WILL ADD A NEW'S!

	// 4a section (method's docs): Message has been successfully sent.
	//noinspection GoNilness (ct and cq change together, cq can't be nil)
	case err == nil:

		// lirester may be disabled
		if s.cleanupRules != nil {
			cq.n++

			// add cleanup rule
			when := uint64(s.consts.Ts[ct.ChatIDT.Type()] + now)
			s.cleanupRules.PushFront(uint64(ct.ChatIDT), when)
		}

		// call backend "sending success" callback
		if s.bridge.SendOK != nil {
			s.bridge.SendOK(ct.Ctx, res)
		}

		// deferring OnSuccess callbacks calls and transactions' finishes
		if len(ct.OnSuccess) > 0 || ct.NeedToFinish() {
			s.mu.Lock()
			s.handledResponses.PushBack(ct.MakeSuccess(res))
			s.mu.Unlock()
		}

	// 4c section (method's docs): Message has not been sent, and error of that
	// is final (can NOT be changed in the future)
	// or there is no more attempts to sending.
	case isFinalErr || ct.RetryAttempts == 0:

		// call backend "sending error" callback
		if s.bridge.SendErr != nil {
			s.bridge.SendErr(ct.Ctx, err)
		}

		// deferring OnError callbacks calls and transactions' finishes
		if len(ct.OnError) > 0 || ct.NeedToFinish() {
			s.mu.Lock()
			s.handledResponses.PushBack(ct.MakeError(err))
			s.mu.Unlock()
		}

	// 4b section (method's docs): Message has not been sent, and error of that
	// is not final (can be changed in the future) and there are attempts.
	case !isFinalErr && ct.RetryAttempts != 0:

		// anyway save error, maybe it's a new error
		ct.MakeError(err)

		// maybe stock of "infinity number of retrying attempts" is over?
		if ct.RetryAttempts == s.consts.retryAttemptsInfMax {
			if s.bridge.SendInfOverflow != nil {
				s.bridge.SendInfOverflow(ct.Ctx, ct.SendingErr, ct.Config)
			}

		} else {
			// decrease counter (for either finite or infinite numbers of attempts),
			// push back to queue
			ct.RetryAttempts--

			//noinspection GoNilness (ct and cq change together, cq can't be nil)
			cq.q.PushFront(ct)
		}
	}

exit:

	// do lirester cleanup (over accumulated rules)
	for i, n := int32(0), s.cleanupRules.Len(); i < n; i++ {
		idBadType, whenBadType := s.cleanupRules.PopFront()

		if int64(whenBadType) > now {
			// the time has not come yet
			s.cleanupRules.PushBack(idBadType, whenBadType)
			continue
		}

		if cq = s.preparedResponses[chat.IDT(idBadType)]; cq != nil {
			cq.n--
		}
	}

	return false
}

// mainLoop starts the infinity loop of iter calls while iter calls allows it.
// This is the loop of B thread.
func (s *Sender) mainLoop() {

	if s.consts.mainLoopDelay != 0 {

		ticker := time.NewTicker(s.consts.mainLoopDelay)
		for tick := range ticker.C {
			if interrupt := s.iter(tick.UnixNano()); interrupt {
				break
			}
		}
		ticker.Stop()

	} else {
		for continu := true; !continu; {
			continu = !s.iter(time.Now().UnixNano())
		}
	}
}

// secondaryLoop starts the infinity loop of calling OnSuccess, OnError finishers
// of completed Tusents and also finishes the transactions.
// This is the loop of C thread (of C thread routines tbh).
func (s *Sender) secondaryLoop() {

	aN, maxN := int16(0), int16(1<<s.consts.handledResponsesLenExp)
	a := make([]*Tusent, maxN)

	for {
		s.mu.Lock()

		aN = s.handledResponses.Len()
		// s.handledResponses can be grown,
		// but anyway don't take more handled tusents than determined by const
		if aN > maxN {
			aN = maxN
		}

		// move N handled tusents to prepared buffer
		// they will be finished after a realising mutex...
		for i := int16(0); i < aN; i++ {
			a[i] = (*Tusent)(s.handledResponses.PopFront())
		}

		// ... but before we must check whether we stop?
		// (restart or full stop is requested)
		if s.isStopped || s.restartRequestedWith != nil {
			break
		}

		s.mu.Unlock()

		// apply all tusents without locking mutex
		for i := int16(0); i < aN; i++ {
			a[i].Call()
			a[i] = nil // TODO: here Tusent will be GC'd (realloc optimise)
		}
	}

	// apply tusents that's still in a, but don't handled
	// because the loop above has been stopped by isStopped condition
	for i := int16(0); i < aN; i++ {
		a[i].Call()
		a[i] = nil // TODO: here Tusent will be GC'd (realloc optimise)
	}

	s.wg.Done()
}

// threadB is the lifecycle of Sender B thread (read docs).
func (s *Sender) threadB() {

	s.bridge.ML.Debug(
		`Serving prepared Tusents successfully started.`,
		logger.KindAsField(logger.Core, logger.Initialization),
		zap.Int("goroutine_numbers", 1),
		zap.Bool("lirester_enabled", s.cleanupRules != nil),
		zap.Duration("lirester_delay_ns", s.consts.mainLoopDelay),
	)

	s.wg.Add(1)
	s.mainLoop()
	s.wg.Done()

	s.bridge.ML.Debug(
		`Serving prepared Tusents successfully stopped.`,
		logger.KindAsField(logger.Core, logger.Initialization),
	)

	s.wg.Wait() // wait thread C goroutines

	// AT THIS CODE POINT IT'S GUARANTEED THAT:
	// - s.preparedResponses is empty, s.handledResponses is empty,
	// - all thread C goroutines are stopped.

	// Do restart if it was requested
	s.mu.Lock()
	if s.restartRequestedWith != nil {

		s.init().applyParams(s.restartRequestedWith).start()

		s.bridge.ML.Debug(
			`Restart of Sender successfully completed.`,
			logger.KindAsField(logger.Core),
		)
	}
	s.mu.Unlock()
}

// threadC is the lifecycle of Sender C thread (and its routines) (read docs).
func (s *Sender) threadC() {

	s.bridge.ML.Debug(
		`Serving handled Tusents successfully started.`,
		logger.KindAsField(logger.Core, logger.Initialization),
		zap.Uint8("goroutine_numbers", s.consts.threadCGN),
	)

	// start additional routins
	for i := uint8(0); i < s.consts.threadCGN-1; i++ {
		s.wg.Add(1)
		go s.secondaryLoop()
	}

	// start main routine of this thread ...
	s.wg.Add(1)
	s.secondaryLoop()
	s.wg.Done()

	s.wg.Wait() // ... and wait for additional C routines and thread B to complete

	s.bridge.ML.Debug(
		`Serving handled Tusents successfully stopped.`,
		logger.KindAsField(logger.Core, logger.Initialization),
	)
}

// start starts a two workers (each in its own goroutine):
//
// - Main worker (B thread).
//   A worker that calls iter method every time Lirester allows it.
//   Also performs restart if it was requested.
//
// - Secondary worker (with its own internal workers) (C thread).
//   Calls OnSuccess, OnError methods, finishes transactions.
//
// threadB leads.
func (s *Sender) start() *Sender {

	go s.threadB()
	go s.threadC()

	// mu is already locked if this is a restart
	// and no lock required if this is a constructor
	s.isStopped = false

	return s
}

// init initializes internal values of Sender by its DEFAULT values and constants.
func (s *Sender) init() *Sender {

	s.consts.cqCap = deque.MinCapacity
	s.consts.cqReuseBufferLen = 256
	s.consts.cqLifetime = int64(10 * time.Minute) // the same as .Nanoseconds() call

	s.consts.retryAttempts = 1
	s.consts.retryAttemptsInfMax = math.MinInt8

	s.consts.undisposedResponsesLenExp = 10 // 1024, because it's 2^10
	s.consts.handledResponsesLenExp = 11    // 2048, because it's 2^11
	s.consts.threadCGN = 1

	return s
}

// applyParams applies passed params to s, enables/disables/restarts lirester,
// allocates mem, etc.
// After that method, Sender object is considered ready to start.
func (s *Sender) applyParams(params []interface{}) *Sender {

	// separate lirester params from sender params
	var liresterParams []interface{}

	for _, pv := range params {

		if typedParam, ok := pv.(param); ok && typedParam != nil {
			typedParam(s)

		} else if paramGen, ok := pv.(func() param); ok && paramGen != nil {
			if typedParam := paramGen(); typedParam != nil {
				typedParam(s)
			}

		} else if liresterParam, ok := pv.(paraml); ok {
			liresterParams = append(liresterParams, liresterParam)

		} else if liresterParam, ok := pv.(func() paraml); ok {
			liresterParams = append(liresterParams, liresterParam)
		}
	}

	// What's about Lirester?
	if s.consts.disableLirester {
		if isEnabledLirester := s.cleanupRules != nil; isEnabledLirester {
			// restart prob: lirester must be disabled but is enabled now
			for _, v := range s.preparedResponses {
				v.n = 0
			}
			s.cleanupRules.Clear()
			s.cleanupRules = nil
		} else {
			// lirester must be disabled, but it's already so.
			// nothing to do
		}
		// lirester will be enabled at the next restart,
		// unless otherwise will specified
		s.consts.disableLirester = false

	} else {
		if isDisabledLirester := s.cleanupRules == nil; isDisabledLirester {
			// constructor prob but maybe restart with disabled lirester:
			// lirester should be created
			s.cleanupRules = deque.For128bit(s.consts.cleanupRulesLenExp)
			// TODO: grow, shrink factors
		} else {
			// restart prob: lirester must be restarted too and is enabled now,
			// but there is no hard restart of lirester
			// and new constants already overwritten.
			// nothing to do.
		}
	}

	// it can be not nil if sender is restarting,
	// and in that case we must save all Tusents that this queue have.
	if s.undisposedResponses == nil {
		s.undisposedResponses = deque.ForPointers(s.consts.undisposedResponsesLenExp)
	}

	// don't reallocate if it's the same size
	if s.handledResponses.Cap() == 1<<s.consts.handledResponsesLenExp {
		// call NewStatic because it's most readable constructor,
		// and minimum capacity == base capacity, but make it grow-possible
		// make it semi-static
		s.handledResponses = deque.ForPointers(s.consts.handledResponsesLenExp)
		// todo: grow, shrink factors ? DONT DISABLE AUTO SHRINK
		// s.handledResponses.DisableAutoGrow = false
	}

	if s.preparedResponses == nil {
		s.preparedResponses = make(map[chat.IDT]*chatQueue)
	}

	switch /* What's about chat queues reuse buffer? */ {

	// restart prob: reusing chatQueue must be disabled but is enabled now
	case s.consts.cqReuseBufferLen == 0 && s.cqReuseBuffer != nil:
		for i, n := 0, len(s.cqReuseBuffer); i < n; i++ {
			s.cqReuseBuffer[i].q.Clear() // GC Tusents in chat queues
			s.cqReuseBuffer[i] = nil     // GC chat queues
		}
		s.cqReuseBuffer = nil

	// constructor prob, but maybe restart with disabled reusing chatQueue:
	// buffer must be created
	case s.consts.cqReuseBufferLen > 0 && s.cqReuseBuffer == nil:
		s.cqReuseBuffer = make([]*chatQueue, 0, s.consts.cqReuseBufferLen)

	// restart prob: reusing chatQueue must be enabled, but it's enabled already
	case s.consts.cqReuseBufferLen > 0:
		switch c := uint16(cap(s.cqReuseBuffer)); {

		// new requested buffer len the same as it was, nothing to do.
		case c == s.consts.cqReuseBufferLen:

		// new requested buffer len smaller than it was,
		// need to reallocate with shrinking (copying)
		case c > s.consts.cqReuseBufferLen:
			// https://github.com/go101/go101/wiki/How-to-efficiently-clone-a-slice
			s.cqReuseBuffer = append(s.cqReuseBuffer[:0:0],
				s.cqReuseBuffer[:s.consts.cqReuseBufferLen]...)

		// new requested buffer len bigger than it was,
		// need to reallocate with growing (copying)
		case c < s.consts.cqReuseBufferLen:
			prev := s.cqReuseBuffer
			s.cqReuseBuffer = make([]*chatQueue, s.consts.cqReuseBufferLen)
			copy(s.cqReuseBuffer, prev)
		}
	}

	// if it was a restart, this is done
	s.restartRequestedWith = nil

	return s
}

// MakeSender creates a new Sender object, the passed params will be applied to.
func MakeSender(b *bridge.Bridge, params []interface{}) *Sender {
	return new(Sender).init().applyParams(params).start()
}

// Wipe completely destroys passed Sender by its double pointer.
// After calling this function, the *s will have a nil instead of valid pointer!
//
// WARNING!
// DO NOT USE PASSED OBJECT AFTER CALLING THIS METHOD, BECAUSE OF UB AND PANIC.
func Wipe(s **Sender) {

	if s == nil || *s == nil {
		return
	}

	ss := *s

	ss.mu.Lock()
	ss.isStopped = true           // request to stop
	ss.restartRequestedWith = nil // cancel request to restart (if it was)
	ss.mu.Unlock()

	ss.wg.Wait() // wait for threads B, C to complete

	// preparedResponses and handledResponses must be empty at this code point,
	// but additional statements are cheaper than errors

	for chatIDT, cq := range ss.preparedResponses {
		cq.q.Clear()
		delete(ss.preparedResponses, chatIDT)
	}

	ss.handledResponses.Clear()

	// it's guaranteed at this code point that
	// s.preparedResponses and s.handledResponses are empty

	ss.preparedResponses = nil
	ss.handledResponses = nil

	for i, n := 0, len(ss.cqReuseBuffer); i < n; i++ {
		ss.cqReuseBuffer[i].q.Clear() // GC Tusents in chat queues
		ss.cqReuseBuffer[i] = nil     // GC chat queues
	}

	ss.undisposedResponses.Clear()
	ss.undisposedResponses = nil

	// done

	*s = nil
}
