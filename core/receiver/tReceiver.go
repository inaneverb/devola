// Copyright © 2018. All rights reserved.
// Author: Alice Qio.
// Contacts: <qioalice@gmail.com>.
// License: https://opensource.org/licenses/MIT

package receiver

import (
	"fmt"
	"log"
	"reflect"
	"unsafe"

	api "github.com/go-telegram-bot-api/telegram-bot-api"

	"../registrator"
)

// -- Receiver --
// 'Receiver' is the part of Telegram TBot UIG Framework.
// 'Receiver' is engaged is receiving incoming messages as bot updates,
// converting updates to the context objects, if it's possible,
// and calls the registered middlewares and handlers.
//
// 'Receiver' must be created by 'NewReceiver' function directly
// (way to specify some 'Receiver' options) or by internal 'TBot' constructor.
//
// As object, 'Receiver' contains logger pointer, bot pointer,
// updates channel (to which the Telegram TBot updates are coming),
// map of registered handlers (specified for event types),
// map of registered middlewares (specified for event types),
// general handlers and middlewares (if special handlers or middlewares
// aren't registered) and executor's channel.
// See 'tEventType', 'tHandlerCallback', 'tMiddlewareCallback' types
// for details, it's important to understand how it works and what is it.
//
// How it works.
// 1. Telegram sends update to the bot server. TBot server receives
// that update and creates 'tgbotapi.Update' object by that update.
// 2. Receiver takes each such update object and tries to find
// executors for it (registered handlers).
// 2.1. Before that, Receiver tries to construct Event object that will
// contain info about event of received update. That object will contain
// info "what kind of event" and "what is body of event".
// 3. Then Receiver tries to create context object for received update.
// (See 'TCtx' type for details). That object will be passed to the
// registered middlewares and handlers.
// 4. Receiver performs all middleware checks for created context object.
// If middlewares allowed the next processing, the context object wraps in
// 'tExecutor' object, that will contain context object and all executors
// which must be called for that context object.
// 5. Receiver pass 'tExecutor' object to the executor's channel and
// some free executor goroutine will take it from channel and execute it.
//
// What is the middleware?
// Middleware is the special callback that can be registered for any event,
// takes context object as argument and must return true or false.
// So, when any event received from Telegram servers, and if
// middleware registered for that event, the middleware(s) will be called
// before handler(s).
// If at least one of middlewares returns 'false' no one next of middlewares
// and handlers will be called!
// Thus with middlewares you can manage calls of handlers.
// Keep in mind, that all middlewares executes in one goroutine for
// each events, and should not be too algorithmically complex,
// otherwise the perfomance may decline.
type Receiver struct {

	// TBot object, this Receiver object associated with.
	endpoint *api.BotAPI

	// Consts section.
	consts struct {
		maxRKBTextLen int

		executorsChanLen int

		executorsCount int
	}

	registrator *registrator.Registrator

	behaviourType tReceiverBehaviourType

	isRun bool

	chUpdates api.UpdatesChannel

	chMustBeExecuted chan tExecutor

	isCtxExtended bool

	ctxExtender reflect.Value

	tCtxExtended reflect.Type
}

//
func (r *Receiver) PushUpdate()

