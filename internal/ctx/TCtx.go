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
	"strings"

	api "github.com/go-telegram-bot-api/telegram-bot-api"

)

//
type TCtx struct {
	bot    *TBot

	Event  tEvent     `json:"occurred_event"`
	Update *api.Update `json:"api_update"`
	From   *api.User   `json:"from_user"`
	Chat   tChatInfo   `json:"from_chat"`

	sess   *tSession   `json:"session"`

	g      *tCtxMessageGenerator
}

//
func (ctx *TCtx) TrSessFinish() error {

}

//
func (ctx *TCtx) TrChatFinish() error {

}

//
func (ctx *TCtx) currentViewIDEncoded() tViewIDEncoded {
	
}


// 'isAllowByMdlws' returns 'true' if all middlewares that can be called
// for event stored in context object 'ctx' with specified 'when' type
// will return 'true' (they will allow futher procession).
// If false returned by at least one of middlewares, the 'isAllowByMdlws'
// will return false immediately.
func (ctx *TCtx) isAllowedByMiddlewares(when session.TStep) bool {
	panicChecker := func(ctx *TCtx) {
		recoverErr := recover()
		if recoverErr == nil {
			return
		}
		r.bot.log.Named("tReceiver")
		if err := recover(); err != nil {
			r.bot.log.Class("tReceiver").Method("isAllowByMdlws").Errorw(
				"Panic occurred and has been restored while trying to pass update "+
					"through middlewares",
				"ctx", ctx, "panic_err", err)
		}
	}
	defer panicChecker(ctx)
	// Try to apply general middlewares
	for _, mdlw := range r.mdlwsMain {
		if !mdlw(ctx) {
			return false
		}
	}
	// Try to apply specified middlewares
	// Apply middlewares for text handler type
	if ctx.Event.Type == HandlerTypeText {
		if when != session.TStepAny {
			for _, mdlw := range r.mdlwsTextWhen[when] {
				if !mdlw(ctx) {
					return false
				}
			}
		} else {
			for _, mdlw := range r.mdlwsTextJust {
				if !mdlw(ctx) {
					return false
				}
			}
		}
		return true
	}
	// Apply middlewares for non text handler types
	if when != session.TStepAny {
		if mdlws := r.mdlwsWhen[ctx.Event.Type][when]; mdlws != nil {
			for _, mdlw := range mdlws[ctx.Event.Body] {
				if !mdlw(ctx) {
					return false
				}
			}
		}
	} else {
		for _, mdlw := range r.mdlwsJust[ctx.Event.Type][ctx.Event.Body] {
			if !mdlw(ctx) {
				return false
			}
		}
	}
	return true
}

// 'OnSuccess' registers handler 'cb' as handler that will be called
// when constructing outgoing message config will be sent successfully.
// It works only for 'New' or 'Edit' finishers.
// If you will use 'Del', 'DelLast', 'DelAll', 'DelAllExceptLast' finishers,
// callback willn't be called.
// Context object and object of new or edit message will be passed to the
// callback.
func (ctx *TCtx) OnSuccess(cb FViewSuccessFinisher) *TCtx {
	if ctx.alloc() == nil {
		return ctx
	}
	if cb != nil {
		ctx.g.onSuccess = append(ctx.g.onSuccess, cb)
	}
	return ctx
}

// 'OnError' registers handler 'cb' as handler that will be called
// when constructing outgoing message config will be sent unsuccessfully
// and some error will occur.
// It works for all finishers.
// Context object and error object will be passed to the callback.
func (ctx *TCtx) OnError(cb FViewErrorFinisher) *TCtx {
	if ctx.alloc() == nil {
		return ctx
	}
	if cb != nil {
		ctx.g.onError = append(ctx.g.onError, cb)
	}
	return ctx
}

// 'Attempts' allows to specify how many attempts to send outgoing message's
// config will be made if any error while trying to send message will occur.
// Note! 'n' must be in the range: [-1..max], where 'max' is
// 'maxSendRetryAttempts' const of used tReceiver.
// If it's not so, this method will do nothing.
// '-1' is reserved value and means that message must be sent anyway.
// It means that Sender will try to send message until it's successfully sent.
func (ctx *TCtx) Attempts(n int) *TCtx {
	if ctx.alloc() == nil {
		return ctx
	}
	if n >= -1 && n <= ctx.sender.consts.maxSendRetryAttempts {
		ctx.g.retryAttempts = int8(n)
	}
	return ctx
}

