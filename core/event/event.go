// Copyright © 2019. All rights reserved.
// Author: Alice Qio.
// Contacts: <qioalice@gmail.com>.
// License: https://opensource.org/licenses/MIT

package event

// Event represents a meta info about occurred event.
// The Event object answers to the following questions:
// - What kind of event is occurred?
// - What data has been received with occurred event?
//
// Thus you can always to get info about what event you're handling
// inside your handler.
type Event struct {

	// The type of occurred event. Is one of predefined backend constants.
	Type Type `json:"type"`

	// The occurred event's data.
	Data Data `json:"data,omitempty"`
}

// String returns a string representation of event.
func (e *Event) String() string {
	if e == nil {
		return ""
	}
	return "Type: " + e.Type.String() + ", Data: \"" + string(e.Data) + "\""
}

// MakeEvent creates a new Event object with passed event type and event data.
func MakeEvent(typ Type, data Data) *Event {
	return &Event{
		Type: typ,
		Data: data,
	}
}
