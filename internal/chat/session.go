// Copyright © 2018. All rights reserved.
// Author: Alice Qio.
// Contacts: <qioalice@gmail.com>.
// License: https://opensource.org/licenses/MIT

package chat

import (
	"time"

	"../view"
)

// Session represents some dialogue between user (Telegram Chat) and
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
type Session struct {

	//
	ID SessionID `json:"id"`

	//
	ViewID view.ID `json:"view_id,omitempty"`

	//
	ViewIDEncoded view.IDEnc `json:"view_id_encoded"`

	//
	SentMessages MessageIDs `json:"sent_messages"`

	//
	ExpirationUnixstamp int64 `json:"expiration_unixstamp"`
}

//
func (s *Session) isEternal() bool {
	panic("implement me")
}

//
func (s *Session) isExpiredAt(unixstamp int64) bool {
	panic("implement me")
}

//
func (s *Session) isExpiredNow() bool {
	return s.isExpiredAt(time.Now().Unix())
}
