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
	"reflect"
	"unsafe"

	"github.com/qioalice/gext/dangerous"
)

// tRegistrator is a part of tReceiver type.
// todo: comment what is it and how it works
// todo: comment why *[]FViewHandler instead []FViewHandler, etc
// todo: comment why unsafe.Pointer as *[]FViewHandler, *[]FViewMiddlewre, *[]unsafe.Pointer
type tRegistrator struct {

	// Link to View ID converter (used in tRegistrator.saveAccumulated).
	converter *tViewIDConverter

	// [ 1 SECTION ]
	// handlers, middlewares storage:
	// occurred event type -> current View ID -> occurred event data,
	// except text events (handlerTextWhen field).
	handlersWhen       map[tEventType]map[tViewIDEncoded]map[tEventData]unsafe.Pointer // T of 3rd map's value is *[]FViewHandler
	handlersWhenExt    map[tEventType]map[tViewIDEncoded]map[tEventData]unsafe.Pointer // T of 3rd map's value is *[]unsafe.Pointer
	middlewaresWhen    map[tEventType]map[tViewIDEncoded]map[tEventData]unsafe.Pointer // T of 3rd map's value is *[]FViewMiddleware
	middlewaresWhenExt map[tEventType]map[tViewIDEncoded]map[tEventData]unsafe.Pointer // T of 3rd map's value is *[]unsafe.Pointer

	// [ 2 SECTION ]
	// handlers, middlewares storage:
	// occurred event type -> occurred event data,
	// except text events (handlerTextJust field).
	handlersJust       map[tEventType]map[tEventData]unsafe.Pointer // T of 2nd map's value is *[]FViewHandler
	handlersJustExt    map[tEventType]map[tEventData]unsafe.Pointer // T of 2nd map's value is *[]unsafe.Pointer
	middlewaresJust    map[tEventType]map[tEventData]unsafe.Pointer // T of 2nd map's value is *[]FViewMiddleware
	middlewaresJustExt map[tEventType]map[tEventData]unsafe.Pointer // T of 2nd map's value is *[]unsafe.Pointer

	// [ 3 SECTION ]
	// main handlers, middlewares storage.
	handlersMain       unsafe.Pointer // T is *[]FViewHandler
	handlersMainExt    unsafe.Pointer // T is *[]unsafe.Pointer
	middlewaresMain    unsafe.Pointer // T is *[]FViewMiddleware
	middlewaresMainExt unsafe.Pointer // T is *[]unsafe.Pointer

	// [ 4 SECTION ]
	// text handlers, middlewares storage by current View ID.
	handlerTextWhen        map[tViewIDEncoded]unsafe.Pointer // T of value is *[]FViewHandler
	handlerTextWhenExt     map[tViewIDEncoded]unsafe.Pointer // T of value is *[]unsafe.Pointer
	middlewaresTextWhen    map[tViewIDEncoded]unsafe.Pointer // T of value is *[]FViewMiddleware
	middlewaresTextWhenExt map[tViewIDEncoded]unsafe.Pointer // T of value is *[]unsafe.Pointer

	// [ 5 SECTION ]
	// text handlers, middlewares.
	handlerTextJust        unsafe.Pointer // T is *[]FViewHandler
	handlerTextJustExt     unsafe.Pointer // T is *[]unsafe.Pointer
	middlewaresTextJust    unsafe.Pointer // T is *[]FViewMiddleware
	middlewaresTextJustExt unsafe.Pointer // T is *[]unsafe.Pointer

	// REGISTERING NEW EVENTS

	// todo: comment
	tCtxExtended reflect.Type

	// todo: comment
	tHandlerCtxExtended reflect.Type

	// todo: comment
	tMiddlewareCtxExtended reflect.Type

	// todo: comment
	onRegistered []tEventRegistered

	// todo: comment
	regErrors []tRegistratorError
}

// Text marks that the callback passed into the next Handler or Middleware
// methods will be called when text event will be occurred.
func (r *tRegistrator) Text(when ...string) *tRegistrator {

	if r != nil {
		event := makeEventRegistered(CEventTypeText, "", when)
		r.onRegistered = append(r.onRegistered, *event)
	}

	return r
}

// Button marks that the callback passed into the next Handler or Middleware
// methods will be called when keyboard button event will be occurred and
// pressed button will be the same as what argument.
func (r *tRegistrator) Button(what string, when ...string) *tRegistrator {

	if r != nil && what != "" {
		event := makeEventRegistered(CEventTypeKeyboardButton, tEventType(what), when)
		r.onRegistered = append(r.onRegistered, *event)
	}

	return r
}

