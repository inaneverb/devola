// Copyright © 2018. All rights reserved.
// Author: Alice Qio.
// Contacts: <qioalice@gmail.com>.
// License: https://opensource.org/licenses/MIT

package tusent

import (
	"unsafe"

	"go.uber.org/zap"

	"github.com/qioalice/devola/core/chat"
	"github.com/qioalice/devola/core/ctxutils"
	"github.com/qioalice/devola/core/logger"
	"github.com/qioalice/devola/core/sys/fn"
)

// Tusent is a data struct that contains a rule of sending some Telegram message.

// Tusent (to send) is a complex type of sending "reaction" (a set of actions)
// as "response" on the some occurred backend event.
//
// Backend depended context objects creates objects of this class, fills it fields
// and use it to sending responses and handle result of sending.
//
//
// So, backend depended context objects creates this object, fills it fields,
// construct internal objects, they will use directly which
type Tusent struct {

	// Pointer to main logger.
	//
	// Tusent does not depend on big parts of Devola SDK, but requires a logger object.
	// That's why the logger ptr passed directly, not a pointer to the some Devola SDK big part,
	// using which Tusent could get access to logger.
	// log *logger.Logger
	ctxer *ctxutils.CtxInfoer

	// ID of chat to which the message configs will be sent.
	chatID chat.ID

	// Determines behaviour of Tusent.
	flags flag

	// Untyped pointer to context object using which messages of that Tusent
	// were created and probably will be sent.
	// Also this pointer will be passed to the onSuccess/onError callbacks.
	ctxPtr unsafe.Pointer

	// Config is a created by backend the config of sendable message.
	// Will be send using backend depended method.
	Config interface{}
	// TODO: Probably refuse to use interface{} ?
	// (use unsafe.Pointer or [N]byte as core of some struct{} instead)

	// Pointer to the info of successfully sent message.
	// Have to be passed to the onSuccess callbacks as the second argument.
	sentSuccess unsafe.Pointer

	// Error of the sending message.
	// Have to be passed to the onError callbacks as the second argument.
	sentError error

	// A set of callbacks that will be called when message is successfully sent.
	onSuccess []fn.Named

	// A set of callbacks that will be called when message is unsent because of error.
	onError []fn.Named

	// TODO: comment
	RetryAttempts int8
}

// MakeSuccess makes ts a success-typed Tusent object and then returns it.
func (ts *Tusent) MakeSuccess(sentMsg unsafe.Pointer) *Tusent {
	ts.sentSuccess, ts.sentError = sentMsg, nil
	return ts
}

// MakeError makes ts an error-typed Tusent object and then returns it.
func (ts *Tusent) MakeError(err error) *Tusent {
	ts.sentSuccess, ts.sentError = nil, err
	return ts
}

// Call calls saved callbacks passing context object and object of sent msg
// or sending message error object to them.
//
// Optionally protect calls by panic guard and tries to finish transactions
// (depends on what flags were passed to the constructor).
func (ts *Tusent) Call() {
	for _, cb := range ts.onSuccess {
		ts.invoke(cb, true)
	}
	for _, cb := range ts.onError {
		ts.invoke(cb, false)
	}
	ts.finishTransactions()
}

// Protect tries to recover panic and if it is so writes the log message
// using log about it.
//
// WARNING!
// This method should be called with "defer" keyword, otherwise there is no-op.
func (ts *Tusent) recoverPanicOf(cb fn.Named) {

	err := recover()
	if err == nil {
		return
	}

	ts.ctxer.ML.Warn(
		`There was a restored panic in the user function.`,
		logger.KindAsField(logger.Core, logger.RecoveredPanic),

		zap.String("ctx", ts.ctxer.ViewFullJSON(ts.ctxPtr)),
		zap.String("fn_name", cb.Name),
		zap.Uint("fn_addr", uint(uintptr(cb.Ptr))),
		zap.Any("recovered_panic", err),
	)
}

// invoke safety (if panic guard is enabled) calls cb as func with 2 args.
//
// Untyped pointer to ctx will always be passed as 1st arg.
// isSuccess == true: 2nd arg is pointer to the object of sent message.
// isSuccess == false: 2nd arg is an error of sending message.
func (ts *Tusent) invoke(cb fn.Named, isSuccess bool) {

	if ts.flags.TestFlag(CEnablePanicGuard) {
		defer ts.recoverPanicOf(cb)
	}

	if isSuccess {
		cbTypedPtr := (*func(_, _ unsafe.Pointer))(cb.Ptr)
		(*cbTypedPtr)(ts.ctxPtr, ts.sentSuccess) // call

	} else {
		cbTypedPtr := (*func(unsafe.Pointer, error))(cb.Ptr)
		(*cbTypedPtr)(ts.ctxPtr, ts.sentError) // call
	}
}

// finishTransactions tries to finish open session and chat transactions if it is need.
//
// WARNING!
// If a session transaction wasn't finished,
// a chat transaction will also not be finished!
func (ts *Tusent) finishTransactions() {

	var err error

	if ts.flags.TestFlag(CFinishSessionTransaction) {
		goto finishChatTr
	}

	err = ts.ctxer.CompleteSessionTransaction(ts.ctxPtr)
	if err == nil {
		goto finishChatTr
	}

	ts.ctxer.ML.Warn(
		`Unable to finish session transaction.`,
		logger.KindAsField(logger.Core, logger.Transaction, logger.SessionTransaction),

		zap.String("ctx", ts.ctxer.ViewFullJSON(ts.ctxPtr)),
		zap.Error(err),
	)

	return

finishChatTr:

	if !ts.flags.TestFlag(CFinishChatTransaction) {
		return
	}

	err = ts.ctxer.CompleteChatTransaction(ts.ctxPtr)
	if err == nil {
		return
	}

	ts.ctxer.ML.Warn(
		`Unable to finish chat transaction.`,
		logger.KindAsField(logger.Core, logger.Transaction, logger.ChatTransaction),

		zap.String("ctx", ts.ctxer.ViewFullJSON(ts.ctxPtr)),
		zap.Error(err),
	)
}

// MakeTusent creates a new Tusent object using passed arguments.
// You should then specify type of Tusent using any of MakeSuccess, MakeError method.
func MakeTusent(flags flag, chatID chat.ID, onSuccess, onError []fn.Named, ctxPtr unsafe.Pointer) *Tusent {

	return &Tusent{
		chatID:    chatID,
		flags:     flags,
		ctxPtr:    ctxPtr,
		onSuccess: onSuccess,
		onError:   onError,
	}
}
