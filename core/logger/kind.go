// Copyright © 2019. All rights reserved.
// Author: Alice Qio.
// Contacts: <qioalice@gmail.com>.
// License: https://opensource.org/licenses/MIT

package logger

import (
	"go.uber.org/zap"
)

//
type kind uint16

//
const (
	Core kind = 1 + iota
	Backend

	RecoveredPanic = 100

	Transaction        = 200
	ChatTransaction    = 201
	SessionTransaction = 202
)

//
func KindAsField(kinds ...kind) zap.Field {

}
