// Copyright © 2019. All rights reserved.
// Author: Alice Qio.
// Contacts: <qioalice@gmail.com>.
// License: https://opensource.org/licenses/MIT

package bridge

import (
	"unsafe"

	"go.uber.org/zap"

	"github.com/qioalice/devola/core/logger"
)

//
type Bridge struct {

	// Main Logger
	ML *logger.Logger

	// CtxNext
	CtxNext func() (ctx unsafe.Pointer)

	// CtxView is an alias to function that takes backend context object
	// and returns its string representation.
	CtxView func(ctx unsafe.Pointer, isFull, asJSON bool) string

	// Completors used to complete (finish, close) session or chat transactions
	// in Tusent objects after all callbacks has been called.
	FinishTr func(ctx unsafe.Pointer, isSessionTr bool) error

	//
	DoSend func(obj interface{}) (res unsafe.Pointer, err error, isFinalErr bool)

	//
	SendOK func(ctx, res unsafe.Pointer)

	//
	SendErr func(ctx unsafe.Pointer, err error)

	// 		// it's prohibited to be a literally "infinity" number of retrying attempts,
	// 		// because of that when a negative decreasing counter will reach its max,
	// 		// we also
	SendInfOverflow func(ctx unsafe.Pointer, err error, config interface{})
}

//
func (b *Bridge) RecoverPanicOf(ctx unsafe.Pointer, fTyp Kind, fPtr unsafe.Pointer, fName string, addInfo interface{}) {

	err := recover()
	if err == nil {
		return
	}

	if fName == "" {
		fName = "Unnamed"
	}

	b.ML.With()

	b.ML.Warn(
		"There was a restored panic in the user function.",
		logger.KindAsField(logger.Core, logger.RecoveredPanic),

		zap.String("ctx", b.CtxView(ctx, true, true)),
		zap.String("fn_name", fName),
		zap.String("fn_kind", fTyp.String()),
		zap.Uintptr("fn_addr", uintptr(fPtr)),
		zap.Any("recovered_panic", err),
		zap.Any("add_info", addInfo),
	)
}