// 'ReplyTo' marks the future message (message that will be created
// by one of finishers and sent) as reply message to the message with
// provided message id 'messageId'.
// If 'messageId' is 0, there is no-op.
// It works only for 'New' finishers.
// If you will use 'Edit', 'Del', 'DelLast', 'DelAll', 'DelAllExceptLast'
// finishers, the message id will be ignored and in the end turns out that this
// method does nothing.
func (ctx *TCtx) ReplyTo(messageId int) *TCtx {
	if ctx.alloc() == nil {
		return ctx
	}
	if messageId != 0 {
		ctx.g.replyTo = messageId
	}
	return ctx
}

// 'HTML' switches parse mode of message's body to the HTML support mode.
// Warning! It disables Markdown support mode!
//
// Note! You can also use 'SenderAlwaysUseHTML' Sender parameter
// while creating TBot or directly Sender object.
// In that case you could not to call always 'HTML' method, 'cause
// it'll be enabled by default, but if you'll want to disable HTML
// or switch to Markdown support for some message, you can call
// HTML(false) or MD(true).
func (ctx *TCtx) HTML(enabled bool) *TCtx {
	if ctx.alloc() == nil {
		return ctx
	}
	ctx.g.isHTML, ctx.g.isMarkdown = enabled, false
	return ctx
}

// 'MD' switches parse mode of message's body to the Markdown support mode.
// Warning! It disables HTML support mode!
//
// Note! You can also use 'SenderAlwaysUseMD' Sender parameter
// while creating TBot or directly Sender object.
// In that case you could not to call always 'MD' method, 'cause
// it'll be enabled by default, but if you'll want to disable MD
// or switch to HTML support for some message, you can call
// MD(false) or HTML(true).
func (ctx *TCtx) MD(enabled bool) *TCtx {
	if ctx.alloc() == nil {
		return ctx
	}
	ctx.g.isMarkdown, ctx.g.isHTML = enabled, false
	return ctx
}

//
func (ctx *TCtx) Keyboard(rules ...interface{}) *tCtxKeyboardGenerator {

}

// // 'Keyboard' assigns 'keyboard' as keyboard that will be sent along with
// // the message.
// // If nil is passed it means that keyboard that was sent already must be
// // deleted. If you don't want to attach any keyboard to the message
// // just do not call this method.
// // It works only for 'New' or 'Edit' finishers.
// // If you will use 'Del', 'DelLast', 'DelAll', 'DelAllExceptLast' finishers,
// // the attached keyboard (even if it's nil) will be ignored.
// // Allowed types:
// // - *tgbotapi.InlineKeyboardMarkup
// // - *tgbotapi.ReplyKeyboardMarkup
// // - *tgbotapi.ReplyKeyboardHide
// // All other types will be ignored.
// func (ctx *TCtx) Keyboard(keyboard interface{}) *TCtx {
// 	if ctx.alloc() == nil { return ctx }
// 	if keyboard == nil {
// 		ctx.g.keyboard, ctx.g.isDeleteKeyboard = nil, true
// 		return ctx
// 	}
// 	switch keyboard.(type) {
// 	case *tgbotapi.ReplyKeyboardMarkup, *tgbotapi.ReplyKeyboardHide,
// 	*tgbotapi.InlineKeyboardMarkup:
// 		ctx.g.keyboard, ctx.g.isDeleteKeyboard = keyboard, false
// 	default:
// 	}
// 	return ctx
// }

// 'Body' assigns 'body' as content of message that will be sent.
// If after space trimming body will be empty, there is no-op.
// It works only for 'New' or 'Edit' finishers.
// If you will use 'Del', 'DelLast', 'DelAll', 'DelAllExceptLast' finishers,
// the content of message will be ignored.
// Use methods 'HTML(true)' or 'MD(true)' if you want to enable
// HTML or Markdown support respectively.
func (ctx *TCtx) Body(body string) *TCtx {
	if ctx.alloc() == nil {
		return ctx
	}
	if body = strings.TrimSpace(body); body != "" {
		ctx.g.text = body
	}
	return ctx
}

// Finisher.
// 'New' tries to put together the all data that has been provided
// before by preparing methods and then tries to create config of
// new outgoing message.
// You must call at least 'Body' method before calling this finisher,
// otherwise this method do nothing.
// Method takes only one bool argument that answers to the question:
// "Shouldn't the successfully sent message's id be added to the session's
// message ids list?"
// If you pass 'true', the message will just send, but the id of successfully
// sent message willn't be added to the session message's ids.
func (ctx *TCtx) New(dontAddToSession ...bool) *TCtx {
	if !ctx.isBotObjectValid() || ctx.g == nil {
		return ctx
	}
	if ctx.g.text == "" {
		return ctx
	}
	// Create Telegram new message config, fill it
	nmsg := tgbotapi.MessageConfig{}
	nmsg.ChatID, nmsg.Text = ctx.Chat.ID, ctx.g.text
	// Realize the parse mode and attach it to the sendable config,
	// if it's not empty
	if mode := ctx.g.genParseMode(); mode != "" {
		nmsg.ParseMode = mode
	}
	// Realize the id of message, this message is reply to, attach it,
	// if it's not zero (reserved value - means that message isn't reply)
	if ctx.g.replyTo != 0 {
		nmsg.ReplyToMessageID = ctx.g.replyTo
	}
	// Generate keyboard and attach it to the sendable config
	if kb := ctx.g.genKb(); kb != nil {
		nmsg.ReplyMarkup = kb
	}
	// Create and defferring send the sendable config
	ctx.sender.deferSend(&tSendableConfig{
		config: nmsg, ctx: ctx,
		onSuccess: ctx.g.onSuccess, onError: ctx.g.onError,
		isUpdateSession: !(len(dontAddToSession) > 0 && dontAddToSession[0]),
		retryAttempts:   ctx.g.retryAttempts,
	})
	return ctx.cleanup2()
}

