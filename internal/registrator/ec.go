// Copyright © 2018. All rights reserved.
// Author: Alice Qio.
// Contacts: <qioalice@gmail.com>.
// License: https://opensource.org/licenses/MIT

package registrator

import (
	"../errors"
)

// Predefined error codes of all registration operations.
/////// These codes may be returned from IDConverter methods. TODO:
const (

	// Bad type of callback that is registering as handler.
	// Returned:
	// - From Registrator.Handler, Registrator.MainHandler methods
	//   if passed callback has type that incompatible with context type.
	ECBadHandler errors.Code = 11

	// Bad type of callback that is registering as middleware.
	// Returned:
	// - From Registrator.Middleware, Registrator.MainMiddleware methods
	//  if passed callback has type that incompatible with context type.
	ECBadMiddleware errors.Code = 12
)
