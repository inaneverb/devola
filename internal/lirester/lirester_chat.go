// Copyright © 2018. All rights reserved.
// Author: Alice Qio.
// Contacts: <qioalice@gmail.com>.
// License: https://opensource.org/licenses/MIT

package lirester

// chat represents a Lirester chat that contains three important things:
//
// - How much messages has been sent at this moment to that chat?
// - Is this chat with user or a channel/group/supergroup?
// - When this chat has been updated last time?
type chat struct {

	// First 7 bits represents the count of messages that have already been sent,
	// and the high bit is set if it is a channel/group/supergroup chat.
	data uint8

	// Timestamp when this chat was updated last time.
	lastUpdated int64
}

// howMuch returns the value of how many messages already sent to the chat.
func (c *chat) howMuch() uint8 {
	// High bit is indicator of group chat - ignore it
	return uint8(c.data) & 0x7F
}

// setHowMuch sets the counter to the v value in the current chat and
// returns changed value.
//
// WARNING!
// v physically must be in the range [-127..127] and logically in [0..127].
// Otherwise, v % 128 will be used as v.
//
// Deprecated: Unneccessary
func (c *chat) setHowMuch(v uint8) *chat {
	uint8(c.data) &= 0x80
	uint8(c.data) |= v & 0x7F
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
func (c *chat) incHowMuch(delta uint8) *chat {

	var nv = (uint8(c.data) & 0x7F) + delta
	if nv&0x80 == 0 {
		uint8(c.data) = (uint8(c.data) & 0x80) | nv
	}

	return c
}

// decHowMuch decreases a counter by the delta value in the current chat
// and returns a changed value.
//
// Uses the same limits and restrictions as incHowMuch method.
// More info: chat.incHowMuch notes, warnings.
func (c *chat) decHowMuch(delta uint8) *chat {
	return c.incHowMuch(-delta)
}

// isUser returns true only if the current chat is a chat with user.
func (c *chat) isUser() bool {
	// Zero high bit means that the current chan is just a chat with user.
	return uint8(c.data)&0x80 == 0
}

// isGroup returns true only if the current chat is a group/supergroup
// not chat with a user.
func (c *chat) isGroup() bool {
	return !c.isUser()
}

// setType changes type of chat (with a user or a group) by isUser flag
// in the current chat object.
func (c *chat) setType(isUser bool) *chat {
	uint8(c.data) &= 0x7F // cleanup prev value of flag
	// if it's not user, set high bit
	if !isUser {
		uint8(c.data) |= 0x80
	}
	return c
}

// setLastUpdated updates the last update time by now value.
func (c *chat) setLastUpdated(now int64) *chat {
	c.lastUpdated = now
	return c
}

// makeChat creates a new chat object with passed chat id.
//
// NOTE.
// After creating, specify whether created chat is a chat with user or a group chat
// using method chat.setType!
func makeChat() *chat {
	return &chat{}
}
