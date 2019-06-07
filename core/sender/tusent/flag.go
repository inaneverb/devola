// Copyright © 2019. All rights reserved.
// Author: Alice Qio.
// Contacts: <qioalice@gmail.com>.
// License: https://opensource.org/licenses/MIT

package tusent

// flag is the internal type to store constants of Tusent's behaviour.
// flag is a bitmask and each bit is meaning something
// (see constants, described below).
//
// More info: Finisher.
type flag uint8

// Predefined constants.
const (

	// Panic Guard.
	// Enabled Panic Guard allows panic in callbacks to be happened without
	// shutdown the whole server - the panic will be recovered.
	CEnablePanicGuard flag = 0x01

	// Finish session transaction.
	// Enabling this option will trigger the session transaction completion process
	// after all callbacks has been called.
	CFinishSessionTransaction flag = 0x04

	// Finish chat transaction.
	// Enabling this option will trigger the chat transaction completion process
	// after all callbacks has been called.
	CFinishChatTransaction flag = 0x08

	// If this bit is set it means that an error occurred while trying to finish
	// session transaction.
	// Error object will be stored in Finisher.Err field.
	CIsSessionTransactionError flag = 0x10

	// If this bit is set it means that an error occurred while trying to finish
	// chat transaction.
	// Error object will be stored in Finisher.Err field.
	CIsChatTransactionError flag = 0x20

	// free HEX values:
	// 0x02, 0x40, 0x80
)

// SetFlag sets each bit in f that set in flag.
func (f *flag) SetFlag(flag flag) {
	*f |= flag
}

// TestFlag returns true only if ALL bits in flag are set in f.
func (f *flag) TestFlag(flag flag) bool {
	return *f&flag == flag
}
