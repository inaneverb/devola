// Copyright © 2018. All rights reserved.
// Author: Alice Qio.
// Contacts: <qioalice@gmail.com>.
// License: https://opensource.org/licenses/MIT

package view

// finisherFlag is the internal type to store constants of Finisher's behaviour.
// finisherFlag is the bitmask and each bit is meaning something
// (see constants, described below).
//
// More info: Finisher Ctx.
type finisherFlag uint8

// Predefined constants.
const (

	// Panic Guard.
	// Enabled Panic Guard allows panic in callbacks to be happened without
	// shutown the whole server - the panic will be recovered.
	CEnablePanicGuard finisherFlag = 0x01

	// Finish session transaction.
	// Enabling this option will trigger the session transaction completion process
	// after all callbacks has been called.
	//
	// More info: Ctx.TrSessionFinish, Session.TrFinish.
	CFinishSessionTransaction finisherFlag = 0x04

	// Finish chat transaction.
	// Enabling this option will trigger the chat transaction completion process
	// after all callbacks has been called.
	//
	// More info: Ctx.TrChatFinish, ChatInfo.TrFinish.
	CFinishChatTransaction finisherFlag = 0x08

	// If this bit is set it means that an error occurred while trying to finish
	// session transaction.
	// Error object will be stored in Finisher.Err field.
	CIsSessionTransactionError finisherFlag = 0x10

	// If this bit is set it means that an error occurred while trying to finish
	// chat transaction.
	// Error object will be stored in Finisher.Err field.
	CIsChatTransactionError finisherFlag = 0x20

	// free HEX values:
	// 0x02, 0x40, 0x80
)

// SetFlag sets each bit in f that set in flag.
func (f *finisherFlag) SetFlag(flag finisherFlag) {
	*f |= flag
}

// TestFlag returns true only if ALL bits in flag are set in f.
func (f *finisherFlag) TestFlag(flag finisherFlag) bool {
	return *f&flag == flag
}
