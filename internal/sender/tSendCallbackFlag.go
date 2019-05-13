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

// tSendCallbackFlag is the internal type to store constants of behaviour of
// tViewBaseFinisher objects and all derived types.
//
// tSendCallbackFlag is the bitmask and each bit is meaning something
// (see constants, described below).
//
// More info: FViewSuccessFinisher, FViewErrorFinisher,
// tViewBaseFinisher, tViewSuccessFinisher, tViewErrorFinisher, TCtx.
type tSendCallbackFlag uint8

// Predefined constants.
const (

	// Panic Guard.
	// Set this flag if you want to enable panic guard in invoke methods.
	//
	// Panic Guard allows panic in callbacks to be happened without
	// shutown the whole server - the panic will be recovered.
	//
	// More info: tViewBaseFinisher.invoke, tViewSuccessFinisher.invoke,
	// tViewErrorFinisher.invoke.
	cSendCallbackEnablePanicGuard tSendCallbackFlag = 0x01

	// Use extended context.
	// Set this flag if you want to use callbacks that receives extended
	// context object instead standard.
	//
	// TBot.ExtendContext allows to extend standard context TCtx type
	// by user defined and use it as context.
	//
	// Finishers can also receive a context of extended type but
	// it should be enabled by this flag.
	cSendCallbackUseExtendedContext tSendCallbackFlag = 0x02

	// Finish session transaction.
	// Set this flag if you want a session transaction to be finished
	// after all callbacks will be executed.
	//
	// More info: TCtx.TrSessionFinish, tSession.TrFinish.
	cSendCallbackFinishSessionTransaction tSendCallbackFlag = 0x04

	// Finish chat transaction.
	// Set this flag if you want a chat transaction to be finished
	// afger all callbacks will be executed.
	//
	// More info: TCtx.TrChatFinish, tChatInfo.TrFinish.
	cSendCallbackFinishChatTransaction tSendCallbackFlag = 0x08

	// An error occurred while trying to finish session transaction.
	// Error object will be stored in afterAllErr field of tViewBaseFinisher.
	cSendCallbackFinishSessionTransactionError tSendCallbackFlag = 0x10

	// An error occurred while trying to finish chat transaction.
	// Error object will be stored in afterAllErr field of tViewBaseFinisher.
	cSendCallbackFinishChatTransactionError tSendCallbackFlag = 0x20
)

// SetFlag sets all each bit in f that set in flag.
func (f *tSendCallbackFlag) SetFlag(flag tSendCallbackFlag) {
	*f |= flag
}

// TestFlag returns true only if ALL bits in flag are set in f.
func (f *tSendCallbackFlag) TestFlag(flag tSendCallbackFlag) bool {
	return *f&flag == flag
}
