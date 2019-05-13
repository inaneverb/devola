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

// 'tReger' is class that simplifies the registration of handlers for
// various events.
// So, first you calling the method with name of type of event that you
// want to register:
// 'Text' for text events, 'Button' for keyboard button, 'InlineButton'
// for callback event (button of inline keyboard) and 'Command' for
// register some command.
// You pass required arguments to one of these methods and then you
// calling 'Handler' or 'Middleware' method and pass callback which
// you want to be called for event that you specified in previous call.
// For example: 'Text().Handler(func...)'



// tBehaviourCreator tries to make a link between some events that will be 
// occurred and callback that should be called as a reaction. 
//
// When you calls Text, Button, or some other "event declaration" method,
// it creates a new tEventRegistered object, pushes it to the events field.
// Next, when you call "callback registration" method (handler or middleware),
// it links ALL accumulated events with the passed callback
type tBehaviourCreator struct {

	// Pointer to receiver in which the new registered event will be applied.
	receiver *tReceiver

	// Events that will be registered as a signal to some handler triggering.
	events []tEventRegistered
}

// 'Text' prepares internal objects to registration
// the handler (or the middleware) for text event, that will be called only
// when session type will be equal to one of 'when' types.
// So, if you want the handler (or the middleware) will be called for
// each text event (regardless of session type) just do not pass anything
// as argument, pass nil or 'AnyType'.
func (rr *tReger) Text(when ...session.TStep) *tReger {
	if rr == nil {
		return nil
	}
	if len(when) != 0 {
		rr.whens = append(rr.whens, when...)
	}
	e := tRegistrableEvent{typ: HandlerTypeText}
	rr.events = append(rr.events, e)
	return rr
}

// 'Button' prepares internal objects to registration
// the handler (or the middleware) for keyboard button event,
// that will be called only for buttons with 'what' text and when session type
// will be equal to one of 'when' types.
// So, if you want the handler (or the middleware) will be called for
// all keyboard buttons with specified 'what' text (regardless of session type)
// just do not pass anything as argument, pass nil or 'AnyType'.
func (rr *tReger) Button(what string, when ...session.TStep) *tReger {
	if rr == nil {
		return nil
	}
	if what = strings.TrimSpace(what); what == "" {
		return nil
	}
	if len(when) != 0 {
		rr.whens = append(rr.whens, when...)
	}
	e := tRegistrableEvent{typ: HandlerTypeKeyboardButton, data: what}
	rr.events = append(rr.events, e)
	return rr
}

// 'Buttons' the same as 'Button' but allows you to register a lot of
// buttons at once in one function call.
// It very useful when you have different languages in your app and
// of course keyboard buttons for each language.
// Just generate button's text slice and pass it to call of 'Buttons' method.
func (rr *tReger) Buttons(whats []string, when ...session.TStep) *tReger {
	if rr == nil || len(whats) == 0 {
		return nil
	}
	if len(when) != 0 {
		rr.whens = append(rr.whens, when...)
	}
	for _, what := range whats {
		if what = strings.TrimSpace(what); what == "" {
			continue
		}
		e := tRegistrableEvent{typ: HandlerTypeKeyboardButton, data: what}
		rr.events = append(rr.events, e)
	}
	return rr
}

// 'InlineButton' prepares internal objects to registration
// the handler (or the middleware) for inline keyboard button event,
// that will be called only for inline buttons with 'what' action and
// when session type will be equal to one of 'when' types.
// So, if you want the handler (or the middleware) will be called for
// all inline keyboard buttons with specified 'what' action
// (regardless of session type) just do not pass anything as argument,
// pass nil or 'AnyType'.
func (rr *tReger) InlineButton(
	what session.Action, when ...session.TStep,
) *tReger {
	if rr == nil {
		return nil
	}
	if len(when) != 0 {
		rr.whens = append(rr.whens, when...)
	}
	e := tRegistrableEvent{typ: HandlerTypeInlineButton, data: what.String()}
	rr.events = append(rr.events, e)
	return rr
}

// 'InlineButtons' the same as 'InlineButton' but allows you to register
// a lot of inline keyboard buttons at once in one function call.
func (rr *tReger) InlineButtons(
	whats []session.Action, when ...session.TStep,
) *tReger {
	if rr == nil || len(whats) == 0 {
		return nil
	}
	if len(when) != 0 {
		rr.whens = append(rr.whens, when...)
	}
	for _, what := range whats {
		e := tRegistrableEvent{typ: HandlerTypeInlineButton, data: what.String()}
		rr.events = append(rr.events, e)
	}
	return rr
}

