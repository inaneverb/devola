// Copyright © 2019. All rights reserved.
// Author: Alice Qio.
// Contacts: <qioalice@gmail.com>.
// License: https://opensource.org/licenses/MIT

package registrator

import (
	"reflect"
	"unsafe"

	"github.com/qioalice/devola/core/event"
	"github.com/qioalice/devola/core/sys/fn"
	"github.com/qioalice/devola/core/view"
)

// Registrator is a special part of Devola project that provides registering,
// storing, keeping, extracting, matching, finding and processing of all
// handlers and middlewares.
//
// You should use Simple or Complex method to accumulate rules to which
// some next callback should be applied and then call Handler or Middleware method
// to flush it.
// You can also use MainHandler or MainMiddleware methods to set up general,
// regular callbacks.
//
// But! You must pass callback with compatible type with your context type!
// If your context type has been changed, use RegenerateRequiredTypes method
// to updated restriction rules.
type Registrator struct {

	// Link to View ID converter (used in Registrator.save method).
	converter *view.IDConv

	// There is 3 important entity using which you can get or store any callback:
	// event type (event.Type), event data (event.Data) and view ID encoded (view.IDEnc),
	// where
	// event type is the type of occurred event,
	// event data is the body of occurred event,
	// view ID encoded is the current View ID in the session associated with occurred event.
	//
	// By default you should determine all these 3 things for each callback
	// at the registering operation and then use them to get registered callback.
	//
	// But it is possible that you don't use sessions (and you don't have View ID)
	// for example, or it's meaningless to determine events by its data
	// because it is empty or too complex and different (for example if it is
	// text message and text message can be absolutely any), or you want to
	// register or get the callback that will be called for all untyped events, etc.
	//
	// Because of that there is 5 sections of storage.

	// [ 1 SECTION ]
	// handlers, middlewares storage:
	// occurred event type -> current View ID -> occurred event data.
	handlersWhen    map[event.Type]map[view.IDEnc]map[event.Data][]unsafe.Pointer
	middlewaresWhen map[event.Type]map[view.IDEnc]map[event.Data][]unsafe.Pointer

	// [ 2 SECTION ]
	// handlers, middlewares storage:
	// occurred event type -> occurred event data.
	handlersJust    map[event.Type]map[event.Data][]unsafe.Pointer
	middlewaresJust map[event.Type]map[event.Data][]unsafe.Pointer

	// [ 3 SECTION ]
	// main handlers, middlewares storage.
	handlersMain    []unsafe.Pointer // T is *[]FViewHandler
	middlewaresMain []unsafe.Pointer // T is *[]FViewMiddleware

	// [ 4 SECTION ]
	// occurred event type -> current View ID.
	handlerTextWhen     map[event.Type]map[view.IDEnc][]unsafe.Pointer
	middlewaresTextWhen map[event.Type]map[view.IDEnc][]unsafe.Pointer

	// [ 5 SECTION ]
	// occurred event type.
	handlerTextJust     map[event.Type][]unsafe.Pointer
	middlewaresTextJust map[event.Type][]unsafe.Pointer

	// Generated type the registered handlers must have.
	handlerTypeRequired reflect.Type

	// Generated type the registered middlewares must have.
	middlewareTypeRequired reflect.Type

	// The set of rules to which next generated callback (handler or middleware)
	// will be applied.
	// Is accumulated by Simple or Complex, is flushed by Handler or Middleware.
	accumulatedRules []rule

	// The function that determines what event type can be considered "simple"
	// and what can not.
	// "Simple" is the type of event that will be handled by callbacks from
	// 4 or 5 sections - independent by event data.
	simplesChecker func(typ event.Type) (isSimple bool)
}

// Text marks that the callback passed into the next Handler or Middleware
// methods will be called when text event will be occurred.
func (r *Registrator) Simple(typ event.Type, when []string) *Registrator {
	return r.Complex(typ, string(event.CDataNil), when)
}

