// Copyright © 2018. All rights reserved.
// Author: Alice Qio.
// Contacts: <qioalice@gmail.com>.
// License: https://opensource.org/licenses/MIT

package lirester

import (
	"../chat"
)

// cleanupConfig represents an object that contains two important things:
// - Chat id of chat, cleanup operation will perform over
// - Unixnano timestamp when cleanup operation will perform
//
// These rules are created by tLirester.Approve method and also passed
// to the its queue.
// Then, using method tLirester.Cleanup, these rules are applied.
type cleanupConfig struct {

	// chat id to which this cleanup rule will be applied
	chatID chat.ID

	// timestamp, when this cleanup rule should be applied
	when int64
}

// makeCleanupConfig creates a new cleanupConfig object with
// passed chat id and timestamp when created cleanup rule should be applied.
func makeCleanupConfig(chatID chat.ID, when int64) cleanupConfig {
	return cleanupConfig{
		chatID: chatID,
		when:   when,
	}
}
