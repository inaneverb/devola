// Copyright © 2018. All rights reserved.
// Author: Alice Qio.
// Contacts: <qioalice@gmail.com>.
// License: https://opensource.org/licenses/MIT

package lirester

// lirchat represents a Lirester chat that contains three important things:
//
// - How much messages has been sent at this moment to that chat?
// - When this chat has been updated last time?
type lirchat struct {

	// NOTE.
	// Type named as lirchat (instead of just "chat") because of
	// avoiding collision with imported "chat" package.

	// WARNING!
	// ALL TIMESTAMPS IN NANO SECONDS!

	n           uint8
	lastUpdated int64
}