// 'isAllowByMdlws' returns 'true' if all middlewares that can be called
// for event stored in context object 'ctx' with specified 'when' type
// will return 'true' (they will allow futher procession).
// If false returned by at least one of middlewares, the 'isAllowByMdlws'
// will return false immediately.
func (ctx *Ctx) isAllowedByMiddlewares(when session.TStep) bool {
	panicChecker := func(ctx *Ctx) {
		recoverErr := recover()
		if recoverErr == nil {
			return
		}
		r.bot.log.Named("Receiver")
		if err := recover(); err != nil {
			r.bot.log.Class("Receiver").Method("isAllowByMdlws").Errorw(
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

//
func (r *Receiver) run() error {

}

//
func (r *Receiver) stop() error {

}

//
func (r *Receiver) restart() error {

}

//
func (r *Receiver) ExtendContext(extender interface{}) error {

}

//
func (r *Receiver) ExtendContext2(f func(c *Ctx) unsafe.Pointer) {

}

//
func (r *Receiver) serveUpdates() {
	log.Println("Receiver.serveUpdates",
		"Serving receiving the Telegram updates successfully started "+
			"at the separated goroutine")
	for upd := range r.chUpdates {
		r.serveUpdate(&upd)
	}
	log.Println("Receiver.serveUpdates",
		"Serving receiving the Telegram updates successfully stopped. "+
			"The separated goroutine has been shutdown.")
}

//
func (r *Receiver) serveUpdate(update *tgbotapi.Update) {
	var (
		ctx  = makeCtx(r.bot, update)
		ssid = session.CIDNil
	)
	switch {

	// If update object is callback query (inline button event)
	case update.CallbackQuery != nil:
		ctx.Event.Type = HandlerTypeInlineButton
		ctx.Chat = update.CallbackQuery.Message.Chat
		ctx.From = update.CallbackQuery.From
		ssid,
			ctx.Event.Body,
			ctx.Event.args = actionDecode(update.CallbackQuery.Data)

	// If updateate object is text message, command or reply keyboard button
	case update.Message != nil:
		ctx.Chat = update.Message.Chat
		ctx.From = update.Message.From
		// If update is command
		if cmd := update.Message.CommandWithAt(); cmd != "" {
			// todo: Add command arguments supporting
			ctx.Event.Type = HandlerTypeCommand
			ctx.Event.Body = cmd
		} else {
			// If update is text or keyboard button, let's solve it in the future
			ctx.Event.Type = handlerTypeKeyboardButtonOrText
			ctx.Event.Body = update.Message.Text
		}

	// Unsupported event type
	default:
		log.Println("Receiver.serveUpdate",
			"Unhandled Telegram updateate",
			"Unsupported or unknown type of updateate",
			update)
		return
	}

	// Try to import session using ssid or chat id (already in ctx)
	// and then get session step (to recognize what handlers should be used)
	ctx.importSession(ssid)
	when := ctx.sessionStep()

	// Try to get handlers
	var hdlr []tHandlerCallback
	switch ctx.Event.Type {

	// If keyboard button or text was set as event type, we'll resolve it here
	case handlerTypeKeyboardButtonOrText:
		// First, check by body length. If it more than allowable limits,
		// we'll think that it's text, not a button
		if len(ctx.Event.Body) > r.consts.maxRKBTextLen {
			goto labelHandleAsText
		}
		// Let's believe that event is button.
		// We'll try to find handlers for button, and if it's not exists,
		// we'll change the mind.
		ctx.Event.Type = HandlerTypeKeyboardButton
		// If map by specified 'when' was initialized, then
		// tries to find handlers over map
		if r.handlersWhen[HandlerTypeKeyboardButton][when] != nil {
			hdlr = r.handlersWhen[HandlerTypeKeyboardButton][when][ctx.Event.Body]
			// If when isn't specified (equals to AnyType), or by previous
			// operations handlers aren't found, tries to find
			// keyboard buttons handlers using event's body as 'what'
		} else if when == session.CStepAny || len(hdlr) == 0 {
			hdlr = r.handlersJust[HandlerTypeKeyboardButton][ctx.Event.Body]
			// If still handlers aren't found, handle it as text
		} else {
			goto labelHandleAsText
		}
		// If handlers aren't registered, also handle it as text
		if len(hdlr) == 0 {
			goto labelHandleAsText
		}
		// At this code point we successfully resolve the handlers for
		// keyboard event
		break

		// Handle event as text
	labelHandleAsText:
		// First, update the event type, init handlers by text 'when' handlers
		// It will be overwritten if 'when' is specified to 'TypeAny' or
		// registered handlers aren't found
		hdlr = r.handlerTextWhen[when]
		ctx.Event.Type = HandlerTypeText
		if when == session.CStepAny || len(hdlr) == 0 {
			hdlr = r.handlerTextJust
		}

	// Try to get handlers for inline button
	case HandlerTypeInlineButton:
		// Get registered handlers by 'when' for inline button if they're
		// registered
		if when != session.CStepAny &&
			r.handlersWhen[HandlerTypeInlineButton][when] != nil {
			hdlr = r.handlersWhen[HandlerTypeInlineButton][when][ctx.Event.Body]
		}
		// If length of handlers array is zero, try to apply general
		// inline keyboard handlers
		if len(hdlr) != 0 {
			break
		}
		hdlr = r.handlersJust[HandlerTypeInlineButton][ctx.Event.Body]

	// Try to get handlers for command
	case HandlerTypeCommand:
		// Get registered handlers by 'when' for command if they're registered
		if when != session.CStepAny &&
			r.handlersWhen[HandlerTypeCommand][when] != nil {
			hdlr = r.handlersWhen[HandlerTypeCommand][when][ctx.Event.Body]
		}
		// If length of handlers array is zero, try to apply general
		// command handlers
		if len(hdlr) != 0 {
			break
		}
		hdlr = r.handlersJust[HandlerTypeCommand][ctx.Event.Body]
	}

	// If still handlers aren't found, try to apply general handlers
	if len(hdlr) == 0 {
		hdlr = r.handlersMain
	}

	// Final check whether handlers are found
	if len(hdlr) == 0 {
		log.Println("Receiver.serveUpdate",
			"Unhandled Telegram updateate",
			"No handlers found for received Telegram updateate",
			ctx.Event, ctx.Sess, update)
		return
	}

	// At this code point we have handlers and resolved event
	// Try to apply all needed middlewares
	if ctx.isAllowedByMiddlewares(when) {
		return
	}

	// Create executor object and pass it to the executor's channel
	// It will be executed in the separated goroutines
	r.chMustBeExecuted <- makeExecutor(ctx, hdlr)
}

//
func (r *Receiver) serveExecutors() {
	r.bot.log.Class("Receiver").Method("serveUpdates").Debugw(fmt.Sprintf(
		"Serving executing registered handlers on the incoming updates "+
			"successfully started at the %d separated goroutine(s)",
		r.consts.executorsCount))
	panicChecker := func(ctx *TCtx) {
		if err := recover(); err != nil {
			r.bot.log.Class("Receiver").Method("serveExecutors").Errorw(
				"Panic occurred and has been restored while trying to apply handlers",
				"ctx", ctx, "panic_err", err)
		}
	}
	safeExecute := func(cb tHandlerCallback, ctx *TCtx) {
		defer panicChecker(ctx)
		cb(ctx)
	}
	executor := func() {
		for mustBeExecuted := range r.chMustBeExecuted {
			for _, executor := range mustBeExecuted.executors {
				safeExecute(executor, mustBeExecuted.ctx)
			}
		}
	}
	for i := 0; i < r.consts.executorsCount; i++ {
		go executor()
	}
	r.bot.log.Class("Receiver").Method("serveUpdates").Debugw(fmt.Sprintf(
		"Serving executing registered handlers on the incoming updates "+
			"successfully stopped. %d separated goroutine(s) has been shutdown.",
		r.consts.executorsCount))
}

//
func (r *Receiver) Serve() error {
	if r == nil {
		return fmt.Errorf("nil receiver object")
	}
	if r.isServed {
		return fmt.Errorf("already serving")
	}
}

// 'makeRegistrator' is auxiliary function that creates the 'tReger' objects,
// saves the pointer to the current receiver to created object and returns it.
// It's used for working registration methods in 'Receiver' class that
// in fact calls the same methods in the 'tReger' class.
func (r *Receiver) makeRegistrator() *tReger {
	if r == nil {
		return nil
	}
	return &tReger{receiver: r}
}

// 'isAnyWhen' returnss true if 'whens' is empty slice of session types
// or if at least one of these types is 'TypeAny'.
// Otherwise false is returned.
func (r *Receiver) isAnyWhen(whens []session.TStep) bool {
	if len(whens) == 0 {
		return true
	}
	for _, when := range whens {
		if when == session.TStepAny {
			return true
		}
	}
	return false
}

// 'regHandlerText' registers 'cb' as handler for text event, that will be
// called only when session type will be one of 'whens'.
// Note. If call 'isAnyWhen' for 'whens' is true, it means that 'cb' will
// be called for each text events, for which another callback with not
// 'TypeAny' 'when' willn't be registered by other registrator calls.
func (r *Receiver) regHandlerText(
	whens []session.TStep, cb tHandlerCallback,
) {
	if !r.isAnyWhen(whens) {
		for _, when := range whens {
			// todo: check whether 'when' is valid session.TStep
			if r.handlerTextWhen == nil {
				r.handlerTextWhen = make(map[session.TStep][]tHandlerCallback)
			}
			r.handlerTextWhen[when] = append(r.handlerTextWhen[when], cb)
		}
	} else {
		r.handlerTextJust = append(r.handlerTextJust, cb)
	}
}

// 'regHandlerOther' registers 'cb' as handler for event with 'typ' type and
// with 'what' body, that will be called only when session type will be one of
// 'whens'.
// Note. If call 'isAnyWhen' for 'whens' is true, it means that 'cb' will
// be called for all events with specified type and body, for which another
// callback with not 'TypeAny' 'when' willn't be registered by other
// registrator calls.
func (r *Receiver) regHandlerOther(
	whens []session.TStep, typ tEventType, what string, cb tHandlerCallback,
) {
	if !r.isAnyWhen(whens) {
		for _, when := range whens {
			// todo: check whether 'when' is valid session.TStep
			if r.handlersWhen == nil {
				r.handlersWhen =
					make(map[tEventType]map[session.TStep]map[string][]tHandlerCallback)
			}
			if r.handlersWhen[typ] == nil {
				r.handlersWhen[typ] =
					make(map[session.TStep]map[string][]tHandlerCallback)
			}
			if r.handlersWhen[typ][when] == nil {
				r.handlersWhen[typ][when] =
					make(map[string][]tHandlerCallback)
			}
			r.handlersWhen[typ][when][what] =
				append(r.handlersWhen[typ][when][what], cb)
		}
	} else {
		if r.handlersJust == nil {
			r.handlersJust =
				make(map[tEventType]map[string][]tHandlerCallback)
		}
		if r.handlersJust[typ] == nil {
			r.handlersJust[typ] =
				make(map[string][]tHandlerCallback)
		}
		r.handlersJust[typ][what] = append(r.handlersJust[typ][what], cb)
	}
}

// 'regMiddlewareText' registers 'cb' as middlewaree for text event,
// that will be called only when session type will be one of 'whens'.
// Note. If call 'isAnyWhen' for 'whens' is true, it means that 'cb' will
// be called for each text events, for which another callback with not
// 'TypeAny' 'when' willn't be registered by other registrator calls.
func (r *Receiver) regMiddlewareText(
	whens []session.TStep, cb tMiddlewareCallback,
) {
	if !r.isAnyWhen(whens) {
		for _, when := range whens {
			// todo: check whether 'when' is valid session.TStep
			if r.mdlwsTextWhen == nil {
				r.mdlwsTextWhen = make(map[session.TStep][]tMiddlewareCallback)
			}
			r.mdlwsTextWhen[when] = append(r.mdlwsTextWhen[when], cb)
		}
	} else {
		r.mdlwsTextJust = append(r.mdlwsTextJust, cb)
	}
}

// 'regMiddlewareOther' registers 'cb' as middleware for event with 'typ' type
// and with 'what' body, that will be called only when session type will be
// one of 'whens'.
// Note. If call 'isAnyWhen' for 'whens' is true, it means that 'cb' will
// be called for all events with specified type and body, for which another
// callback with not 'TypeAny' 'when' willn't be registered by other
// registrator calls.
func (r *Receiver) regMiddlewareOther(
	whens []session.TStep, typ tEventType, what string, cb tMiddlewareCallback,
) {
	if !r.isAnyWhen(whens) {
		for _, when := range whens {
			// todo: check whether 'when' is valid session.TStep
			if r.mdlwsWhen == nil {
				r.mdlwsWhen =
					make(map[tEventType]map[session.TStep]map[string][]tMiddlewareCallback)
			}
			if r.mdlwsWhen[typ] == nil {
				r.mdlwsWhen[typ] =
					make(map[session.TStep]map[string][]tMiddlewareCallback)
			}
			if r.mdlwsWhen[typ][when] == nil {
				r.mdlwsWhen[typ][when] =
					make(map[string][]tMiddlewareCallback)
			}
			r.mdlwsWhen[typ][when][what] = append(r.mdlwsWhen[typ][when][what], cb)
		}
	} else {
		if r.mdlwsJust == nil {
			r.mdlwsJust =
				make(map[tEventType]map[string][]tMiddlewareCallback)
		}
		if r.mdlwsJust[typ] == nil {
			r.mdlwsJust[typ] =
				make(map[string][]tMiddlewareCallback)
		}
		r.mdlwsJust[typ][what] = append(r.mdlwsJust[typ][what], cb)
	}
}

// 'Text' the same as 'tReger.Text'.
func (r *Receiver) Text(when ...session.TStep) *tReger {
	return r.makeRegistrator().Text(when...)
}

// 'Button' the same as 'tReger.Button'.
func (r *Receiver) Button(what string, when ...session.TStep) *tReger {
	return r.makeRegistrator().Button(what, when...)
}

// 'Buttons' the same as 'tReger.Buttons'.
func (r *Receiver) Buttons(what []string, when ...session.TStep) *tReger {
	return r.makeRegistrator().Buttons(what, when...)
}

// 'InlineButton' the same as 'tReger.InlineButton'.
func (r *Receiver) InlineButton(
	what session.Action, when ...session.TStep,
) *tReger {
	return r.makeRegistrator().InlineButton(what, when...)
}

// 'InlineButtons' the same as 'tReger.InlineButtons'.
func (r *Receiver) InlineButtons(
	what []session.Action, when ...session.TStep,
) *tReger {
	return r.makeRegistrator().InlineButtons(what, when...)
}

// 'Command' the same as 'tReger.Command'.
func (r *Receiver) Command(what string, when ...session.TStep) *tReger {
	return r.makeRegistrator().Command(what, when...)
}

// 'Commands' the same as 'tReger.Commands'.
func (r *Receiver) Commands(what []string, when ...session.TStep) *tReger {
	return r.makeRegistrator().Commands(what, when...)
}

// 'params' might be only 'tReceiverParam'
func makeReceiver(parent *TBot, params ...interface{}) (*Receiver, error) {
	r := &Receiver{bot: parent}
	// The default values of receiver params
	const (
		defMaxRKBTextLen    = 50
		defExecutorsChanLen = 1024 * 16
		defExecutorsCount   = 1
	)
	// Apply default receiver params.
	// It'll be overwritten later by passed receiver params.
	r.consts.maxRKBTextLen = defMaxRKBTextLen
	r.consts.executorsChanLen = defExecutorsChanLen
	r.consts.executorsCount = defExecutorsCount
	// Apply receiver params
	for _, param := range params {
		if param, ok := param.(tReceiverParam); ok && param != nil {
			param(r)
		}
	}
	// Allocate mem for handlers that will be called for special events
	// with some session type
	r.handlersWhen =
		make(map[tEventType]map[session.TStep]map[string][]tHandlerCallback)
	r.handlersWhen[HandlerTypeCommand] =
		make(map[session.TStep]map[string][]tHandlerCallback)
	r.handlersWhen[HandlerTypeKeyboardButton] =
		make(map[session.TStep]map[string][]tHandlerCallback)
	r.handlersWhen[HandlerTypeInlineButton] =
		make(map[session.TStep]map[string][]tHandlerCallback)
	// Allocate mem for handlers that will be called for events
	// excluding events with some session types
	r.handlersJust = make(map[tEventType]map[string][]tHandlerCallback)
	r.handlersJust[HandlerTypeCommand] = make(map[string][]tHandlerCallback)
	r.handlersJust[HandlerTypeKeyboardButton] = make(map[string][]tHandlerCallback)
	r.handlersJust[HandlerTypeInlineButton] = make(map[string][]tHandlerCallback)
	// Allocate mem for special text handlers
	r.handlerTextWhen = make(map[session.TStep][]tHandlerCallback)
	// Allocate mem for middlewares that will be called for special events
	// with some session type
	r.mdlwsWhen =
		make(map[tEventType]map[session.TStep]map[string][]tMiddlewareCallback)
	r.mdlwsWhen[HandlerTypeCommand] =
		make(map[session.TStep]map[string][]tMiddlewareCallback)
	r.mdlwsWhen[HandlerTypeKeyboardButton] =
		make(map[session.TStep]map[string][]tMiddlewareCallback)
	r.mdlwsWhen[HandlerTypeInlineButton] =
		make(map[session.TStep]map[string][]tMiddlewareCallback)
	// Allocate mem for middlewares that will be called for events
	// excluding events with some session types
	r.mdlwsJust = make(map[tEventType]map[string][]tMiddlewareCallback)
	r.mdlwsJust[HandlerTypeCommand] = make(map[string][]tMiddlewareCallback)
	r.mdlwsJust[HandlerTypeKeyboardButton] = make(map[string][]tMiddlewareCallback)
	r.mdlwsJust[HandlerTypeInlineButton] = make(map[string][]tMiddlewareCallback)
	// Allocate mem for special text middlewares
	r.mdlwsTextWhen = make(map[session.TStep][]tMiddlewareCallback)
	// Create executor's channel
	r.chMustBeExecuted = make(chan tExecutor, r.consts.executorsChanLen)
	// Receiver successfully created
	parent.receiver = r
	// Start serving if it's forced
	if r.consts.forceServe {
		return r, r.Serve()
	}
	return r, nil
}
