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

// tViewBaseFinisher is the internal type and represents the common base of
// tViewSuccessFinisher and tViewErrorFinisher types.
// This is the base type for them.
//
// Moreover, while these types represents success/error callbacks that
// receives standard context object, this type already has a storage
// for success/error callbacks that receives extended context object
// and the storage in which both of standard and extended context objects
// can be stored.
//
// More info: tViewSuccessFinisher, tViewErrorFinisher.
type tViewBaseFinisher struct {

	// flags determines the behaviour Call method in derived types
	// (tViewSuccessFinisher, tViewErrorFinisher)
	// such as enable panic guard, enable session transaction finisher,
	// enable chat transaction finisher, etc).
	flags tSendCallbackFlag

	// ctx is context by which an successfully sent or unsent Telegram message
	// was created.
	// Can be used to store both standard and extended contexts.
	ctx iCtx

	// All occurred panics in callbacks will be stored here
	// and panic won't shutdown server (if panic guard enabled).
	recoveredPanics []interface{}

	// If session transaction finisher or chat transaction finisher
	// will be completed with error, this error will be here,
	// and type error (session transaction finisher of chat transaction finisher)
	// can be figure out by flags field.
	afterAllErr error

	// Extended context object's callbacks (user defined type) casted to
	// reflect.Value to speed up.
	cbs []reflect.Value
}

// protectFromPanic tries to recover panic, and if it was successfull,
// saves the recovered panic info to the panics field in current cb object
// to analyse it in the caller code.
func (cb *tViewBaseFinisher) protectFromPanic() {

	if recoverInfo := recover(); recoverInfo != nil {
		cb.recoveredPanics = append(cb.recoveredPanics, recoverInfo)
	}
}

// invoke safety (if panic checker is enabled) calls cbo (a view finisher which
// receives an extended context) passing args as call arguments using reflect.Call.
func (cb *tViewBaseFinisher) invoke(cbo reflect.Value, args []reflect.Value) {

	if cb.flags.TestFlag(cSendCallbackEnablePanicGuard) {
		defer cb.protectFromPanic()
	}
	cbo.Call(args)
}

// trFinish tries to finish open session and chat transactions if it is need.
//
// WARNING!
// If a session transaction wasn't finished,
// a chat transaction will also not be finished!
func (cb *tViewBaseFinisher) trFinish() {

	// Finish session transaction (if it's need)
	// Stop doing next things if error is occurred
	if cb.flags.TestFlag(cSendCallbackFinishSessionTransaction) {
		if err := cb.ctx.TrSessFinish(); err != nil {

			cb.afterAllErr = err
			cb.flags.SetFlag(cSendCallbackFinishSessionTransactionError)
			return
		}
	}

	// Finish chat transaction (if it's need)
	if cb.flags.TestFlag(cSendCallbackFinishChatTransaction) {
		if err := cb.ctx.TrChatFinish(); err != nil {

			cb.afterAllErr = err
			cb.flags.SetFlag(cSendCallbackFinishChatTransactionError)
		}
	}
}
