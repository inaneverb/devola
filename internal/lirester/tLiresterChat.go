// Copyright © 2018. All rights reserved.
// Author: Alice Qio.
// Contacts: <qioalice@gmail.com>.
// License: https://opensource.org/licenses/MIT
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"),
// to deal in the Software without restriction, including without limitation
// the rights to use, copy, modify, merge, publish, distribute, sublicense,
// and/or sell copies of the Software, and to permit persons to whom
// the Software is furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS
// IN THE SOFTWARE.

package tgbot

// tLiresterChat represents a Lirester chat that contains three important things:
//
// - How much messages has been sent at this moment to that chat?
// - Is this chat with user or a channel/group/supergroup?
// - When this chat has been updated last time?
type tLiresterChat struct {

	// First 7 bits represents the count of messages that have already been sent,
	// and the high bit is set if it is a channel/group/supergroup chat.
	data uint8

	// Timestamp when this chat was updated last time.
	lastUpdated int64
}

// howMuch returns the value of how many messages already sent to the chat.
func (c *tLiresterChat) howMuch() uint8 { 

	// High bit is indicator of group chat - ignore it
	return uint8(c.data) & 0x7F 
}

// setHowMuch sets the counter to the 'v' value in the current chat and
// returns changed value.
// Warning! 'v' physically must be in the range [-127..127].
// Warning! 'v' logically must be in the range [0..127].
// Otherwise, 'v mod 128' will be used as 'v'.
// Deprecated: Unneccessary
func (c *tLiresterChat) setHowMuch(v uint8) *tLiresterChat {
	uint8(c.data) &= 0x80
	uint8(c.data) |= v & 0x7F
	return c
}

// incHowMuch increases a counter by the delta value in the current chat
// and returns a changed value.
// 
// NOTE!
// If you want to decrease a counter, just use decHowMuch method,
// or pass the negative delta to the current method.
//
// WARNING! 
// c.howMuch() + delta logically must be in the range [0..127] and
// physically (you can pass negative values) must be in the range [-127..127].
// Otherwise, there is no-op.
func (c *tLiresterChat) incHowMuch(delta uint8) *tLiresterChat {

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
// More info: tLiresterChat.incHowMuch notes, warnings.
func (c *tLiresterChat) decHowMuch(delta uint8) *tLiresterChat {

	return c.incHowMuch(-delta)
}

// isUser returns true only if the current chat is a chat with user.
func (c *tLiresterChat) isUser() bool { 

	// Zero high bit means that the current chan is just a chat with user.
	return uint8(c.data)&0x80 == 0 
}

// isGroup returns true only if the current chat is a group/supergroup
// not chat with a user.
func (c *tLiresterChat) isGroup() bool { 

	return !c.isUser()
}

// setType changes type of chat (with a user or a group) by isUser flag
// in the current chat object.
func (c *tLiresterChat) setType(isUser bool) *tLiresterChat {

	uint8(c.data) &= 0x7F // cleanup prev value of flag

	// if it's not user, set high bit
	if !isUser {
		uint8(c.data) |= 0x80
	} 

	return c
}

// setLastUpdated updates the last update time by now value.
func (c *tLiresterChat) setLastUpdated(now int64) *tLiresterChat {

	c.lastUpdated = now
	return c
}

// makeLiresterChat creates a new tLiresterChat object with passed chat id.
//
// NOTE! 
// After creating, specify whether created chat is a chat with user or a group chat
// using method tLiresterChat.setType!
func makeLiresterChat() *tLiresterChat {

	return &tLiresterChat{}
}