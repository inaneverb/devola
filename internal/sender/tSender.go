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
	"fmt"
	"log"
	"strings"
	// api "github.com/go-telegram-bot-api/telegram-bot-api"
)

// -- Sender --
// Sender is the part of Telegram TBot UIG Framework.
// Sender is enagaged in sending outgoing messages, analysing results,
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
// sendable config, which contains sending entity as object of class implemented
// tgbotapi.Chattable interface, TCtx object connected to that entity,
// onSuccess, onError callbacks and some internal parts.
// (see tSendableConfig and TCtx types for details).
// 2. TCtx methods sends created tSendableConfig object to the common
// channel of all prepared to sending configs (chPreparedSendable)
// using method deferSend.
// 3. One of Sender goroutines when time will come, extract one of
// ready to sending tSendableConfig object and perform sending operation.
type tSender struct {

	// TBot object, this tSender object associated with.
	parent *TBot

	// Consts section.
	consts struct {
		preparedSendableChanLen int
		ctxChanLen              int
		ctxChanLifetime         int64
		maxSendRetryAttempts    int
		callbacksChanLen        int
		isAlwaysUseHTML         bool
		isAlwaysUseMD           bool
	}

	lirester *tLirester

	// todo: Delete chPreparedSendable
	chPreparedSendable chan *tSendableConfig

	sendable map[tChatID]chan *tSendableConfig

	chSendSuccessCallbacks chan *tViewSuccessFinisher

	chSendErrorCallbacks chan *tViewErrorFinisher
}

//
func (s *tSender) Start() {

}

// chchGet returns the channel of tSendableConfig's associated with the
// specified chat id chatId.
// If interenal sender object don't have channel for that chat id,
// it will be created.
// todo: Optimise (reuse allocated RAM for channels).
func (s *tSender) chchGet(chatID tChatID) chan *tSendableConfig {
	if _, ok := s.sendable[chatId]; !ok {
		s.sendable[chatId] = make(chan *tSendableConfig, s.consts.ctxChanLen)
	}
	return s.sendable[chatId]
}

// chchCleanup (believing that now is the current unixnano timestamp)
// receives the some chat id and checks whether lifetime of the associated
// tSendableConfig's channel with that chat id is over.
// If it's true, starts the delete context channel procedure.
func (s *tSender) chchCleanup(now int64, chatID tChatID) {
	if s.lirester.LastUpdated(chatId)+s.consts.ctxChanLifetime < now {
		s.chchForceDelete(chatId)
	}
}

// chchForceDelete just deletes tSendableConfig's channel associated with the
// specified chat id chatId.
// There are three important operations:
// 1. Close the channel.
// 2. Set channel pointer to nil (for GC).
// 3. Delete channel from map by chat id key.
func (s *tSender) chchForceDelete(chatID tChatID) {
	close(s.sendable[chatId])
	s.sendable[chatId] = nil // todo: probably unnecessary (added for GC)?
	delete(s.sendable, chatId)
}

// deferSend moves the tSendableConfig object to the common queue.
// This sendable config will be proceed at the next main sender loop iteration.
func (s *tSender) deferSend(cfg *tSendableConfig) {
	if cfg.validate() {
		s.chPreparedSendable <- cfg
	}
}

// serveMessages starts the infinity loop, which each new iteration starts
// when its allowed by Lirester MainLoop object (see Lirester type for details).
// (see serveMessageIteration for details).
func (s *tSender) serveMessages() {
	// Write debug log message, that serving is started
	log.Println("tSender.serveMessages",
		"Serving sendable configs successfully started at the separated goroutine.")
	// Really start serving until channel is closed
	for tick := range s.lirester.MainLoop.C {
		s.serveMessageIteration(tick.UnixNano())
	}
	// Code below will be executed when s.lirester.MainLoop.C channel
	// will be closed (it will be when lirester will be wiped (deleted)).
	// So, we write debug log message about it
	log.Println("tSender.serveMessages",
		"Serving sendable configs successfully stopped. "+
			"The separated goroutine has been shutdown.")
}

