// Copyright © 2018. All rights reserved.
// Author: Alice Qio.
// Contacts: <qioalice@gmail.com>.
// License: https://opensource.org/licenses/MIT

package view

import (
	"unsafe"

	api "github.com/go-telegram-bot-api/telegram-bot-api"
)

// Finisher is the internal type that represents actions were occurred
// after some Telegram message has been sent or not sent using Telegram API.
//
// Finisher types.
// Finisher can be of two types: Success finisher or Error finisher.
// Success finisher (constructed by MakeFinisherSuccess) calls
// func(ctx *Ctx, msg *api.Message) callbacks, where msg is the object of sent message.
// Error finisher (constructed by MakeFinisherError) calls
// func (ctx *Ctx, err error) callbacks, where err is the reason why message was not sent.
//
// Panic guard.
// Finisher can protect calling of these callbacks by Panic Guard feature.
// In that case if any callback panicking, it will be recovered and stored into
// RecoveredPanics field.
// To enable that feature, use CEnablePanicGuard flag at the constructor.
//
// Transactions finishing.
// Finisher can finish session or chat transactions after all callbacks were called.
// Read more about in Ctx.TrSessionFinish, Session.TrFinish and
// Ctx.TrChatFinish, ChatInfo.TrFinish.
// To enable that feature, use CFinishSessionTransaction or/and CFinishChatTransaction.
// To figure out whether some error of finisher is occured, use TrSessionError, TrChatError methods.
type Finisher struct {

	// TODO: Add tests

	// Determines behaviour of finisher.
	// More info: finisherFlag.
	Flags finisherFlag

	// Untyped pointer to ORIGINAL context object using which message of that Finisher
	// is created and probably sent.
	// ALWAYS POINTS TO *Ctx EVEN IF CONTEXT IS EXTENDED!
	originalCtx unsafe.Pointer

	// Untyped pointer to context object using which message of that Finisher
	// is created and probably sent.
	passCallbacksCtx unsafe.Pointer

	// Data of successfully sent message.
	// Is not nil if this is a Success finisher.
	sentMsg *api.Message

	// A reason why message was not sent.
	// Is not nil if this is an Error finisher.
	sentErr error

	// Callbacks which should be called as finishing action.
	// Can be empty (when only transaction finishing is required, for example).
	callbacks []unsafe.Pointer

	// Slice of all recovered panics from callbacks.
	RecoveredPanics []interface{}

	// There is session transaction error or chat transaction error is placed.
	Err error
}

// Completors (finishers) that used to complete (finish, close) session or chat
// transactions in Finisher objects after all callbacks has been called.
// Once arument, ctx should be *Ctx, not a pointer to some extended context object!
var (
	fCompletorSessionTransaction func(ctx unsafe.Pointer) error
	fCompletorChatTransaction    func(ctx unsafe.Pointer) error
)

// InitCompletors initializes transaction complete functions (completors).
func InitCompletors(cSessTr, cChatTr func(ctx unsafe.Pointer) error) {
	fCompletorSessionTransaction = cSessTr
	fCompletorChatTransaction = cChatTr
}

// Call calls saved callbacks passing context object and object of sent msg
// or sending message error object to them.
// Optionally protect calls by panic guard and tries to finish transactions
// (depends on what flags were passed to the constructor).
func (f *Finisher) Call() {
	for _, cb := range f.callbacks {
		f.invoke(cb)
	}
	f.trFinish()
}

// TrSessionError returns an error object of finishing session transaction.
// It returns nil if that operation was not required.
func (f *Finisher) TrSessionError() error {
	if f.Flags.TestFlag(CIsSessionTransactionError) && f.Err != nil {
		return f.Err
	}
	return nil
}

// TrChatError returns an error object of finishing chat transaction.
// It returns nil if that operation was not required.
func (f *Finisher) TrChatError() error {
	if f.Flags.TestFlag(CIsChatTransactionError) && f.Err != nil {
		return f.Err
	}
	return nil
}

// protectFromPanic tries to recover panic, and if it was successfull,
// saves the recovered panic info to the panics field in current cb object
// to analyse it in the caller code.
func (f *Finisher) protectFromPanic() {
	if recoverInfo := recover(); recoverInfo != nil {
		f.RecoveredPanics = append(f.RecoveredPanics, recoverInfo)
	}
}

// invoke safety (if panic guard is enabled) calls cb,
// passing untyped pointer to ctx as 1st argument and object of sent message
// or sending message error (depends on which of them is not a nil).
func (f *Finisher) invoke(cb unsafe.Pointer) {

	if f.Flags.TestFlag(CEnablePanicGuard) {
		defer f.protectFromPanic()
	}

	switch {
	case f.sentMsg != nil:
		cbTypedPtr := (*func(unsafe.Pointer, *api.Message))(cb)
		(*cbTypedPtr)(f.passCallbacksCtx, f.sentMsg)

	case f.sentErr != nil:
		cbTypedPtr := (*func(unsafe.Pointer, error))(cb)
		(*cbTypedPtr)(f.passCallbacksCtx, f.sentErr)
	}
}

// trFinish tries to finish open session and chat transactions if it is need.
//
// WARNING!
// If a session transaction wasn't finished,
// a chat transaction will also not be finished!
func (f *Finisher) trFinish() {

	// Finish session transaction (if it's need)
	// Stop doing next things if error is occurred
	if f.Flags.TestFlag(CFinishSessionTransaction) {
		if err := fCompletorSessionTransaction(f.originalCtx); err != nil {
			f.Err = err
			f.Flags.SetFlag(CIsSessionTransactionError)
			return
		}
	}

	// Finish chat transaction (if it's need)
	if f.Flags.TestFlag(CFinishChatTransaction) {
		if err := fCompletorChatTransaction(f.originalCtx); err != nil {
			f.Err = err
			f.Flags.SetFlag(CIsChatTransactionError)
		}
	}
}

// MakeFinisherSuccess creates a new Success finisher using passed arguments.
func MakeFinisherSuccess(flags finisherFlag, cbs []unsafe.Pointer, originalCtx, pass2callbacksCtx unsafe.Pointer, sentMsg *api.Message) *Finisher {
	return &Finisher{
		Flags:            flags,
		originalCtx:      originalCtx,
		passCallbacksCtx: pass2callbacksCtx,
		sentMsg:          sentMsg,
	}
}

// MakeFinisherError creates a new Error finisher using passed arguments.
func MakeFinisherError(flags finisherFlag, cbs []unsafe.Pointer, originalCtx, pass2callbacksCtx unsafe.Pointer, err error) *Finisher {
	return &Finisher{
		Flags:            flags,
		originalCtx:      originalCtx,
		passCallbacksCtx: pass2callbacksCtx,
		sentErr:          err,
	}
}
