// Copyright © 2018. All rights reserved.
// Author: Alice Qio.
// Contacts: <qioalice@gmail.com>.
// License: https://opensource.org/licenses/MIT

package sender

import (
	"log"

	api "github.com/go-telegram-bot-api/telegram-bot-api"

	"../chat"
	"../lirester"
	"../view"
)

// Sender is the part of Telegram TBot UIG Framework.
// Sender is engaged in sending outgoing messages, analysing results,
// calls deferring callbacks of sending result, and more.
//
// Sender must be created only by newSender function.
//
// As object, Sender contains TBot object pointer, Lirester object by pointer
// (see Lirester type for details, it's important to understand how sender works),
// map from chat id to associated channel
// (in which all connected context object will be stored),
// and channels for deferred onSuccess and onError callbacks.
//
// How it works.
// 1. TCtx methods calls makeSendableConfig consturctor to create
// chch config, which contains sending entity as object of class implemented
// tgbotapi.Chattable interface, TCtx object connected to that entity,
// onSuccess, onError callbacks and some internal parts.
// (see tSendableConfig and TCtx types for details).
// 2. TCtx methods sends created tSendableConfig object to the common
// channel of all prepared to sending configs (chPreparedSendable)
// using method deferSend.
// 3. One of Sender goroutines when time will come, extract one of
// ready to sending tSendableConfig object and perform sending operation.
type Sender struct {

	// TBot object, this Sender object associated with.
	endpoint *api.BotAPI

	lirester *lirester.Lirester

	// Consts section.
	consts struct {

		// Base capacity of chatchan queue. It is passed to its constructor.
		chchCap int16

		// Const length of chchUnused field.
		chchUnusedCap int16

		// Delay in ns (as standard time.Duration) that must pass since
		// the last changing of some chatchan queue, before it will be deleted.
		chchLifetime int64

		// How many attempts to send a message will be made
		// before it is considered to be completely unsent.
		retryAttempts int8

		// Buffer size of chFinishers chan.
		finisherChanLen int

		// Used in ctx.Ctx // TODO: what methods?
		isAlwaysUseHTML bool

		// Used in ctx.Ctx // TODO: what methods?
		isAlwaysUseMD bool
	}

	// Core of Sender object.
	// All chatchan queues are stored here and can be accessed by chat ID
	// they're represents.
	chch map[chat.ID]*chatchan

	// Reallocation optimization.
	// All chatchan that should be deleted from chch field go here.
	// When a new chatchan object is requested it can be used from this storage.
	chchUnused []*chatchan

	// Success and error finishers bufferized channel.
	// Read more: view.Finisher.
	chFinishers chan *view.Finisher
}

// Predefined default values of some important Sender constants.
const (

	// Default const length of chchUnused field.
	cChchUnusedCap int16 = 0

	// Default delay in ns that must pass since the last changing
	// of some chatchan queue, before it will be deleted.
	cChchLifetime int64 = 0

	// Default retry attempts value.
	cRetryAttempts int8 = 0

	// Default buffer size of chFinishers chan.
	cFinisherChanLen int = 0
)

//
func (s *Sender) Start(level int) {
	panic("implement me")
}

//
func (s *Sender) Stop(level int) {
	panic("implement me")
}

//
func (s *Sender) Restart() {
	panic("implement me")
}

// chchGet returns chatchan object associated with passed chat ID
// from chch field of s.
// It will be created before return if it doesn't exist.
func (s *Sender) chchGet(chatID chat.ID) *chatchan {

	// if chatchan of requested chat ID already exists
	if chch, ok := s.chch[chatID]; ok {
		return chch
	}

	var chch *chatchan

	// have to create a new chatchan object for requested chat ID
	// or reuse already created (reallocation optimization).

	// reallocation optimization
	if l := len(s.chchUnused); l != 0 {
		chch = s.chchUnused[l-1]
		// avoid leaving ptr to reused chch in reuse storage
		// even in inaccessible place (slice will be shrinked)
		s.chchUnused[l-1] = nil
		s.chchUnused = s.chchUnused[:l-1]
		chch.Clear()

	} else {
		// nothing to reuse, allocate a new instance
		chch = makeChatchan(s.consts.chchCap)
	}

	// associate with requested chat ID in main DB and return it
	s.chch[chatID] = chch
	return chch
}

