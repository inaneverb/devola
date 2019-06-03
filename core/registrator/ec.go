// Copyright © 2018. All rights reserved.
// Author: Alice Qio.
// Contacts: <qioalice@gmail.com>.
// License: https://opensource.org/licenses/MIT

package registrator

import (
	"github.com/qioalice/devola/core/errors"
)

// Predefined error codes of all registration operations.
//
// Registrator.Handler, Registrator.Middleware,
// Registrator.MainHandler, Registrator.MainMiddleware may return an error object
// EBadCallback, that implements errors.Error interface.
//
// So, method EBadCallback.Code returns one of these constant
// (or errors.ECOK if all is good).
const (

	// Bad type of callback that is registering as handler
	// (passed handler has incompatible type with context type).
	ECBadHandler errors.Code = 11

	// Bad type of callback that is registering as middleware.
	// (passed middleware has incompatible type with context type).
	ECBadMiddleware errors.Code = 12
)