// Buttons marks that the callback passed into the next Handler or Middleware
// methods will be called when keyboard button event will be occurred and
// pressed button will be the same as any from whats argument.
func (r *tRegistrator) Buttons(whats []string, when ...string) *tRegistrator {

	for _, what := range whats {
		r = r.Button(what, when)
	}

	return r
}

// InlineButton marks that the callback passed into the next Handler or Middleware
// methods will be called when inline keyboard button event will be occurred and
// pressed inline button will be the same as what argument.
func (r *tRegistrator) InlineButton(what string, when ...string) *tRegistrator {

	if r != nil && what != "" {
		event := makeEventRegistered(CEventTypeInlineKeyboardButton, tEventType(what), when)
		r.onRegistered = append(r.onRegistered, *event)
	}

	return r
}

// InlineButtons marks that the callback passed into the next Handler or Middleware
// methods will be called when inline keyboard button event will be occurred and
// pressed inline button will be the same as any from whats argument.
func (r *tRegistrator) InlineButtons(whats []string, when ...string) *tRegistrator {

	for _, what := range whats {
		r = r.InlineButton(what, when)
	}

	return r
}

// Command marks that the callback passed into the next Handler or Middleware
// methods will be called when command event will be occurred and
// requested command will be the same as what argument.
func (r *tRegistrator) Command(what string, when ...string) *tRegistrator {

	if what != "" && (what[0] == '/' || what[0] == '\\') {
		what = what[1:]
	}

	if r != nil && what != "" {
		event := makeEventRegistered(CEventTypeInlineKeyboardButton, tEventType(what), when)
		r.onRegistered = append(r.onRegistered, *event)
	}

	return r
}

// Commands marks that the callback passed into the next Handler or Middleware
// methods will be called when command event will be occurred and
// requested command will be the same as any from whats argument.
func (r *tRegistrator) Commands(whats []string, when ...string) *tRegistrator {

	for _, what := range whats {
		r = r.Command(what, when)
	}

	return r
}

// Handler links all accumulated events by Text, Button, InlineButton, Command
// methods with passed handler (and then list of all accumulated events will be cleared).
//
// If there are no accumulated events, the handler will be registered as
// main handler (handles all events).
//
// handler should be FViewHandler type or the same type as your context extender
// returns, if you extends the context.
//
// ATTENTION!
// If any error will occur while trying to register handler,
// list of all accumulated events won't be cleared!
func (r *tRegistrator) Handler(handler interface{}) *TBot {

	if r == nil {
		return nil
	}

	var err *tRegistratorError

	switch {
	case handler == nil:
		err = Errors.EventRegistration.NilHandler.cp()
		err = err.tw(reflect.TypeOf(FViewHandler(nil)).String())

	case !r.parent.isCtxExtended:
		if typedHandler, ok := handler.(FViewHandler); ok {
			if typedHandler != nil {
				err = r.saveWhenDefaultCtx(typedHandler, nil)
			} else {
				err = Errors.EventRegistration.NilHandler.cp()
				err = err.tw(reflect.TypeOf(FViewHandler(nil)).String())
			}
		} else {
			err = Errors.EventRegistration.IncompatibleContextType.cp()
			err = err.th(reflect.TypeOf(handler).String())
			err = err.tw(reflect.TypeOf(FViewHandler(nil)).String())
		}

	case r.parent.isCtxExtended:
		err = r.saveWhenExtendedCtx(handler, false)
	}

	if err != nil {
		r.regErrors = append(r.regErrors, *err.e(r.onRegistered))
	}

	// if r != nil, tRegistrator is private type and it guarantees,
	// that parent != nil and parent.parent also.
	return r.parent.parent
}

