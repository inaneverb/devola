// Copyright © 2018. All rights reserved.
// Author: Alice Qio.
// Contacts: <qioalice@gmail.com>.
// License: https://opensource.org/licenses/MIT

package lirester

// lirchat represents a Lirester chat that contains three important things:
//
// - How much messages has been sent at this moment to that chat?
// - Is this chat with user or a channel/group/supergroup?
// - When this chat has been updated last time?
type lirchat struct {

	// NOTE.
	// Type named as lirchat (instead of just "chat") because of
	// avoiding collision with imported "chat" package.

	// First 7 bits represents the count of messages that have already been sent,
	// and the high bit is set if it is a channel/group/supergroup chat.
	data uint8

	// Timestamp when this chat was updated last time.
	lastUpdated int64
}

// howMuch returns the value of how many messages already sent to the chat.
func (c *lirchat) howMuch() uint8 {
	// High bit is indicator of group chat - ignore it
	return c.data & 0x7F
}

// setHowMuch sets the counter to the v value in the current chat and
// returns changed value.
//
// WARNING!
// v physically must be in the range [-127..127] and logically in [0..127].
// Otherwise, v % 128 will be used as v.
//
// Deprecated: Unneccessary
func (c *lirchat) setHowMuch(v uint8) *lirchat {
	c.data &= 0x80
	c.data |= v & 0x7F
	return c
}

// incHowMuch increases the counter by the delta value in the current chat
// and returns a changed value.
//
// NOTE.
// If you want to decrease a counter, just use decHowMuch method,
// or pass the negative delta to the current method.
//
// WARNING!
// c.howMuch() + delta logically must be in the range [0..127] and
// physically (you can pass negative values) must be in the range [-127..127].
// Otherwise, there is no-op.
func (c *lirchat) incHowMuch(delta uint8) *lirchat {

	var nv = (c.data & 0x7F) + delta
	if nv&0x80 == 0 {
		c.data = c.data & 0x80 | nv
	}

	return c
}

// decHowMuch decreases a counter by the delta value in the current chat
// and returns a changed value.
//
// Uses the same limits and restrictions as incHowMuch method.
// More info: chat.incHowMuch notes, warnings.
func (c *lirchat) decHowMuch(delta uint8) *lirchat {
	return c.incHowMuch(-delta)
}

// isUser returns true only if the current chat is a chat with user.
func (c *lirchat) isUser() bool {
	// Zero high bit means that the current chan is just a chat with user.
	return c.data & 0x80 == 0
}

// isGroup returns true only if the current chat is a group/supergroup
// not chat with a user.
func (c *lirchat) isGroup() bool {
	return !c.isUser()
}

// setType changes type of chat (with a user or a group) by isUser flag
// in the current chat object.
func (c *lirchat) setType(isUser bool) *lirchat {
	c.data &= 0x7F // cleanup prev value of flag
	// if it's not user, set high bit
	if !isUser {
		c.data |= 0x80
	}
	return c
}

// setLastUpdated updates the last update time by now value.
func (c *lirchat) setLastUpdated(now int64) *lirchat {
	c.lastUpdated = now
	return c
}

// makeChat creates a new chat object with passed chat id.
//
// NOTE.
// After creating, specify whether created chat is a chat with user or a group chat
// using method chat.setType!
func makeChat() *lirchat {
	return &lirchat{}
}