// InlineButton marks that the callback passed into the next Handler or Middleware
// methods will be called when inline keyboard button event will be occurred and
// pressed inline button will be the same as what argument.
func (r *Registrator) Complex(typ event.Type, what string, when []string) *Registrator {
	// *(*[]view.ID)(unsafe.Pointer(&when)) is
	// []string -> []view.ID conversion without memory reallocation
	e := makeRule(typ, event.Data(what), *(*[]view.ID)(unsafe.Pointer(&when)))
	r.accumulatedRules = append(r.accumulatedRules, *e)
	return r
}

// Handler links all accumulated events by Simple or Complex methods
// with passed handler (and then list of all accumulated events will be cleared).
//
// If there are no accumulated events, the handler will be registered as
// main handler (handles all events).
//
// handler type should be "func(*T)" where T is the same type as backend Ctx
// or the same as your context extender returns, if you extends the context.
// Otherwise a not nil EBadCallback error object is returned.
//
// ATTENTION!
// If any error will occur while trying to register handler,
// list of all accumulated events will not be cleared!
func (r *Registrator) Handler(handler interface{}) *EBadCallback {
	return r.save(handler, false, r.handlerTypeRequired)
}

// Middleware links all accumulated events by Simple or Complex methods
// with passed middleware (and then list of all accumulated events will be cleared).
//
// If there are no accumulated events, the middleware will be registered as
// main middleware (checks all events).
//
// middleware type should be "func(*T) bool" where T is the same type as backend Ctx
// or the same as your context extender returns, if you extends the context.
// Otherwise a not nil EBadCallback error object is returned.
//
// ATTENTION!
// If any error will occur while trying to register middleware,
// list of all accumulated events will not be cleared!
func (r *Registrator) Middleware(middleware interface{}) *EBadCallback {
	return r.save(middleware, true, r.middlewareTypeRequired)
}

// MainHandler register handler as main handler (handles all events).
// You could just use Handler method instead with no accumulated events,
// but using this method improves code readability.
//
// handler type should be "func(*T)" where T is the same type as backend Ctx
// or the same as your context extender returns, if you extends the context.
// Otherwise a not nil EBadCallback error object is returned.
func (r *Registrator) MainHandler(handler interface{}) *EBadCallback {

	// Handler will be registered as main handler only when accumulatedRules
	// field is empty.

	currentRules := r.accumulatedRules
	r.accumulatedRules = nil
	returned := r.Handler(handler)
	r.accumulatedRules = currentRules
	return returned
}

// MainMiddleware register middleware as main middleware (checks all events).
// You could just use Middleware method instead with no accumulated events,
// but using this method improves code readability.
//
// middleware type should be "func(*T) bool" where T is the same type as backend Ctx
// or the same as your context extender returns, if you extends the context.
// Otherwise a not nil EBadCallback error object is returned.
func (r *Registrator) MainMiddleware(middleware interface{}) *EBadCallback {

	// Middleware will be registered as main middleware only when accumulatedRules
	// field is empty.

	currentRules := r.accumulatedRules
	r.accumulatedRules = nil
	returned := r.Middleware(middleware)
	r.accumulatedRules = currentRules
	return returned
}

// RegenerateRequiredTypes updates handlerTypeRequired and middlewareTypeRequired
// fields by new generated types that will be depended from ctxType.
// It allows to register handler or middlewares with a new signature
// after context has been extended.
func (r *Registrator) RegenerateRequiredTypes(ctxType reflect.Type) {

	// It excepts ctxType is not nil

	inBoth := []reflect.Type{ctxType}
	outMiddleware := []reflect.Type{reflect.TypeOf(true)}

	r.handlerTypeRequired = reflect.FuncOf(inBoth, nil, false)
	r.middlewareTypeRequired = reflect.FuncOf(inBoth, outMiddleware, false)
}