// Middleware links all accumulated events by Text, Button, InlineButton, Command
// methods with passed middleware (and then list of all accumulated events will be cleared).
//
// If there are no accumulated events, the middleware will be registered as
// main middleware (checks all events).
//
// middleware should be FViewMiddleware type or the same type as your context extender
// returns, if you extends the context.
//
// ATTENTION!
// If any error will occur while trying to register middleware,
// list of all accumulated events won't be cleared!
func (r *tRegistrator) Middleware(middleware interface{}) *TBot {

	if r == nil {
		return nil
	}

	var err *tRegistratorError

	switch {
	case middleware == nil:
		err = Errors.EventRegistration.NilMiddleware.cp()
		err = err.tw(reflect.TypeOf(FViewMiddleware(nil)).String())

	case !r.parent.isCtxExtended:
		if typedMiddleware, ok := middleware.(FViewMiddleware); ok {
			if typedMiddleware != nil {
				err = r.saveWhenDefaultCtx(nil, typedMiddleware)
			} else {
				err = Errors.EventRegistration.NilMiddleware.cp()
				err = err.tw(reflect.TypeOf(FViewMiddleware(nil)).String())
			}
		} else {
			err = Errors.EventRegistration.IncompatibleContextType.cp()
			err = err.th(reflect.TypeOf(middleware).String())
			err = err.tw(reflect.TypeOf(FViewMiddleware(nil)).String())
		}

	case r.parent.isCtxExtended:
		err = r.saveWhenExtendedCtx(middleware, true)
	}

	if err != nil {
		r.regErrors = append(r.regErrors, *err.e(r.onRegistered))
	}

	// if r != nil, tRegistrator is private type and it guarantees,
	// that parent != nil and parent.parent also.
	return r.parent.parent
}

// MainHandler register handler as main handler (handles all events).
// You could just use Handler method instead with no accumulated events,
// but using this method improves code readability.
//
// handler should be FViewHandler type or the same type as your context extender
// returns, if you extends the context.
func (r *tRegistrator) MainHandler(handler interface{}) *TBot {

	if r == nil {
		return nil
	}

	// Handler will be registered as main handler only when onRegistered
	// field is empty.

	onRegistered := r.onRegistered
	r.onRegistered = nil
	returned := r.Handler(handler)
	r.onRegistered = onRegistered
	return returned
}

// MainMiddleware register middleware as main middleware (checks all events).
// You could just use Middleware method instead with no accumulated events,
// but using this method improves code readability.
//
// middleware should be FViewMiddleware type or the same type as your context extender
// returns, if you extends the context.
func (r *tRegistrator) MainMiddleware(middleware interface{}) *TBot {

	if r == nil {
		return nil
	}

	// Middleware will be registered as main middleware only when onRegistered
	// field is empty.

	onRegistered := r.onRegistered
	r.onRegistered = nil
	returned := r.Middleware(middleware)
	r.onRegistered = onRegistered
	return returned
}

// saveWhenDefaultCtx performs linking all accumulated events with passed
// default view handler h and default view middleware m
// (one of them can be nil, but not both at the same time).
func (r *tRegistrator) saveWhenDefaultCtx(h FViewHandler, m FViewMiddleware) *tRegistratorError {

	// At this code point h or m is valid (not nil) callback.

	switch {
	case h != nil:
		cbptr = unsafe.Pointer(&h)
	case m != nil:
		cbptr = unsafe.Pointer(&m)
	}

	return r.saveAccumulated(cbptr, m != nil, false)
}

// saveWhenExtendedCtx performs linking all accumulated events with passed
// view handler or view middleware user callback cb, when context is extended.
// isMiddleware reports whether cb is view handler or view middleware.
func (r *tRegistrator) saveWhenExtendedCtx(cb interface{}, isMiddleware bool) *tRegistratorError {

	// cb is not nil interface at this code point
	// (checked by Handler/Middleware methods).

	cbType := reflect.TypeOf(cb)

	if !isMiddleware && cbType != r.tHandlerCtxExtended {
		err := Errors.EventRegistration.IncompatibleContextType.cp()
		err = err.th(cbType.String())
		err = err.tw(r.tHandlerCtxExtended.String())
		return err
	}

	if isMiddleware && cbType != r.tMiddlewareCtxExtended {
		err := Errors.EventRegistration.IncompatibleContextType.cp()
		err = err.th(cbType.String())
		err = err.tw(r.tMiddlewareCtxExtended.String())
		return err
	}

	// if cbType.Kind() != reflect.Func {
	// 	// return err - incompatible type
	// }

	// // cb must take only one argument
	// if cbType.NumIn() != 1 {
	// 	// return err - incompatible type
	// }

	// // cb must take argument by pointer
	// if cbArgType := cbType.In(0); cbArgType.Kind() == reflect.Ptr {

	// 	// and type of argument passed by its pointer should be the same
	// 	// as the type returned by user context extender
	// 	if cbArgType = cbArgType.Elem(); cbArgType != r.tCtxExtended {
	// 		// return err - incompatible type
	// 	}
	// } else {
	// 	// return err - incompatible type
	// }

	// // cb must not have return arguments if cb is handler
	// if !isMiddleware && cbType.NumOut() != 0 {
	// 	// return err - incompatible type
	// }

	// if isMiddleware {
	// 	// cb must have only one return argument if cb is middleware
	// 	if cbType.NumOut() != 1 {
	// 		// return err - incompatible type
	// 	}

	// 	// cb must have bool as type of returned argument if cb is middleware
	// 	if retArgType := cbType.Out(0); retArgType.Kind() != reflect.Bool {
	// 		// return err - incompatible type
	// 	}
	// }

	cbPtr := dangerous.FnPtrCallable(cb)
	if cbPtr == nil {
		err := Errors.EventRegistration.InternalError.cp()
		err = err.msg("gext.dangerous.FnPtrCallable in saveWhenExtendedCtx return nil")
		return err
	}

	// At this code point cb is an untyped pointer to valid handler or middleware
	// (not nil, arguments types are valid, return types are valid, etc).
	return r.saveAccumulated(cbPtr, isMiddleware, true)
}