// chchCleanup deletes the chatchan object associated with passed chat ID
// from chch field of s if that chat was updated a long ago than its allowed
// by chchLifetime const.
// If that chat is relatively new, it won't be deleted.
//
// Forcing cleanup.
// You can avoid last update time checks and forcing cleanup
// by passing a zero or any negative value (-1 mostly) as now argument.
func (s *Sender) chchCleanup(now int64, chatID chat.ID) {

	if now >= 0 && s.lirester.LastUpdated(chatID)+s.consts.chchLifetime > now {
		return
	}

	// reallocation optimization (if it's enabled)
	if s.chchUnused != nil {
		if int16(len(s.chchUnused)) < s.consts.chchUnusedCap {
			s.chchUnused = append(s.chchUnused, s.chch[chatID])
			s.chch[chatID] = nil // TODO: unnecessary statement ?
		}
	}

	delete(s.chch, chatID)
}

// DeferSend moves passed ToSend config to the chatchan queue of the chat
// whose ID is in passed config.
func (s *Sender) DeferSend(cfg *ToSend) {
	// chchGet always return not nil chatchan
	s.chchGet(cfg.chatID).PushBack(cfg)
}

//
func (s *Sender) DeferFinisher(finisher *view.Finisher) {

}

// applyConfig performs one sending iteration that contains:
// 1. Attempt to get some ready chch config from common chPreparedSendable
// channel.
// 2. Attempt to send ready outgoing Telegram messages from tSendableConfig
// from 1st step using Telegram API.
// 3. Deferring onSuccess and onError callbacks, if it's need.
// 4. Writing log messages, if it's need.
// 5. Mark unsent messages as messages that must be sent later
// (return it to the queue).
//
// For more details, following the code and read submethods docs.
func (s *Sender) applyConfig(now int64) {

	var (
		handlingConfig      *ToSend
		handlingConfigQueue *chatchan
	)

	// Try to find non empty chat and check if it's allow to send message to.
	// Golang maps guarantees that order of this loop will be randomized
	// each time.
	for chatID, chQueue := range s.chch {

		// maybe we should to cleanup chat if it's empty?
		if chQueue.IsEmpty() {
			s.chchCleanup(now, chatID)
			continue
		}

		// do lirester allow us to send message to that chat?
		if !s.lirester.Try(now, chatID) {
			continue
		}

		// This could be never happen, but anyway this check is here
		// just for the future, 'cause in other way there is will be
		// a whole full iteration instead just go to next chat's channel.
		if handlingConfig = chQueue.PopFront(); handlingConfig == nil {
			continue
		}

		// handling found config
		handlingConfigQueue = chQueue
		break
	}

	// If after prev loop ToSend config still is nil, there is nothing to do.
	if handlingConfig == nil {
		return
	}

	// Really send message.
	if sentMsg := s.reallySend(handlingConfig); sentMsg != nil {
		// Successfully sent

		s.lirester.Approve(now, handlingConfig.chatID, handlingConfig.isUserChat)

		// Update session if it's need (save info about sent message)
		if handlingConfig.isUpdateSession {
			// TODO: save sent message ID
		}
	} else {
		// Message has not been sent (some error occurred).

		// do we need an another try?
		if handlingConfig.retryAttempts != 0 {

			// Negative value means infinity number of tries
			if handlingConfig.retryAttempts > 0 {
				handlingConfig.retryAttempts--
			}

			handlingConfigQueue.PushFront(handlingConfig)
		}
	}

	// Perform Lirester cleanup
	// TODO: Too often, decrease calls count
	s.lirester.Cleanup(now)
}

//
func (s *Sender) applyFinisher(finisher *view.Finisher) {
	finisher.Call()
	// TODO: Handle finisher errors (session tr, chat tr, recovered panics)
}