// Match returns a slice of handlers or slice of middlewares (isMiddleware flag)
// associated with an event, identification signs of which are passed.
// A real type of return value will be *[]unsafe.Pointer.
func (r *Registrator) Match(typ event.Type, data event.Data, viewID view.IDEnc, isMiddleware bool) []unsafe.Pointer {
	return r.access(nil, typ, data, viewID, isMiddleware)
}

// save performs linking all accumulated events with passed
// handler or middleware and returns nil if it was successfully.
//
// cb is handler or middleware functor,
// isMiddleware reports whether cb is handler or middleware,
// wantType is one of r.handlerTypeRequired or r.middlewareTypeRequired.
func (r *Registrator) save(cb interface{}, isMiddleware bool, wantType reflect.Type) *EBadCallback {

	// cb == nil also handlers here
	if haveType := reflect.TypeOf(cb); haveType != wantType {
		return makeEBadCallback(isMiddleware, haveType == nil, wantType.String(), haveType.String(), r.accumulatedRules)
	}

	// FnPtrCallable may return nil only if cb == nil, there is no need a nil check
	cbPtr := fn.TakeCallableAddr(cb)

	// Reg as general handler/middleware
	if r.accumulatedRules == nil {
		r.access(cbPtr, event.CTypeInvalid, event.CDataNil, view.CIDEncNil, isMiddleware)
		return nil
	}

	for _, rule := range r.accumulatedRules {
		if len(rule.When) == 0 {
			r.access(cbPtr, rule.Type, rule.Data, view.CIDEncNil, isMiddleware)
			continue
		}
		for _, when := range rule.When {
			viewID := r.converter.Encode(when)
			r.access(cbPtr, rule.Type, rule.Data, viewID, isMiddleware)
		}
	}

	return nil
}