// saveAccumulated links all accumulated events with passed untyped callback
// (may be handler, middleware; for default or extended context).
func (r *tRegistrator) saveAccumulated(cb unsafe.Pointer, isMiddleware, isCtxExtended bool) []tRegistratorError {

	// Reg as general handler/middleware
	if r.onRegistered == nil {
		r.save(cb, cEventTypeInvalid, cEvDa, viewID, isMiddleware, isCtxExtended)
		return nil
	}

	var (
		converter = r.parent.parent.converter
		errs      []tRegistratorError
	)

	for _, eventRegister := range r.onRegistered {

		if len(eventRegister.When) != 0 {
			for _, when := range eventRegister.When {

				if whenEncoded, encodeErr := converter.Encode(when); encodeErr == nil {
					r.save(cb, eventRegister.Type, eventRegister.Data, whenEncoded, isMiddleware, isCtxExtended)

				} else {
					// todo: Add error grouping
					err := Errors.EventRegistration.InternalError.cp()
					err = err.msg(encodeErr.What)
					err = err.e2(*makeEventRegistered(eventRegister.Type, eventRegister.Data, []tViewID{when}))
					errs = append(errs, err)
				}
			}
		} else {
			r.save(cb, eventRegister.Type, eventRegister.Data, cViewIDEncodedNull, isMiddleware, isCtxExtended)
		}

	}

	return errs
}

// match returns a slice of handlers or slice of middlewares (isMiddleware flag)
// associated with an event, identification signs of which are passed.
// A real type of return value will be one of:
// *[]FViewHandler, *[]FViewMiddleware, *[]unsafe.Pointer.
// More info: tRegistrator.access.
func (r *tRegistrator) match(typ tEventType, data tEventData, viewID tViewIDEncoded, isMiddleware, isCtxExtended bool) unsafe.Pointer {
	return r.access(nil, typ, data, viewID, isMiddleware, isCtxExtended)
}

// save links cb as handler or middleware (isMiddleware flag)
// of default or extended context (isCtxExtended flat)
// with an event identification signs of which are passed.
// More info: tRegistrator.access.
func (r *tRegistrator) save(cb unsafe.Pointer, typ tEventType, data tEventData, viewID tViewIDEncoded, isMiddleware, isCtxExtended bool) {
	r.access(cb, typ, data, viewID, isMiddleware, isCtxExtended)
}