// reallySend tries to send Telegram messages from cfg and if it's need
// register success or error finishers.
func (s *Sender) reallySend(cfg *ToSend) *api.Message {

	// cfg2 is additional message config can contain only delete message config
	// and if it's so, we have to sent it first
	if cfg.cfg2 != nil {

		if _, err := s.endpoint.Send(cfg.cfg2); err == nil {
			// message by cfg2 has been succesfully deleted
			cfg.cfg2 = nil
		} else {
			// We can't process sending process (can't beginning to sent cfg1
			// if cfg2 hasn't been sent)
			// TODO: figure out what kind of error occurred, handle different errors
			return nil
		}
	}
	var (
		// successfully sent message
		sent *api.Message

		// finisher constructor arguments
		ff    = view.CEnablePanicGuard
		ctxp1 = cfg.originalCtx
		ctxp2 = cfg.pass2finisherCtx
	)

	// try to send a primary message config (cfg)
	switch sentObj, err := s.endpoint.Send(cfg.cfg); {

	case err == nil && len(cfg.onSuccess) != 0:
		sent = &sentObj
		s.DeferFinisher(view.MakeFinisherSuccess(ff, cfg.onSuccess, ctxp1, ctxp2, sent))

	case err != nil && len(cfg.onError) != 0:
		s.DeferFinisher(view.MakeFinisherError(ff, cfg.onError, ctxp1, ctxp2, err))
	}

	// TODO: log
	return sent
}

// serveConfigs starts the infinity loop, which each new iteration starts
// when its allowed by Lirester MainLoop object.
// More info: lirester.Lirester, Sender.applyConfig.
func (s *Sender) serveConfigs() {

	// Write debug log message, that serving is started
	log.Println("Sender.serveConfigs",
		"Serving configs successfully started at the separated goroutine.")

	// Really start serving until channel is closed
	for tick := range s.lirester.MainLoop.C {
		s.applyConfig(tick.UnixNano())
	}

	// Code below will be executed when s.lirester.MainLoop.C channel
	// will be closed.
	log.Println("Sender.serveConfigs",
		"Serving configs successfully stopped. The separated goroutine has been shutdown.")
}

// serveFinishers writes debug log messages before loop over
// deferred onSuccess callbacks has been started and when that loop is completed.
// Also, starts loop over deferred callbacks until callbacks channel is closed.
func (s *Sender) serveFinishers() {

	// Write debug log message, that serving is started
	log.Println("Sender.serveFinishers",
		"Serving finishers successfully started at the separated goroutine.")

	// Really start serving until channel is closed
	for finisher := range s.chFinishers {
		s.applyFinisher(finisher)
	}

	// Code below will be executed when s.chFinishers
	// will be closed.
	log.Println("Sender.serveFinishers",
		"Serving finishers successfully stopped. The separated goroutine has been shutdown.")
}

// applyParams applies each param from params slice to the current Sender.
// It overwrites alreay applied parameters.
func (s *Sender) applyParams(params []interface{}) {
	for _, param := range params {

		if typedParam, ok := param.(Param); ok && typedParam != nil {
			typedParam(s)

		} else if paramGen, ok := param.(func() Param); ok && paramGen != nil {
			if typedParam := paramGen(); typedParam != nil {
				typedParam(s)
			}
		}
	}
}

// MakeSender creates a new Sender object, the passed params will be applied to.
func MakeSender(params []interface{}) *Sender {

	var s Sender

	// Apply default sender params.
	// It'll be overwritten later by passed lirester params.
	s.consts.chchCap = cChChMinCapacity
	s.consts.chchUnusedCap = cChchUnusedCap
	s.consts.chchLifetime = cChchLifetime
	s.consts.retryAttempts = cRetryAttempts
	s.consts.finisherChanLen = cFinisherChanLen
	s.consts.isAlwaysUseMD = false
	s.consts.isAlwaysUseHTML = false

	s.applyParams(params)

	// Allocate mem for main core, finisher's chan and reused chatchan's ptr
	// storage.
	s.chch = make(map[chat.ID]*chatchan)
	if s.consts.chchUnusedCap > 0 {
		// negative chchUnusedCap value means that a feature
		// of chatchan's reallocation optimization will be disabled.
		s.chchUnused = make([]*chatchan, 0, s.consts.chchUnusedCap)
	}
	s.chFinishers = make(chan *view.Finisher, s.consts.finisherChanLen)

	return &s
}
