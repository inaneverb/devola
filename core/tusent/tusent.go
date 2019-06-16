// Copyright © 2018. All rights reserved.
// Author: Alice Qio.
// Contacts: <qioalice@gmail.com>.
// License: https://opensource.org/licenses/MIT

package tusent

import (
	"unsafe"

	"go.uber.org/zap"

	"github.com/qioalice/devola/core/chat"
	"github.com/qioalice/devola/core/errors"
	"github.com/qioalice/devola/core/logger"
	"github.com/qioalice/devola/core/sys/flag"
	"github.com/qioalice/devola/core/sys/fn"

	"github.com/qioalice/devola/modules/brdige"
)

// Tusent (to send) is a complex type of sending "reaction" (a set of actions)
// as "response" on the some occurred backend event.
//
// Backend depended context object creates objects of this class, fills it fields
// and use it to sending responses and handle result of sending.
type Tusent struct {

	// HOW IT WORKS.
	//
	// It works in tandem with core.sender/Sender, but Sender is leading.
	// Almost all fields are public, because Sender must have access to them,
	// but there are some methods that should be used
	// instead of direct access to the fields.
	//
	// This object can be in one of two states:
	// - Send operation (depended on the backend) has not yet been performed
	//   (and this object is in any of
	//   Sender.undisposedResponses or Sender.preparedResponses field).
	// - Finish actions has not yet been performed
	//   (and this object is in Sender.handledResponses field).
	//
	// When state is 1, it waits for send operation (depended on the backend)
	// to be performed while executing of which the state changes from 1 to 2
	// (and objects goes to Sender.handledResponses).
	// Sender performs send operation using direct access to the fields.
	//
	// When state is 2, it waits for calling OnSuccess or OnError callbacks
	// (depended on the result of send operation) and finishing chat's or session's
	// transactions.
	// Sender calls callbacks and finishes transactions using Call method
	// (for both of these operations).
	//
	// When response(s) is/are sent, callback(s) is/are called and
	// transaction(s) is/are finished, the purpose of this object is completed
	// and it will be GC'ed.
	//
	// Logging and finishing transaction uses core/bridge.Bridge module to call
	// backend depended functions.

	bridge *brdige.Bridge
	flags  flag.F8

	// ChatIDT is combined backend chat's ID and chat's type.
	ChatIDT chat.IDT

	// Ctx is an untyped pointer to context object using which messages
	// of that Tusent were created and probably will be sent.
	//
	// Also this pointer will be passed to the OnSuccess/OnError callbacks.
	Ctx unsafe.Pointer

	// Config is a created config of sendable message by backend.
	// Will be send using backend depended method.
	Config interface{}
	// TODO: Probably refuse to use interface{} ?
	// (use unsafe.Pointer or [N]byte as core of some struct{} instead)

	// SentObj is an untyped pointer to the info of successfully sent message.
	// Have to be passed to the OnSuccess callbacks as the second argument.
	SentObj unsafe.Pointer

	// SendingErr is an error of sending message.
	// Have to be passed to the OnError callbacks as the second argument.
	SendingErr error

	// OnSuccess is a set of callbacks that will be called
	// when message is successfully sent.
	OnSuccess []fn.Named

	// OnError is a set of callbacks that will be called
	// when message is unsent because of error.
	OnError []fn.Named

	// RetryAttempts is the counter of additional attempts of sending message.
	// It's in [Sender.consts.retryAttemptsInfMax .. max(int8)] / {0}.
	//
	// Negative value means that there is an "infinity" number of attempts.
	// (see Sender.consts.retryAttemptsInfMax docs).
	//
	// This counter is decremented each time an error occurred sending a message
	// and until it reached 0.
	// When it reaches zero, the message is considered completely unsent.
	//
	// WARNING!
	// OnError callbacks are not called until this counter becomes 0!
	RetryAttempts int8
}