// 'Command' prepares internal objects to registration
// the handler (or the middleware) for command event,
// that will be called only for command with 'what' text as body of command
// and when session type will be equal to one of 'when' types.
// So, if you want the handler (or the middleware) will be called for
// all commands with specified 'what' body (regardless of session type)
// just do not pass anything as argument, pass nil or 'AnyType'.
func (rr *tReger) Command(what string, when ...session.TStep) *tReger {
	if rr == nil {
		return nil
	}
	if what = strings.TrimSpace(what); what == "" {
		return nil
	}
	if len(when) != 0 {
		rr.whens = append(rr.whens, when...)
	}
	e := tRegistrableEvent{typ: HandlerTypeCommand, data: what}
	rr.events = append(rr.events, e)
	return rr
}

// 'Commands' the same as 'Command' but allows you to register a lot of
// commands at once in one function call.
// It very userful when you have some command and aliases for it.
// Just generate command's body slice and pass it to call of 'Command' method.
func (rr *tReger) Commands(whats []string, when ...session.TStep) *tReger {
	if rr == nil || len(whats) == 0 {
		return nil
	}
	if len(when) != 0 {
		rr.whens = append(rr.whens, when...)
	}
	for _, what := range whats {
		if what = strings.TrimSpace(what); what == "" {
			continue
		}
		e := tRegistrableEvent{typ: HandlerTypeCommand, data: what}
		rr.events = append(rr.events, e)
	}
	return rr
}

// 'Handler' registers the 'cb' as callback for all events that was
// specified early by calling 'Text', 'Button', 'InlineButton' or 'Command'
// methods or by combination of them.
// So, if no one call of these methods was, it means that 'cb' will
// registered as global handler callback.
// Global handler callback is callback that will be called for each
// events which doesn't have any other specified registered callbacks.
// Note! After calling 'Handler' or 'Middleware' method, all stored metadata
// about events for which next call of 'Handler' or 'Middleware' must
// register callback, will be deleted.
// It means, that chaining:
// rr.Text().Handler(cb1).Button(btn).Handler(cb2)
// will register 'cb1' as callback for all text events and 'cb2' as callback
// for keyboard button with text 'btn'.
func (rr *tReger) Handler(cb tHandlerCallback) *tReger {
	if rr == nil || cb == nil {
		return nil
	}
	// If there are no saved events, it means that registered callback
	// must be one of global handlers
	if len(rr.events) == 0 {
		rr.receiver.handlersMain = append(rr.receiver.handlersMain, cb)
		return rr.cleanup()
	}
	// Perform registration 'cb' callback for each saved event
	for _, event := range rr.events {
		if event.typ == HandlerTypeText {
			rr.receiver.regHandlerText(rr.whens, cb)
		} else {
			rr.receiver.regHandlerOther(rr.whens, event.typ, event.data, cb)
		}
	}
	return rr.cleanup()
}

// 'Middleware' is the same as 'Handler' but registers a middleware,
// not a handler.
// So you can read what is middleware in the docs for 'Receiver' class.
// Also, see 'Handler' docs to understand all the subtleties of this method.
// ('Handler' and 'Middleware' are very similar methods and works the same way,
// but one registers the handlers and the other - middlewares).
func (rr *tReger) Middleware(cb tMiddlewareCallback) *tReger {
	if rr == nil || cb == nil {
		return nil
	}
	// If there are no saved events, it means that registered callback
	// must be one of global handlers
	if len(rr.events) == 0 {
		rr.receiver.mdlwsMain = append(rr.receiver.mdlwsMain, cb)
		return rr.cleanup()
	}
	// Perform registration 'cb' callback for each saved event
	for _, event := range rr.events {
		if event.typ == HandlerTypeText {
			rr.receiver.regMiddlewareText(rr.whens, cb)
		} else {
			rr.receiver.regMiddlewareOther(rr.whens, event.typ, event.data, cb)
		}
	}
	return rr.cleanup()
}

// 'cleanup' performs the cleaning process of 'tReger' object after
// calling 'Handler' or 'Middleware' method.
func (rr *tReger) cleanup() *tReger {
	rr.events, rr.whens = nil, nil // slices will be deleted by GC
	return rr
}