// serveMessageIteration performs one sending iteration that contains:
// 1. Attempt to get some ready sendable config from common chPreparedSendable
// channel.
// 2. Attempt to send ready outgoing Telegram messages from tSendableConfig
// from 1st step using Telegram API.
// 3. Deferring onSuccess and onError callbacks, if it's need.
// 4. Writing log messages, if it's need.
// 5. Mark unsent messages as messages that must be sent later
// (return it to the queue).
//
// For more details, following the code and read submethods docs.
func (s *tSender) serveMessageIteration(now int64) {
	var handlingConfig *tSendableConfig
	// Extract all prepared to send context objects and moves them
	// to the separated channels
	for i, n := 0, len(s.chPreparedSendable); i < n; i++ {
		preparedSendable := <-s.chPreparedSendable
		s.chchGet(preparedSendable.ctx.Chat.ID) <- preparedSendable
	}
	// Loop over all online chats, try to find chat with non empty
	// outgoing message channel and check if it's allow to send
	// message to the associated chat.
	// Golang maps guarantees that order of this loop will be randomized
	// each time.
	for chatId, chSendable := range s.sendable {
		if len(chSendable) == 0 {
			s.chchCleanup(now, chatId)
			continue
		}
		if !s.lirester.Try(now, chatId) {
			continue
		}
		// This could be never happen, but anyway this check is here
		// just for the future, 'cause in other way there is will be
		// a whole full iteration instead just go to next chat's channel.
		if handlingConfig = <-chSendable; handlingConfig == nil {
			continue
		}
		// handling config found
		break
	}
	// If after prev loop the sendable config still is nil, no one
	// sendable config is exists. Just go to next iteration.
	if handlingConfig == nil {
		return
	}
	// Really send message, and if it was successfull,
	// approve sending in lirester, update session.
	// Otherwise register any next attempt if it's necessary.
	if sentMsg := s.send(handlingConfig); sentMsg != nil {
		chat := handlingConfig.ctx.Chat
		isUser := chat.IsChannel() || chat.IsGroup() || chat.IsSuperGroup()
		s.lirester.Approve(now, handlingConfig.ctx.Chat.ID, isUser)
		// Update session if it's need (save info about sent message)
		if handlingConfig.isUpdateSession {
			// handlingConfig.ctx.
		}
	} else {
		// Need at least one more try if retryAttempts isn't zero
		if handlingConfig.retryAttempts != 0 {
			// Negative value means infinity number of tries
			if handlingConfig.retryAttempts > 0 {
				handlingConfig.retryAttempts--
			}
			s.sendable[handlingConfig.ctx.Chat.ID] <- handlingConfig
		}
	}
	// Perform Lirester cleanup
	s.lirester.Cleanup(now)
}

// serveSuccessCallbacks writes debug log messages before loop over
// deferred onSuccess callbacks has been started and when that loop is completed.
// Also, starts loop over deferred callbacks until callbacks channel is closed.
func (s *tSender) serveSuccessCallbacks() {
	// Write debug log message, that serving is started
	log.Println("tSender.serveSuccessCallbacks",
		"Serving successfull callbacks of the sendable configs "+
			"successfully started at the separated goroutine.")
	// Really start serving until channel is closed
	for cb := range s.chSendSuccessCallbacks {
		cb.Call(true) // todo: make depend panicCheck from consts
	}
	// Code below will be executed when s.chSendSuccessCallbacks
	// will be closed. So, we write debug log message about it
	log.Println("tSender.serveSuccessCallbacks",
		"Serving successfull callbacks of the sendable configs "+
			"successfully stopped. The separated goroutine has been shutdown.")
}

// serveErrorCallbacks writes debug log messages before loop over
// deferred onError callbacks has been started and when that loop is completed.
// Also, starts loop over deferred callbacks until callbacks channel is closed.
func (s *tSender) serveErrorCallbacks() {
	// Write debug log message, that serving is started
	log.Println("tSender.serveErrorCallbacks",
		"Serving error callbacks of the sendable configs successfully started "+
			"at the separated goroutine.")
	// Really start serving until channel is closed
	for cb := range s.chSendErrorCallback {
		cb.Call(true) // todo: make depend panicCheck from consts
	}
	// Code below will be executed when s.chSendErrorCallback
	// will be closed. So, we write debug log message about it
	log.Println("tSender.serveErrorCallbacks",
		"Serving error callbacks of the sendable configs successfully stopped. "+
			"The separated goroutine has been shutdown.")
}

