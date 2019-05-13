// Copyright Â© 2018. All rights reserved.
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

import (
	"time"
)

// tSession represents some dialogue between user (Telegram Chat) and
// server (Telegram Bot).
//
// Moreover session allows to represent the dialogue as the "wizard" of
// some business entity - a step-oriented dialogue, when each next step
// is edited (or deleted) message about prev step.
//
// Session object keeps in itself internal or external information
// about one continious action using bot.
//
// More info: TCtx, tChatInfoManager, tSender, tReceiver.
type tSession struct {
	
	//
	ID tSessionID `json:"id"`
	
	//
	ViewID tViewID `json:"view_id,omitempty"`
	
	//
	ViewIDEncoded tViewIDEncoded `json:"view_id_encoded"`
	
	//
	SentMessages tChatMessageIDs `json:"sent_messages"`
	
	//
	ExpirationUnixstamp int64 `json:"expiration_unixstamp"`
}

//
func (s *tSession) isEternal() bool {
	
}

//
func (s *tSession) isExpiredAt(unixstamp int64) bool {
	
}

//
func (s *tSession) isExpiredNow() bool {
	return s.isExpiredAt(time.Now().Unix())
}