// access does the one of two things Registrator.save and Registrator.Match describes.
// It depends on whether cb is nil or not. Returns nil if works in "save" mode
// or if requested callbacks not found in "Match" mode.
// Detailed description inside.
func (r *Registrator) access(cb unsafe.Pointer, typ event.Type, data event.Data, viewID view.IDEnc, isMiddleware bool) []unsafe.Pointer {

	// Typedefs described below are created for a more compact way
	// to describe read/write operations with Registrator's storages.
	//
	// S means "section", L - "level".
	// Thus, S1L3 - section 1, level 3 - an alias to the map type of the last
	// "view" in 1st storage section callbacks.

	type tS1L3 map[event.Data][]unsafe.Pointer
	type tS1L2 map[view.IDEnc]tS1L3
	type tS1L1 map[event.Type]tS1L2

	type tS2L2 map[event.Data][]unsafe.Pointer
	type tS2L1 map[event.Type]tS2L2

	type tS4L2 map[view.IDEnc][]unsafe.Pointer
	type tS4L1 map[event.Type]tS4L2

	type tS5L1 map[event.Type][]unsafe.Pointer

	// ptrField is a pointer to some Registrator field.
	// In all switch cases presented below the first action is initializing
	// ptrField.
	// It allows to write one code-snippet for both different cases:
	// handlers storage and middlewares storage.
	var ptrField unsafe.Pointer

	isReg := cb != nil
	isSimpleType := typ != event.CTypeInvalid && r.simplesChecker(typ)

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
	// in the Registator's fields description.

	switch {

	// 3 SECTION: "GLOBAL CALLBACKS".
	case typ == event.CTypeInvalid:

		switch {
		case !isMiddleware:
			ptrField = unsafe.Pointer(&r.handlersMain)

		case isMiddleware:
			ptrField = unsafe.Pointer(&r.middlewaresMain)
		}

		if isReg {
			storage := *(*[]unsafe.Pointer)(ptrField)
			storage = append(storage, cb)
			*(*[]unsafe.Pointer)(ptrField) = storage
		} else {
			return *(*[]unsafe.Pointer)(ptrField)
		}

	// 5 SECTION: "TEXT CALLBACKS W/O VIEW ID".
	case isSimpleType && viewID == view.CIDEncNil:

		switch {
		case !isMiddleware:
			ptrField = unsafe.Pointer(&r.handlerTextJust)

		case isMiddleware:
			ptrField = unsafe.Pointer(&r.middlewaresTextJust)
		}

		if isReg {
			if *(*tS5L1)(ptrField) == nil {
				*(*tS5L1)(ptrField) = make(tS5L1)
			}
			storage := (*(*tS5L1)(ptrField))[typ]
			storage = append(storage, cb)
			(*(*tS5L1)(ptrField))[typ] = storage
		} else {
			return (*(*tS5L1)(ptrField))[typ]
		}

	// 4 SECTION "TEXT CALLBACKS WITH VIEW ID"
	// Condition "viewID != view.CIDEncNil" is omitted.
	case isSimpleType:

		switch {
		case !isMiddleware:
			ptrField = unsafe.Pointer(&r.handlerTextWhen)

		case isMiddleware:
			ptrField = unsafe.Pointer(&r.middlewaresTextWhen)
		}

		if isReg {
			if *(*tS4L1)(ptrField) == nil {
				*(*tS4L1)(ptrField) = make(tS4L1)
			}
			if (*(*tS4L1)(ptrField))[typ] == nil {
				(*(*tS4L1)(ptrField))[typ] = make(tS4L2)
			}
			storage := (*(*tS4L1)(ptrField))[typ][viewID]
			storage = append(storage, cb)
			(*(*tS4L1)(ptrField))[typ][viewID] = storage
		} else {
			return (*(*tS4L1)(ptrField))[typ][viewID]
		}

	// 2 SECTION "ALL OTHER CALLBACKS W/O VIEW ID"
	// Condition "isSimpleType" is omitted.
	// Conditions "typ == <typ1> || ... || typ == <typN>" are omitted.
	case viewID == view.CIDEncNil:

		switch {
		case !isMiddleware:
			ptrField = unsafe.Pointer(&r.handlersJust)

		case isMiddleware:
			ptrField = unsafe.Pointer(&r.middlewaresJust)
		}

		if isReg {
			if *(*tS2L1)(ptrField) == nil {
				*(*tS2L1)(ptrField) = make(tS2L1)
			}
			if (*(*tS2L1)(ptrField))[typ] == nil {
				(*(*tS2L1)(ptrField))[typ] = make(tS2L2)
			}
			storage := (*(*tS2L1)(ptrField))[typ][data]
			storage = append(storage, cb)
			(*(*tS2L1)(ptrField))[typ][data] = storage
		} else {
			return (*(*tS2L1)(ptrField))[typ][data]
		}

	// 1 SECTION "ALL OTHER CALLBACKS WITH VIEW ID"
	// Condition "isSimpleType" is omitted.
	// Conditions "typ == <typ1> || ... || typ == <typN>" are omitted.
	// Condition "viewID != view.CIDEncNil" is omitted.
	default:

		switch {
		case !isMiddleware:
			ptrField = unsafe.Pointer(&r.handlersWhen)

		case isMiddleware:
			ptrField = unsafe.Pointer(&r.middlewaresWhen)
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
			storage = append(storage, cb)
			(*(*tS1L1)(ptrField))[typ][viewID][data] = storage
		} else {
			return (*(*tS1L1)(ptrField))[typ][viewID][data]
		}

	}

	// All registration operations presented above return nothing and lead here.
	// We shouldn't return something if it was a registration.
	return nil
}

// MakeRegistrator creates a new Registrator object, initializes it with passed
// view ID converter object and simples event type's checker.
// It also sets that registered handlers or middlewares should be compatible
// with passed context type.
func MakeRegistrator(converter *view.IDConv, ctxType reflect.Type, simplesChecker func(typ event.Type) (isSimple bool)) *Registrator {

	var r Registrator

	r.converter = converter
	r.simplesChecker = simplesChecker
	r.RegenerateRequiredTypes(ctxType)

	return &r
}
