// Copyright © 2019. All rights reserved.
// Author: Alice Qio.
// Contacts: <qioalice@gmail.com>.
// License: https://opensource.org/licenses/MIT

package ctxutils

import (
	"unsafe"

	"github.com/qioalice/devola/core/logger"
)

//
type CtxInfoer struct {
	ML *logger.Logger

	ViewShort ctxStringer

	ViewShortJSON ctxStringer

	ViewFull ctxStringer

	ViewFullJSON ctxStringer

	// Completors used to complete (finish, close) session or chat transactions
	// in Tusent objects after all callbacks has been called.
	CompleteSessionTransaction ctxTransactionFinisher

	CompleteChatTransaction ctxTransactionFinisher
}

// InitCompletors initializes transaction complete functions (completors).
func InitCompletors(cSessTr, cChatTr func(ctx unsafe.Pointer) error) {
	fCompletorSessionTransaction = cSessTr
	fCompletorChatTransaction = cChatTr
}