// Predefined flags that determines the behaviour of Tusent.
const (

	// CEnablePanicGuard enables Panic Guard that allows callback panicking
	// without shutdown the whole server because the panic will be recovered.
	CEnablePanicGuard flag.F8 = 0x01

	// CFinishSessionTransaction will trigger the session transaction completion
	// process after all callbacks has been called.
	CFinishSessionTransaction flag.F8 = 0x04

	// CFinishChatTransaction will trigger the chat transaction completion
	// process after all callbacks has been called.
	CFinishChatTransaction flag.F8 = 0x08

	// free HEX values:
	// 0x02, 0x40, 0x80
)

// MakeSuccess makes ts a success-typed Tusent object and then returns it.
func (ts *Tusent) MakeSuccess(sentMsg unsafe.Pointer) *Tusent {
	ts.SentObj, ts.SendingErr = sentMsg, nil
	return ts
}

// MakeError makes ts an error-typed Tusent object and then returns it.
func (ts *Tusent) MakeError(err error) *Tusent {

	if ts.SentObj = nil; ts.SendingErr == nil {
		ts.SendingErr = err

	} else if ts.SendingErr != err {
		if errsTyped, ok := ts.SendingErr.(errors.Set); ok {
			for _, errAlready := range errsTyped {
				if errAlready == err {
					return ts
				}
			}
			ts.SendingErr = append(errsTyped, err)
		}
	}
	return ts
}

// NeedToFinish returns true if current Tusent object should to close (finish)
// chat's or session's transaction of self.
func (ts *Tusent) NeedToFinish() bool {
	return ts.flags.AnyFlag(CFinishChatTransaction | CFinishSessionTransaction)
}

// Call calls saved callbacks passing context object and object of sent msg
// or sending message error object to them.
//
// Optionally protect calls by panic guard and tries to finish transactions
// (depends on what flags were passed to the constructor).
func (ts *Tusent) Call() {
	for _, cb := range ts.OnSuccess {
		ts.invoke(cb, true)
	}
	for _, cb := range ts.OnError {
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

	ts.bridge.ML.Warn(
		"There was a restored panic in the user function.",
		logger.KindAsField(logger.Core, logger.RecoveredPanic),

		zap.String("ctx", ts.bridge.CtxView(ts.Ctx, true, true)),
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
		(*cbTypedPtr)(ts.Ctx, ts.SentObj) // call

	} else {
		cbTypedPtr := (*func(unsafe.Pointer, error))(cb.Ptr)
		(*cbTypedPtr)(ts.Ctx, ts.SendingErr) // call
	}
}

// finishTransactions tries to finish open session and chat transactions if it is need.
//
// WARNING!
// If a session transaction wasn't finished,
// a chat transaction will also not be finished!
func (ts *Tusent) finishTransactions() {

	var err error

	if !ts.flags.TestFlag(CFinishSessionTransaction) {
		goto finishChatTr
	}

	err = ts.bridge.FinishTr(ts.Ctx, true)
	if err == nil {
		goto finishChatTr
	}

	ts.bridge.ML.Warn(
		"Unable to finish session transaction.",
		logger.KindAsField(logger.Core, logger.Transaction, logger.SessionTransaction),

		zap.String("ctx", ts.bridge.CtxView(ts.Ctx, true, true)),
		zap.Error(err),
	)

	return

finishChatTr:

	if !ts.flags.TestFlag(CFinishChatTransaction) {
		return
	}

	err = ts.bridge.FinishTr(ts.Ctx, false)
	if err == nil {
		return
	}

	ts.bridge.ML.Warn(
		"Unable to finish chat transaction.",
		logger.KindAsField(logger.Core, logger.Transaction, logger.ChatTransaction),

		zap.String("ctx", ts.bridge.CtxView(ts.Ctx, true, true)),
		zap.Error(err),
	)
}

// New creates a new Tusent object using passed arguments.
// You should then specify type of Tusent using any of MakeSuccess, MakeError method.
func New(
	bridge *brdige.Bridge,
	flags flag.F8,
	chatID chat.IDT,
	onSuccess, onError []fn.Named,
	ctxPtr unsafe.Pointer,
) *Tusent {
	return &Tusent{
		bridge:    bridge,
		flags:     flags,
		ChatIDT:   chatID,
		Ctx:       ctxPtr,
		OnSuccess: onSuccess,
		OnError:   onError,
	}
}