// 'send' checks whether context object contains valid fields,
// and if it's true, try to send outgoing Telegram messages.
func (s *tSender) send(cfg *tSendableConfig) *tgbotapi.Message {
	if s.bot == nil {
		s.bot.log.Class("tSender").Method("send").Errorw("Message can't be sent",
			"err", "Wrong usage of tSender object: Nil tSender object")
		return nil
	}
	// Try to send outgoing Telegram message using API
	if sentMsg, err := s.bot.bot.Send(cfg.config); err == nil {
		// Sent successfully. Defer onSuccess callbacks, if it's need.
		// todo: CONSTRUCTOR NEVER RETURNS NIL! HANDLE IT MANUALLY! because of tr finishers
		if cfg := makeSendSuccessCallbackConfig(cfg, &sentMsg); cfg != nil {
			s.chSendSuccessCallbacks <- cfg
		}
		return &sentMsg
	} else {
		// Has not been sent. Defer onError callback, extract body of outgoing
		// unsent message and log it.
		// todo: CONSTRUCTOR NEVER RETURNS NIL! HANDLE IT MANUALLY! because of tr finishers
		// todo: EXECUTE FINISHERS ONLY IF NO ONE ATTEMPT LEAVE
		if cfg := makeSendErrorCallbackConfig(cfg, err); cfg != nil {
			s.chSendErrorCallback <- cfg
		}
		messageBody := fmt.Sprintf("%+v", cfg.config)
		messageBody = strings.Replace(messageBody, "\n", "\\n", -1)
		log.Println("tSender.send", "Message can't be sent",
			err, cfg.ctx.Event, cfg.ctx.From, cfg.ctx.Sess, messageBody)
		return nil
	}
}

// 'params' might be 'tSenderParam' or 'tLiresterParam'
func newSender(parent *TBot, params ...interface{}) *tSender {
	s := &tSender{bot: parent}
	// The default values of sender params
	const (
		defCtxChanLen           = int(0)
		defCtxChanLifetime      = int64(0)
		defMaxSendRetryAttempts = int(0)
		defCallbackChanLen      = int(0)
		defIsAlwaysUseMD        = false
		defIsAlwaysUseHTML      = false
	)
	panic("overwrite consts") // todo: above
	// Apply default sender params.
	// It'll be overwritten later by passed sender params.
	s.consts.ctxChanLen = defCtxChanLen
	s.consts.ctxChanLifetime = defCtxChanLifetime
	s.consts.maxSendRetryAttempts = defMaxSendRetryAttempts
	s.consts.callbacksChanLen = defCallbackChanLen
	s.consts.isAlwaysUseMD = defIsAlwaysUseMD
	s.consts.isAlwaysUseHTML = defIsAlwaysUseHTML
	// Declare slices for sender and lirester params
	vSenderParams := []tSenderParam(nil)
	vLiresterParams := []tLiresterParam(nil)
	// Split all params to the Sender params and Lirester params separately
	for _, param := range params {
		switch param.(type) {
		case tSenderParam:
			vSenderParams = append(vSenderParams, param.(tSenderParam))
		case tLiresterParam:
			vLiresterParams = append(vLiresterParams, param.(tLiresterParam))
		}
	}
	// Create Lirester object and store it to the Sender
	s.lirester = newLirester(vLiresterParams...)
	// Apply sender params
	for _, param := range vSenderParams {
		if param != nil {
			param(s)
		}
	}
	// Allocate mem for sendable configs channel and map
	s.chPreparedSendable = make(chan *tSendableConfig,
		s.consts.preparedSendableChanLen)
	s.sendable = make(map[int64]chan *tSendableConfig)
	// Allocate mem for callback's channels
	s.chSendSuccessCallbacks = make(chan *tSendSuccessCallbackConfig,
		s.consts.callbacksChanLen)
	s.chSendErrorCallback = make(chan *tViewErrorFinisher,
		s.consts.callbacksChanLen)
	// Sender successfully created
	parent.sender = s
	return s
}
