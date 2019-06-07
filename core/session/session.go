// Copyright © 2019. All rights reserved.
// Author: Alice Qio.
// Contacts: <qioalice@gmail.com>.
// License: https://opensource.org/licenses/MIT

package session

import (
	"time"

	"github.com/qioalice/devola/core/chat"
	"github.com/qioalice/devola/core/view"
)

// Session represents some dialogue between user and server.
//
// Moreover session allows to represent the dialogue as the "wizard" of
// some business entity - a step-oriented dialogue, when each next step
// is edited (or deleted) message about prev step.
//
// Session object keeps in itself internal or external information
// about one continuous action using bot.
type Session struct {

	//
	ID SessionID `json:"id"`

	//
	ViewID view.ID `json:"view_id,omitempty"`

	//
	ViewIDEncoded view.IDEnc `json:"view_id_encoded"`

	//
	SentMessages chat.MessageIDs `json:"sent_messages"`

	//
	ExpirationUnixstamp int64 `json:"expiration_unixstamp"`
}

// isEternal returns true if s is eternal (infinity) session.
func (s *Session) isEternal() bool {
	return s.ExpirationUnixstamp == -1
}

// isExpiredAt returns true if s will be expired at the stamp (should be unix timestamp).
func (s *Session) isExpiredAt(stamp int64) bool {
	return !s.isEternal() && s.ExpirationUnixstamp <= stamp
}

// isExpired returns true if s is already expired.
func (s *Session) isExpired() bool {
	return s.isExpiredAt(time.Now().Unix())
}
