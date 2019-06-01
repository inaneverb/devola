// Copyright © 2018. All rights reserved.
// Author: Alice Qio.
// Contacts: <qioalice@gmail.com>.
// License: https://opensource.org/licenses/MIT

package registrator

import (
	"github.com/qioalice/devola/core/event"
	"github.com/qioalice/devola/core/view"
)

// rule is a special type that describes a set of conditions to which
// some callback (handler or middleware) will be applied using Registrator type.
// Combined event description (type, data) and view ID conditions under which
// the callback that will be linked with, should be called.
type rule struct {

	// A base type.
	event.Event `json:",inline"`

	// A set of "current" ViewIDs when this registering event should be reacted.
	// If empty, registering event will be reacted anytime, but if it's not,
	// the registering event will be handled ONLY WHEN current session's View ID
	// is the same as any View ID from this field.
	When []view.ID `json:"when,omitempty"`
}

// makeRule creates a new rule object, that directly calls event.Event's
// constructor and saves when as When field.
func makeRule(typ event.Type, data event.Data, when []view.ID) *rule {
	var r rule
	r.Event.Type, r.Event.Data = typ, data
	r.When = when
	return &r
}