// access does the one of two things tRegistrator.save and tRegistrator.match describes.
// It depends on whether cb is nil or not. Returns nil if works in "save" mode
// or if requested callbacks not found in "match" mode.
// Detailed description inside.
func (r *tRegistrator) access(cb unsafe.Pointer, typ tEventType, data tEventData, viewID tViewIDEncoded, isMiddleware, isCtxExtended bool) unsafe.Pointer {

	// save infers storage type using isMiddleware, isCtxExtended flags,
	// checks whether storage is nil, and if it is so, allocates memory,
	// and then after all saves cb to storage, returns storage.
	save := func(storage, cb unsafe.Pointer, isMiddleware, isCtxExtended bool) unsafe.Pointer {

		switch {
		case !isMiddleware && !isCtxExtended:
			storageTypedPtr := (*[]FViewHandler)(storage)
			if storageTypedPtr == nil {
				storageData := make([]FViewHandler, 0, 1)
				storageTypedPtr = &storageData
			}
			*storageTypedPtr = append(*storageTypedPtr, *(*FViewHandler)(cb))
			storage = unsafe.Pointer(storageTypedPtr)

		case isMiddleware && !isCtxExtended:
			storageTypedPtr := (*[]FViewMiddleware)(storage)
			if storageTypedPtr == nil {
				storageData := make([]FViewMiddleware, 0, 1)
				storageTypedPtr = &storageData
			}
			*storageTypedPtr = append(*storageTypedPtr, *(*FViewMiddleware)(cb))
			storage = unsafe.Pointer(storageTypedPtr)

		case isCtxExtended:
			storageTypedPtr := (*[]unsafe.Pointer)(storage)
			if storageTypedPtr == nil {
				storageData := make([]unsafe.Pointer, 0, 1)
				storageTypedPtr = &storageData
			}
			*storageTypedPtr = append(*storageTypedPtr, cb)
			storage = unsafe.Pointer(storageTypedPtr)
		}
		return storage
	}

	// Typedefs described below are created for a more compact way
	// to describe read/write operations with tRegistrator's storages.
	//
	// S means "section", L - "level".
	// Thus, S1L3 - section 1, level 3 - an alias to the map type of the last
	// "view" in 1st storage section callbacks.

	type tS1L3 map[tEventData]unsafe.Pointer
	type tS1L2 map[tViewIDEncoded]tS1L3
	type tS1L1 map[tEventType]tS1L2

	type tS2L2 map[tEventData]unsafe.Pointer
	type tS2L1 map[tEventType]tS2L2

	type tS4L1 map[tViewIDEncoded]unsafe.Pointer

	// ptrField is a pointer to some tRegistrator field.
	// In all switch cases presented below the first action is initializing
	// ptrField.
	// It allows to write one code-snippet for 4 different cases:
	// handlers storage, handlers for extended context storage,
	// middlewares storage, middlewares for extended context storage.
	var ptrField unsafe.Pointer

	isReg := cb != nil

	// WARNING!
	// All algorithms and operations presented below (in switch cases)
	// are performed with an untyped pointers (like void* in C) to decrease
	// amount of code and avoid unnecessary type-checking in some cases.
	// BE CAREFUL BEFORE MAKING ANY CHANGES!

	// NOTE.
	// A nested switch presented below is a readability way to represent
	// a many nested if-else statements.
	// Golang guarantees that all cases are checked left-right, top-bottom.
	// Using that feature, we can omit some bool expressions as "else" branch
	// if their positive pairs has been checked by branches above.

	// NOTE.
	// All switch cases are numbered and have the same number as the section number
	// in the tRegistator's fields description.

	switch {

	// 3 SECTION: "GLOBAL CALLBACKS".
	case typ == cEventTypeInvalid:

		switch {
		case !isMiddleware && !isCtxExtended:
			ptrField = unsafe.Pointer(&r.handlersMain)

		case !isMiddleware && isCtxExtended:
			ptrField = unsafe.Pointer(&r.handlersMainExt)

		case isMiddleware && !isCtxExtended:
			ptrField = unsafe.Pointer(&r.middlewaresMain)

		case isMiddleware && isCtxExtended:
			ptrField = unsafe.Pointer(&r.middlewaresMainExt)
		}

		if isReg {
			storage := *(*unsafe.Pointer)(ptrField)
			storage = save(storage, cb, isMiddleware, isCtxExtended)
			*(*unsafe.Pointer)(ptrField) = storage
		} else {
			return *(*unsafe.Pointer)(ptrField)
		}

	// 5 SECTION: "TEXT CALLBACKS W/O VIEW ID".
	case typ == CEventTypeText && viewID == cViewIDEncodedNull:

		switch {
		case !isMiddleware && !isCtxExtended:
			ptrField = unsafe.Pointer(&r.handlerTextJust)

		case !isMiddleware && isCtxExtended:
			ptrField = unsafe.Pointer(&r.handlerTextJustExt)

		case isMiddleware && !isCtxExtended:
			ptrField = unsafe.Pointer(&r.middlewaresTextJust)

		case isMiddleware && isCtxExtended:
			ptrField = unsafe.Pointer(&r.middlewaresTextJustExt)
		}

		if isReg {
			storage := *(*unsafe.Pointer)(ptrField)
			storage = save(storage, cb, isMiddleware, isCtxExtended)
			*(*unsafe.Pointer)(ptrField) = storage
		} else {
			return *(*unsafe.Pointer)(ptrField)
		}

	// 4 SECTION "TEXT CALLBACKS WITH VIEW ID"
	// Condition "viewID != cViewIDEncodedNull" is omitted.
	case typ == CEventTypeText:

		switch {
		case !isMiddleware && !isCtxExtended:
			ptrField = unsafe.Pointer(&r.handlerTextWhen)

		case !isMiddleware && isCtxExtended:
			ptrField = unsafe.Pointer(&r.handlerTextWhenExt)

		case isMiddleware && !isCtxExtended:
			ptrField = unsafe.Pointer(&r.middlewaresTextWhen)

		case isMiddleware && isCtxExtended:
			ptrField = unsafe.Pointer(&r.middlewaresTextWhenExt)
		}

		if isReg {
			if *(*tS4L1)(ptrField) == nil {
				*(*tS4L1)(ptrField) = make(tS4L1)
			}
			storage := (*(*tS4L1)(ptrField))[viewID]
			storage = save(storage, cb, isMiddleware, isCtxExtended)
			(*(*tS4L1)(ptrField))[viewID] = storage
		} else {
			return (*(*tS4L1)(ptrField))[viewID]
		}

	// 2 SECTION "ALL OTHER CALLBACKS W/O VIEW ID"
	// Condition "typ != CEventTypeText" is omitted.
	// Conditions "typ == <typ1> || ... || typ == <typN>" are omitted.
	case viewID == cViewIDEncodedNull:

		switch {
		case !isMiddleware && !isCtxExtended:
			ptrField = unsafe.Pointer(&r.handlersJust)

		case !isMiddleware && isCtxExtended:
			ptrField = unsafe.Pointer(&r.handlersJustExt)

		case isMiddleware && !isCtxExtended:
			ptrField = unsafe.Pointer(&r.middlewaresJust)

		case isMiddleware && isCtxExtended:
			ptrField = unsafe.Pointer(&r.middlewaresJustExt)
		}

		if isReg {
			if *(*tS2L1)(ptrField) == nil {
				*(*tS2L1)(ptrField) = make(tS2L1)
			}
			if (*(*tS2L1)(ptrField))[typ] == nil {
				(*(*tS1L1)(ptrField))[typ] = make(tS2L2)
			}
			storage := (*(*tS2L1)(ptrField))[typ][data]
			storage = save(storage, cb, isMiddleware, isCtxExtended)
			(*(*tS2L1)(ptrField))[typ][data] = storage
		} else {
			return (*(*tS2L1)(ptrField))[typ][data]
		}

	// 1 SECTION "ALL OTHER CALLBACKS WITH VIEW ID"
	// Condition "typ != CEventTypeText" is omitted.
	// Conditions "typ == <typ1> || ... || typ == <typN>" are omitted.
	// Condition "viewID != cViewIDEncodedNull" is omitted.
	default:

		switch {
		case !isMiddleware && !isCtxExtended:
			ptrField = unsafe.Pointer(&r.handlersWhen)

		case !isMiddleware && isCtxExtended:
			ptrField = unsafe.Pointer(&r.handlersWhenExt)

		case isMiddleware && !isCtxExtended:
			ptrField = unsafe.Pointer(&r.middlewaresWhen)

		case isMiddleware && isCtxExtended:
			ptrField = unsafe.Pointer(&r.middlewaresWhenExt)
		}

		if isReg {
			if *(*tS1L1)(ptrField) == nil {
				*(*tS1L1)(ptrField) = make(tS1L1)
			}
			if (*(*tS1L1)(ptrField))[typ] == nil {
				(*(*tS1L1)(ptrField))[typ] = make(tS1L2)
			}
			if (*(*tS1L1)(ptrField))[typ][viewID] == nil {
				(*(*tS1L1)(ptrField))[typ][viewID] = make(tS1L3)
			}
			storage := (*(*tS1L1)(ptrField))[typ][viewID][data]
			storage = save(storage, cb, isMiddleware, isCtxExtended)
			(*(*tS1L1)(ptrField))[typ][viewID][data] = storage
		} else {
			return (*(*tS1L1)(ptrField))[typ][viewID][data]
		}

	}

	// All registration operations presented above do not return something
	// and lead here.
	return nil
}

// todo: comment me, implement me
func makeRegistrator() *tRegistrator {
	panic("implement me")
}