// Finisher.
// 'EditLast' tries to put together the all data that has been provided
// before by preparing methods and then tries to create config of
// editing the last message from bot in the current session
// (if it is or the most last message from bot).
// You must call at least 'Body' method before calling this finisher,
// otherwise this method do nothing.
// Warning! If session object is nil (depends by registered Sessioner)
// or session object doesn't have any stored message id,
// the 'New' method (w/o args) will be called.
// It means that in this case the new message will be created (instead
// of editing previous, 'cause "previous" doesn't exists) and
// when it will be sent successfully, the message id of sent message will be
// added to the session object (if it's not nil).
func (ctx *TCtx) EditLast() *TCtx {
	if !ctx.isBotObjectValid() || ctx.g == nil {
		return ctx
	}
	if ctx.g.text == "" {
		return ctx
	}
	// If no one message id is stored in session, just send new message
	if ctx.Sess == nil || ctx.Sess.MessagerEmpty() {
		return ctx.New()
	}
	// Try to extract last bot's message from the session
	msgId := ctx.Sess.MessagerPeek()
	// Create Telegram edit message config, fill it
	emsg := tgbotapi.EditMessageTextConfig{}
	emsg.ChatID, emsg.MessageID, emsg.Text = ctx.Chat.ID, msgId, ctx.g.text
	// Realize the parse mode and attach it to the sendable config,
	// if it's not empty
	if mode := ctx.g.genParseMode(); mode != "" {
		emsg.ParseMode = mode
	}
	// Generate keyboard and attach it to the sendable config
	// (only if it's inline keyboard)
	if ikb, ok := ctx.g.genKb().(*tgbotapi.InlineKeyboardMarkup); ok {
		emsg.ReplyMarkup = ikb
	}
	// Create and deferring send the sendable config
	ctx.sender.deferSend(&tSendableConfig{
		config: emsg, ctx: ctx,
		onSuccess: ctx.g.onSuccess, onError: ctx.g.onError,
		// tSession must not be updated, 'cause 'Edit' finisher assumes that
		// message id of message that must be edited already stored
		// in the session object (from which the message id was extracted)
		isUpdateSession: false,
		retryAttempts:   ctx.g.retryAttempts,
	})
	return ctx.cleanup2()
}

// Finisher.
// 'DelLast' tries to delete the last sent message by bot in chat for
// the current session.
// Note! If the current session doesn't have any message from bot,
// this method will do nothing.
func (ctx *TCtx) DelLast() *TCtx {
	if !ctx.isBotObjectValid() {
		return ctx
	}
	// If no one message id is stored in session, just do nothing
	if ctx.Sess == nil || ctx.Sess.MessagerEmpty() {
		return ctx
	}
	// Try to get last bot's message from the session
	msgId := ctx.Sess.MessagerPop()
	// Create Telegram delete message config, fill it
	dmsg := tgbotapi.DeleteMessageConfig{}
	dmsg.ChatID, dmsg.MessageID = ctx.Chat.ID, msgId
	// Create and deferring send the sendable config
	ctx.sender.deferSend(&tSendableConfig{
		config: dmsg, ctx: ctx,
		onSuccess: ctx.g.onSuccess, onError: ctx.g.onError,
		// tSession must not be updated, 'cause 'Del' finisher already
		// pops (not peeks) the deleting message id from session object
		isUpdateSession: false,
		retryAttempts:   ctx.g.retryAttempts,
	})
	return ctx.cleanup2()
}

