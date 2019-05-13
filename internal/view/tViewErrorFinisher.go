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
)

// todo: Add disabling panic guard possibility using constructor or smth else
// todo: Add tests

// tViewErrorFinisher is the internal type and represents all error callbacks
// (and all stuff linked with it) of one unsent Telegram message.
//
// More info: FViewErrorFinisher, TCtx, tSender.
type tViewErrorFinisher struct {

	// Enable inheritance from tViewBaseFinisher.
	tViewBaseFinisher

	// err is sending error
	err error

	// Standard context object's callbacks (predefined type).
	cbs []FViewErrorFinisher
}

// invoke safety (if panic guard is enabled) calls cbo (an error view finisher
// which receives a standard context) passing ctx.
func (cb *tViewErrorFinisher) invoke(cbo FViewErrorFinisher, ctx *TCtx) {

	if cb.flags.TestFlag(cSendCallbackEnablePanicGuard) {
		defer cb.protectFromPanic()
	}
	cbo(ctx, cb.err)
}

// Call calls all stored error finishers (predefined type finishers and
// extended finishers) passing standard or extended context object and
// error object to them.
//
// Each call will be panic protected if panic guard is enabled.
//
// After all calls
// (regardless of whether there was a panic in one or in several of them)
// a session transaction will be finished (if it is enabled) and then
// a chat transaction will be finished (if it is enabled).
//
// WARNING!
// If a session transaction wasn't finished,
// a chat transaction will also not be finished!
func (cb *tViewErrorFinisher) Call() {

	// Call callbacks with splitting behaviour for standard context
	// and extended context.
	if cb.flags.TestFlag(cSendCallbackUseExtendedContext) {

		// Use extended context

		// You may think that if statement is redundant (it is in the loop below),
		// but it is to don't create args if it doesn't need.
		if len(cb.tViewBaseFinisher.cbs) > 0 {

			args := []reflect.Value{
				reflect.ValueOf(cb.ctx),
				reflect.ValueOf(cb.err),
			}

			for _, cboExtended := range cb.tViewBaseFinisher.cbs {
				cb.tViewBaseFinisher.invoke(cboExtended, args)
			}
		}

	} else {

		// Use standard context

		// You may think that if statement is redundant (it is in the loop below),
		// but it is to don't convert ctx object if it doesn't need.
		if len(cb.cbs) > 0 {

			ctx := cb.ctx.(*TCtx)
			for _, cboStandard := range cb.cbs {
				cb.invoke(cboStandard, ctx)
			}
		}
	}

	cb.trFinish()
}

// makeViewErrorFinisher creates a new tViewErrorFinisher object
// using sendable config cfg and err as occurred error.
func makeViewErrorFinisher(cfg *tSendableConfig, err error) *tViewErrorFinisher {

	// Enable panic guard by default, enable use extended context if
	// callbacks that it receives are.
	flags := cSendCallbackEnablePanicGuard
	if len(cfg.errorExtendedFinishers) == 0 {
		flags.SetFlag(cSendCallbackUseExtendedContext)
	}

	finisher := &tViewErrorFinisher{}

	finisher.flags = flags
	finisher.err = err
	finisher.ctx = cfg.ctx

	finisher.cbs = cfg.errorStandardFinishers
	finisher.tViewBaseFinisher.cbs = cfg.errorExtendedFinishers

	return finisher
}
