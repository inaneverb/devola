// Copyright © 2018. All rights reserved.
// Author: Alice Qio.
// Contacts: <qioalice@gmail.com>.
// License: https://opensource.org/licenses/MIT

package ctx

import (
	"github.com/qioalice/devola/core/event"
	"github.com/qioalice/devola/core/session"
)

// BaseCtx is a part of Devola SDK core that is an abstract over absolutely
// all backend-depended context types.
type BaseCtx struct {

	// In Devola SDK each backend can declare its own context type.
	// Objects of that context type are created when events are occurred
	// and contains all important info about these events.
	//
	// MOST IMPORTANTLY,
	// ABSOLUTELY EACH BACKEND-DEPENDED CONTEXT TYPE
	// MUST BE LIKE THIS TYPE!
	//
	// It means that if we have T as backend-depended context type,
	// and v as *T object,
	// the ASSIGN WITH CASTING like
	//
	//   var v *BaseCtx = (*BaseCtx)(unsafe.Pointer(v))
	//
	// MUST BE ABSOLUTELY LEGIT AND NORMAL, and moreover,
	// the ACCESS to any of BaseCtx FIELDS using dereference of v like
	//
	//   eventType, sessID := v.Event.Type, v.Session.ID
	//
	// MUST BE ABSOLUTELY LEGIT AND NORMAL TOO!
	//
	// THE THINGS DESCRIBED ABOVE IS THE BASE AXIOMS AND PRINCIPLES
	// OF DEVOLA SDK AND ALL DEVOLA BACKENDS MUST FOLLOW THEM!
	//
	// You can achieve the behaviour described above by following one of these
	// two ways:
	//
	// 1. Embed BaseCtx (not a *BaseCtx, it's important!)
	//    to your backend-depended context type.
	//
	// 2. Your backend-depended context type must have:
	//    - event.Event as 1st embedded type,
	//    - session.Session as 2nd embedded type.

	Event   event.Event
	Session session.Session
}