// Finisher.
// 'DelAll' tries to delete absolutely all bot's messages in chat for the
// current session.
// Warning! The more messages sent by bot, the more delete configs must
// be sent, but Telegram has limits (See 'Lirester' for details).
// IT IS NOT RECOMMENDED TO USE THIS FEATURE!
// Let Telegram take care about it.
// Note! If the current session doesn't have any message from bot,
// this method will do nothing.
func (ctx *TCtx) DelAll() *TCtx {
	if !ctx.isBotObjectValid() {
		return ctx
	}
	// If no one message id is stored in session, just do nothing
	if ctx.Sess == nil || ctx.Sess.MessagerEmpty() {
		return ctx
	}
	// Try to get all bot's messages from the session
	msgIds := ctx.Sess.MessagerClean()
	// Create Telegram delete message config
	dmsg := tgbotapi.DeleteMessageConfig{}
	dmsg.ChatID = ctx.Chat.ID
	// Loop over all message ids which messages must be deleted
	// Assign the message id to the delete config and deferring sending it
	for _, msgId := range msgIds {
		dmsg.MessageID = msgId
		// Create and deferring send the sendable config
		ctx.sender.deferSend(&tSendableConfig{
			config: dmsg, ctx: ctx,
			onSuccess: ctx.g.onSuccess, onError: ctx.g.onError,
			// tSession must not be updated, 'cause 'DelAll' finisher already
			// deleted all message ids.
			isUpdateSession: false,
			retryAttempts:   ctx.g.retryAttempts,
		})
	}
	return ctx.cleanup2()
}

// Finisher.
// 'DelAllExceptLast' tries to delete absolutely all bot's messages in chat
// for the current session except the last message.
// Warning! The more messages sent by bot, the more delete configs must
// be sent, but Telegram has limits (See 'Lirester' for details).
// IT IS NOT RECOMMENDED TO USE THIS FEATURE!
// Let Telegram take care about it.
// Note! If the current session contains only one message from bot,
// this method will do nothing.
func (ctx *TCtx) DelAllExceptLast() *TCtx {
	if !ctx.isBotObjectValid() {
		return ctx
	}
	// If no one message id is stored in session, just do nothing
	if ctx.Sess == nil || ctx.Sess.MessagerEmpty() {
		return ctx
	}
	// Try to get all bot's messages from the session
	msgIds := ctx.Sess.MessagerClean()
	// If only one message is stored, do nothing
	if len(msgIds) == 1 {
		return ctx
	}
	// Store back the last message
	ctx.Sess.MessagerPush(msgIds[len(msgIds)-1])
	// All messages except last must be deleted
	msgIds = msgIds[:len(msgIds)-1]
	// Create Telegram delete message config
	dmsg := tgbotapi.DeleteMessageConfig{}
	dmsg.ChatID = ctx.Chat.ID
	// Loop over all message ids which messages must be deleted
	// Assign the message id to the delete config and deferring sending it
	for _, msgId := range msgIds {
		dmsg.MessageID = msgId
		// Create and deferring send the sendable config
		ctx.sender.deferSend(&tSendableConfig{
			config: dmsg, ctx: ctx,
			onSuccess: ctx.g.onSuccess, onError: ctx.g.onError,
			// tSession must not be updated, 'cause 'DelAll' finisher already
			// deleted all message ids.
			isUpdateSession: false,
			retryAttempts:   ctx.g.retryAttempts,
		})
	}
	return ctx.cleanup2()
}

// 'isBotObjectValid' checks whether the current context object isn't nil
// and contains valid bot object, that contains valid Sender object.
// Returns bool result of that checks.
func (ctx *TCtx) isBotObjectValid() bool {
	return ctx != nil && ctx.sender != nil && ctx.Chat != nil
}

//
func (ctx *TCtx) alloc() *TCtx {
	if !ctx.isBotObjectValid() {
		return nil
	}
	if ctx.g == nil {
		ctx.g = &tCtxMessageGenerator{}
	}
	return ctx
}

// // 'cleanup' performs cleanup context operation after one of finisher call.
// // Cleanup operation means that all stored prepared data
// // (stored by prepared methods) will be deleted.
// func (ctx *TCtx) cleanup() *TCtx {
// 	ctx.g.onSuccess, ctx.g.onError = nil, nil
// 	ctx.g.retryAttempts, ctx.g.replyTo, ctx.g.text = 0, 0, ""
// 	ctx.g.isHTML = ctx.bot.sender.consts.isAlwaysUseHTML
// 	ctx.g.isMarkdown = ctx.bot.sender.consts.isAlwaysUseMD
// 	ctx.g.keyboard, ctx.g.isDeleteKeyboard = nil, false
// 	return ctx
// }

//
func (ctx *TCtx) cleanup2() *TCtx {
	ctx.g = nil
	return ctx
}



//
func (c *tCtxCreator) CtxFor(chatId int64) *TCtx {

}

//
func (c *tCtxCreator) CtxForCb(chatId int64, cb tHandlerCallback) {

}

//
func makeCtx(bot *TBot, update *tgbotapi.Update) *TCtx {

}
