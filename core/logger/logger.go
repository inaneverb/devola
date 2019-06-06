// Copyright © 2019. All rights reserved.
// Author: Alice Qio.
// Contacts: <qioalice@gmail.com>.
// License: https://opensource.org/licenses/MIT

package logger

import (
	"go.uber.org/zap"
)

type Logger struct {
	*zap.Logger
}
